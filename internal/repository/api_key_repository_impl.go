package repository

import (
	"context"
	"database/sql"
	"go-rest-api-template/internal/domain/entity"
	"go-rest-api-template/internal/domain/repository"
	"go-rest-api-template/internal/model"

	"github.com/jmoiron/sqlx"
)

// apiKeyRepositoryImpl - Infrastructure implementation (read-only)
type apiKeyRepositoryImpl struct {
	db *sqlx.DB
}

// NewApiKeyRepository creates repository implementation
func NewApiKeyRepository(db *sqlx.DB) repository.ApiKeyRepository {
	return &apiKeyRepositoryImpl{db: db}
}

func (r *apiKeyRepositoryImpl) GetByApiKey(ctx context.Context, apiKey string) (*entity.ApiKey, error) {
	var apiKeyModel model.ApiKeyModel

	query := `SELECT * FROM api_key WHERE api_key = ? AND status = 'active'`
	err := r.db.GetContext(ctx, &apiKeyModel, query, apiKey)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return r.modelToEntity(&apiKeyModel), nil
}

func (r *apiKeyRepositoryImpl) GetByAuthKey(ctx context.Context, authKey string) (*entity.ApiKey, error) {
	var apiKeyModel model.ApiKeyModel

	query := `SELECT * FROM api_key WHERE auth_key = ? AND status = 'active'`
	err := r.db.GetContext(ctx, &apiKeyModel, query, authKey)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return r.modelToEntity(&apiKeyModel), nil
}

func (r *apiKeyRepositoryImpl) GetByID(ctx context.Context, id int) (*entity.ApiKey, error) {
	var apiKeyModel model.ApiKeyModel

	query := `SELECT * FROM api_key WHERE id = ?`
	err := r.db.GetContext(ctx, &apiKeyModel, query, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return r.modelToEntity(&apiKeyModel), nil
}

func (r *apiKeyRepositoryImpl) GetAll(ctx context.Context, limit, offset int) ([]*entity.ApiKey, error) {
	var apiKeyModels []model.ApiKeyModel

	query := `SELECT * FROM api_key ORDER BY created_at DESC LIMIT ? OFFSET ?`
	err := r.db.SelectContext(ctx, &apiKeyModels, query, limit, offset)
	if err != nil {
		return nil, err
	}

	return r.modelsToEntities(apiKeyModels), nil
}

func (r *apiKeyRepositoryImpl) GetByStatus(ctx context.Context, status string, limit, offset int) ([]*entity.ApiKey, error) {
	var apiKeyModels []model.ApiKeyModel

	query := `SELECT * FROM api_key WHERE status = ? ORDER BY created_at DESC LIMIT ? OFFSET ?`
	err := r.db.SelectContext(ctx, &apiKeyModels, query, status, limit, offset)
	if err != nil {
		return nil, err
	}

	return r.modelsToEntities(apiKeyModels), nil
}

func (r *apiKeyRepositoryImpl) GetCount(ctx context.Context) (int, error) {
	var count int
	query := `SELECT COUNT(*) FROM api_key`
	err := r.db.GetContext(ctx, &count, query)
	return count, err
}

func (r *apiKeyRepositoryImpl) UpdateLastAccess(ctx context.Context, id int) error {
	query := `UPDATE api_key SET last_access = NOW() WHERE id = ?`
	_, err := r.db.ExecContext(ctx, query, id)
	return err
}

// Helper methods for model conversion
func (r *apiKeyRepositoryImpl) modelToEntity(model *model.ApiKeyModel) *entity.ApiKey {
	return &entity.ApiKey{
		ID:          model.ID,
		Name:        model.Name,
		Description: model.Description,
		ApiKey:      model.ApiKey,
		AuthKey:     model.AuthKey,
		Status:      model.Status,
		H2H:         model.H2H,
		LastAccess:  model.LastAccess,
		IPWhitelist: model.IPWhitelist,
		CreatedAt:   model.CreatedAt,
		CreatedBy:   model.CreatedBy,
		UpdatedAt:   model.UpdatedAt,
		UpdatedBy:   model.UpdatedBy,
	}
}

func (r *apiKeyRepositoryImpl) modelsToEntities(models []model.ApiKeyModel) []*entity.ApiKey {
	entities := make([]*entity.ApiKey, len(models))
	for i, model := range models {
		entities[i] = r.modelToEntity(&model)
	}
	return entities
}
