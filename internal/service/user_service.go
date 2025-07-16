package service

import (
	"context"
	"database/sql"
	"errors"
	"go-rest-api-template/internal/domain/entity"
	"go-rest-api-template/internal/domain/repository"
	"go-rest-api-template/internal/domain/usecase"
)

type userService struct {
	userRepo repository.UserRepository
}

// NewUserService creates a new user service
func NewUserService(userRepo repository.UserRepository) usecase.UserUsecase {
	return &userService{
		userRepo: userRepo,
	}
}

func (s *userService) CreateUser(ctx context.Context, user *entity.User) error {
	// Business validation
	if err := user.ValidateForCreate(); err != nil {
		return err
	}

	// Check if username already exists
	existingUser, err := s.userRepo.GetByUsername(ctx, user.Username)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return err
	}
	if existingUser != nil {
		return errors.New("username already exists")
	}

	// Check if email already exists
	existingUser, err = s.userRepo.GetByEmail(ctx, user.Email)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return err
	}
	if existingUser != nil {
		return errors.New("email already exists")
	}

	// Create user
	return s.userRepo.Create(ctx, user)
}

func (s *userService) GetUserByID(ctx context.Context, id int) (*entity.User, error) {
	return s.userRepo.GetByID(ctx, id)
}

func (s *userService) GetUserByUsername(ctx context.Context, username string) (*entity.User, error) {
	return s.userRepo.GetByUsername(ctx, username)
}

func (s *userService) GetUserByEmail(ctx context.Context, email string) (*entity.User, error) {
	return s.userRepo.GetByEmail(ctx, email)
}

func (s *userService) GetAllUsers(ctx context.Context, limit, offset int) ([]*entity.User, error) {
	return s.userRepo.GetAll(ctx, limit, offset)
}

func (s *userService) GetUserCount(ctx context.Context) (int, error) {
	return s.userRepo.GetCount(ctx)
}

func (s *userService) UpdateUser(ctx context.Context, user *entity.User) error {
	// Check if user exists
	existingUser, err := s.userRepo.GetByID(ctx, user.ID)
	if err != nil {
		return err
	}
	if existingUser == nil {
		return errors.New("user not found")
	}

	// Check if username already exists (if changed)
	if user.Username != existingUser.Username {
		existingUser, err := s.userRepo.GetByUsername(ctx, user.Username)
		if err != nil && !errors.Is(err, sql.ErrNoRows) {
			return err
		}
		if existingUser != nil {
			return errors.New("username already exists")
		}
	}

	// Check if email already exists (if changed)
	if user.Email != existingUser.Email {
		existingUser, err := s.userRepo.GetByEmail(ctx, user.Email)
		if err != nil && !errors.Is(err, sql.ErrNoRows) {
			return err
		}
		if existingUser != nil {
			return errors.New("email already exists")
		}
	}

	// Update user
	return s.userRepo.Update(ctx, user)
}

func (s *userService) DeleteUser(ctx context.Context, id int) error {
	// Check if user exists
	existingUser, err := s.userRepo.GetByID(ctx, id)
	if err != nil {
		return err
	}
	if existingUser == nil {
		return errors.New("user not found")
	}

	// Soft delete user
	return s.userRepo.Delete(ctx, id)
}

func (s *userService) ChangePassword(ctx context.Context, userID int, currentPassword, newPassword string) error {
	// Get user
	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return err
	}
	if user == nil {
		return errors.New("user not found")
	}

	// Verify current password
	if !user.CheckPassword(currentPassword) {
		return errors.New("invalid current password")
	}

	// Hash new password
	if err := user.HashPassword(newPassword); err != nil {
		return err
	}

	// Update user
	return s.userRepo.Update(ctx, user)
}

func (s *userService) ForgotPassword(ctx context.Context, email string) error {
	// Get user by email
	user, err := s.userRepo.GetByEmail(ctx, email)
	if err != nil {
		return err
	}
	if user == nil {
		return errors.New("user not found")
	}

	// Generate reset token
	if err := user.GenerateResetPasswordToken(); err != nil {
		return err
	}

	// Update user with reset token
	if err := s.userRepo.Update(ctx, user); err != nil {
		return err
	}

	// In a real application, you would send an email with the reset token here
	// For now, we just return success
	return nil
}

func (s *userService) ResetPassword(ctx context.Context, token, newPassword string) error {
	// Get user by reset token
	user, err := s.userRepo.GetByResetPasswordToken(ctx, token)
	if err != nil {
		return err
	}
	if user == nil {
		return errors.New("invalid or expired reset token")
	}

	// Validate reset token
	if !user.IsResetPasswordTokenValid(token) {
		return errors.New("reset token has expired")
	}

	// Hash new password
	if err := user.HashPassword(newPassword); err != nil {
		return err
	}

	// Clear reset token
	user.ClearResetPasswordToken()

	// Update user
	return s.userRepo.Update(ctx, user)
}
