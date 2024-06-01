package api

import (
	"go-admin/internal/modules/auth/model"
	"go-admin/internal/modules/auth/service"
	"go-admin/pkg/util"

	"github.com/gin-gonic/gin"
)

type Auth struct {
	AuthService *service.Auth
}

// @Tags AuthAPI
// @Summary Get captcha ID
// @Success 200 {object} util.ResponseResult{data=model.Captcha}
// @Router /api/v1/captcha/id [get]
func (a *Auth) GetCaptcha(c *gin.Context) {
	ctx := c.Request.Context()
	data, err := a.AuthService.GetCaptcha(ctx)
	if err != nil {
		util.ResError(c, err)
		return
	}
	util.ResSuccess(c, data, "Get Captcha Successfully")
}

// @Tags AuthAPI
// @Summary Response captcha image
// @Param id query string true "Captcha ID"
// @Param reload query number false "Reload captcha image (reload=1)"
// @Produce image/png
// @Success 200 "Captcha image"
// @Failure 404 {object} util.ResponseResult
// @Router /api/v1/captcha/image [get]
func (a *Auth) ResponseCaptcha(c *gin.Context) {
	ctx := c.Request.Context()
	err := a.AuthService.ResponseCaptcha(ctx, c.Writer, c.Query("id"), c.Query("reload") == "1")
	if err != nil {
		util.ResError(c, err)
	}
}

// @Tags AuthAPI
// @Summary Login system with username and password
// @Param body body model.LoginForm true "Request body"
// @Success 200 {object} util.ResponseResult{data=model.LoginToken}
// @Failure 400 {object} util.ResponseResult
// @Failure 500 {object} util.ResponseResult
// @Router /api/v1/login [post]
func (a *Auth) Login(c *gin.Context) {
	ctx := c.Request.Context()
	item := new(model.LoginForm)
	if err := util.ParseJSON(c, item); err != nil {
		util.ResError(c, err)
		return
	}

	data, err := a.AuthService.Login(ctx, item.Trim())
	if err != nil {
		util.ResError(c, err)
		return
	}
	util.ResSuccess(c, data, "Login Successfully")
}

// @Tags AuthAPI
// @Summary Login system with username and password
// @Param body body model.RegisterForm true "Request body"
// @Success 200 {object} util.ResponseResult{data=model.LoginToken}
// @Failure 400 {object} util.ResponseResult
// @Failure 500 {object} util.ResponseResult
// @Router /api/v1/register [post]
func (a *Auth) Register(c *gin.Context) {
	ctx := c.Request.Context()
	item := new(model.RegisterForm)
	if err := util.ParseJSON(c, item); err != nil {
		util.ResError(c, err)
		return
	}
	data, err := a.AuthService.Register(ctx, item.Trim())
	if err != nil {
		util.ResError(c, err)
		return
	}
	util.ResSuccess(c, data, "Register Successfully")
}

// @Tags AuthAPI
// @Security ApiKeyAuth
// @Summary Logout system
// @Success 200 {object} util.ResponseResult
// @Failure 500 {object} util.ResponseResult
// @Router /api/v1/current/logout [post]
func (a *Auth) Logout(c *gin.Context) {
	ctx := c.Request.Context()
	err := a.AuthService.Logout(ctx)
	if err != nil {
		util.ResError(c, err)
		return
	}
	util.ResOK(c)
}

// @Tags AuthAPI
// @Security ApiKeyAuth
// @Summary Refresh current access token
// @Success 200 {object} util.ResponseResult{data=model.LoginToken}
// @Failure 401 {object} util.ResponseResult
// @Failure 500 {object} util.ResponseResult
// @Router /api/v1/current/refresh-token [post]
func (a *Auth) RefreshToken(c *gin.Context) {
	ctx := c.Request.Context()
	data, err := a.AuthService.RefreshToken(ctx)
	if err != nil {
		util.ResError(c, err)
		return
	}
	util.ResSuccess(c, data, "")
}

// @Tags AuthAPI
// @Security ApiKeyAuth
// @Summary Get current user info
// @Success 200 {object} util.ResponseResult{data=model.User}
// @Failure 401 {object} util.ResponseResult
// @Failure 500 {object} util.ResponseResult
// @Router /api/v1/current/user [get]
func (a *Auth) GetUserInfo(c *gin.Context) {
	ctx := c.Request.Context()
	data, err := a.AuthService.GetUserInfo(ctx)
	if err != nil {
		util.ResError(c, err)
		return
	}
	util.ResSuccess(c, data, "Get User Info Successfully")
}

// @Tags AuthAPI
// @Security ApiKeyAuth
// @Summary Change current user password
// @Param body body model.UpdateLoginPassword true "Request body"
// @Success 200 {object} util.ResponseResult
// @Failure 400 {object} util.ResponseResult
// @Failure 401 {object} util.ResponseResult
// @Failure 500 {object} util.ResponseResult
// @Router /api/v1/current/password [put]
func (a *Auth) UpdatePassword(c *gin.Context) {
	ctx := c.Request.Context()
	item := new(model.UpdateLoginPassword)
	if err := util.ParseJSON(c, item); err != nil {
		util.ResError(c, err)
		return
	}

	err := a.AuthService.UpdatePassword(ctx, item)
	if err != nil {
		util.ResError(c, err)
		return
	}
	util.ResOK(c)
}

// @Tags AuthAPI
// @Security ApiKeyAuth
// @Summary Query current user menus based on the current user role
// @Success 200 {object} util.ResponseResult{data=[]model.Menu}
// @Failure 401 {object} util.ResponseResult
// @Failure 500 {object} util.ResponseResult
// @Router /api/v1/current/menus [get]
func (a *Auth) QueryMenus(c *gin.Context) {
	ctx := c.Request.Context()
	data, err := a.AuthService.QueryMenus(ctx)
	if err != nil {
		util.ResError(c, err)
		return
	}
	util.ResSuccess(c, data, "Get Menus Successfully")
}

// @Tags AuthAPI
// @Security ApiKeyAuth
// @Summary Update current user info
// @Param body body model.UpdateCurrentUser true "Request body"
// @Success 200 {object} util.ResponseResult
// @Failure 400 {object} util.ResponseResult
// @Failure 401 {object} util.ResponseResult
// @Failure 500 {object} util.ResponseResult
// @Router /api/v1/current/user [put]
func (a *Auth) UpdateUser(c *gin.Context) {
	ctx := c.Request.Context()
	item := new(model.UpdateCurrentUser)
	if err := util.ParseJSON(c, item); err != nil {
		util.ResError(c, err)
		return
	}

	err := a.AuthService.UpdateUser(ctx, item)
	if err != nil {
		util.ResError(c, err)
		return
	}
	util.ResOK(c)
}
