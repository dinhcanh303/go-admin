package api

import (
	"go-admin/internal/modules/auth/model"
	"go-admin/internal/modules/auth/service"
	"go-admin/pkg/util"

	"github.com/gin-gonic/gin"
)

// Role management for RBAC
type Role struct {
	RoleService *service.Role
}

// @Tags RoleAPI
// @Security ApiKeyAuth
// @Summary Query role list
// @Param current query int true "pagination index" default(1)
// @Param pageSize query int true "pagination size" default(10)
// @Param name query string false "Display name of role"
// @Param status query string false "Status of role (disabled, enabled)"
// @Success 200 {object} util.ResponseResult{data=[]model.Role}
// @Failure 401 {object} util.ResponseResult
// @Failure 500 {object} util.ResponseResult
// @Router /api/v1/roles [get]
func (a *Role) Query(c *gin.Context) {
	ctx := c.Request.Context()
	var params model.RoleQueryParam
	if err := util.ParseQuery(c, &params); err != nil {
		util.ResError(c, err)
		return
	}

	result, err := a.RoleService.Query(ctx, params)
	if err != nil {
		util.ResError(c, err)
		return
	}
	util.ResPage(c, result.Data, result.PageResult)
}

// @Tags RoleAPI
// @Security ApiKeyAuth
// @Summary Get role record by ID
// @Param id path string true "unique id"
// @Success 200 {object} util.ResponseResult{data=model.Role}
// @Failure 401 {object} util.ResponseResult
// @Failure 500 {object} util.ResponseResult
// @Router /api/v1/roles/{id} [get]
func (a *Role) Get(c *gin.Context) {
	ctx := c.Request.Context()
	item, err := a.RoleService.Get(ctx, c.Param("id"))
	if err != nil {
		util.ResError(c, err)
		return
	}
	util.ResSuccess(c, item, "Get Role Successfully")
}

// @Tags RoleAPI
// @Security ApiKeyAuth
// @Summary Create role record
// @Param body body model.RoleForm true "Request body"
// @Success 200 {object} util.ResponseResult{data=model.Role}
// @Failure 400 {object} util.ResponseResult
// @Failure 401 {object} util.ResponseResult
// @Failure 500 {object} util.ResponseResult
// @Router /api/v1/roles [post]
func (a *Role) Create(c *gin.Context) {
	ctx := c.Request.Context()
	item := new(model.RoleForm)
	if err := util.ParseJSON(c, item); err != nil {
		util.ResError(c, err)
		return
	} else if err := item.Validate(); err != nil {
		util.ResError(c, err)
		return
	}

	result, err := a.RoleService.Create(ctx, item)
	if err != nil {
		util.ResError(c, err)
		return
	}
	util.ResSuccess(c, result, "")
}

// @Tags RoleAPI
// @Security ApiKeyAuth
// @Summary Update role record by ID
// @Param id path string true "unique id"
// @Param body body model.RoleForm true "Request body"
// @Success 200 {object} util.ResponseResult
// @Failure 400 {object} util.ResponseResult
// @Failure 401 {object} util.ResponseResult
// @Failure 500 {object} util.ResponseResult
// @Router /api/v1/roles/{id} [put]
func (a *Role) Update(c *gin.Context) {
	ctx := c.Request.Context()
	item := new(model.RoleForm)
	if err := util.ParseJSON(c, item); err != nil {
		util.ResError(c, err)
		return
	} else if err := item.Validate(); err != nil {
		util.ResError(c, err)
		return
	}

	err := a.RoleService.Update(ctx, c.Param("id"), item)
	if err != nil {
		util.ResError(c, err)
		return
	}
	util.ResOK(c)
}

// @Tags RoleAPI
// @Security ApiKeyAuth
// @Summary Delete role record by ID
// @Param id path string true "unique id"
// @Success 200 {object} util.ResponseResult
// @Failure 401 {object} util.ResponseResult
// @Failure 500 {object} util.ResponseResult
// @Router /api/v1/roles/{id} [delete]
func (a *Role) Delete(c *gin.Context) {
	ctx := c.Request.Context()
	err := a.RoleService.Delete(ctx, c.Param("id"))
	if err != nil {
		util.ResError(c, err)
		return
	}
	util.ResOK(c)
}
