package repository

import (
	"context"
	"go-rest-api-template/internal/domain/entity"
)

// ApiKeyRepository defines the interface for API key data operations (read-only)
type ApiKeyRepository interface {
	// Core read operations for JWT middleware
	GetByApiKey(ctx context.Context, apiKey string) (*entity.ApiKey, error)
	GetByAuthKey(ctx context.Context, authKey string) (*entity.ApiKey, error)
	GetByID(ctx context.Context, id int) (*entity.ApiKey, error)

	// List operations for admin purposes
	GetAll(ctx context.Context, limit, offset int) ([]*entity.ApiKey, error)
	GetByStatus(ctx context.Context, status string, limit, offset int) ([]*entity.ApiKey, error)
	GetCount(ctx context.Context) (int, error)

	// Update last access for logging (minimal write operation)
	UpdateLastAccess(ctx context.Context, id int) error
}
