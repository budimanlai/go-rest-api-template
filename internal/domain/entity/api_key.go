package entity

import (
	"time"
)

// ApiKey represents an API key entity (read-only for JWT middleware integration)
type ApiKey struct {
	ID          int        `json:"id"`
	Name        string     `json:"name"`
	Description *string    `json:"description"`
	ApiKey      string     `json:"api_key"`
	AuthKey     string     `json:"auth_key"`
	Status      string     `json:"status"`
	H2H         string     `json:"h2h"`
	LastAccess  *time.Time `json:"last_access"`
	IPWhitelist *string    `json:"ip_whitelist"`
	CreatedAt   *time.Time `json:"created_at"`
	CreatedBy   int        `json:"created_by"`
	UpdatedAt   *time.Time `json:"updated_at"`
	UpdatedBy   *int       `json:"updated_by"`
}

// IsActive checks if the API key is active
func (a *ApiKey) IsActive() bool {
	return a.Status == "active"
}

// IsH2HEnabled checks if H2H (Host-to-Host) is enabled
func (a *ApiKey) IsH2HEnabled() bool {
	return a.H2H == "Y"
}

// IsIPWhitelisted checks if the given IP is whitelisted
func (a *ApiKey) IsIPWhitelisted(ip string) bool {
	if a.IPWhitelist == nil || *a.IPWhitelist == "" {
		return true // No whitelist means all IPs are allowed
	}
	// Simple implementation - in production you might want more sophisticated IP matching
	return *a.IPWhitelist == ip
}

// UpdateLastAccess updates the last access time (for logging purposes)
func (a *ApiKey) UpdateLastAccess() {
	now := time.Now()
	a.LastAccess = &now
}
