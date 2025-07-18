package usecase

import (
	"context"
	"go-rest-api-template/internal/domain/entity"
)

// UserUsecase defines business logic interface for user operations
type UserUsecase interface {
	// Authentication methods
	Login(ctx context.Context, username, password string) (*entity.User, string, error) // returns user, token, error
	RefreshToken(ctx context.Context, tokenString string) (string, error)

	// User management
	CreateUser(ctx context.Context, user *entity.User) error
	GetUserByID(ctx context.Context, id int) (*entity.User, error)
	GetUserByEmail(ctx context.Context, email string) (*entity.User, error)
	GetUserByUsername(ctx context.Context, username string) (*entity.User, error)
	UpdateUser(ctx context.Context, user *entity.User) error
	DeleteUser(ctx context.Context, id int) error
	GetAllUsers(ctx context.Context, limit, offset int) ([]*entity.User, error)
	GetUserCount(ctx context.Context) (int, error)

	// Password management
	ForgotPassword(ctx context.Context, email string) error
	ResetPassword(ctx context.Context, token, newPassword string) error
	ChangePassword(ctx context.Context, userID int, currentPassword, newPassword string) error
}
