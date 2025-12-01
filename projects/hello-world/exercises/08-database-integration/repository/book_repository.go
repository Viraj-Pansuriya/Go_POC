package repository

import (
	"errors"

	"github.com/viraj/go-mono-repo/projects/hello-world/exercises/08-database-integration/model"
	"gorm.io/gorm"
)

// =============================================================================
// BOOK REPOSITORY - Demonstrates Many-to-Many with Tags
// =============================================================================

// BookRepository defines operations for Book entity
type BookRepository interface {
	// Basic CRUD
	Create(book *model.Book) error
	FindByID(id uint) (*model.Book, error)
	FindByIDWithTags(id uint) (*model.Book, error)      // Eager load tags
	FindByIDWithAuthor(id uint) (*model.Book, error)    // Eager load author
	FindByIDFull(id uint) (*model.Book, error)          // Load everything
	FindAll() ([]model.Book, error)
	Update(book *model.Book) error
	Delete(id uint) error

	// Association management (Many-to-Many)
	AddTag(bookID uint, tag *model.Tag) error
	RemoveTag(bookID uint, tagID uint) error
	ReplaceAllTags(bookID uint, tags []model.Tag) error
	GetTags(bookID uint) ([]model.Tag, error)

	// Custom queries
	FindByAuthorID(authorID uint) ([]model.Book, error)
	FindByTagName(tagName string) ([]model.Book, error)
	FindByPriceRange(minPrice, maxPrice float64) ([]model.Book, error)
	SearchByTitle(title string) ([]model.Book, error)
}

type bookRepository struct {
	db *gorm.DB
}

// NewBookRepository creates a new BookRepository
func NewBookRepository(db *gorm.DB) BookRepository {
	return &bookRepository{db: db}
}

// =============================================================================
// BASIC CRUD
// =============================================================================

// Create inserts a new book
func (r *bookRepository) Create(book *model.Book) error {
	// If book.Tags is populated, GORM will:
	// 1. Create the book
	// 2. Create entries in the join table (book_tags)
	// This is like CascadeType.PERSIST in JPA
	return r.db.Create(book).Error
}

// FindByID retrieves a book by ID (without associations)
func (r *bookRepository) FindByID(id uint) (*model.Book, error) {
	var book model.Book
	err := r.db.First(&book, id).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &book, err
}

// FindByIDWithTags retrieves book with tags eagerly loaded
// Java: @EntityGraph(attributePaths = {"tags"})
func (r *bookRepository) FindByIDWithTags(id uint) (*model.Book, error) {
	var book model.Book

	// Preload Many-to-Many association
	// GORM will execute:
	// 1. SELECT * FROM books WHERE id = ?
	// 2. SELECT * FROM tags INNER JOIN book_tags ON tags.id = book_tags.tag_id WHERE book_tags.book_id = ?
	err := r.db.Preload("Tags").First(&book, id).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &book, err
}

// FindByIDWithAuthor retrieves book with author eagerly loaded
func (r *bookRepository) FindByIDWithAuthor(id uint) (*model.Book, error) {
	var book model.Book

	// For Many-to-One, we use Preload with the association name
	err := r.db.Preload("Author").First(&book, id).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &book, err
}

// FindByIDFull retrieves book with all associations
func (r *bookRepository) FindByIDFull(id uint) (*model.Book, error) {
	var book model.Book

	// Chain multiple Preload() for multiple associations
	// Like @EntityGraph with multiple attributePaths
	err := r.db.
		Preload("Tags").
		Preload("Author").
		First(&book, id).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &book, err
}

// FindAll retrieves all books
func (r *bookRepository) FindAll() ([]model.Book, error) {
	var books []model.Book
	err := r.db.Find(&books).Error
	return books, err
}

// Update saves changes to a book
func (r *bookRepository) Update(book *model.Book) error {
	return r.db.Save(book).Error
}

// Delete soft-deletes a book
func (r *bookRepository) Delete(id uint) error {
	return r.db.Delete(&model.Book{}, id).Error
}

// =============================================================================
// MANY-TO-MANY ASSOCIATION MANAGEMENT
// =============================================================================
// In JPA, you'd modify the collection and save the entity
// GORM provides explicit association methods

