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
	if err != nil && err != sql.ErrNoRows {
		return err
	}
	if existingUser != nil {
		return errors.New("username already exists")
	}

	// Check if email already exists
	existingUser, err = s.userRepo.GetByEmail(ctx, user.Email)
	if err != nil && err != sql.ErrNoRows {
		return err
	}
	if existingUser != nil {
		return errors.New("email already exists")
	}

	// Hash password
	if err := user.HashPassword(); err != nil {
		return err
	}

	// Set default status
	user.Status = "active"

	return s.userRepo.Create(ctx, user)
}

func (s *userService) GetUserByID(ctx context.Context, id int) (*entity.User, error) {
	if id <= 0 {
		return nil, errors.New("invalid user ID")
	}

	user, err := s.userRepo.GetByID(ctx, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("user not found")
		}
		return nil, err
	}

	return user, nil
}

func (s *userService) GetUserByEmail(ctx context.Context, email string) (*entity.User, error) {
	if email == "" {
		return nil, errors.New("email is required")
	}

	user, err := s.userRepo.GetByEmail(ctx, email)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("user not found")
		}
		return nil, err
	}

	return user, nil
}

func (s *userService) GetUserByUsername(ctx context.Context, username string) (*entity.User, error) {
	if username == "" {
		return nil, errors.New("username is required")
	}

	user, err := s.userRepo.GetByUsername(ctx, username)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("user not found")
		}
		return nil, err
	}

	return user, nil
}

func (s *userService) UpdateUser(ctx context.Context, user *entity.User) error {
	if user.ID <= 0 {
		return errors.New("invalid user ID")
	}

	// Check if user exists
	existingUser, err := s.userRepo.GetByID(ctx, user.ID)
	if err != nil {
		if err == sql.ErrNoRows {
			return errors.New("user not found")
		}
		return err
	}

	// Check if new username is taken by another user
	if user.Username != existingUser.Username {
		userByUsername, err := s.userRepo.GetByUsername(ctx, user.Username)
		if err != nil && err != sql.ErrNoRows {
			return err
		}
		if userByUsername != nil && userByUsername.ID != user.ID {
			return errors.New("username already exists")
		}
	}

	// Check if new email is taken by another user
	if user.Email != existingUser.Email {
		userByEmail, err := s.userRepo.GetByEmail(ctx, user.Email)
		if err != nil && err != sql.ErrNoRows {
			return err
		}
		if userByEmail != nil && userByEmail.ID != user.ID {
			return errors.New("email already exists")
		}
	}

	// If password is provided, hash it
	if user.Password != "" {
		if err := user.HashPassword(); err != nil {
			return err
		}
	} else {
		// Keep existing password if not provided
		user.Password = existingUser.Password
	}

	return s.userRepo.Update(ctx, user)
}

func (s *userService) DeleteUser(ctx context.Context, id int) error {
	if id <= 0 {
		return errors.New("invalid user ID")
	}

	// Check if user exists
	_, err := s.userRepo.GetByID(ctx, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return errors.New("user not found")
		}
		return err
	}

	return s.userRepo.Delete(ctx, id)
}

func (s *userService) GetAllUsers(ctx context.Context, limit, offset int) ([]*entity.User, error) {
	if limit <= 0 {
		limit = 10 // default limit
	}
	if offset < 0 {
		offset = 0
	}

	return s.userRepo.GetAll(ctx, limit, offset)
}

func (s *userService) GetUserCount(ctx context.Context) (int, error) {
	return s.userRepo.GetCount(ctx)
}
