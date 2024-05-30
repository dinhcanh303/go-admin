package tests

import (
	"net/http"
	"testing"

	"go-admin/internal/modules/auth/model"
	"go-admin/pkg/util"

	"github.com/stretchr/testify/assert"
)

func TestRole(t *testing.T) {
	e := tester(t)

	menuFormItem := model.MenuForm{
		Code:        "role",
		Name:        "Role management",
		Description: "Role management",
		Type:        "page",
		Path:        "/system/role",
		Properties:  `{"icon":"role"}`,
		Status:      model.MenuStatusEnabled,
	}

	var menu model.Menu
	e.POST(baseAPI + "/menus").WithJSON(menuFormItem).
		Expect().Status(http.StatusOK).JSON().Decode(&util.ResponseResult{Data: &menu})

	assert := assert.New(t)
	assert.NotEmpty(menu.ID)
	assert.Equal(menuFormItem.Code, menu.Code)
	assert.Equal(menuFormItem.Name, menu.Name)
	assert.Equal(menuFormItem.Description, menu.Description)
	assert.Equal(menuFormItem.Type, menu.Type)
	assert.Equal(menuFormItem.Path, menu.Path)
	assert.Equal(menuFormItem.Properties, menu.Properties)
	assert.Equal(menuFormItem.Status, menu.Status)

	roleFormItem := model.RoleForm{
		Code: "admin",
		Name: "Administrator",
		Menus: model.RoleMenus{
			{MenuID: menu.ID},
		},
		Description: "Administrator",
		Status:      model.RoleStatusEnabled,
	}

	var role model.Role
	e.POST(baseAPI + "/roles").WithJSON(roleFormItem).Expect().Status(http.StatusOK).JSON().Decode(&util.ResponseResult{Data: &role})
	assert.NotEmpty(role.ID)
	assert.Equal(roleFormItem.Code, role.Code)
	assert.Equal(roleFormItem.Name, role.Name)
	assert.Equal(roleFormItem.Description, role.Description)
	assert.Equal(roleFormItem.Status, role.Status)
	assert.Equal(len(roleFormItem.Menus), len(role.Menus))

	var roles model.Roles
	e.GET(baseAPI + "/roles").Expect().Status(http.StatusOK).JSON().Decode(&util.ResponseResult{Data: &roles})
	assert.GreaterOrEqual(len(roles), 1)

	newName := "Administrator 1"
	newStatus := model.RoleStatusDisabled
	role.Name = newName
	role.Status = newStatus
	e.PUT(baseAPI + "/roles/" + role.ID).WithJSON(role).Expect().Status(http.StatusOK)

	var getRole model.Role
	e.GET(baseAPI + "/roles/" + role.ID).Expect().Status(http.StatusOK).JSON().Decode(&util.ResponseResult{Data: &getRole})
	assert.Equal(newName, getRole.Name)
	assert.Equal(newStatus, getRole.Status)

	e.DELETE(baseAPI + "/roles/" + role.ID).Expect().Status(http.StatusOK)
	e.GET(baseAPI + "/roles/" + role.ID).Expect().Status(http.StatusNotFound)

	e.DELETE(baseAPI + "/menus/" + menu.ID).Expect().Status(http.StatusOK)
	e.GET(baseAPI + "/menus/" + menu.ID).Expect().Status(http.StatusNotFound)
}
