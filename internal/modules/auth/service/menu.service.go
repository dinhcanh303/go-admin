package service

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"go-admin/internal/config"
	"go-admin/internal/modules/auth/model"
	"go-admin/internal/modules/auth/repo"
	"go-admin/pkg/cachex"
	"go-admin/pkg/encoding/json"
	"go-admin/pkg/encoding/yaml"
	"go-admin/pkg/errors"
	"go-admin/pkg/logging"
	"go-admin/pkg/util"

	"go.uber.org/zap"
)

// Menu management for RBAC
type Menu struct {
	Cache            cachex.Cacher
	Trans            *util.Trans
	MenuRepo         *repo.Menu
	MenuResourceRepo *repo.MenuResource
	RoleMenuRepo     *repo.RoleMenu
}

func (a *Menu) InitFromFile(ctx context.Context, menuFile string) error {
	f, err := os.ReadFile(menuFile)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			logging.Context(ctx).Warn("Menu data file not found, skip init menu data from file", zap.String("file", menuFile))
			return nil
		}
		return err
	}

	var menus model.Menus
	if ext := filepath.Ext(menuFile); ext == ".json" {
		if err := json.Unmarshal(f, &menus); err != nil {
			return errors.Wrapf(err, "Unmarshal JSON file '%s' failed", menuFile)
		}
	} else if ext == ".yaml" || ext == ".yml" {
		if err := yaml.Unmarshal(f, &menus); err != nil {
			return errors.Wrapf(err, "Unmarshal YAML file '%s' failed", menuFile)
		}
	} else {
		return errors.Errorf("Unsupported file type '%s'", ext)
	}

	return a.Trans.Exec(ctx, func(ctx context.Context) error {
		return a.createInBatchByParent(ctx, menus, nil)
	})
}

func (a *Menu) createInBatchByParent(ctx context.Context, items model.Menus, parent *model.Menu) error {
	total := len(items)
	for i, item := range items {
		var parentID string
		if parent != nil {
			parentID = parent.ID
		}

		exist := false

		if item.ID != "" {
			exists, err := a.MenuRepo.Exists(ctx, item.ID)
			if err != nil {
				return err
			} else if exists {
				exist = true
			}
		} else if item.Code != "" {
			exists, err := a.MenuRepo.ExistsCodeByParentID(ctx, item.Code, parentID)
			if err != nil {
				return err
			} else if exists {
				exist = true
				existItem, err := a.MenuRepo.GetByCodeAndParentID(ctx, item.Code, parentID)
				if err != nil {
					return err
				}
				if existItem != nil {
					item.ID = existItem.ID
				}
			}
		} else if item.Name != "" {
			exists, err := a.MenuRepo.ExistsNameByParentID(ctx, item.Name, parentID)
			if err != nil {
				return err
			} else if exists {
				exist = true
				existItem, err := a.MenuRepo.GetByNameAndParentID(ctx, item.Name, parentID)
				if err != nil {
					return err
				}
				if existItem != nil {
					item.ID = existItem.ID
				}
			}
		}

		if !exist {
			if item.ID == "" {
				item.ID = util.NewXID()
			}
			if item.Status == "" {
				item.Status = model.MenuStatusEnabled
			}
			if item.Sequence == 0 {
				item.Sequence = total - i
			}

			item.ParentID = parentID
			if parent != nil {
				item.ParentPath = parent.ParentPath + parentID + util.TreePathDelimiter
			}
			item.CreatedAt = time.Now()

			if err := a.MenuRepo.Create(ctx, item); err != nil {
				return err
			}
		}

		for _, res := range item.Resources {
			if res.ID != "" {
				exists, err := a.MenuResourceRepo.Exists(ctx, res.ID)
				if err != nil {
					return err
				} else if exists {
					continue
				}
			}

			if res.Path != "" {
				exists, err := a.MenuResourceRepo.ExistsMethodPathByMenuID(ctx, res.Method, res.Path, item.ID)
				if err != nil {
					return err
				} else if exists {
					continue
				}
			}

			if res.ID == "" {
				res.ID = util.NewXID()
			}
			res.MenuID = item.ID
			res.CreatedAt = time.Now()
			if err := a.MenuResourceRepo.Create(ctx, res); err != nil {
				return err
			}
		}

		if item.Children != nil {
			if err := a.createInBatchByParent(ctx, *item.Children, item); err != nil {
				return err
			}
		}
	}
	return nil
}

