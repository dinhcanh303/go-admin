package service

import (
	"context"
	"fmt"
	"time"

	"go-admin/internal/config"
	"go-admin/internal/modules/auth/model"
	"go-admin/internal/modules/auth/repo"
	"go-admin/pkg/cachex"
	"go-admin/pkg/errors"
	"go-admin/pkg/util"
)

// Permission management for RBAC
type Permission struct {
	Cache              cachex.Cacher
	Trans              *util.Trans
	PermissionRepo     *repo.Permission
	RolePermissionRepo *repo.RolePermission
}

// Query roles from the data access object based on the provided parameters and options.
func (a *Permission) Query(ctx context.Context, params model.PermissionQueryParam) (*model.PermissionQueryResult, error) {
	params.Pagination = true

	var selectFields []string
	if params.ResultType == model.PermissionResultTypeSelect {
		params.Pagination = false
		selectFields = []string{"id", "name"}
	}

	result, err := a.PermissionRepo.Query(ctx, params, model.PermissionQueryOptions{
		QueryOptions: util.QueryOptions{
			OrderFields: []util.OrderByParam{
				// {Field: "sequence", Direction: util.DESC},
				{Field: "created_at", Direction: util.DESC},
			},
			SelectFields: selectFields,
		},
	})
	if err != nil {
		return nil, err
	}
	return result, nil
}

// Get the specified role from the data access object.
func (a *Permission) Get(ctx context.Context, id string) (*model.Permission, error) {
	role, err := a.PermissionRepo.Get(ctx, id)
	if err != nil {
		return nil, err
	} else if role == nil {
		return nil, errors.NotFound("", "Permission not found")
	}
	return role, nil
}

// Create a new role in the data access object.
func (a *Permission) Create(ctx context.Context, formItem *model.PermissionForm) (*model.Permission, error) {
	if exists, err := a.PermissionRepo.ExistsCode(ctx, formItem.Code); err != nil {
		return nil, err
	} else if exists {
		return nil, errors.BadRequest("", "Permission code already exists")
	}

	role := &model.Permission{
		ID:        util.NewXID(),
		CreatedAt: time.Now(),
	}
	if err := formItem.FillTo(role); err != nil {
		return nil, err
	}

	err := a.Trans.Exec(ctx, func(ctx context.Context) error {
		if err := a.PermissionRepo.Create(ctx, role); err != nil {
			return err
		}
		return a.syncToCasbin(ctx)
	})
	if err != nil {
		return nil, err
	}
	return role, nil
}

// Update the specified role in the data access object.
func (a *Permission) Update(ctx context.Context, id string, formItem *model.PermissionForm) error {
	role, err := a.PermissionRepo.Get(ctx, id)
	if err != nil {
		return err
	} else if role == nil {
		return errors.NotFound("", "Permission not found")
	} else if role.Code != formItem.Code {
		if exists, err := a.PermissionRepo.ExistsCode(ctx, formItem.Code); err != nil {
			return err
		} else if exists {
			return errors.BadRequest("", "Permission code already exists")
		}
	}

	if err := formItem.FillTo(role); err != nil {
		return err
	}
	role.UpdatedAt = time.Now()
	return a.Trans.Exec(ctx, func(ctx context.Context) error {
		if err := a.PermissionRepo.Update(ctx, role); err != nil {
			return err
		}
		return a.syncToCasbin(ctx)
	})
}

// Delete the specified role from the data access object.
func (a *Permission) Delete(ctx context.Context, id string) error {
	exists, err := a.PermissionRepo.Exists(ctx, id)
	if err != nil {
		return err
	} else if !exists {
		return errors.NotFound("", "Permission not found")
	}

	return a.Trans.Exec(ctx, func(ctx context.Context) error {
		if err := a.PermissionRepo.Delete(ctx, id); err != nil {
			return err
		}
		return a.syncToCasbin(ctx)
	})
}

func (a *Permission) syncToCasbin(ctx context.Context) error {
	return a.Cache.Set(ctx, config.CacheNSForPermission, config.CacheKeyForSyncToCasbin, fmt.Sprintf("%d", time.Now().Unix()))
}
