package repository

import (
	"context"
	"go-rest-api-template/internal/constant"
	"go-rest-api-template/internal/domain/entity"
	"go-rest-api-template/internal/domain/repository"
	"go-rest-api-template/internal/model"
	"time"

	common "github.com/budimanlai/go-common"
	"github.com/jmoiron/sqlx"
)

// userRepositoryImpl - Infrastructure implementation
type userRepositoryImpl struct {
	db *sqlx.DB
}

// NewUserRepository creates repository implementation
func NewUserRepository(db *sqlx.DB) repository.UserRepository {
	return &userRepositoryImpl{db: db}
}

func (r *userRepositoryImpl) Create(ctx context.Context, user *entity.User) error {
	now := time.Now()
	// Convert domain entity to database model
	userModel := &model.UserModel{
		Username:          user.Username,
		Email:             user.Email,
		PasswordHash:      user.PasswordHash,
		VerificationToken: user.VerificationToken,
		AuthKey:           common.GenerateRandomString(constant.AuthKeyLength),
		// Set default values
		CreatedAt: now,
		UpdatedAt: &now,
		DeletedAt: nil, // Not deleted initially
		DeletedBy: nil, // Not deleted initially
		// Set status to active by default
		Status:    constant.UserStatusActive,
		CreatedBy: user.CreatedBy,
	}

	query := `INSERT INTO user (username, auth_key, email, password_hash, status, created_by, created_at, updated_at) 
			  VALUES (:username, :auth_key, :email, :password_hash, :status, :created_by, NOW(), NOW())`

	result, err := r.db.NamedExecContext(ctx, query, userModel)
	if err != nil {
		return err
	}

	// Set the ID back to domain entity
	id, _ := result.LastInsertId()
	user.ID = int(id)

	return nil
}

func (r *userRepositoryImpl) GetByID(ctx context.Context, id int) (*entity.User, error) {
	var userModel model.UserModel

	query := `SELECT * FROM user WHERE id = ? AND deleted_at IS NULL`
	err := r.db.GetContext(ctx, &userModel, query, id)
	if err != nil {
		return nil, err
	}

	// Convert database model to domain entity
	return &entity.User{
		ID:                 userModel.ID,
		Username:           userModel.Username,
		Email:              userModel.Email,
		PasswordHash:       userModel.PasswordHash,
		PasswordResetToken: userModel.PasswordResetToken,
		Status:             userModel.Status,
		VerificationToken:  userModel.VerificationToken,
		AuthKey:            userModel.AuthKey,
		CreatedAt:          userModel.CreatedAt,
		UpdatedAt:          userModel.UpdatedAt,
		DeletedAt:          userModel.DeletedAt,
		CreatedBy:          userModel.CreatedBy,
		UpdatedBy:          userModel.UpdatedBy,
		DeletedBy:          userModel.DeletedBy,
	}, nil
}

func (r *userRepositoryImpl) GetByEmail(ctx context.Context, email string) (*entity.User, error) {
	var userModel model.UserModel

	query := `SELECT * FROM user WHERE email = ? AND deleted_at IS NULL`
	err := r.db.GetContext(ctx, &userModel, query, email)
	if err != nil {
		return nil, err
	}

	return &entity.User{
		ID:                 userModel.ID,
		Username:           userModel.Username,
		Email:              userModel.Email,
		PasswordHash:       userModel.PasswordHash,
		PasswordResetToken: userModel.PasswordResetToken,
		Status:             userModel.Status,
		VerificationToken:  userModel.VerificationToken,
		AuthKey:            userModel.AuthKey,
		CreatedAt:          userModel.CreatedAt,
		UpdatedAt:          userModel.UpdatedAt,
		DeletedAt:          userModel.DeletedAt,
		CreatedBy:          userModel.CreatedBy,
		UpdatedBy:          userModel.UpdatedBy,
		DeletedBy:          userModel.DeletedBy,
	}, nil
}

func (r *userRepositoryImpl) GetByUsername(ctx context.Context, username string) (*entity.User, error) {
	var userModel model.UserModel

	query := `SELECT * FROM user WHERE username = ? AND deleted_at IS NULL`
	err := r.db.GetContext(ctx, &userModel, query, username)
	if err != nil {
		return nil, err
	}

	return &entity.User{
		ID:                 userModel.ID,
		Username:           userModel.Username,
		Email:              userModel.Email,
		PasswordHash:       userModel.PasswordHash,
		PasswordResetToken: userModel.PasswordResetToken,
		Status:             userModel.Status,
		VerificationToken:  userModel.VerificationToken,
		AuthKey:            userModel.AuthKey,
		CreatedAt:          userModel.CreatedAt,
		UpdatedAt:          userModel.UpdatedAt,
		DeletedAt:          userModel.DeletedAt,
		CreatedBy:          userModel.CreatedBy,
		UpdatedBy:          userModel.UpdatedBy,
		DeletedBy:          userModel.DeletedBy,
	}, nil
}

