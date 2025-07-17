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
	userRepo   repository.UserRepository
	jwtService JWTService
}

// NewUserService creates a new user service
func NewUserService(userRepo repository.UserRepository, jwtService JWTService) usecase.UserUsecase {
	return &userService{
		userRepo:   userRepo,
		jwtService: jwtService,
	}
}

// Login authenticates a user and returns user data with JWT token
func (s *userService) Login(ctx context.Context, username, password string) (*entity.User, string, error) {
	// Get user by username or email
	var user *entity.User
	var err error

	// Try to find user by username first
	user, err = s.userRepo.GetByUsername(ctx, username)
	if err != nil || user == nil {
		// If not found, try by email
		user, err = s.userRepo.GetByEmail(ctx, username)
		if err != nil {
			return nil, "", errors.New("invalid credentials")
		}
	}

	if user == nil {
		return nil, "", errors.New("invalid credentials")
	}

	// Check if user is active
	if !user.IsActive() {
		return nil, "", errors.New("account is not active")
	}

	// Verify password
	if !user.CheckPassword(password) {
		return nil, "", errors.New("invalid credentials")
	}

	// Note: For login, we don't generate JWT token here anymore
	// The token generation should be handled by the login handler
	// which will have access to the API key information
	// This method now only validates user credentials

	return user, "", nil
}

// RefreshToken is no longer available in the new JWT system
// Token refresh should be handled at the handler level with API key validation
func (s *userService) RefreshToken(ctx context.Context, tokenString string) (string, error) {
	return "", errors.New("refresh token not supported in new JWT system")
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

	// Generate verification token
	if err := user.GenerateVerificationToken(); err != nil {
		return err
	}

	// Update user with verification token
	if err := s.userRepo.UpdateVerificationToken(ctx, user); err != nil {
		return err
	}

	// In a real application, you would send an email with the reset token here
	// For now, we just return success
	return nil
}

func (s *userService) ResetPassword(ctx context.Context, token, newPassword string) error {
	// Get user by verification token
	user, err := s.userRepo.GetByVerificationToken(ctx, token)
	if err != nil {
		return err
	}
	if user == nil {
		return errors.New("invalid verification token")
	}

	// Validate verification token
	if !user.IsVerificationTokenValid(token) {
		return errors.New("invalid verification token")
	}

	// Hash new password
	if err := user.HashPassword(newPassword); err != nil {
		return err
	}

	// Clear verification token
	user.ClearVerificationToken()

	// Update user
	return s.userRepo.Update(ctx, user)
}
