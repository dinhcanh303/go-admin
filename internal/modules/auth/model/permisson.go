package model

import (
	"time"

	"go-admin/internal/config"
	"go-admin/pkg/util"
)

// Permission management for RBAC
type Permission struct {
	ID          string    `json:"id" gorm:"size:20;primarykey;"` // Unique ID
	Name        string    `json:"name" gorm:"size:128;index"`    // Display name of permission
	Description string    `json:"description" gorm:"size:1024"`  // Details about permission
	Slug        string    `json:"slug" gorm:"size:255"`          //Slug
	HttpMethod  string    `json:"http_method" gorm:"size:255"`   // HTTP method
	HttpPath    string    `json:"http_path" gorm:"size:1024"`    // HTTP path
	CreatedAt   time.Time `json:"created_at" gorm:"index;"`      // Create time
	UpdatedAt   time.Time `json:"updated_at" gorm:"index;"`      // Update time
}

func (a *Permission) TableName() string {
	return config.C.FormatTableName("permissions")
}

// Defining the query parameters for the `Role` struct.
type PermissionQueryParam struct {
	util.PaginationParam
	LikeName    string     `form:"name"`                                       // Display name of role
	Status      string     `form:"status" binding:"oneof=disabled enabled ''"` // Status of role (disabled, enabled)
	ResultType  string     `form:"resultType"`                                 // Result type (options: select)
	InIDs       []string   `form:"-"`                                          // ID list
	GtUpdatedAt *time.Time `form:"-"`                                          // Update time is greater than
}

// Defining the query options for the `Role` struct.
type PermissionQueryOptions struct {
	util.QueryOptions
}

// Defining the query result for the `Permission` struct.
type PermissionQueryResult struct {
	Data       Permissions
	PageResult *util.PaginationResult
}

// Defining the slice of `Permission` struct.
type Permissions []*Permission

// Defining the data structure for creating a `Permission` struct.
type PermissionForm struct {
	Name        string `json:"name" binding:"required,max=128"`                  // Display name of Permission
	Description string `json:"description"`                                      // Details about Permission
	Sequence    int    `json:"sequence"`                                         // Sequence for sorting
	Status      string `json:"status" binding:"required,oneof=disabled enabled"` // Status of Permission (enabled, disabled)
}

// A validation function for the `PermissionForm` struct.
func (a *PermissionForm) Validate() error {
	return nil
}

func (a *PermissionForm) FillTo(p *Permission) error {
	p.Name = a.Name
	p.Description = a.Description
	// p.Slug = a.Slug
	return nil
}
