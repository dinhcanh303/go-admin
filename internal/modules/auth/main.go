package auth

import (
	"context"
	"path/filepath"

	"go-admin/internal/config"
	"go-admin/internal/modules/auth/api"
	"go-admin/internal/modules/auth/model"
	"go-admin/pkg/logging"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type Auth struct {
	DB            *gorm.DB
	MenuAPI       *api.Menu
	RoleAPI       *api.Role
	UserAPI       *api.User
	AuthAPI       *api.Auth
	PermissionAPI *api.Permission
	Casbinx       *Casbinx
}

func (a *Auth) AutoMigrate(ctx context.Context) error {
	return a.DB.AutoMigrate(
		new(model.Menu),
		new(model.MenuResource),
		new(model.Role),
		new(model.RoleMenu),
		new(model.User),
		new(model.UserRole),
		new(model.Permission),
	)
}

func (a *Auth) Init(ctx context.Context) error {
	if config.C.Storage.DB.AutoMigrate {
		if err := a.AutoMigrate(ctx); err != nil {
			return err
		}
	}

	if err := a.Casbinx.Load(ctx); err != nil {
		return err
	}

	if name := config.C.General.MenuFile; name != "" {
		fullPath := filepath.Join(config.C.General.WorkDir, name)
		if err := a.MenuAPI.MenuService.InitFromFile(ctx, fullPath); err != nil {
			logging.Context(ctx).Error("failed to init menu data", zap.Error(err), zap.String("file", fullPath))
		}
	}

	return nil
}

func (a *Auth) RegisterV1Routers(ctx context.Context, v1 *gin.RouterGroup) error {
	captcha := v1.Group("captcha")
	{
		captcha.GET("id", a.AuthAPI.GetCaptcha)
		captcha.GET("image", a.AuthAPI.ResponseCaptcha)
	}
	v1.POST("login", a.AuthAPI.Login)
	v1.POST("register", a.AuthAPI.Register)

	current := v1.Group("current")
	{
		current.POST("refresh-token", a.AuthAPI.RefreshToken)
		current.GET("user", a.AuthAPI.GetUserInfo)
		current.GET("menus", a.AuthAPI.QueryMenus)
		current.PUT("password", a.AuthAPI.UpdatePassword)
		current.PUT("user", a.AuthAPI.UpdateUser)
		current.POST("logout", a.AuthAPI.Logout)
	}
	menu := v1.Group("menus")
	{
		menu.GET("", a.MenuAPI.Query)
		menu.GET(":id", a.MenuAPI.Get)
		menu.POST("", a.MenuAPI.Create)
		menu.PUT(":id", a.MenuAPI.Update)
		menu.DELETE(":id", a.MenuAPI.Delete)
	}
	permission := v1.Group("permissions")
	{
		permission.GET("", a.PermissionAPI.Query)
		permission.GET(":id", a.PermissionAPI.Get)
		permission.POST("", a.PermissionAPI.Create)
		permission.PUT(":id", a.PermissionAPI.Update)
		permission.DELETE(":id", a.PermissionAPI.Delete)
	}
	role := v1.Group("roles")
	{
		role.GET("", a.RoleAPI.Query)
		role.GET(":id", a.RoleAPI.Get)
		role.POST("", a.RoleAPI.Create)
		role.PUT(":id", a.RoleAPI.Update)
		role.DELETE(":id", a.RoleAPI.Delete)
	}
	user := v1.Group("users")
	{
		user.GET("", a.UserAPI.Query)
		user.GET(":id", a.UserAPI.Get)
		user.POST("", a.UserAPI.Create)
		user.PUT(":id", a.UserAPI.Update)
		user.DELETE(":id", a.UserAPI.Delete)
		user.PATCH(":id/reset-pwd", a.UserAPI.ResetPassword)
	}
	return nil
}

func (a *Auth) Release(ctx context.Context) error {
	if err := a.Casbinx.Release(ctx); err != nil {
		return err
	}
	return nil
}
