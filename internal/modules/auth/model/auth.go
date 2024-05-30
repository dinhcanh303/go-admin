package model

import (
	"go-admin/pkg/crypto/hash"
	"go-admin/pkg/errors"
	"strings"
)

type Captcha struct {
	CaptchaID string `json:"captcha_id"` // Captcha ID
}

type LoginForm struct {
	Email       string `json:"email" binding:"required"`        // Login name
	Password    string `json:"password" binding:"required"`     // Login password (md5 hash)
	CaptchaID   string `json:"captcha_id" binding:"required"`   // Captcha verify id
	CaptchaCode string `json:"captcha_code" binding:"required"` // Captcha verify code
}

func (a *LoginForm) Trim() *LoginForm {
	a.Email = strings.TrimSpace(a.Email)
	a.CaptchaCode = strings.TrimSpace(a.CaptchaCode)
	return a
}

type RegisterForm struct {
	FirstName       string `json:"first_name" binding:"required"`
	LastName        string `json:"last_name" binding:"required"`
	Email           string `json:"email" binding:"required"`
	Password        string `json:"password" binding:"required"`
	PasswordConfirm string `json:"password_confirm" binding:"required"`
	CaptchaID       string `json:"captcha_id" binding:"required"`
	CaptchaCode     string `json:"captcha_code" binding:"required"`
}

func (a *RegisterForm) Trim() *RegisterForm {
	a.Email = strings.TrimSpace(a.Email)
	a.FirstName = strings.TrimSpace(a.FirstName)
	a.LastName = strings.TrimSpace(a.LastName)
	a.CaptchaCode = strings.TrimSpace(a.CaptchaCode)
	return a
}

// Convert `UserForm` to `User` object.
func (a *RegisterForm) FillTo(user *User) error {
	user.Email = a.Email
	user.FirstName = a.FirstName
	user.LastName = a.LastName
	user.FullName = a.FirstName + " " + a.LastName

	if pass := a.Password; pass != "" {
		hashPass, err := hash.GeneratePassword(pass)
		if err != nil {
			return errors.BadRequest("", "Failed to generate hash password: %s", err.Error())
		}
		user.Password = hashPass
	}

	return nil
}

type UpdateLoginPassword struct {
	OldPassword string `json:"old_password" binding:"required"` // Old password (md5 hash)
	NewPassword string `json:"new_password" binding:"required"` // New password (md5 hash)
}

type LoginToken struct {
	AccessToken string `json:"access_token"` // Access token (JWT)
	TokenType   string `json:"token_type"`   // Token type (Usage: Authorization=${token_type} ${access_token})
	ExpiresAt   int64  `json:"expires_at"`   // Expired time (Unit: second)
}

type UpdateCurrentUser struct {
	FirstName string `json:"first_name" binding:"required,max=64"` // Name of user
	LastName  string `json:"last_name" binding:"required,max=64"`  // Name of user
	Phone     string `json:"phone" binding:"max=32"`               // Phone number of user
	Remark    string `json:"remark" binding:"max=1024"`            // Remark of user
}
