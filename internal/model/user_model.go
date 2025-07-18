package model

import (
	"go-rest-api-template/pkg/validator"
	"time"
)

// UserModel - Database model (infrastructure concern)
type UserModel struct {
	ID                 int        `db:"id" json:"id"`
	Username           string     `db:"username" json:"username"`
	AuthKey            string     `db:"auth_key" json:"-"`
	PasswordHash       string     `db:"password_hash" json:"-"`
	PasswordResetToken *string    `db:"password_reset_token" json:"-"`
	Email              string     `db:"email" json:"email"`
	Status             string     `db:"status" json:"status"`
	CreatedAt          time.Time  `db:"created_at" json:"created_at"`
	CreatedBy          *int       `db:"created_by" json:"created_by,omitempty"`
	UpdatedAt          *time.Time `db:"updated_at" json:"updated_at,omitempty"`
	UpdatedBy          *int       `db:"updated_by" json:"updated_by,omitempty"`
	DeletedAt          *time.Time `db:"deleted_at" json:"deleted_at,omitempty"`
	DeletedBy          *int       `db:"deleted_by" json:"deleted_by,omitempty"`
	VerificationToken  *string    `db:"verification_token" json:"-"`
}

// UserCreateRequest - DTO for HTTP requests with comprehensive validation
type UserCreateRequest struct {
	Username string `json:"username" validate:"required,min=3,max=50,alphanum"`
	Email    string `json:"email" validate:"required,email,max=100"`
	Password string `json:"password" validate:"required,min=8,max=100"`
}

// UserLoginRequest - DTO for login requests
type UserLoginRequest struct {
	Username string `json:"username" validate:"required,min=3,max=50"`
	Password string `json:"password" validate:"required,min=8,max=100"`
}

// UserUpdateRequest - DTO for HTTP update requests
type UserUpdateRequest struct {
	Username string `json:"username" validate:"omitempty,min=3,max=50,alphanum"`
	Email    string `json:"email" validate:"omitempty,email,max=100"`
	Password string `json:"password" validate:"omitempty,min=8,max=100"`
	Status   string `json:"status" validate:"omitempty,oneof=active inactive suspended"`
}

// ForgotPasswordRequest - DTO for forgot password requests
type ForgotPasswordRequest struct {
	Email string `json:"email" validate:"required,email,max=100"`
}

// ResetPasswordRequest - DTO for reset password requests
type ResetPasswordRequest struct {
	Token       string `json:"token" validate:"required,min=1,max=255"`
	NewPassword string `json:"new_password" validate:"required,min=8,max=100"`
}

// ChangePasswordRequest - DTO for change password requests
type ChangePasswordRequest struct {
	CurrentPassword string `json:"current_password" validate:"required,min=1"`
	NewPassword     string `json:"new_password" validate:"required,min=8,max=100,nefield=CurrentPassword"`
}

// UserResponse - DTO for HTTP responses
type UserResponse struct {
	ID        int        `json:"id"`
	Username  string     `json:"username"`
	Email     string     `json:"email"`
	Status    string     `json:"status"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt *time.Time `json:"updated_at"`
}

// Validate validates UserCreateRequest
func (r *UserCreateRequest) Validate() error {
	return validator.ValidateStruct(r)
}

// Validate validates UserLoginRequest
func (r *UserLoginRequest) Validate() error {
	return validator.ValidateStruct(r)
}

// Validate validates UserUpdateRequest
func (r *UserUpdateRequest) Validate() error {
	return validator.ValidateStruct(r)
}

// Validate validates ForgotPasswordRequest
func (r *ForgotPasswordRequest) Validate() error {
	return validator.ValidateStruct(r)
}

// Validate validates ResetPasswordRequest
func (r *ResetPasswordRequest) Validate() error {
	return validator.ValidateStruct(r)
}

// Validate validates ChangePasswordRequest
func (r *ChangePasswordRequest) Validate() error {
	return validator.ValidateStruct(r)
}
