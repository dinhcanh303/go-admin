package model

import (
	"time"

	"go-admin/internal/config"
	"go-admin/pkg/util"
)

const (
	PermissionResultTypeSelect = "select" // Select
)

// Permission management for RBAC
type Permission struct {
	ID          string    `json:"id" gorm:"size:20;primarykey;"` // Unique ID
	Code        string    `json:"code" gorm:"size:32;index"`     // Display name of permission
	Name        string    `json:"name" gorm:"size:255;index"`    // Display name of permission
	Description string    `json:"description" gorm:"size:1024"`  // Details about permission
	HttpMethod  string    `json:"http_method" gorm:"size:255"`   // HTTP method
	HttpPath    string    `json:"http_path" gorm:"size:1024"`    // HTTP path
	Sequence    int       `json:"sequence" gorm:"index"`         // Sequence for sorting
	CreatedAt   time.Time `json:"created_at" gorm:"index;"`      // Create time
	UpdatedAt   time.Time `json:"updated_at" gorm:"index;"`      // Update time
}

func (a *Permission) TableName() string {
	return config.C.FormatTableName("permissions")
}

// Defining the query parameters for the `Permission` struct.
type PermissionQueryParam struct {
	util.PaginationParam
	LikeName    string     `form:"name"`       // Display name of permission
	ResultType  string     `form:"resultType"` // Result type (options: select)
	InIDs       []string   `form:"-"`          // ID list
	GtUpdatedAt *time.Time `form:"-"`          // Update time is greater than
}

// Defining the query options for the `Permission` struct.
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
	Code        string `json:"code" binding:"required,max=32"`         // Code of permission (unique)
	Name        string `json:"name" binding:"required,max=128"`        // Display name of permission
	Description string `json:"description"`                            // Details about permission
	Sequence    int    `json:"sequence"`                               // Sequence for sorting
	HttpMethod  string `json:"http_method" binding:"required,max=255"` // HTTP method
	HttpPath    string `json:"http_path" binding:"required,max=1024"`  // HTTP path
}

// A validation function for the `PermissionForm` struct.
func (a *PermissionForm) Validate() error {
	return nil
}

func (a *PermissionForm) FillTo(permission *Permission) error {
	permission.Code = a.Code
	permission.Name = a.Name
	permission.Description = a.Description
	permission.Sequence = a.Sequence
	permission.HttpMethod = a.HttpMethod
	permission.HttpPath = a.HttpPath
	return nil
}
