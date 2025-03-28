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

// Role management for RBAC
type Role struct {
	Cache              cachex.Cacher
	Trans              *util.Trans
	RoleRepo           *repo.Role
	RoleMenuRepo       *repo.RoleMenu
	RolePermissionRepo *repo.RolePermission
	UserRoleRepo       *repo.UserRole
}

// Query roles from the data access object based on the provided parameters and options.
func (a *Role) Query(ctx context.Context, params model.RoleQueryParam) (*model.RoleQueryResult, error) {
	params.Pagination = true

	var selectFields []string
	if params.ResultType == model.RoleResultTypeSelect {
		params.Pagination = false
		selectFields = []string{"id", "name"}
	}

	result, err := a.RoleRepo.Query(ctx, params, model.RoleQueryOptions{
		QueryOptions: util.QueryOptions{
			OrderFields: []util.OrderByParam{
				{Field: "sequence", Direction: util.DESC},
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
func (a *Role) Get(ctx context.Context, id string) (*model.Role, error) {
	role, err := a.RoleRepo.Get(ctx, id)
	if err != nil {
		return nil, err
	} else if role == nil {
		return nil, errors.NotFound("", "Role not found")
	}

	roleMenuResult, err := a.RoleMenuRepo.Query(ctx, model.RoleMenuQueryParam{
		RoleID: id,
	})
	if err != nil {
		return nil, err
	}
	role.Menus = roleMenuResult.Data

	return role, nil
}

// Create a new role in the data access object.
func (a *Role) Create(ctx context.Context, formItem *model.RoleForm) (*model.Role, error) {
	if exists, err := a.RoleRepo.ExistsCode(ctx, formItem.Code); err != nil {
		return nil, err
	} else if exists {
		return nil, errors.BadRequest("", "Role code already exists")
	}

	role := &model.Role{
		ID:        util.NewXID(),
		CreatedAt: time.Now(),
	}
	if err := formItem.FillTo(role); err != nil {
		return nil, err
	}

	err := a.Trans.Exec(ctx, func(ctx context.Context) error {
		if err := a.RoleRepo.Create(ctx, role); err != nil {
			return err
		}

		for _, roleMenu := range formItem.Menus {
			roleMenu.ID = util.NewXID()
			roleMenu.RoleID = role.ID
			roleMenu.CreatedAt = time.Now()
			if err := a.RoleMenuRepo.Create(ctx, roleMenu); err != nil {
				return err
			}
		}
		for _, rolePermission := range formItem.Permissions {
			rolePermission.ID = util.NewXID()
			rolePermission.RoleID = role.ID
			rolePermission.CreatedAt = time.Now()
			if err := a.RolePermissionRepo.Create(ctx, rolePermission); err != nil {
				return err
			}
		}
		return a.syncToCasbin(ctx)
	})
	if err != nil {
		return nil, err
	}
	role.Menus = formItem.Menus

	return role, nil
}

// Update the specified role in the data access object.
func (a *Role) Update(ctx context.Context, id string, formItem *model.RoleForm) error {
	role, err := a.RoleRepo.Get(ctx, id)
	if err != nil {
		return err
	} else if role == nil {
		return errors.NotFound("", "Role not found")
	} else if role.Code != formItem.Code {
		if exists, err := a.RoleRepo.ExistsCode(ctx, formItem.Code); err != nil {
			return err
		} else if exists {
			return errors.BadRequest("", "Role code already exists")
		}
	}

	if err := formItem.FillTo(role); err != nil {
		return err
	}
	role.UpdatedAt = time.Now()

	return a.Trans.Exec(ctx, func(ctx context.Context) error {
		if err := a.RoleRepo.Update(ctx, role); err != nil {
			return err
		}
		if err := a.RoleMenuRepo.DeleteByRoleID(ctx, id); err != nil {
			return err
		}
		for _, roleMenu := range formItem.Menus {
			if roleMenu.ID == "" {
				roleMenu.ID = util.NewXID()
			}
			roleMenu.RoleID = role.ID
			if roleMenu.CreatedAt.IsZero() {
				roleMenu.CreatedAt = time.Now()
			}
			roleMenu.UpdatedAt = time.Now()
			if err := a.RoleMenuRepo.Create(ctx, roleMenu); err != nil {
				return err
			}
		}
		if err := a.RolePermissionRepo.DeleteByRoleID(ctx, id); err != nil {
			return err
		}
		for _, rolePermission := range formItem.Permissions {
			if rolePermission.ID == "" {
				rolePermission.ID = util.NewXID()
			}
			rolePermission.RoleID = role.ID
			if rolePermission.CreatedAt.IsZero() {
				rolePermission.CreatedAt = time.Now()
			}
			rolePermission.UpdatedAt = time.Now()
			if err := a.RolePermissionRepo.Create(ctx, rolePermission); err != nil {
				return err
			}
		}
		return a.syncToCasbin(ctx)
	})
}

// Delete the specified role from the data access object.
func (a *Role) Delete(ctx context.Context, id string) error {
	exists, err := a.RoleRepo.Exists(ctx, id)
	if err != nil {
		return err
	} else if !exists {
		return errors.NotFound("", "Role not found")
	}

	return a.Trans.Exec(ctx, func(ctx context.Context) error {
		if err := a.RoleRepo.Delete(ctx, id); err != nil {
			return err
		}
		if err := a.RoleMenuRepo.DeleteByRoleID(ctx, id); err != nil {
			return err
		}
		if err := a.UserRoleRepo.DeleteByRoleID(ctx, id); err != nil {
			return err
		}

		return a.syncToCasbin(ctx)
	})
}

func (a *Role) syncToCasbin(ctx context.Context) error {
	return a.Cache.Set(ctx, config.CacheNSForRole, config.CacheKeyForSyncToCasbin, fmt.Sprintf("%d", time.Now().Unix()))
}
