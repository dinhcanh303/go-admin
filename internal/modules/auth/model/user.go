package model

import (
	"time"

	"go-admin/internal/config"
	"go-admin/pkg/crypto/hash"
	"go-admin/pkg/errors"
	"go-admin/pkg/util"

	"github.com/go-playground/validator/v10"
)

const (
	UserStatusActive   = "active"
	UserStatusInactive = "inactive"
)

// User management for RBAC
type User struct {
	ID        string    `json:"id" gorm:"size:20;primarykey;"`    // Unique ID
	Email     string    `json:"email" gorm:"size:255;index"`      // Email for login
	FirstName string    `json:"first_name" gorm:"size:100;index"` // First Name of user
	LastName  string    `json:"last_name" gorm:"size:100;index"`  // Last Name of user
	FullName  string    `json:"full_name" gorm:"size:255;index"`  // Full Name of user
	Password  string    `json:"-" gorm:"size:255;"`               // Password for login (encrypted)
	Phone     string    `json:"phone" gorm:"size:32;"`            // Phone number of user
	Remark    string    `json:"remark" gorm:"size:1024;"`         // Remark of user
	Status    string    `json:"status" gorm:"size:20;index"`      // Status of user (active, inactive)
	CreatedAt time.Time `json:"created_at" gorm:"index;"`         // Create time
	UpdatedAt time.Time `json:"updated_at" gorm:"index;"`         // Update time
	Roles     UserRoles `json:"roles" gorm:"-"`                   // Roles of user
}

func (a *User) TableName() string {
	return config.C.FormatTableName("users")
}

// Defining the query parameters for the `User` struct.
type UserQueryParam struct {
	util.PaginationParam
	LikeEmail    string `form:"email"`                                  // Email for login
	LikeFullName string `form:"full_name"`                              // Full Name of user
	Status       string `form:"status" binding:"oneof=active inactive"` // Status of user (active, inactive)
}

// Defining the query options for the `User` struct.
type UserQueryOptions struct {
	util.QueryOptions
}

// Defining the query result for the `User` struct.
type UserQueryResult struct {
	Data       Users
	PageResult *util.PaginationResult
}

// Defining the slice of `User` struct.
type Users []*User

func (a Users) ToIDs() []string {
	var ids []string
	for _, item := range a {
		ids = append(ids, item.ID)
	}
	return ids
}

// Defining the data structure for creating a `User` struct.
type UserForm struct {
	Email     string    `json:"email" binding:"required,max=128"`                // Username for login
	FirstName string    `json:"first_name" binding:"required,max=64"`            // First Name of user
	LastName  string    `json:"last_name" binding:"required,max=64"`             // Last Name of user
	Password  string    `json:"password" binding:"required,max=64"`              // Password for login (md5 hash)
	Phone     string    `json:"phone" binding:"max=32"`                          // Phone number of user
	Remark    string    `json:"remark" binding:"max=1024"`                       // Remark of user
	Status    string    `json:"status" binding:"required,oneof=active inactive"` // Status of user (active, inactive)
	Roles     UserRoles `json:"roles"`                                           // Roles of user
}

// A validation function for the `UserForm` struct.
func (a *UserForm) Validate() error {
	if a.Email != "" && validator.New().Var(a.Email, "email") != nil {
		return errors.BadRequest("", "Invalid email address")
	}
	return nil
}

// Convert `UserForm` to `User` object.
func (a *UserForm) FillTo(user *User) error {
	user.Email = a.Email
	user.FirstName = a.FirstName
	user.LastName = a.LastName
	user.FullName = a.FirstName + " " + a.LastName
	user.Phone = a.Phone
	user.Remark = a.Remark
	user.Status = a.Status

	if pass := a.Password; pass != "" {
		hashPass, err := hash.GeneratePassword(pass)
		if err != nil {
			return errors.BadRequest("", "Failed to generate hash password: %s", err.Error())
		}
		user.Password = hashPass
	}

	return nil
}
