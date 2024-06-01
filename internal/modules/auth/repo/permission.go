package repo

import (
	"context"

	"go-admin/internal/modules/auth/model"
	"go-admin/pkg/errors"
	"go-admin/pkg/util"

	"gorm.io/gorm"
)

// Get Permission storage instance
func GetPermissionDB(ctx context.Context, defDB *gorm.DB) *gorm.DB {
	return util.GetDB(ctx, defDB).Model(new(model.Permission))
}

// Permission management for auth
type Permission struct {
	DB *gorm.DB
}

// Query Permissions from the database based on the provided parameters and options.
func (a *Permission) Query(ctx context.Context, params model.PermissionQueryParam, opts ...model.PermissionQueryOptions) (*model.PermissionQueryResult, error) {
	var opt model.PermissionQueryOptions
	if len(opts) > 0 {
		opt = opts[0]
	}

	db := GetPermissionDB(ctx, a.DB)
	if v := params.InIDs; len(v) > 0 {
		db = db.Where("id IN (?)", v)
	}
	if v := params.LikeName; len(v) > 0 {
		db = db.Where("name LIKE ?", "%"+v+"%")
	}
	if v := params.GtUpdatedAt; v != nil {
		db = db.Where("updated_at > ?", v)
	}

	var list model.Permissions
	pageResult, err := util.WrapPageQuery(ctx, db, params.PaginationParam, opt.QueryOptions, &list)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	queryResult := &model.PermissionQueryResult{
		PageResult: pageResult,
		Data:       list,
	}
	return queryResult, nil
}

// Get the specified Permission from the database.
func (a *Permission) Get(ctx context.Context, id string, opts ...model.PermissionQueryOptions) (*model.Permission, error) {
	var opt model.PermissionQueryOptions
	if len(opts) > 0 {
		opt = opts[0]
	}

	item := new(model.Permission)
	ok, err := util.FindOne(ctx, GetPermissionDB(ctx, a.DB).Where("id=?", id), opt.QueryOptions, item)
	if err != nil {
		return nil, errors.WithStack(err)
	} else if !ok {
		return nil, nil
	}
	return item, nil
}

// Exist checks if the specified Permission exists in the database.
func (a *Permission) Exists(ctx context.Context, id string) (bool, error) {
	ok, err := util.Exists(ctx, GetPermissionDB(ctx, a.DB).Where("id=?", id))
	return ok, errors.WithStack(err)
}

func (a *Permission) ExistsCode(ctx context.Context, code string) (bool, error) {
	ok, err := util.Exists(ctx, GetPermissionDB(ctx, a.DB).Where("code=?", code))
	return ok, errors.WithStack(err)
}

// Create a new Permission.
func (a *Permission) Create(ctx context.Context, item *model.Permission) error {
	result := GetPermissionDB(ctx, a.DB).Create(item)
	return errors.WithStack(result.Error)
}

// Update the specified Permission in the database.
func (a *Permission) Update(ctx context.Context, item *model.Permission) error {
	result := GetPermissionDB(ctx, a.DB).Where("id=?", item.ID).Select("*").Omit("created_at").Updates(item)
	return errors.WithStack(result.Error)
}

// Delete the specified Permission from the database.
func (a *Permission) Delete(ctx context.Context, id string) error {
	result := GetPermissionDB(ctx, a.DB).Where("id=?", id).Delete(new(model.Permission))
	return errors.WithStack(result.Error)
}
