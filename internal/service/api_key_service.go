package service

import (
	"context"
	"go-rest-api-template/internal/domain/entity"
	"go-rest-api-template/internal/domain/repository"
)

// ApiKeyService handles API key business logic (read-only for JWT middleware)
type ApiKeyService interface {
	// Core validation methods for JWT middleware
	ValidateApiKey(ctx context.Context, apiKey string) (*entity.ApiKey, error)
	ValidateAuthKey(ctx context.Context, authKey string) (*entity.ApiKey, error)

	// Admin methods for management
	GetApiKeyByID(ctx context.Context, id int) (*entity.ApiKey, error)
	GetAllApiKeys(ctx context.Context, limit, offset int) ([]*entity.ApiKey, error)
	GetActiveApiKeys(ctx context.Context, limit, offset int) ([]*entity.ApiKey, error)

	// Logging method
	LogApiKeyAccess(ctx context.Context, id int) error
}

type apiKeyService struct {
	apiKeyRepo repository.ApiKeyRepository
}

// NewApiKeyService creates a new API key service
func NewApiKeyService(apiKeyRepo repository.ApiKeyRepository) ApiKeyService {
	return &apiKeyService{
		apiKeyRepo: apiKeyRepo,
	}
}

func (s *apiKeyService) ValidateApiKey(ctx context.Context, apiKey string) (*entity.ApiKey, error) {
	key, err := s.apiKeyRepo.GetByApiKey(ctx, apiKey)
	if err != nil {
		return nil, err
	}

	if key == nil || !key.IsActive() {
		return nil, nil // Invalid or inactive key
	}

	return key, nil
}

func (s *apiKeyService) ValidateAuthKey(ctx context.Context, authKey string) (*entity.ApiKey, error) {
	key, err := s.apiKeyRepo.GetByAuthKey(ctx, authKey)
	if err != nil {
		return nil, err
	}

	if key == nil || !key.IsActive() {
		return nil, nil // Invalid or inactive key
	}

	return key, nil
}

func (s *apiKeyService) GetApiKeyByID(ctx context.Context, id int) (*entity.ApiKey, error) {
	return s.apiKeyRepo.GetByID(ctx, id)
}

func (s *apiKeyService) GetAllApiKeys(ctx context.Context, limit, offset int) ([]*entity.ApiKey, error) {
	return s.apiKeyRepo.GetAll(ctx, limit, offset)
}

func (s *apiKeyService) GetActiveApiKeys(ctx context.Context, limit, offset int) ([]*entity.ApiKey, error) {
	return s.apiKeyRepo.GetByStatus(ctx, "active", limit, offset)
}

func (s *apiKeyService) LogApiKeyAccess(ctx context.Context, id int) error {
	return s.apiKeyRepo.UpdateLastAccess(ctx, id)
}
