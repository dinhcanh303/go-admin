package tests

import (
	"net/http"
	"testing"

	"go-admin/internal/modules/auth/model"
	"go-admin/pkg/crypto/hash"
	"go-admin/pkg/util"

	"github.com/stretchr/testify/assert"
)

func TestUser(t *testing.T) {
	e := tester(t)

	menuFormItem := model.MenuForm{
		Code:        "user",
		Name:        "User management",
		Description: "User management",
		Sequence:    7,
		Type:        "page",
		Path:        "/system/user",
		Properties:  `{"icon":"user"}`,
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

	roleFormItem := model.RoleForm{
		Code: "user",
		Name: "Normal",
		Menus: model.RoleMenus{
			{MenuID: menu.ID},
		},
		Description: "Normal",
		Sequence:    8,
		Status:      model.RoleStatusEnabled,
	}

	var role model.Role
	e.POST(baseAPI + "/roles").WithJSON(roleFormItem).Expect().Status(http.StatusOK).JSON().Decode(&util.ResponseResult{Data: &role})
	assert.NotEmpty(role.ID)
	assert.Equal(roleFormItem.Code, role.Code)
	assert.Equal(roleFormItem.Name, role.Name)
	assert.Equal(roleFormItem.Description, role.Description)
	assert.Equal(roleFormItem.Sequence, role.Sequence)
	assert.Equal(roleFormItem.Status, role.Status)
	assert.Equal(len(roleFormItem.Menus), len(role.Menus))

	userFormItem := model.UserForm{
		Email:     "test@test.com",
		FirstName: "Test",
		Password:  hash.MD5String("test"),
		Phone:     "0720",
		Remark:    "test user",
		Status:    model.UserStatusActive,
		Roles:     model.UserRoles{{RoleID: role.ID}},
	}

	var user model.User
	e.POST(baseAPI + "/users").WithJSON(userFormItem).Expect().Status(http.StatusOK).JSON().Decode(&util.ResponseResult{Data: &user})
	assert.NotEmpty(user.ID)
	assert.Equal(userFormItem.Email, user.Email)
	assert.Equal(userFormItem.FirstName, user.FirstName)
	assert.Equal(userFormItem.Phone, user.Phone)
	assert.Equal(userFormItem.Email, user.Email)
	assert.Equal(userFormItem.Remark, user.Remark)
	assert.Equal(userFormItem.Status, user.Status)
	assert.Equal(len(userFormItem.Roles), len(user.Roles))

	var users model.Users
	e.GET(baseAPI+"/users").WithQuery("email", userFormItem.Email).Expect().Status(http.StatusOK).JSON().Decode(&util.ResponseResult{Data: &users})
	assert.GreaterOrEqual(len(users), 1)

	newName := "Test 1"
	newStatus := model.UserStatusInactive
	user.FirstName = newName
	user.Status = newStatus
	e.PUT(baseAPI + "/users/" + user.ID).WithJSON(user).Expect().Status(http.StatusOK)

	var getUser model.User
	e.GET(baseAPI + "/users/" + user.ID).Expect().Status(http.StatusOK).JSON().Decode(&util.ResponseResult{Data: &getUser})
	assert.Equal(newName, getUser.FirstName)
	assert.Equal(newStatus, getUser.Status)

	e.DELETE(baseAPI + "/users/" + user.ID).Expect().Status(http.StatusOK)
	e.GET(baseAPI + "/users/" + user.ID).Expect().Status(http.StatusNotFound)

	e.DELETE(baseAPI + "/roles/" + role.ID).Expect().Status(http.StatusOK)
	e.GET(baseAPI + "/roles/" + role.ID).Expect().Status(http.StatusNotFound)

	e.DELETE(baseAPI + "/menus/" + menu.ID).Expect().Status(http.StatusOK)
	e.GET(baseAPI + "/menus/" + menu.ID).Expect().Status(http.StatusNotFound)
}
