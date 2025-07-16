package repository

import (
	"context"
	"go-rest-api-template/internal/domain/entity"
)

// UserRepository defines the interface for user data operations
type UserRepository interface {
	Create(ctx context.Context, user *entity.User) error
	GetByID(ctx context.Context, id int) (*entity.User, error)
	GetByEmail(ctx context.Context, email string) (*entity.User, error)
	GetByUsername(ctx context.Context, username string) (*entity.User, error)
	Update(ctx context.Context, user *entity.User) error
	Delete(ctx context.Context, id int) error
	GetAll(ctx context.Context, limit, offset int) ([]*entity.User, error)
	GetCount(ctx context.Context) (int, error)
}