// Query menus from the data access object based on the provided parameters and options.
func (a *Menu) Query(ctx context.Context, params model.MenuQueryParam) (*model.MenuQueryResult, error) {
	params.Pagination = false

	if err := a.fillQueryParam(ctx, &params); err != nil {
		return nil, err
	}

	result, err := a.MenuRepo.Query(ctx, params, model.MenuQueryOptions{
		QueryOptions: util.QueryOptions{
			OrderFields: model.MenusOrderParams,
		},
	})
	if err != nil {
		return nil, err
	}

	if params.LikeName != "" || params.CodePath != "" {
		result.Data, err = a.appendChildren(ctx, result.Data)
		if err != nil {
			return nil, err
		}
	}

	if params.IncludeResources {
		for i, item := range result.Data {
			resResult, err := a.MenuResourceRepo.Query(ctx, model.MenuResourceQueryParam{
				MenuID: item.ID,
			})
			if err != nil {
				return nil, err
			}
			result.Data[i].Resources = resResult.Data
		}
	}

	result.Data = result.Data.ToTree()
	return result, nil
}

func (a *Menu) fillQueryParam(ctx context.Context, params *model.MenuQueryParam) error {
	if params.CodePath != "" {
		var (
			codes    []string
			lastMenu model.Menu
		)
		for _, code := range strings.Split(params.CodePath, util.TreePathDelimiter) {
			if code == "" {
				continue
			}
			codes = append(codes, code)
			menu, err := a.MenuRepo.GetByCodeAndParentID(ctx, code, lastMenu.ParentID, model.MenuQueryOptions{
				QueryOptions: util.QueryOptions{
					SelectFields: []string{"id", "parent_id", "parent_path"},
				},
			})
			if err != nil {
				return err
			} else if menu == nil {
				return errors.NotFound("", "Menu not found by code '%s'", strings.Join(codes, util.TreePathDelimiter))
			}
			lastMenu = *menu
		}
		params.ParentPathPrefix = lastMenu.ParentPath + lastMenu.ID + util.TreePathDelimiter
	}
	return nil
}

func (a *Menu) appendChildren(ctx context.Context, data model.Menus) (model.Menus, error) {
	if len(data) == 0 {
		return data, nil
	}

	existsInData := func(id string) bool {
		for _, item := range data {
			if item.ID == id {
				return true
			}
		}
		return false
	}

	for _, item := range data {
		childResult, err := a.MenuRepo.Query(ctx, model.MenuQueryParam{
			ParentPathPrefix: item.ParentPath + item.ID + util.TreePathDelimiter,
		})
		if err != nil {
			return nil, err
		}
		for _, child := range childResult.Data {
			if existsInData(child.ID) {
				continue
			}
			data = append(data, child)
		}
	}

	if parentIDs := data.SplitParentIDs(); len(parentIDs) > 0 {
		parentResult, err := a.MenuRepo.Query(ctx, model.MenuQueryParam{
			InIDs: parentIDs,
		})
		if err != nil {
			return nil, err
		}
		for _, p := range parentResult.Data {
			if existsInData(p.ID) {
				continue
			}
			data = append(data, p)
		}
	}
	sort.Sort(data)

	return data, nil
}

// Get the specified menu from the data access object.
func (a *Menu) Get(ctx context.Context, id string) (*model.Menu, error) {
	menu, err := a.MenuRepo.Get(ctx, id)
	if err != nil {
		return nil, err
	} else if menu == nil {
		return nil, errors.NotFound("", "Menu not found")
	}

	menuResResult, err := a.MenuResourceRepo.Query(ctx, model.MenuResourceQueryParam{
		MenuID: menu.ID,
	})
	if err != nil {
		return nil, err
	}
	menu.Resources = menuResResult.Data

	return menu, nil
}

// Create a new menu in the data access object.
func (a *Menu) Create(ctx context.Context, formItem *model.MenuForm) (*model.Menu, error) {
	menu := &model.Menu{
		ID:        util.NewXID(),
		CreatedAt: time.Now(),
	}

	if parentID := formItem.ParentID; parentID != "" {
		parent, err := a.MenuRepo.Get(ctx, parentID)
		if err != nil {
			return nil, err
		} else if parent == nil {
			return nil, errors.NotFound("", "Parent not found")
		}
		menu.ParentPath = parent.ParentPath + parent.ID + util.TreePathDelimiter
	}

	if exists, err := a.MenuRepo.ExistsCodeByParentID(ctx, formItem.Code, formItem.ParentID); err != nil {
		return nil, err
	} else if exists {
		return nil, errors.BadRequest("", "Menu code already exists at the same level")
	}

	if err := formItem.FillTo(menu); err != nil {
		return nil, err
	}

	err := a.Trans.Exec(ctx, func(ctx context.Context) error {
		if err := a.MenuRepo.Create(ctx, menu); err != nil {
			return err
		}

		for _, res := range formItem.Resources {
			res.ID = util.NewXID()
			res.MenuID = menu.ID
			res.CreatedAt = time.Now()
			if err := a.MenuResourceRepo.Create(ctx, res); err != nil {
				return err
			}
		}

		return nil
	})
	if err != nil {
		return nil, err
	}
	return menu, nil
}

