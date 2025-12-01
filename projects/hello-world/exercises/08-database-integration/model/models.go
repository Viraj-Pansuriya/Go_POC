package model

import (
	"time"

	"gorm.io/gorm"
)

// =============================================================================
// BASE MODEL - Like JPA's @MappedSuperclass
// =============================================================================
// gorm.Model provides: ID, CreatedAt, UpdatedAt, DeletedAt (soft delete)
// In Java: @Id @GeneratedValue + @CreatedDate + @LastModifiedDate + @Where(clause="deleted_at IS NULL")

// =============================================================================
// ONE-TO-ONE RELATIONSHIP: User ↔ Profile
// =============================================================================
// Java equivalent:
// @Entity class User {
//     @OneToOne(mappedBy = "user", cascade = CascadeType.ALL)
//     private Profile profile;
// }

// User represents a system user
// Think of it like @Entity in JPA/Hibernate
type User struct {
	gorm.Model                  // Embeds ID, CreatedAt, UpdatedAt, DeletedAt
	Name    string  `gorm:"size:100;not null" json:"name"`                  // VARCHAR(100) NOT NULL
	Email   string  `gorm:"size:100;uniqueIndex;not null" json:"email"`     // UNIQUE INDEX
	Age     int     `gorm:"default:0" json:"age"`                           // DEFAULT 0
	Profile Profile `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE;"` // ONE-TO-ONE: User has one Profile

	// ONE-TO-MANY: User has many Posts
	// Java: @OneToMany(mappedBy = "user", cascade = CascadeType.ALL)
	Posts []Post `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE;" json:"posts,omitempty"`
}

// Profile - ONE-TO-ONE with User (the "owned" side)
// Java equivalent:
// @Entity class Profile {
//     @OneToOne @JoinColumn(name = "user_id")
//     private User user;
// }
type Profile struct {
	gorm.Model
	UserID    uint   `gorm:"uniqueIndex;not null" json:"user_id"` // Foreign key to User (unique = one-to-one)
	Bio       string `gorm:"type:text" json:"bio"`                // TEXT type for longer content
	AvatarURL string `gorm:"size:255" json:"avatar_url"`
	Website   string `gorm:"size:255" json:"website"`
}

// =============================================================================
// ONE-TO-MANY RELATIONSHIP: User → Posts
// =============================================================================
// Java equivalent:
// @Entity class Post {
//     @ManyToOne @JoinColumn(name = "user_id")
//     private User user;
// }

// Post belongs to a User (Many-to-One from Post's perspective)
type Post struct {
	gorm.Model
	Title   string `gorm:"size:200;not null" json:"title"`
	Content string `gorm:"type:text" json:"content"`
	UserID  uint   `gorm:"index;not null" json:"user_id"` // Foreign key (indexed for performance)

	// Self-referential: Post can have parent (for comments/replies)
	// Java: @ManyToOne @JoinColumn(name = "parent_id")
	ParentID *uint  `gorm:"index" json:"parent_id,omitempty"` // Nullable (pointer = optional)
	Replies  []Post `gorm:"foreignKey:ParentID" json:"replies,omitempty"`
}

// =============================================================================
// ONE-TO-MANY RELATIONSHIP: Author → Books (Another example)
// =============================================================================

// Author can write many books
type Author struct {
	gorm.Model
	Name        string `gorm:"size:100;not null" json:"name"`
	Country     string `gorm:"size:50" json:"country"`
	BirthDate   *time.Time `json:"birth_date,omitempty"` // Pointer = nullable

	// One Author has Many Books
	Books []Book `gorm:"foreignKey:AuthorID" json:"books,omitempty"`
}

// Book belongs to an Author
type Book struct {
	gorm.Model
	Title       string  `gorm:"size:200;not null" json:"title"`
	ISBN        string  `gorm:"size:20;uniqueIndex" json:"isbn"` // Unique ISBN
	PublishedAt *time.Time `json:"published_at,omitempty"`
	Price       float64 `gorm:"type:decimal(10,2);default:0" json:"price"` // DECIMAL(10,2)
	AuthorID    uint    `gorm:"index;not null" json:"author_id"`           // Foreign key to Author

	// MANY-TO-MANY: Book has many Tags, Tag has many Books
	// Java: @ManyToMany @JoinTable(name = "book_tags", ...)
	Tags []Tag `gorm:"many2many:book_tags;" json:"tags,omitempty"`
}