func (r *userRepositoryImpl) Update(ctx context.Context, user *entity.User) error {
	userModel := &model.UserModel{
		ID:           user.ID,
		Username:     user.Username,
		Email:        user.Email,
		PasswordHash: user.PasswordHash,
		Status:       user.Status,
		UpdatedBy:    user.UpdatedBy,
	}

	query := `UPDATE user SET username = :username, email = :email, password_hash = :password_hash, 
			  status = :status, updated_by = :updated_by, updated_at = NOW() WHERE id = :id AND deleted_at IS NULL`

	_, err := r.db.NamedExecContext(ctx, query, userModel)
	return err
}

func (r *userRepositoryImpl) Delete(ctx context.Context, id int) error {
	query := `UPDATE user SET deleted_at = NOW(), deleted_by = ? WHERE id = ? AND deleted_at IS NULL`
	_, err := r.db.ExecContext(ctx, query, constant.DefaultUpdatedBy, id)
	return err
}

func (r *userRepositoryImpl) GetAll(ctx context.Context, limit, offset int) ([]*entity.User, error) {
	var userModels []model.UserModel

	query := `SELECT * FROM user WHERE deleted_at IS NULL ORDER BY created_at DESC LIMIT ? OFFSET ?`
	err := r.db.SelectContext(ctx, &userModels, query, limit, offset)
	if err != nil {
		return nil, err
	}

	users := make([]*entity.User, len(userModels))
	for i, userModel := range userModels {
		users[i] = &entity.User{
			ID:                 userModel.ID,
			Username:           userModel.Username,
			Email:              userModel.Email,
			PasswordHash:       userModel.PasswordHash,
			PasswordResetToken: userModel.PasswordResetToken,
			Status:             userModel.Status,
			VerificationToken:  userModel.VerificationToken,
			AuthKey:            userModel.AuthKey,
			CreatedAt:          userModel.CreatedAt,
			UpdatedAt:          userModel.UpdatedAt,
			DeletedAt:          userModel.DeletedAt,
			CreatedBy:          userModel.CreatedBy,
			UpdatedBy:          userModel.UpdatedBy,
			DeletedBy:          userModel.DeletedBy,
		}
	}

	return users, nil
}

func (r *userRepositoryImpl) GetCount(ctx context.Context) (int, error) {
	var count int
	query := `SELECT COUNT(*) FROM user WHERE deleted_at IS NULL`
	err := r.db.GetContext(ctx, &count, query)
	return count, err
}

func (r *userRepositoryImpl) GetByVerificationToken(ctx context.Context, token string) (*entity.User, error) {
	var userModel model.UserModel

	query := `SELECT * FROM user WHERE verification_token = ? AND deleted_at IS NULL`
	err := r.db.GetContext(ctx, &userModel, query, token)
	if err != nil {
		return nil, err
	}

	return &entity.User{
		ID:                 userModel.ID,
		Username:           userModel.Username,
		Email:              userModel.Email,
		PasswordHash:       userModel.PasswordHash,
		PasswordResetToken: userModel.PasswordResetToken,
		Status:             userModel.Status,
		VerificationToken:  userModel.VerificationToken,
		AuthKey:            userModel.AuthKey,
		CreatedAt:          userModel.CreatedAt,
		UpdatedAt:          userModel.UpdatedAt,
		DeletedAt:          userModel.DeletedAt,
		CreatedBy:          userModel.CreatedBy,
		UpdatedBy:          userModel.UpdatedBy,
		DeletedBy:          userModel.DeletedBy,
	}, nil
}

func (r *userRepositoryImpl) UpdateVerificationToken(ctx context.Context, user *entity.User) error {
	userModel := &model.UserModel{
		ID:                user.ID,
		VerificationToken: user.VerificationToken,
		UpdatedBy:         user.UpdatedBy,
	}

	query := `UPDATE user SET verification_token = :verification_token, 
			  updated_by = :updated_by, updated_at = NOW() 
			  WHERE id = :id AND deleted_at IS NULL`

	_, err := r.db.NamedExecContext(ctx, query, userModel)
	return err
}
