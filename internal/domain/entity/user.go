package entity

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"time"

	"golang.org/x/crypto/bcrypt"
)

// User represents a user entity
type User struct {
	ID                     int        `json:"id"`
	Username               string     `json:"username"`
	Email                  string     `json:"email"`
	PasswordHash           string     `json:"-"`
	Status                 string     `json:"status"`
	ResetPasswordToken     *string    `json:"-"`
	ResetPasswordExpiresAt *time.Time `json:"-"`
	CreatedAt              time.Time  `json:"created_at"`
	UpdatedAt              time.Time  `json:"updated_at"`
	CreatedBy              *int       `json:"created_by,omitempty"`
	UpdatedBy              *int       `json:"updated_by,omitempty"`
	DeletedBy              *int       `json:"deleted_by,omitempty"`
}

// Business validation rules
func (u *User) ValidateForCreate() error {
	if len(u.Username) < 3 {
		return errors.New("username must be at least 3 characters")
	}
	return nil
}

// ValidatePassword validates password before hashing
func (u *User) ValidatePassword(password string) error {
	if len(password) < 6 {
		return errors.New("password must be at least 6 characters")
	}
	return nil
}

// Hash password using bcrypt
func (u *User) HashPassword(password string) error {
	if err := u.ValidatePassword(password); err != nil {
		return err
	}
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	u.PasswordHash = string(hashedPassword)
	return nil
}

// Check if password matches
func (u *User) CheckPassword(password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(u.PasswordHash), []byte(password))
	return err == nil
}

// Check if user is active
func (u *User) IsActive() bool {
	return u.Status == "active"
}

// GenerateResetPasswordToken creates a new reset password token
func (u *User) GenerateResetPasswordToken() error {
	// Generate 32 random bytes
	bytes := make([]byte, 32)
	if _, err := rand.Read(bytes); err != nil {
		return err
	}

	// Convert to hex string
	token := hex.EncodeToString(bytes)
	u.ResetPasswordToken = &token

	// Set expiration to 1 hour from now
	expiresAt := time.Now().Add(1 * time.Hour)
	u.ResetPasswordExpiresAt = &expiresAt

	return nil
}

// IsResetPasswordTokenValid checks if reset token is valid and not expired
func (u *User) IsResetPasswordTokenValid(token string) bool {
	if u.ResetPasswordToken == nil || *u.ResetPasswordToken != token {
		return false
	}

	if u.ResetPasswordExpiresAt == nil || time.Now().After(*u.ResetPasswordExpiresAt) {
		return false
	}

	return true
}

// ClearResetPasswordToken removes the reset password token
func (u *User) ClearResetPasswordToken() {
	u.ResetPasswordToken = nil
	u.ResetPasswordExpiresAt = nil
}

// SetAuditFields sets the audit fields for create/update operations
func (u *User) SetCreatedBy(userID int) {
	u.CreatedBy = &userID
}

func (u *User) SetUpdatedBy(userID int) {
	u.UpdatedBy = &userID
}

func (u *User) SetDeletedBy(userID int) {
	u.DeletedBy = &userID
}
