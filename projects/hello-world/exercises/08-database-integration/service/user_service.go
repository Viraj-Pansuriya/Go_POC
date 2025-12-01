package service

import (
	"errors"
	"fmt"

	"github.com/viraj/go-mono-repo/projects/hello-world/exercises/08-database-integration/model"
	"github.com/viraj/go-mono-repo/projects/hello-world/exercises/08-database-integration/repository"
)

// =============================================================================
// USER SERVICE - Business Logic Layer
// =============================================================================
// In Java/Spring: @Service class with @Autowired repository
// Go: Struct with repository dependency injected via constructor

// UserService defines business operations for users
type UserService interface {
	Register(name, email string, age int) (*model.User, error)
	GetByID(id uint) (*model.User, error)
	GetByIDWithProfile(id uint) (*model.User, error)
	GetByEmail(email string) (*model.User, error)
	GetAll() ([]model.User, error)
	GetAllPaginated(page, pageSize int) ([]model.User, int64, error)
	UpdateProfile(userID uint, bio, avatarURL, website string) error
	Delete(id uint) error

	// Business operations
	SearchUsers(query string) ([]model.User, error)
	GetAdults() ([]model.User, error)
}

// userService implements UserService
type userService struct {
	repo repository.UserRepository
}

// NewUserService creates a UserService with injected dependencies
// Java equivalent: @Service class with @Autowired constructor
func NewUserService(repo repository.UserRepository) UserService {
	return &userService{repo: repo}
}

// =============================================================================
// SERVICE METHODS
// =============================================================================

// Register creates a new user with validation
// Java: @Transactional public User register(RegisterRequest request) { ... }
func (s *userService) Register(name, email string, age int) (*model.User, error) {
	// Business validation (not just database constraints)
	if name == "" {
		return nil, errors.New("name is required")
	}
	if email == "" {
		return nil, errors.New("email is required")
	}
	if age < 0 {
		return nil, errors.New("age cannot be negative")
	}

	// Check if email already exists (business rule)
	exists, err := s.repo.ExistsByEmail(email)
	if err != nil {
		return nil, fmt.Errorf("failed to check email: %w", err)
	}
	if exists {
		return nil, errors.New("email already registered")
	}

	// Create user with default profile
	user := &model.User{
		Name:  name,
		Email: email,
		Age:   age,
		Profile: model.Profile{
			Bio: "New user",
		},
	}

	if err := s.repo.Create(user); err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	return user, nil
}

// GetByID retrieves a user by ID
func (s *userService) GetByID(id uint) (*model.User, error) {
	user, err := s.repo.FindByID(id)
	if err != nil {
		return nil, fmt.Errorf("failed to find user: %w", err)
	}
	if user == nil {
		return nil, errors.New("user not found")
	}
	return user, nil
}

// GetByIDWithProfile retrieves user with profile loaded
func (s *userService) GetByIDWithProfile(id uint) (*model.User, error) {
	user, err := s.repo.FindByIDWithProfile(id)
	if err != nil {
		return nil, fmt.Errorf("failed to find user: %w", err)
	}
	if user == nil {
		return nil, errors.New("user not found")
	}
	return user, nil
}

// GetByEmail retrieves a user by email
func (s *userService) GetByEmail(email string) (*model.User, error) {
	if email == "" {
		return nil, errors.New("email is required")
	}
	return s.repo.FindByEmail(email)
}

// GetAll retrieves all users
func (s *userService) GetAll() ([]model.User, error) {
	return s.repo.FindAll()
}

// GetAllPaginated retrieves users with pagination
func (s *userService) GetAllPaginated(page, pageSize int) ([]model.User, int64, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 10
	}
	return s.repo.FindAllWithPagination(page, pageSize)
}

// UpdateProfile updates a user's profile
// Demonstrates working with One-to-One relationship
func (s *userService) UpdateProfile(userID uint, bio, avatarURL, website string) error {
	user, err := s.repo.FindByIDWithProfile(userID)
	if err != nil {
		return fmt.Errorf("failed to find user: %w", err)
	}
	if user == nil {
		return errors.New("user not found")
	}

	// Update profile fields
	user.Profile.Bio = bio
	user.Profile.AvatarURL = avatarURL
	user.Profile.Website = website

	return s.repo.Update(user)
}

// Delete removes a user (soft delete)
func (s *userService) Delete(id uint) error {
	// Check if user exists first
	user, err := s.repo.FindByID(id)
	if err != nil {
		return fmt.Errorf("failed to find user: %w", err)
	}
	if user == nil {
		return errors.New("user not found")
	}

	return s.repo.Delete(id)
}

// SearchUsers searches users by name
func (s *userService) SearchUsers(query string) ([]model.User, error) {
	if query == "" {
		return s.repo.FindAll()
	}
	return s.repo.FindByNameContaining(query)
}

// GetAdults retrieves all users 18 and older
func (s *userService) GetAdults() ([]model.User, error) {
	return s.repo.FindByAgeGreaterThan(17)
}

