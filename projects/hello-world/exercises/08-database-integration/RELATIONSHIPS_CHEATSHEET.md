# 🗃️ Database Relationships Cheatsheet

Quick reference for all relationship patterns in GORM with Java/Hibernate comparisons.

---

## 1️⃣ ONE-TO-ONE: User ↔ Profile

```go
// OWNER side (User) - has the Profile
type User struct {
    gorm.Model
    Name    string
    Profile Profile `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE;"`
}

// OWNED side (Profile) - has the foreign key
type Profile struct {
    gorm.Model
    UserID uint   `gorm:"uniqueIndex"` // Foreign key (unique = one-to-one)
    Bio    string
}
```

**Java equivalent:**
```java
@Entity class User {
    @OneToOne(mappedBy = "user", cascade = CascadeType.ALL)
    private Profile profile;
}

@Entity class Profile {
    @OneToOne @JoinColumn(name = "user_id", unique = true)
    private User user;
}
```

**Query with eager loading:**
```go
db.Preload("Profile").First(&user, 1)
```

---

## 1️⃣➡️🔢 ONE-TO-MANY: Author → Books

```go
// ONE side (Author) - has many Books
type Author struct {
    gorm.Model
    Name  string
    Books []Book `gorm:"foreignKey:AuthorID"` // Has Many
}

// MANY side (Book) - belongs to Author
type Book struct {
    gorm.Model
    Title    string
    AuthorID uint `gorm:"index"` // Foreign key
}
```

**Java equivalent:**
```java
@Entity class Author {
    @OneToMany(mappedBy = "author", cascade = CascadeType.ALL)
    private List<Book> books;
}

@Entity class Book {
    @ManyToOne @JoinColumn(name = "author_id")
    private Author author;
}
```

**Query with eager loading:**
```go
db.Preload("Books").First(&author, 1)
```

---

## 🔢↔️🔢 MANY-TO-MANY: Books ↔ Tags (Simple)

```go
// Both sides reference each other
type Book struct {
    gorm.Model
    Title string
    Tags  []Tag `gorm:"many2many:book_tags;"` // Join table: book_tags
}

type Tag struct {
    gorm.Model
    Name  string
    Books []Book `gorm:"many2many:book_tags;"` // Back reference
}
```

**Java equivalent:**
```java
@Entity class Book {
    @ManyToMany
    @JoinTable(name = "book_tags",
        joinColumns = @JoinColumn(name = "book_id"),
        inverseJoinColumns = @JoinColumn(name = "tag_id"))
    private Set<Tag> tags;
}
```

**Association operations:**
```go
// Add tag to book
db.Model(&book).Association("Tags").Append(&tag)

// Remove tag from book
db.Model(&book).Association("Tags").Delete(&tag)

// Replace all tags
db.Model(&book).Association("Tags").Replace([]Tag{tag1, tag2})

// Get all tags
db.Model(&book).Association("Tags").Find(&tags)
```

---

## 🔢↔️🔢 MANY-TO-MANY: Students ↔ Courses (With Extra Fields)

When your join table needs extra attributes (grade, enrolled_at), use explicit join entity:

```go
type Student struct {
    gorm.Model
    Name        string
    Enrollments []Enrollment `gorm:"foreignKey:StudentID"`
}

type Course struct {
    gorm.Model
    Name        string
    Enrollments []Enrollment `gorm:"foreignKey:CourseID"`
}

// Explicit join table with extra fields
type Enrollment struct {
    gorm.Model
    StudentID  uint      `gorm:"uniqueIndex:idx_student_course"`
    CourseID   uint      `gorm:"uniqueIndex:idx_student_course"`
    EnrolledAt time.Time
    Grade      *string   // Nullable
    Completed  bool
    
    Student Student `gorm:"foreignKey:StudentID"`
    Course  Course  `gorm:"foreignKey:CourseID"`
}
```

**Java equivalent:**
```java
@Entity class Enrollment {
    @ManyToOne private Student student;
    @ManyToOne private Course course;
    private LocalDate enrolledAt;
    private String grade;
}
```

**Query enrollments with related data:**
```go
db.Preload("Student").Preload("Course").Where("student_id = ?", 1).Find(&enrollments)
```

---

## 🔄 SELF-REFERENTIAL: Category Tree

```go
type Category struct {
    gorm.Model
    Name     string
    ParentID *uint      `gorm:"index"`                      // Nullable (root = no parent)
    Parent   *Category  `gorm:"foreignKey:ParentID"`        // Parent reference
    Children []Category `gorm:"foreignKey:ParentID"`        // Children references
}
```

**Java equivalent:**
```java
@Entity class Category {
    @ManyToOne @JoinColumn(name = "parent_id")
    private Category parent;
    
    @OneToMany(mappedBy = "parent")
    private List<Category> children;
}
```

**Query tree:**
```go
// Get category with children
db.Preload("Children").First(&category, 1)

// Get category with parent
db.Preload("Parent").First(&category, 5)
```

---

## 🎯 Quick Reference Table

| Relationship | GORM Tag | Foreign Key Location |
|--------------|----------|---------------------|
| One-to-One | `gorm:"foreignKey:UserID"` | On owned side (Profile) |
| One-to-Many | `gorm:"foreignKey:AuthorID"` | On "many" side (Book) |
| Many-to-Many | `gorm:"many2many:join_table"` | Join table |
| Self-ref | `gorm:"foreignKey:ParentID"` | Same table |

---

## 📝 Common Operations

```go
// CREATE with associations
user := User{
    Name: "Viraj",
    Profile: Profile{Bio: "Developer"},  // Creates both
}
db.Create(&user)

// EAGER LOAD (like @EntityGraph)
db.Preload("Profile").Preload("Posts").First(&user, 1)

// LAZY LOAD (separate query)
var posts []Post
db.Model(&user).Association("Posts").Find(&posts)

// UPDATE association
db.Model(&book).Association("Tags").Append(&newTag)

// DELETE association (not the record)
db.Model(&book).Association("Tags").Delete(&tag)

// CLEAR all associations
db.Model(&book).Association("Tags").Clear()

// COUNT associations
count := db.Model(&book).Association("Tags").Count()
```

---

## ⚠️ Common Gotchas

1. **Preload vs Joins**: `Preload` = separate query, `Joins` = SQL JOIN
2. **Circular references**: Use pointers (`*Category`) to avoid infinite loops
3. **N+1 problem**: Always use `Preload` when you need associations
4. **Soft delete**: `DeletedAt` affects all queries automatically
5. **Unique indexes**: Use for one-to-one foreign keys

---

*This file is part of Exercise 08 - Database Integration*

