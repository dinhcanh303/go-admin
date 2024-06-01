package model

import (
	"time"

	"go-admin/internal/config"
	"go-admin/pkg/util"
)

// User role permission for RBAC
type RolePermission struct {
	ID           string    `json:"id" gorm:"size:20;primarykey"`                 // Unique ID
	RoleID       string    `json:"role_id" gorm:"size:20;index"`                 // From Role.ID
	PermissionID string    `json:"permission_id" gorm:"size:20;index"`           // From Permission.ID
	CreatedAt    time.Time `json:"created_at" gorm:"index;"`                     // Create time
	UpdatedAt    time.Time `json:"updated_at" gorm:"index;"`                     // Update time
	RoleName     string    `json:"permission_name" gorm:"<-:false;-:migration;"` // From Role.Name
}

func (a *RolePermission) TableName() string {
	return config.C.FormatTableName("role_permissions")
}

// Defining the query parameters for the `RolePermission` struct.
type RolePermissionQueryParam struct {
	util.PaginationParam
	InUserIDs    []string `form:"-"` // From User.ID
	RoleID       string   `form:"-"` // From Role.ID
	PermissionID string   `form:"-"` // From Permission.ID
}

// Defining the query options for the `RolePermission` struct.
type RolePermissionQueryOptions struct {
	util.QueryOptions
	JoinRole bool // Join role table
}

// Defining the query result for the `RolePermission` struct.
type RolePermissionQueryResult struct {
	Data       RolePermissions
	PageResult *util.PaginationResult
}

// Defining the slice of `RolePermission` struct.
type RolePermissions []*RolePermission

func (a RolePermissions) ToUserIDMap() map[string]RolePermissions {
	m := make(map[string]RolePermissions)
	for _, rolePermission := range a {
		m[rolePermission.RoleID] = append(m[rolePermission.RoleID], rolePermission)
	}
	return m
}

func (a RolePermissions) ToPermissionIDs() []string {
	var ids []string
	for _, item := range a {
		ids = append(ids, item.PermissionID)
	}
	return ids
}

// Defining the data structure for creating a `RolePermission` struct.
type RolePermissionForm struct {
}

// A validation function for the `RolePermissionForm` struct.
func (a *RolePermissionForm) Validate() error {
	return nil
}

func (a *RolePermissionForm) FillTo(RolePermission *RolePermission) error {
	return nil
}
