package service

import (
	"context"
	"net/http"
	"sort"
	"time"

	"go-admin/internal/config"
	"go-admin/internal/modules/auth/model"
	"go-admin/internal/modules/auth/repo"
	"go-admin/pkg/cachex"
	"go-admin/pkg/crypto/hash"
	"go-admin/pkg/errors"
	"go-admin/pkg/jwtx"
	"go-admin/pkg/logging"
	"go-admin/pkg/util"

	"github.com/LyricTian/captcha"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// Login management for RBAC
type Auth struct {
	Cache        cachex.Cacher
	Auth         jwtx.Auther
	UserRepo     *repo.User
	UserRoleRepo *repo.UserRole
	MenuRepo     *repo.Menu
	UserService  *User
	Trans        *util.Trans
}

func (a *Auth) ParseUserID(c *gin.Context) (string, error) {
	rootID := config.C.General.Root.ID
	if config.C.Middleware.Auth.Disable {
		return rootID, nil
	}

	invalidToken := errors.Unauthorized(config.ErrInvalidTokenID, "Invalid access token")
	token := util.GetToken(c)
	if token == "" {
		return "", invalidToken
	}

	ctx := c.Request.Context()
	ctx = util.NewUserToken(ctx, token)

	userID, err := a.Auth.ParseSubject(ctx, token)
	if err != nil {
		if err == jwtx.ErrInvalidToken {
			return "", invalidToken
		}
		return "", err
	} else if userID == rootID {
		c.Request = c.Request.WithContext(util.NewIsRootUser(ctx))
		return userID, nil
	}

	userCacheVal, ok, err := a.Cache.Get(ctx, config.CacheNSForUser, userID)
	if err != nil {
		return "", err
	} else if ok {
		userCache := util.ParseUserCache(userCacheVal)
		c.Request = c.Request.WithContext(util.NewUserCache(ctx, userCache))
		return userID, nil
	}

	// Check user status, if not activated, force to logout
	user, err := a.UserRepo.Get(ctx, userID, model.UserQueryOptions{
		QueryOptions: util.QueryOptions{SelectFields: []string{"status"}},
	})
	if err != nil {
		return "", err
	} else if user == nil || user.Status != model.UserStatusActivated {
		return "", invalidToken
	}

	roleIDs, err := a.UserService.GetRoleIDs(ctx, userID)
	if err != nil {
		return "", err
	}

	userCache := util.UserCache{
		RoleIDs: roleIDs,
	}
	err = a.Cache.Set(ctx, config.CacheNSForUser, userID, userCache.String())
	if err != nil {
		return "", err
	}

	c.Request = c.Request.WithContext(util.NewUserCache(ctx, userCache))
	return userID, nil
}

// This function generates a new captcha ID and returns it as a `model.Captcha` struct. The length of
// the captcha is determined by the `config.C.Util.Captcha.Length` configuration value.
func (a *Auth) GetCaptcha(ctx context.Context) (*model.Captcha, error) {
	return &model.Captcha{
		CaptchaID: captcha.NewLen(config.C.Util.Captcha.Length),
	}, nil
}

// Response captcha image
func (a *Auth) ResponseCaptcha(ctx context.Context, w http.ResponseWriter, id string, reload bool) error {
	if reload && !captcha.Reload(id) {
		return errors.NotFound("", "Captcha id not found")
	}

	err := captcha.WriteImage(w, id, config.C.Util.Captcha.Width, config.C.Util.Captcha.Height)
	if err != nil {
		if err == captcha.ErrNotFound {
			return errors.NotFound("", "Captcha id not found")
		}
		return err
	}

	w.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
	w.Header().Set("Pragma", "no-cache")
	w.Header().Set("Expires", "0")
	w.Header().Set("Content-Type", "image/png")
	return nil
}

func (a *Auth) genUserToken(ctx context.Context, userID string) (*model.LoginToken, error) {
	token, err := a.Auth.GenerateToken(ctx, userID)
	if err != nil {
		return nil, err
	}

	tokenBuf, err := token.EncodeToJSON()
	if err != nil {
		return nil, err
	}
	logging.Context(ctx).Info("Generate user token", zap.Any("token", string(tokenBuf)))

	return &model.LoginToken{
		AccessToken: token.GetAccessToken(),
		TokenType:   token.GetTokenType(),
		ExpiresAt:   token.GetExpiresAt(),
	}, nil
}

func (a *Auth) Login(ctx context.Context, formItem *model.LoginForm) (*model.LoginToken, error) {
	// verify captcha
	if !captcha.VerifyString(formItem.CaptchaID, formItem.CaptchaCode) {
		return nil, errors.BadRequest(config.ErrInvalidCaptchaID, "Incorrect captcha")
	}

	ctx = logging.NewTag(ctx, logging.TagKeyLogin)

	// login by root
	if formItem.Email == config.C.General.Root.Email {
		if formItem.Password != config.C.General.Root.Password {
			return nil, errors.BadRequest(config.ErrInvalidUsernameOrPassword, "Incorrect username or password")
		}

		userID := config.C.General.Root.ID
		ctx = logging.NewUserID(ctx, userID)
		logging.Context(ctx).Info("Login by root")
		return a.genUserToken(ctx, userID)
	}

	// get user info
	user, err := a.UserRepo.GetByEmail(ctx, formItem.Email, model.UserQueryOptions{
		QueryOptions: util.QueryOptions{
			SelectFields: []string{"id", "password", "status"},
		},
	})
	if err != nil {
		return nil, err
	} else if user == nil {
		return nil, errors.BadRequest(config.ErrInvalidUsernameOrPassword, "Incorrect username or password")
	} else if user.Status != model.UserStatusActivated {
		return nil, errors.BadRequest("", "User status is not activated, please contact the administrator")
	}

	// check password
	if err := hash.CompareHashAndPassword(user.Password, formItem.Password); err != nil {
		return nil, errors.BadRequest(config.ErrInvalidUsernameOrPassword, "Incorrect username or password")
	}

	userID := user.ID
	ctx = logging.NewUserID(ctx, userID)

	// set user cache with role ids
	roleIDs, err := a.UserService.GetRoleIDs(ctx, userID)
	if err != nil {
		return nil, err
	}

	userCache := util.UserCache{RoleIDs: roleIDs}
	err = a.Cache.Set(ctx, config.CacheNSForUser, userID, userCache.String(),
		time.Duration(config.C.Dictionary.UserCacheExp)*time.Hour)
	if err != nil {
		logging.Context(ctx).Error("Failed to set cache", zap.Error(err))
	}
	logging.Context(ctx).Info("Login success", zap.String("email", formItem.Email))

	// generate token
	return a.genUserToken(ctx, userID)
}

func (a *Auth) Register(ctx context.Context, formItem *model.RegisterForm) error {
	// verify captcha
	if !captcha.VerifyString(formItem.CaptchaID, formItem.CaptchaCode) {
		return errors.BadRequest(config.ErrInvalidCaptchaID, "Incorrect captcha")
	}
	ctx = logging.NewTag(ctx, logging.TagKeyRegister)
	existsEmail, err := a.UserRepo.ExistsEmail(ctx, formItem.Email)
	if err != nil {
		return err
	} else if existsEmail {
		return errors.BadRequest("", "Email already exists")
	}
	user := &model.User{
		ID:        util.NewXID(),
		CreatedAt: time.Now(),
	}
	if err := formItem.FillTo(user); err != nil {
		return err
	}
	err = a.Trans.Exec(ctx, func(ctx context.Context) error {
		if err := a.UserRepo.Create(ctx, user); err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return err
	}
	ctx = logging.NewUserID(ctx, user.ID)
	logging.Context(ctx).Info("Register success", zap.String("email", formItem.Email))
	return nil
}

func (a *Auth) RefreshToken(ctx context.Context) (*model.LoginToken, error) {
	userID := util.FromUserID(ctx)

	user, err := a.UserRepo.Get(ctx, userID, model.UserQueryOptions{
		QueryOptions: util.QueryOptions{
			SelectFields: []string{"status"},
		},
	})
	if err != nil {
		return nil, err
	} else if user == nil {
		return nil, errors.BadRequest("", "Incorrect user")
	} else if user.Status != model.UserStatusActivated {
		return nil, errors.BadRequest("", "User status is not activated, please contact the administrator")
	}

	return a.genUserToken(ctx, userID)
}

func (a *Auth) Logout(ctx context.Context) error {
	userToken := util.FromUserToken(ctx)
	if userToken == "" {
		return nil
	}

	ctx = logging.NewTag(ctx, logging.TagKeyLogout)
	if err := a.Auth.DestroyToken(ctx, userToken); err != nil {
		return err
	}

	userID := util.FromUserID(ctx)
	err := a.Cache.Delete(ctx, config.CacheNSForUser, userID)
	if err != nil {
		logging.Context(ctx).Error("Failed to delete user cache", zap.Error(err))
	}
	logging.Context(ctx).Info("Logout success")

	return nil
}

// Get user info
func (a *Auth) GetUserInfo(ctx context.Context) (*model.User, error) {
	if util.FromIsRootUser(ctx) {
		return &model.User{
			ID:        config.C.General.Root.ID,
			Email:     config.C.General.Root.Email,
			FirstName: config.C.General.Root.FirstName,
			LastName:  config.C.General.Root.LastName,
			Status:    model.UserStatusActivated,
		}, nil
	}

	userID := util.FromUserID(ctx)
	user, err := a.UserRepo.Get(ctx, userID, model.UserQueryOptions{
		QueryOptions: util.QueryOptions{
			OmitFields: []string{"password"},
		},
	})
	if err != nil {
		return nil, err
	} else if user == nil {
		return nil, errors.NotFound("", "User not found")
	}

	userRoleResult, err := a.UserRoleRepo.Query(ctx, model.UserRoleQueryParam{
		UserID: userID,
	}, model.UserRoleQueryOptions{
		JoinRole: true,
	})
	if err != nil {
		return nil, err
	}
	user.Roles = userRoleResult.Data

	return user, nil
}

// Change login password
func (a *Auth) UpdatePassword(ctx context.Context, updateItem *model.UpdateLoginPassword) error {
	if util.FromIsRootUser(ctx) {
		return errors.BadRequest("", "Root user cannot change password")
	}

	userID := util.FromUserID(ctx)
	user, err := a.UserRepo.Get(ctx, userID, model.UserQueryOptions{
		QueryOptions: util.QueryOptions{
			SelectFields: []string{"password"},
		},
	})
	if err != nil {
		return err
	} else if user == nil {
		return errors.NotFound("", "User not found")
	}

	// check old password
	if err := hash.CompareHashAndPassword(user.Password, updateItem.OldPassword); err != nil {
		return errors.BadRequest("", "Incorrect old password")
	}

	// update password
	newPassword, err := hash.GeneratePassword(updateItem.NewPassword)
	if err != nil {
		return err
	}
	return a.UserRepo.UpdatePasswordByID(ctx, userID, newPassword)
}

// Query menus based on user permissions
func (a *Auth) QueryMenus(ctx context.Context) (model.Menus, error) {
	menuQueryParams := model.MenuQueryParam{
		Status: model.MenuStatusEnabled,
	}

	isRoot := util.FromIsRootUser(ctx)
	if !isRoot {
		menuQueryParams.UserID = util.FromUserID(ctx)
	}
	menuResult, err := a.MenuRepo.Query(ctx, menuQueryParams, model.MenuQueryOptions{
		QueryOptions: util.QueryOptions{
			OrderFields: model.MenusOrderParams,
		},
	})
	if err != nil {
		return nil, err
	} else if isRoot {
		return menuResult.Data.ToTree(), nil
	}

	// fill parent menus
	if parentIDs := menuResult.Data.SplitParentIDs(); len(parentIDs) > 0 {
		var missMenusIDs []string
		menuIDMapper := menuResult.Data.ToMap()
		for _, parentID := range parentIDs {
			if _, ok := menuIDMapper[parentID]; !ok {
				missMenusIDs = append(missMenusIDs, parentID)
			}
		}
		if len(missMenusIDs) > 0 {
			parentResult, err := a.MenuRepo.Query(ctx, model.MenuQueryParam{
				InIDs: missMenusIDs,
			})
			if err != nil {
				return nil, err
			}
			menuResult.Data = append(menuResult.Data, parentResult.Data...)
			sort.Sort(menuResult.Data)
		}
	}

	return menuResult.Data.ToTree(), nil
}

// Update current user info
func (a *Auth) UpdateUser(ctx context.Context, updateItem *model.UpdateCurrentUser) error {
	if util.FromIsRootUser(ctx) {
		return errors.BadRequest("", "Root user cannot update")
	}

	userID := util.FromUserID(ctx)
	user, err := a.UserRepo.Get(ctx, userID)
	if err != nil {
		return err
	} else if user == nil {
		return errors.NotFound("", "User not found")
	}

	user.FirstName = updateItem.FirstName
	user.LastName = updateItem.LastName
	user.FullName = updateItem.FirstName + " " + updateItem.LastName
	user.Phone = updateItem.Phone
	user.Remark = updateItem.Remark
	return a.UserRepo.Update(ctx, user, "name", "phone", "email", "remark")
}
