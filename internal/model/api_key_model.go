package model

import (
	"time"
)

// ApiKeyModel - Database model (infrastructure concern) - Read-only for JWT middleware
type ApiKeyModel struct {
	ID          int        `db:"id" json:"id"`
	Name        string     `db:"name" json:"name"`
	Description *string    `db:"description" json:"description"`
	ApiKey      string     `db:"api_key" json:"api_key"`
	AuthKey     string     `db:"auth_key" json:"auth_key"`
	Status      string     `db:"status" json:"status"`
	H2H         string     `db:"h2h" json:"h2h"`
	LastAccess  *time.Time `db:"last_access" json:"last_access"`
	IPWhitelist *string    `db:"ip_whitelist" json:"ip_whitelist"`
	CreatedAt   *time.Time `db:"created_at" json:"created_at"`
	CreatedBy   int        `db:"created_by" json:"created_by"`
	UpdatedAt   *time.Time `db:"updated_at" json:"updated_at"`
	UpdatedBy   *int       `db:"updated_by" json:"updated_by"`
}

// ApiKeyResponse - DTO for HTTP responses (read-only)
type ApiKeyResponse struct {
	ID          int        `json:"id"`
	Name        string     `json:"name"`
	Description *string    `json:"description"`
	Status      string     `json:"status"`
	H2H         string     `json:"h2h"`
	LastAccess  *time.Time `json:"last_access"`
	IPWhitelist *string    `json:"ip_whitelist"`
	CreatedAt   *time.Time `json:"created_at"`
	CreatedBy   int        `json:"created_by"`
	UpdatedAt   *time.Time `json:"updated_at"`
	UpdatedBy   *int       `json:"updated_by"`
	// Note: ApiKey and AuthKey are intentionally excluded for security
}