// Update the specified menu in the data access object.
func (a *Menu) Update(ctx context.Context, id string, formItem *model.MenuForm) error {
	menu, err := a.MenuRepo.Get(ctx, id)
	if err != nil {
		return err
	} else if menu == nil {
		return errors.NotFound("", "Menu not found")
	}

	oldParentPath := menu.ParentPath
	oldStatus := menu.Status
	var childData model.Menus
	if menu.ParentID != formItem.ParentID {
		if parentID := formItem.ParentID; parentID != "" {
			parent, err := a.MenuRepo.Get(ctx, parentID)
			if err != nil {
				return err
			} else if parent == nil {
				return errors.NotFound("", "Parent not found")
			}
			menu.ParentPath = parent.ParentPath + parent.ID + util.TreePathDelimiter
		} else {
			menu.ParentPath = ""
		}

		childResult, err := a.MenuRepo.Query(ctx, model.MenuQueryParam{
			ParentPathPrefix: oldParentPath + menu.ID + util.TreePathDelimiter,
		}, model.MenuQueryOptions{
			QueryOptions: util.QueryOptions{
				SelectFields: []string{"id", "parent_path"},
			},
		})
		if err != nil {
			return err
		}
		childData = childResult.Data
	}

	if menu.Code != formItem.Code {
		if exists, err := a.MenuRepo.ExistsCodeByParentID(ctx, formItem.Code, formItem.ParentID); err != nil {
			return err
		} else if exists {
			return errors.BadRequest("", "Menu code already exists at the same level")
		}
	}

	if err := formItem.FillTo(menu); err != nil {
		return err
	}

	return a.Trans.Exec(ctx, func(ctx context.Context) error {
		if oldStatus != formItem.Status {
			oldPath := oldParentPath + menu.ID + util.TreePathDelimiter
			if err := a.MenuRepo.UpdateStatusByParentPath(ctx, oldPath, formItem.Status); err != nil {
				return err
			}
		}

		for _, child := range childData {
			oldPath := oldParentPath + menu.ID + util.TreePathDelimiter
			newPath := menu.ParentPath + menu.ID + util.TreePathDelimiter
			err := a.MenuRepo.UpdateParentPath(ctx, child.ID, strings.Replace(child.ParentPath, oldPath, newPath, 1))
			if err != nil {
				return err
			}
		}

		if err := a.MenuRepo.Update(ctx, menu); err != nil {
			return err
		}

		if err := a.MenuResourceRepo.DeleteByMenuID(ctx, id); err != nil {
			return err
		}
		for _, res := range formItem.Resources {
			if res.ID == "" {
				res.ID = util.NewXID()
			}
			res.MenuID = id
			if res.CreatedAt.IsZero() {
				res.CreatedAt = time.Now()
			}
			res.UpdatedAt = time.Now()
			if err := a.MenuResourceRepo.Create(ctx, res); err != nil {
				return err
			}
		}

		return a.syncToCasbin(ctx)
	})
}

// Delete the specified menu from the data access object.
func (a *Menu) Delete(ctx context.Context, id string) error {
	if config.C.General.DenyDeleteMenu {
		return errors.BadRequest("", "Menu deletion is not allowed")
	}

	menu, err := a.MenuRepo.Get(ctx, id)
	if err != nil {
		return err
	} else if menu == nil {
		return errors.NotFound("", "Menu not found")
	}

	childResult, err := a.MenuRepo.Query(ctx, model.MenuQueryParam{
		ParentPathPrefix: menu.ParentPath + menu.ID + util.TreePathDelimiter,
	}, model.MenuQueryOptions{
		QueryOptions: util.QueryOptions{
			SelectFields: []string{"id"},
		},
	})
	if err != nil {
		return err
	}

	return a.Trans.Exec(ctx, func(ctx context.Context) error {
		if err := a.delete(ctx, id); err != nil {
			return err
		}

		for _, child := range childResult.Data {
			if err := a.delete(ctx, child.ID); err != nil {
				return err
			}
		}

		return a.syncToCasbin(ctx)
	})
}

func (a *Menu) delete(ctx context.Context, id string) error {
	if err := a.MenuRepo.Delete(ctx, id); err != nil {
		return err
	}
	if err := a.MenuResourceRepo.DeleteByMenuID(ctx, id); err != nil {
		return err
	}
	if err := a.RoleMenuRepo.DeleteByMenuID(ctx, id); err != nil {
		return err
	}
	return nil
}

func (a *Menu) syncToCasbin(ctx context.Context) error {
	return a.Cache.Set(ctx, config.CacheNSForRole, config.CacheKeyForSyncToCasbin, fmt.Sprintf("%d", time.Now().Unix()))
}
