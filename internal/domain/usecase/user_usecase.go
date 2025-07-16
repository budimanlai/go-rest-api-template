package usecase

import (
	"context"
	"go-rest-api-template/internal/domain/entity"
)

// UserUsecase defines business logic interface for user operations
type UserUsecase interface {
	CreateUser(ctx context.Context, user *entity.User) error
	GetUserByID(ctx context.Context, id int) (*entity.User, error)
	GetUserByEmail(ctx context.Context, email string) (*entity.User, error)
	GetUserByUsername(ctx context.Context, username string) (*entity.User, error)
	UpdateUser(ctx context.Context, user *entity.User) error
	DeleteUser(ctx context.Context, id int) error
	GetAllUsers(ctx context.Context, limit, offset int) ([]*entity.User, error)
	GetUserCount(ctx context.Context) (int, error)
}
