package repository

import (
	"errors"

	"github.com/viraj/go-mono-repo/projects/hello-world/exercises/08-database-integration/model"
	"gorm.io/gorm"
)

// =============================================================================
// USER REPOSITORY
// =============================================================================
// In Java/Spring: public interface UserRepository extends JpaRepository<User, Long>
// Go: Define interface + implementation separately

// UserRepository defines the contract for user data access
// Like Spring Data JPA interface, but we implement methods ourselves
type UserRepository interface {
	// Basic CRUD - like JpaRepository methods
	Create(user *model.User) error
	FindByID(id uint) (*model.User, error)
	FindByIDWithProfile(id uint) (*model.User, error)  // Eager load profile
	FindByIDWithPosts(id uint) (*model.User, error)    // Eager load posts
	FindByEmail(email string) (*model.User, error)
	FindAll() ([]model.User, error)
	FindAllWithPagination(page, pageSize int) ([]model.User, int64, error)
	Update(user *model.User) error
	Delete(id uint) error
	HardDelete(id uint) error  // Permanent delete

	// Custom queries - like @Query in Spring Data
	FindByAgeGreaterThan(age int) ([]model.User, error)
	FindByNameContaining(name string) ([]model.User, error)
	CountByAge(age int) (int64, error)
	ExistsByEmail(email string) (bool, error)
}

// userRepository implements UserRepository
// Private struct - only expose through interface
type userRepository struct {
	db *gorm.DB
}

// NewUserRepository creates a new UserRepository instance
// Factory function pattern - like Spring's bean creation
func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{db: db}
}

// =============================================================================
// BASIC CRUD OPERATIONS
// =============================================================================

// Create inserts a new user into the database
// Java: repository.save(user) for new entity
func (r *userRepository) Create(user *model.User) error {
	// GORM's Create() is like JPA's persist()
	// It will also create associated entities if present (cascade)
	return r.db.Create(user).Error
}

// FindByID retrieves a user by primary key
// Java: repository.findById(id).orElse(null)
func (r *userRepository) FindByID(id uint) (*model.User, error) {
	var user model.User

	// First() finds first record matching condition, ordered by primary key
	// Returns ErrRecordNotFound if not found
	err := r.db.First(&user, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil // Return nil, nil for "not found" (like Optional.empty())
		}
		return nil, err
	}
	return &user, nil
}

// FindByIDWithProfile retrieves user with profile eagerly loaded
// Java: @EntityGraph(attributePaths = {"profile"}) or Hibernate.initialize()
func (r *userRepository) FindByIDWithProfile(id uint) (*model.User, error) {
	var user model.User

	// Preload() is like JPA's eager fetch or @EntityGraph
	// It executes a separate query for the association
	err := r.db.Preload("Profile").First(&user, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}

// FindByIDWithPosts retrieves user with all posts eagerly loaded
// Java: @EntityGraph(attributePaths = {"posts"})
func (r *userRepository) FindByIDWithPosts(id uint) (*model.User, error) {
	var user model.User

	// Multiple Preload() calls for multiple associations
	err := r.db.Preload("Posts").Preload("Profile").First(&user, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}

// FindByEmail finds a user by email address
// Java: Optional<User> findByEmail(String email);
func (r *userRepository) FindByEmail(email string) (*model.User, error) {
	var user model.User

	// Where() is like JPQL: SELECT u FROM User u WHERE u.email = :email
	// Use ? placeholder to prevent SQL injection (like prepared statements)
	err := r.db.Where("email = ?", email).First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}

// FindAll retrieves all users
// Java: repository.findAll()
func (r *userRepository) FindAll() ([]model.User, error) {
	var users []model.User

	// Find() with no conditions returns all records
	// Note: Soft-deleted records are automatically excluded (WHERE deleted_at IS NULL)
	err := r.db.Find(&users).Error
	return users, err
}

// FindAllWithPagination retrieves users with pagination
// Java: repository.findAll(PageRequest.of(page, size))
func (r *userRepository) FindAllWithPagination(page, pageSize int) ([]model.User, int64, error) {
	var users []model.User
	var total int64

	// Count total records first
	r.db.Model(&model.User{}).Count(&total)

	// Offset = (page - 1) * pageSize for 1-based page numbers
	// Limit = pageSize
	offset := (page - 1) * pageSize
	err := r.db.Offset(offset).Limit(pageSize).Find(&users).Error

	return users, total, err
}

// Update saves changes to an existing user
// Java: repository.save(user) for existing entity
func (r *userRepository) Update(user *model.User) error {
	// Save() updates all fields (like merge in JPA)
	// Use Updates() for partial updates
	return r.db.Save(user).Error
}

// Delete performs soft delete (sets deleted_at)
// Java: Custom implementation with @Where clause
func (r *userRepository) Delete(id uint) error {
	// Delete() with gorm.Model performs soft delete
	// The record stays in DB but has deleted_at set
	return r.db.Delete(&model.User{}, id).Error
}

// HardDelete permanently removes a user from database
// Java: @Query with native delete or custom implementation
func (r *userRepository) HardDelete(id uint) error {
	// Unscoped() bypasses soft delete and actually removes the record
	return r.db.Unscoped().Delete(&model.User{}, id).Error
}

// =============================================================================
// CUSTOM QUERY METHODS
// =============================================================================

// FindByAgeGreaterThan finds users older than specified age
// Java: List<User> findByAgeGreaterThan(int age);
func (r *userRepository) FindByAgeGreaterThan(age int) ([]model.User, error) {
	var users []model.User

	// Where() with comparison operator
	err := r.db.Where("age > ?", age).Find(&users).Error
	return users, err
}

// FindByNameContaining finds users whose name contains the search string
// Java: List<User> findByNameContaining(String name); or @Query with LIKE
func (r *userRepository) FindByNameContaining(name string) ([]model.User, error) {
	var users []model.User

	// LIKE query with wildcards
	// %name% matches anywhere in the string
	err := r.db.Where("name LIKE ?", "%"+name+"%").Find(&users).Error
	return users, err
}

// CountByAge counts users with specific age
// Java: long countByAge(int age);
func (r *userRepository) CountByAge(age int) (int64, error) {
	var count int64

	// Model() specifies the table, Count() executes SELECT COUNT(*)
	err := r.db.Model(&model.User{}).Where("age = ?", age).Count(&count).Error
	return count, err
}

// ExistsByEmail checks if a user with given email exists
// Java: boolean existsByEmail(String email);
func (r *userRepository) ExistsByEmail(email string) (bool, error) {
	var count int64
	err := r.db.Model(&model.User{}).Where("email = ?", email).Count(&count).Error
	return count > 0, err
}