// =============================================================================
// MANY-TO-MANY RELATIONSHIP: Books ↔ Tags
// =============================================================================
// Java equivalent:
// @ManyToMany
// @JoinTable(name = "book_tags",
//     joinColumns = @JoinColumn(name = "book_id"),
//     inverseJoinColumns = @JoinColumn(name = "tag_id"))
// private Set<Tag> tags;

// Tag can be applied to many books
type Tag struct {
	gorm.Model
	Name  string `gorm:"size:50;uniqueIndex;not null" json:"name"` // Tag names are unique
	Color string `gorm:"size:7" json:"color"`                      // Hex color like #FF5733

	// Back-reference (optional - for bidirectional navigation)
	Books []Book `gorm:"many2many:book_tags;" json:"books,omitempty"`
}

// =============================================================================
// MANY-TO-MANY RELATIONSHIP: Students ↔ Courses (with join table attributes)
// =============================================================================
// When you need extra fields in the join table, use explicit join model
// Java: @Entity class Enrollment with @ManyToOne to both Student and Course

// Student can enroll in many courses
type Student struct {
	gorm.Model
	Name        string `gorm:"size:100;not null" json:"name"`
	StudentCode string `gorm:"size:20;uniqueIndex" json:"student_code"` // Like "STU001"
	Email       string `gorm:"size:100;uniqueIndex" json:"email"`

	// Many-to-Many through Enrollment (explicit join table)
	Enrollments []Enrollment `gorm:"foreignKey:StudentID" json:"enrollments,omitempty"`
}

// Course can have many students
type Course struct {
	gorm.Model
	Name        string  `gorm:"size:100;not null" json:"name"`
	Code        string  `gorm:"size:20;uniqueIndex" json:"code"` // Like "CS101"
	Credits     int     `gorm:"default:3" json:"credits"`
	MaxStudents int     `gorm:"default:30" json:"max_students"`
	Price       float64 `gorm:"type:decimal(10,2)" json:"price"`

	// Many-to-Many through Enrollment
	Enrollments []Enrollment `gorm:"foreignKey:CourseID" json:"enrollments,omitempty"`
}

// Enrollment - Join table with extra attributes
// This is like having an @Entity for the join table in JPA
// Java: @Entity class Enrollment { @ManyToOne Student; @ManyToOne Course; Date enrolledAt; ... }
type Enrollment struct {
	gorm.Model
	StudentID  uint       `gorm:"uniqueIndex:idx_student_course;not null" json:"student_id"`
	CourseID   uint       `gorm:"uniqueIndex:idx_student_course;not null" json:"course_id"`
	EnrolledAt time.Time  `gorm:"not null" json:"enrolled_at"`
	Grade      *string    `gorm:"size:2" json:"grade,omitempty"` // Nullable: A, B, C, D, F
	Completed  bool       `gorm:"default:false" json:"completed"`
	
	// Navigation properties
	Student Student `gorm:"foreignKey:StudentID" json:"student,omitempty"`
	Course  Course  `gorm:"foreignKey:CourseID" json:"course,omitempty"`
}

// =============================================================================
// SELF-REFERENTIAL RELATIONSHIP: Category Tree
// =============================================================================
// Useful for hierarchical data like categories, org charts, etc.

// Category with parent-child relationship (tree structure)
type Category struct {
	gorm.Model
	Name        string     `gorm:"size:100;not null" json:"name"`
	Description string     `gorm:"type:text" json:"description"`
	ParentID    *uint      `gorm:"index" json:"parent_id,omitempty"` // Nullable (root has no parent)
	
	// Self-referential relationships
	Parent   *Category  `gorm:"foreignKey:ParentID" json:"parent,omitempty"`
	Children []Category `gorm:"foreignKey:ParentID" json:"children,omitempty"`
}

// =============================================================================
// POLYMORPHIC RELATIONSHIP: Comments on multiple entities
// =============================================================================
// Java: This is tricky in JPA, often done with @Any or discriminator columns
// GORM makes it easy with polymorphic tags

// Comment can belong to different types of entities (Post, Book, Course, etc.)
type Comment struct {
	gorm.Model
	Content       string `gorm:"type:text;not null" json:"content"`
	CommentableID uint   `gorm:"index" json:"commentable_id"` // ID of the parent entity
	CommentableType string `gorm:"size:50;index" json:"commentable_type"` // "posts", "books", "courses"
	
	// Author of comment
	UserID uint `gorm:"index" json:"user_id"`
	User   User `gorm:"foreignKey:UserID" json:"user,omitempty"`
}

