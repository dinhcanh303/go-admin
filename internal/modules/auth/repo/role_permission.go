package repo

import (
	"context"
	"fmt"

	"go-admin/internal/modules/auth/model"
	"go-admin/pkg/errors"
	"go-admin/pkg/util"

	"gorm.io/gorm"
)

// Get user role storage instance
func GetRolePermissionDB(ctx context.Context, defDB *gorm.DB) *gorm.DB {
	return util.GetDB(ctx, defDB).Model(new(model.RolePermission))
}

// User roles for auth
type RolePermission struct {
	DB *gorm.DB
}

// Query user roles from the database based on the provided parameters and options.
func (a *RolePermission) Query(ctx context.Context, params model.RolePermissionQueryParam, opts ...model.RolePermissionQueryOptions) (*model.RolePermissionQueryResult, error) {
	var opt model.RolePermissionQueryOptions
	if len(opts) > 0 {
		opt = opts[0]
	}

	db := a.DB.Table(fmt.Sprintf("%s AS a", new(model.RolePermission).TableName()))
	if opt.JoinRole {
		db = db.Joins(fmt.Sprintf("left join %s b on a.role_id=b.id", new(model.Role).TableName()))
		db = db.Select("a.*,b.name as role_name")
	}
	if v := params.InUserIDs; len(v) > 0 {
		db = db.Where("a.user_id IN (?)", v)
	}
	if v := params.RoleID; len(v) > 0 {
		db = db.Where("a.role_id = ?", v)
	}
	if v := params.PermissionID; len(v) > 0 {
		db = db.Where("a.permission_id = ?", v)
	}

	var list model.RolePermissions
	pageResult, err := util.WrapPageQuery(ctx, db, params.PaginationParam, opt.QueryOptions, &list)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	queryResult := &model.RolePermissionQueryResult{
		PageResult: pageResult,
		Data:       list,
	}
	return queryResult, nil
}

// Get the specified user role from the database.
func (a *RolePermission) Get(ctx context.Context, id string, opts ...model.RolePermissionQueryOptions) (*model.RolePermission, error) {
	var opt model.RolePermissionQueryOptions
	if len(opts) > 0 {
		opt = opts[0]
	}

	item := new(model.RolePermission)
	ok, err := util.FindOne(ctx, GetRolePermissionDB(ctx, a.DB).Where("id=?", id), opt.QueryOptions, item)
	if err != nil {
		return nil, errors.WithStack(err)
	} else if !ok {
		return nil, nil
	}
	return item, nil
}

// Exist checks if the specified user role exists in the database.
func (a *RolePermission) Exists(ctx context.Context, id string) (bool, error) {
	ok, err := util.Exists(ctx, GetRolePermissionDB(ctx, a.DB).Where("id=?", id))
	return ok, errors.WithStack(err)
}

// Create a new user role.
func (a *RolePermission) Create(ctx context.Context, item *model.RolePermission) error {
	result := GetRolePermissionDB(ctx, a.DB).Create(item)
	return errors.WithStack(result.Error)
}

// Update the specified user role in the database.
func (a *RolePermission) Update(ctx context.Context, item *model.RolePermission) error {
	result := GetRolePermissionDB(ctx, a.DB).Where("id=?", item.ID).Select("*").Omit("created_at").Updates(item)
	return errors.WithStack(result.Error)
}

// Delete the specified user role from the database.
func (a *RolePermission) Delete(ctx context.Context, id string) error {
	result := GetRolePermissionDB(ctx, a.DB).Where("id=?", id).Delete(new(model.RolePermission))
	return errors.WithStack(result.Error)
}

func (a *RolePermission) DeleteByUserID(ctx context.Context, userID string) error {
	result := GetRolePermissionDB(ctx, a.DB).Where("user_id=?", userID).Delete(new(model.RolePermission))
	return errors.WithStack(result.Error)
}

func (a *RolePermission) DeleteByRoleID(ctx context.Context, roleID string) error {
	result := GetRolePermissionDB(ctx, a.DB).Where("role_id=?", roleID).Delete(new(model.RolePermission))
	return errors.WithStack(result.Error)
}
