package service

import (
	"context"
	"time"

	"go-admin/internal/config"
	"go-admin/internal/modules/auth/model"
	"go-admin/internal/modules/auth/repo"
	"go-admin/pkg/cachex"
	"go-admin/pkg/crypto/hash"
	"go-admin/pkg/errors"
	"go-admin/pkg/util"
)

// User management for RBAC
type User struct {
	Cache        cachex.Cacher
	Trans        *util.Trans
	UserRepo     *repo.User
	UserRoleRepo *repo.UserRole
}

// Query users from the data access object based on the provided parameters and options.
func (a *User) Query(ctx context.Context, params model.UserQueryParam) (*model.UserQueryResult, error) {
	params.Pagination = true

	result, err := a.UserRepo.Query(ctx, params, model.UserQueryOptions{
		QueryOptions: util.QueryOptions{
			OrderFields: []util.OrderByParam{
				{Field: "created_at", Direction: util.DESC},
			},
			OmitFields: []string{"password"},
		},
	})
	if err != nil {
		return nil, err
	}

	if userIDs := result.Data.ToIDs(); len(userIDs) > 0 {
		userRoleResult, err := a.UserRoleRepo.Query(ctx, model.UserRoleQueryParam{
			InUserIDs: userIDs,
		}, model.UserRoleQueryOptions{
			JoinRole: true,
		})
		if err != nil {
			return nil, err
		}
		userRolesMap := userRoleResult.Data.ToUserIDMap()
		for _, user := range result.Data {
			user.Roles = userRolesMap[user.ID]
		}
	}

	return result, nil
}

// Get the specified user from the data access object.
func (a *User) Get(ctx context.Context, id string) (*model.User, error) {
	user, err := a.UserRepo.Get(ctx, id, model.UserQueryOptions{
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
		UserID: id,
	})
	if err != nil {
		return nil, err
	}
	user.Roles = userRoleResult.Data

	return user, nil
}

// Create a new user in the data access object.
func (a *User) Create(ctx context.Context, formItem *model.UserForm) (*model.User, error) {
	existsEmail, err := a.UserRepo.ExistsEmail(ctx, formItem.Email)
	if err != nil {
		return nil, err
	} else if existsEmail {
		return nil, errors.BadRequest("", "Email already exists")
	}

	user := &model.User{
		ID:        util.NewXID(),
		CreatedAt: time.Now(),
	}

	if formItem.Password == "" {
		formItem.Password = config.C.General.DefaultLoginPwd
	}

	if err := formItem.FillTo(user); err != nil {
		return nil, err
	}

	err = a.Trans.Exec(ctx, func(ctx context.Context) error {
		if err := a.UserRepo.Create(ctx, user); err != nil {
			return err
		}

		for _, userRole := range formItem.Roles {
			userRole.ID = util.NewXID()
			userRole.UserID = user.ID
			userRole.CreatedAt = time.Now()
			if err := a.UserRoleRepo.Create(ctx, userRole); err != nil {
				return err
			}
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	user.Roles = formItem.Roles

	return user, nil
}

// Update the specified user in the data access object.
func (a *User) Update(ctx context.Context, id string, formItem *model.UserForm) error {
	user, err := a.UserRepo.Get(ctx, id)
	if err != nil {
		return err
	} else if user == nil {
		return errors.NotFound("", "User not found")
	} else if user.Email != formItem.Email {
		existsEmail, err := a.UserRepo.ExistsEmail(ctx, formItem.Email)
		if err != nil {
			return err
		} else if existsEmail {
			return errors.BadRequest("", "Email already exists")
		}
	}

	if err := formItem.FillTo(user); err != nil {
		return err
	}
	user.UpdatedAt = time.Now()

	return a.Trans.Exec(ctx, func(ctx context.Context) error {
		if err := a.UserRepo.Update(ctx, user); err != nil {
			return err
		}

		if err := a.UserRoleRepo.DeleteByUserID(ctx, id); err != nil {
			return err
		}
		for _, userRole := range formItem.Roles {
			if userRole.ID == "" {
				userRole.ID = util.NewXID()
			}
			userRole.UserID = user.ID
			if userRole.CreatedAt.IsZero() {
				userRole.CreatedAt = time.Now()
			}
			userRole.UpdatedAt = time.Now()
			if err := a.UserRoleRepo.Create(ctx, userRole); err != nil {
				return err
			}
		}

		return a.Cache.Delete(ctx, config.CacheNSForUser, id)
	})
}

// Delete the specified user from the data access object.
func (a *User) Delete(ctx context.Context, id string) error {
	exists, err := a.UserRepo.Exists(ctx, id)
	if err != nil {
		return err
	} else if !exists {
		return errors.NotFound("", "User not found")
	}

	return a.Trans.Exec(ctx, func(ctx context.Context) error {
		if err := a.UserRepo.Delete(ctx, id); err != nil {
			return err
		}
		if err := a.UserRoleRepo.DeleteByUserID(ctx, id); err != nil {
			return err
		}
		return a.Cache.Delete(ctx, config.CacheNSForUser, id)
	})
}

func (a *User) ResetPassword(ctx context.Context, id string) error {
	exists, err := a.UserRepo.Exists(ctx, id)
	if err != nil {
		return err
	} else if !exists {
		return errors.NotFound("", "User not found")
	}

	hashPass, err := hash.GeneratePassword(config.C.General.DefaultLoginPwd)
	if err != nil {
		return errors.BadRequest("", "Failed to generate hash password: %s", err.Error())
	}

	return a.Trans.Exec(ctx, func(ctx context.Context) error {
		if err := a.UserRepo.UpdatePasswordByID(ctx, id, hashPass); err != nil {
			return err
		}
		return nil
	})
}

func (a *User) GetRoleIDs(ctx context.Context, id string) ([]string, error) {
	userRoleResult, err := a.UserRoleRepo.Query(ctx, model.UserRoleQueryParam{
		UserID: id,
	}, model.UserRoleQueryOptions{
		QueryOptions: util.QueryOptions{
			SelectFields: []string{"role_id"},
		},
	})
	if err != nil {
		return nil, err
	}
	return userRoleResult.Data.ToRoleIDs(), nil
}
