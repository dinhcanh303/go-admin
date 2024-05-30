package repo

import (
	"context"

	"go-admin/internal/modules/auth/model"
	"go-admin/pkg/errors"
	"go-admin/pkg/util"

	"gorm.io/gorm"
)

// Get user storage instance
func GetUserDB(ctx context.Context, defDB *gorm.DB) *gorm.DB {
	return util.GetDB(ctx, defDB).Model(new(model.User))
}

// User management for auth
type User struct {
	DB *gorm.DB
}

// Query users from the database based on the provided parameters and options.
func (a *User) Query(ctx context.Context, params model.UserQueryParam, opts ...model.UserQueryOptions) (*model.UserQueryResult, error) {
	var opt model.UserQueryOptions
	if len(opts) > 0 {
		opt = opts[0]
	}

	db := GetUserDB(ctx, a.DB)
	if v := params.LikeUsername; len(v) > 0 {
		db = db.Where("username LIKE ?", "%"+v+"%")
	}
	if v := params.LikeName; len(v) > 0 {
		db = db.Where("name LIKE ?", "%"+v+"%")
	}
	if v := params.Status; len(v) > 0 {
		db = db.Where("status = ?", v)
	}

	var list model.Users
	pageResult, err := util.WrapPageQuery(ctx, db, params.PaginationParam, opt.QueryOptions, &list)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	queryResult := &model.UserQueryResult{
		PageResult: pageResult,
		Data:       list,
	}
	return queryResult, nil
}

// Get the specified user from the database.
func (a *User) Get(ctx context.Context, id string, opts ...model.UserQueryOptions) (*model.User, error) {
	var opt model.UserQueryOptions
	if len(opts) > 0 {
		opt = opts[0]
	}

	item := new(model.User)
	ok, err := util.FindOne(ctx, GetUserDB(ctx, a.DB).Where("id=?", id), opt.QueryOptions, item)
	if err != nil {
		return nil, errors.WithStack(err)
	} else if !ok {
		return nil, nil
	}
	return item, nil
}

func (a *User) GetByEmail(ctx context.Context, email string, opts ...model.UserQueryOptions) (*model.User, error) {
	var opt model.UserQueryOptions
	if len(opts) > 0 {
		opt = opts[0]
	}

	item := new(model.User)
	ok, err := util.FindOne(ctx, GetUserDB(ctx, a.DB).Where("email=?", email), opt.QueryOptions, item)
	if err != nil {
		return nil, errors.WithStack(err)
	} else if !ok {
		return nil, nil
	}
	return item, nil
}

// Exist checks if the specified user exists in the database.
func (a *User) Exists(ctx context.Context, id string) (bool, error) {
	ok, err := util.Exists(ctx, GetUserDB(ctx, a.DB).Where("id=?", id))
	return ok, errors.WithStack(err)
}

func (a *User) ExistsEmail(ctx context.Context, email string) (bool, error) {
	ok, err := util.Exists(ctx, GetUserDB(ctx, a.DB).Where("email=?", email))
	return ok, errors.WithStack(err)
}

// Create a new user.
func (a *User) Create(ctx context.Context, item *model.User) error {
	result := GetUserDB(ctx, a.DB).Create(item)
	return errors.WithStack(result.Error)
}

// Update the specified user in the database.
func (a *User) Update(ctx context.Context, item *model.User, selectFields ...string) error {
	db := GetUserDB(ctx, a.DB).Where("id=?", item.ID)
	if len(selectFields) > 0 {
		db = db.Select(selectFields)
	} else {
		db = db.Select("*").Omit("created_at")
	}
	result := db.Updates(item)
	return errors.WithStack(result.Error)
}

// Delete the specified user from the database.
func (a *User) Delete(ctx context.Context, id string) error {
	result := GetUserDB(ctx, a.DB).Where("id=?", id).Delete(new(model.User))
	return errors.WithStack(result.Error)
}

func (a *User) UpdatePasswordByID(ctx context.Context, id string, password string) error {
	result := GetUserDB(ctx, a.DB).Where("id=?", id).Select("password").Updates(model.User{Password: password})
	return errors.WithStack(result.Error)
}
