package service

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/johnmerga/realgaming-marketplace-backend/marketplace-backend/internal/models"
	"github.com/johnmerga/realgaming-marketplace-backend/marketplace-backend/internal/repository"
	"golang.org/x/crypto/bcrypt"
)

type UserService interface {
	CreateUser(ctx context.Context, req *models.CreateUserRequest) (*models.UserResponse, error)
	GetUserByID(ctx context.Context, id uuid.UUID) (*models.UserResponse, error)
	GetUserByEmail(ctx context.Context, email string) (*models.UserResponse, error)
	UpdateUser(ctx context.Context, id uuid.UUID, req *models.UpdateUserRequest) (*models.UserResponse, error)
	UpdateUserStatus(ctx context.Context, id uuid.UUID, status models.UserStatus) error
	ListUsers(ctx context.Context, role *models.UserRole, status *models.UserStatus, page, limit int) ([]*models.UserResponse, error)
	Login(ctx context.Context, req *models.LoginRequest) (*models.UserResponse, error)
}

type userService struct {
	userRepo repository.UserRepository
}

func NewUserService(userRepo repository.UserRepository) UserService {
	return &userService{
		userRepo: userRepo,
	}
}

func (s *userService) CreateUser(ctx context.Context, req *models.CreateUserRequest) (*models.UserResponse, error) {
	// Check if user already exists
	existingUser, err := s.userRepo.GetByEmail(ctx, req.Email)
	if err != nil {
		return nil, fmt.Errorf("error checking existing user: %w", err)
	}
	if existingUser != nil {
		return nil, errors.New("user with this email already exists")
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("error hashing password: %w", err)
	}

	// Create user model
	user := &models.User{
		Email:        req.Email,
		PasswordHash: string(hashedPassword),
		FirstName:    req.FirstName,
		LastName:     req.LastName,
		Role:         req.Role,
		Status:       models.StatusActive,
	}

	if req.Phone != "" {
		user.Phone = &req.Phone
	}

	// Create user in database
	createdUser, err := s.userRepo.Create(ctx, user)
	if err != nil {
		return nil, fmt.Errorf("error creating user: %w", err)
	}

	return s.userToResponse(createdUser), nil
}

func (s *userService) GetUserByID(ctx context.Context, id uuid.UUID) (*models.UserResponse, error) {
	user, err := s.userRepo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("error getting user: %w", err)
	}
	if user == nil {
		return nil, errors.New("user not found")
	}

	return s.userToResponse(user), nil
}

func (s *userService) GetUserByEmail(ctx context.Context, email string) (*models.UserResponse, error) {
	user, err := s.userRepo.GetByEmail(ctx, email)
	if err != nil {
		return nil, fmt.Errorf("error getting user: %w", err)
	}
	if user == nil {
		return nil, errors.New("user not found")
	}

	return s.userToResponse(user), nil
}

func (s *userService) UpdateUser(ctx context.Context, id uuid.UUID, req *models.UpdateUserRequest) (*models.UserResponse, error) {
	// Get existing user
	existingUser, err := s.userRepo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("error getting user: %w", err)
	}
	if existingUser == nil {
		return nil, errors.New("user not found")
	}

	// Update user fields
	existingUser.FirstName = req.FirstName
	existingUser.LastName = req.LastName

	if req.Phone != "" {
		existingUser.Phone = &req.Phone
	}
	if req.AvatarURL != "" {
		existingUser.AvatarURL = &req.AvatarURL
	}

	// Update user in database
	updatedUser, err := s.userRepo.Update(ctx, existingUser)
	if err != nil {
		return nil, fmt.Errorf("error updating user: %w", err)
	}

	return s.userToResponse(updatedUser), nil
}

func (s *userService) UpdateUserStatus(ctx context.Context, id uuid.UUID, status models.UserStatus) error {
	// Check if user exists
	user, err := s.userRepo.GetByID(ctx, id)
	if err != nil {
		return fmt.Errorf("error getting user: %w", err)
	}
	if user == nil {
		return errors.New("user not found")
	}

	return s.userRepo.UpdateStatus(ctx, id, status)
}

func (s *userService) ListUsers(ctx context.Context, role *models.UserRole, status *models.UserStatus, page, limit int) ([]*models.UserResponse, error) {
	offset := (page - 1) * limit

	users, err := s.userRepo.List(ctx, role, status, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("error listing users: %w", err)
	}

	responses := make([]*models.UserResponse, len(users))
	for i, user := range users {
		responses[i] = s.userToResponse(user)
	}

	return responses, nil
}

func (s *userService) Login(ctx context.Context, req *models.LoginRequest) (*models.UserResponse, error) {
	user, err := s.userRepo.GetByEmail(ctx, req.Email)
	if err != nil {
		return nil, fmt.Errorf("error getting user: %w", err)
	}
	if user == nil {
		return nil, errors.New("invalid credentials")
	}

	// Check password
	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password))
	if err != nil {
		return nil, errors.New("invalid credentials")
	}

	// Check if user is active
	if user.Status != models.StatusActive {
		return nil, errors.New("user account is inactive")
	}

	return s.userToResponse(user), nil
}

// Helper function to convert user model to response
func (s *userService) userToResponse(user *models.User) *models.UserResponse {
	return &models.UserResponse{
		ID:        user.ID,
		Email:     user.Email,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Role:      user.Role,
		Status:    user.Status,
		AvatarURL: user.AvatarURL,
		Phone:     user.Phone,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}
}
