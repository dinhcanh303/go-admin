package tests

import (
	"net/http"
	"testing"

	"go-admin/internal/modules/auth/model"
	"go-admin/pkg/util"

	"github.com/stretchr/testify/assert"
)

func TestMenu(t *testing.T) {
	e := tester(t)

	menuFormItem := model.MenuForm{
		Code:        "menu",
		Name:        "Menu management",
		Description: "Menu management",
		Sequence:    9,
		Type:        "page",
		Path:        "/system/menu",
		Properties:  `{"icon":"menu"}`,
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
	assert.Equal(menuFormItem.Sequence, menu.Sequence)
	assert.Equal(menuFormItem.Type, menu.Type)
	assert.Equal(menuFormItem.Path, menu.Path)
	assert.Equal(menuFormItem.Properties, menu.Properties)
	assert.Equal(menuFormItem.Status, menu.Status)

	var menus model.Menus
	e.GET(baseAPI + "/menus").Expect().Status(http.StatusOK).JSON().Decode(&util.ResponseResult{Data: &menus})
	assert.GreaterOrEqual(len(menus), 1)

	newName := "Menu management 1"
	newStatus := model.MenuStatusDisabled
	menu.Name = newName
	menu.Status = newStatus
	e.PUT(baseAPI + "/menus/" + menu.ID).WithJSON(menu).Expect().Status(http.StatusOK)

	var getMenu model.Menu
	e.GET(baseAPI + "/menus/" + menu.ID).Expect().Status(http.StatusOK).JSON().Decode(&util.ResponseResult{Data: &getMenu})
	assert.Equal(newName, getMenu.Name)
	assert.Equal(newStatus, getMenu.Status)

	e.DELETE(baseAPI + "/menus/" + menu.ID).Expect().Status(http.StatusOK)
	e.GET(baseAPI + "/menus/" + menu.ID).Expect().Status(http.StatusNotFound)
}