// AddTag adds a tag to a book (creates join table entry)
// Java equivalent: book.getTags().add(tag); repository.save(book);
func (r *bookRepository) AddTag(bookID uint, tag *model.Tag) error {
	var book model.Book
	if err := r.db.First(&book, bookID).Error; err != nil {
		return err
	}

	// Association() returns an Association object for the named association
	// Append() adds to the Many-to-Many relationship
	return r.db.Model(&book).Association("Tags").Append(tag)
}

// RemoveTag removes a tag from a book (deletes join table entry)
// Java equivalent: book.getTags().remove(tag); repository.save(book);
func (r *bookRepository) RemoveTag(bookID uint, tagID uint) error {
	var book model.Book
	if err := r.db.First(&book, bookID).Error; err != nil {
		return err
	}

	var tag model.Tag
	if err := r.db.First(&tag, tagID).Error; err != nil {
		return err
	}

	// Delete() removes the association (not the tag itself)
	return r.db.Model(&book).Association("Tags").Delete(&tag)
}

// ReplaceAllTags replaces all tags of a book
// Removes existing associations and creates new ones
func (r *bookRepository) ReplaceAllTags(bookID uint, tags []model.Tag) error {
	var book model.Book
	if err := r.db.First(&book, bookID).Error; err != nil {
		return err
	}

	// Replace() clears existing and adds new associations
	return r.db.Model(&book).Association("Tags").Replace(tags)
}

// GetTags retrieves all tags for a book
func (r *bookRepository) GetTags(bookID uint) ([]model.Tag, error) {
	var book model.Book
	if err := r.db.First(&book, bookID).Error; err != nil {
		return nil, err
	}

	var tags []model.Tag
	err := r.db.Model(&book).Association("Tags").Find(&tags)
	return tags, err
}

// =============================================================================
// CUSTOM QUERIES
// =============================================================================

// FindByAuthorID finds all books by a specific author
// Java: List<Book> findByAuthorId(Long authorId);
func (r *bookRepository) FindByAuthorID(authorID uint) ([]model.Book, error) {
	var books []model.Book
	err := r.db.Where("author_id = ?", authorID).Find(&books).Error
	return books, err
}

// FindByTagName finds all books that have a specific tag
// Java: @Query("SELECT b FROM Book b JOIN b.tags t WHERE t.name = :tagName")
func (r *bookRepository) FindByTagName(tagName string) ([]model.Book, error) {
	var books []model.Book

	// Join through the many-to-many relationship
	// This creates: SELECT books.* FROM books 
	//               INNER JOIN book_tags ON book_tags.book_id = books.id
	//               INNER JOIN tags ON tags.id = book_tags.tag_id
	//               WHERE tags.name = ?
	err := r.db.
		Joins("INNER JOIN book_tags ON book_tags.book_id = books.id").
		Joins("INNER JOIN tags ON tags.id = book_tags.tag_id").
		Where("tags.name = ?", tagName).
		Find(&books).Error

	return books, err
}

// FindByPriceRange finds books within a price range
// Java: List<Book> findByPriceBetween(BigDecimal min, BigDecimal max);
func (r *bookRepository) FindByPriceRange(minPrice, maxPrice float64) ([]model.Book, error) {
	var books []model.Book

	// Multiple Where() conditions are ANDed together
	err := r.db.
		Where("price >= ?", minPrice).
		Where("price <= ?", maxPrice).
		Order("price ASC").
		Find(&books).Error

	return books, err
}

// SearchByTitle searches books by title (case-insensitive)
// Java: @Query("SELECT b FROM Book b WHERE LOWER(b.title) LIKE LOWER(CONCAT('%', :title, '%'))")
func (r *bookRepository) SearchByTitle(title string) ([]model.Book, error) {
	var books []model.Book

	// SQLite uses LIKE which is case-insensitive by default
	// For other databases, use ILIKE (Postgres) or LOWER()
	err := r.db.
		Where("title LIKE ?", "%"+title+"%").
		Preload("Tags").
		Preload("Author").
		Find(&books).Error

	return books, err
}

