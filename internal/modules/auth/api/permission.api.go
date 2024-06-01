package api

import (
	"go-admin/internal/modules/auth/model"
	"go-admin/internal/modules/auth/service"
	"go-admin/pkg/util"

	"github.com/gin-gonic/gin"
)

// Permission management for RBAC
type Permission struct {
	PermissionService *service.Permission
}

// @Tags PermissionAPI
// @Security ApiKeyAuth
// @Summary Query role list
// @Param current query int true "pagination index" default(1)
// @Param pageSize query int true "pagination size" default(10)
// @Param name query string false "Display name of role"
// @Param status query string false "Status of role (disabled, enabled)"
// @Success 200 {object} util.ResponseResult{data=[]model.Permission}
// @Failure 401 {object} util.ResponseResult
// @Failure 500 {object} util.ResponseResult
// @Router /api/v1/permissions [get]
func (a *Permission) Query(c *gin.Context) {
	ctx := c.Request.Context()
	var params model.PermissionQueryParam
	if err := util.ParseQuery(c, &params); err != nil {
		util.ResError(c, err)
		return
	}

	result, err := a.PermissionService.Query(ctx, params)
	if err != nil {
		util.ResError(c, err)
		return
	}
	util.ResPage(c, result.Data, result.PageResult)
}

// @Tags PermissionAPI
// @Security ApiKeyAuth
// @Summary Get role record by ID
// @Param id path string true "unique id"
// @Success 200 {object} util.ResponseResult{data=model.Permission}
// @Failure 401 {object} util.ResponseResult
// @Failure 500 {object} util.ResponseResult
// @Router /api/v1/permissions/{id} [get]
func (a *Permission) Get(c *gin.Context) {
	ctx := c.Request.Context()
	item, err := a.PermissionService.Get(ctx, c.Param("id"))
	if err != nil {
		util.ResError(c, err)
		return
	}
	util.ResSuccess(c, item, "Get Permission Successfully")
}

// @Tags PermissionAPI
// @Security ApiKeyAuth
// @Summary Create role record
// @Param body body model.PermissionForm true "Request body"
// @Success 200 {object} util.ResponseResult{data=model.Permission}
// @Failure 400 {object} util.ResponseResult
// @Failure 401 {object} util.ResponseResult
// @Failure 500 {object} util.ResponseResult
// @Router /api/v1/permissions [post]
func (a *Permission) Create(c *gin.Context) {
	ctx := c.Request.Context()
	item := new(model.PermissionForm)
	if err := util.ParseJSON(c, item); err != nil {
		util.ResError(c, err)
		return
	} else if err := item.Validate(); err != nil {
		util.ResError(c, err)
		return
	}

	result, err := a.PermissionService.Create(ctx, item)
	if err != nil {
		util.ResError(c, err)
		return
	}
	util.ResSuccess(c, result, "")
}

// @Tags PermissionAPI
// @Security ApiKeyAuth
// @Summary Update role record by ID
// @Param id path string true "unique id"
// @Param body body model.PermissionForm true "Request body"
// @Success 200 {object} util.ResponseResult
// @Failure 400 {object} util.ResponseResult
// @Failure 401 {object} util.ResponseResult
// @Failure 500 {object} util.ResponseResult
// @Router /api/v1/permissions/{id} [put]
func (a *Permission) Update(c *gin.Context) {
	ctx := c.Request.Context()
	item := new(model.PermissionForm)
	if err := util.ParseJSON(c, item); err != nil {
		util.ResError(c, err)
		return
	} else if err := item.Validate(); err != nil {
		util.ResError(c, err)
		return
	}

	err := a.PermissionService.Update(ctx, c.Param("id"), item)
	if err != nil {
		util.ResError(c, err)
		return
	}
	util.ResOK(c)
}

// @Tags PermissionAPI
// @Security ApiKeyAuth
// @Summary Delete role record by ID
// @Param id path string true "unique id"
// @Success 200 {object} util.ResponseResult
// @Failure 401 {object} util.ResponseResult
// @Failure 500 {object} util.ResponseResult
// @Router /api/v1/permissions/{id} [delete]
func (a *Permission) Delete(c *gin.Context) {
	ctx := c.Request.Context()
	err := a.PermissionService.Delete(ctx, c.Param("id"))
	if err != nil {
		util.ResError(c, err)
		return
	}
	util.ResOK(c)
}
