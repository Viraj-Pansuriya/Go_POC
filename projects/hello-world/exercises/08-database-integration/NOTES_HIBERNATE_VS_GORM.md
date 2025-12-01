# Hibernate/JPA vs GORM - Complete Comparison

## 📚 Table of Contents
1. [Entity/Model Definition](#1-entitymodel-definition)
2. [Column Annotations](#2-column-annotations)
3. [Relationships](#3-relationships)
4. [Cascade Operations](#4-cascade-operations)
5. [Fetch Types (Lazy/Eager)](#5-fetch-types-lazyeager)
6. [Queries](#6-queries)
7. [Transactions](#7-transactions)
8. [Lifecycle Hooks](#8-lifecycle-hooks)
9. [Migrations](#9-migrations)
10. [Common Patterns](#10-common-patterns)

---

## 1. Entity/Model Definition

### Hibernate/JPA
```java
@Entity
@Table(name = "users")
public class User {
    @Id
    @GeneratedValue(strategy = GenerationType.IDENTITY)
    private Long id;
    
    @Column(name = "full_name", nullable = false, length = 100)
    private String name;
    
    @Column(unique = true)
    private String email;
    
    @CreationTimestamp
    private LocalDateTime createdAt;
    
    @UpdateTimestamp
    private LocalDateTime updatedAt;
}
```

### GORM
```go
type User struct {
    ID        uint           `gorm:"primaryKey"`
    Name      string         `gorm:"column:full_name;not null;size:100"`
    Email     string         `gorm:"unique"`
    CreatedAt time.Time      // Auto-managed by GORM
    UpdatedAt time.Time      // Auto-managed by GORM
    DeletedAt gorm.DeletedAt `gorm:"index"` // Soft delete support
}

// Or use gorm.Model for common fields:
type User struct {
    gorm.Model        // Includes ID, CreatedAt, UpdatedAt, DeletedAt
    Name  string
    Email string `gorm:"unique"`
}
```

---

## 2. Column Annotations

### Complete Mapping Table

| Hibernate/JPA | GORM | Description |
|---------------|------|-------------|
| `@Id` | `gorm:"primaryKey"` | Primary key |
| `@GeneratedValue(IDENTITY)` | `gorm:"autoIncrement"` | Auto increment |
| `@Column(name = "x")` | `gorm:"column:x"` | Column name |
| `@Column(nullable = false)` | `gorm:"not null"` | Not null constraint |
| `@Column(unique = true)` | `gorm:"unique"` | Unique constraint |
| `@Column(length = 100)` | `gorm:"size:100"` | String length |
| `@Column(precision, scale)` | `gorm:"precision:10;scale:2"` | Decimal precision |
| `@Column(columnDefinition)` | `gorm:"type:varchar(100)"` | Custom SQL type |
| `@Lob` | `gorm:"type:text"` | Large object |
| `@Transient` | `gorm:"-"` | Ignore field |
| `@Enumerated(STRING)` | Custom type or string | Enum handling |
| `@Column(insertable=false)` | `gorm:"->;migration"` | Read-only column |
| `@Column(updatable=false)` | `gorm:"<-:create"` | Insert only |

### GORM Tag Examples
```go
type Product struct {
    ID          uint    `gorm:"primaryKey;autoIncrement"`
    Code        string  `gorm:"column:product_code;size:50;uniqueIndex"`
    Name        string  `gorm:"not null;size:200"`
    Price       float64 `gorm:"type:decimal(10,2);default:0"`
    Description string  `gorm:"type:text"`
    Stock       int     `gorm:"default:0;check:stock >= 0"`
    IsActive    bool    `gorm:"default:true"`
    Metadata    string  `gorm:"-"` // Ignored by GORM
    
    // Composite unique index
    SKU         string  `gorm:"uniqueIndex:idx_sku_vendor"`
    VendorID    uint    `gorm:"uniqueIndex:idx_sku_vendor"`
    
    // Timestamps
    CreatedAt   time.Time
    UpdatedAt   time.Time
    DeletedAt   gorm.DeletedAt `gorm:"index"`
}
```

---

## 3. Relationships

### 3.1 One-to-One

#### Hibernate/JPA
```java
@Entity
public class User {
    @Id
    @GeneratedValue
    private Long id;
    
    @OneToOne(mappedBy = "user", cascade = CascadeType.ALL)
    private Profile profile;
}

@Entity
public class Profile {
    @Id
    @GeneratedValue
    private Long id;
    
    @OneToOne
    @JoinColumn(name = "user_id")
    private User user;
}
```

#### GORM
```go
type User struct {
    gorm.Model
    Name    string
    Profile Profile  // Has One
}

type Profile struct {
    gorm.Model
    UserID  uint    // Foreign key (convention: <Parent>ID)
    Bio     string
    Avatar  string
}

// Query with relation
db.Preload("Profile").Find(&user)
```

### 3.2 One-to-Many / Many-to-One

#### Hibernate/JPA
```java
@Entity
public class User {
    @Id
    @GeneratedValue
    private Long id;
    
    @OneToMany(mappedBy = "author", cascade = CascadeType.ALL)
    private List<Post> posts;
}

@Entity
public class Post {
    @Id
    @GeneratedValue
    private Long id;
    
    @ManyToOne(fetch = FetchType.LAZY)
    @JoinColumn(name = "author_id")
    private User author;
}
```

#### GORM
```go
type User struct {
    gorm.Model
    Name  string
    Posts []Post  // Has Many
}

type Post struct {
    gorm.Model
    Title    string
    Content  string
    AuthorID uint  // Foreign key
    Author   User  // Belongs To (optional, for eager loading)
}

// Query with relation
db.Preload("Posts").Find(&users)

// Custom foreign key
type Post struct {
    gorm.Model
    Title     string
    WriterID  uint `gorm:"column:writer_id"`
    Writer    User `gorm:"foreignKey:WriterID"`
}
```

### 3.3 Many-to-Many

#### Hibernate/JPA
```java
@Entity
public class Student {
    @Id
    @GeneratedValue
    private Long id;
    
    @ManyToMany
    @JoinTable(
        name = "student_courses",
        joinColumns = @JoinColumn(name = "student_id"),
        inverseJoinColumns = @JoinColumn(name = "course_id")
    )
    private Set<Course> courses;
}

@Entity
public class Course {
    @Id
    @GeneratedValue
    private Long id;
    
    @ManyToMany(mappedBy = "courses")
    private Set<Student> students;
}
```

#### GORM
```go
type Student struct {
    gorm.Model
    Name    string
    Courses []Course `gorm:"many2many:student_courses;"`
}

type Course struct {
    gorm.Model
    Title    string
    Students []Student `gorm:"many2many:student_courses;"`
}

// Adding association
db.Model(&student).Association("Courses").Append(&course)

// Query with relation
db.Preload("Courses").Find(&students)

// Custom join table
type Student struct {
    gorm.Model
    Name    string
    Courses []Course `gorm:"many2many:enrollments;joinForeignKey:StudentID;joinReferences:CourseID"`
}
```

### 3.4 Self-Referential (Tree/Hierarchy)

#### Hibernate/JPA
```java
@Entity
public class Category {
    @Id
    @GeneratedValue
    private Long id;
    
    @ManyToOne
    @JoinColumn(name = "parent_id")
    private Category parent;
    
    @OneToMany(mappedBy = "parent")
    private List<Category> children;
}
```

#### GORM
```go
type Category struct {
    gorm.Model
    Name     string
    ParentID *uint       // Nullable for root categories
    Parent   *Category   `gorm:"foreignKey:ParentID"`
    Children []Category  `gorm:"foreignKey:ParentID"`
}

// Query tree
db.Preload("Children").Where("parent_id IS NULL").Find(&rootCategories)
```

---

## 4. Cascade Operations

### Hibernate/JPA
```java
@OneToMany(cascade = CascadeType.ALL, orphanRemoval = true)
private List<Post> posts;

// Cascade types:
// CascadeType.PERSIST - save children when parent saved
// CascadeType.MERGE - update children when parent updated
// CascadeType.REMOVE - delete children when parent deleted
// CascadeType.ALL - all of above
// orphanRemoval - delete child when removed from collection
```

### GORM
```go
// GORM doesn't have built-in cascade like JPA
// You handle it manually or use hooks

// Option 1: Use hooks
func (u *User) BeforeDelete(tx *gorm.DB) error {
    return tx.Where("user_id = ?", u.ID).Delete(&Post{}).Error
}

// Option 2: Use Association mode
db.Select("Posts").Delete(&user)  // Delete user and posts

// Option 3: Database-level cascade (in migration)
db.Exec(`ALTER TABLE posts ADD CONSTRAINT fk_user 
         FOREIGN KEY (user_id) REFERENCES users(id) 
         ON DELETE CASCADE`)

// Option 4: Manual cascade
func DeleteUserWithPosts(db *gorm.DB, userID uint) error {
    return db.Transaction(func(tx *gorm.DB) error {
        if err := tx.Where("user_id = ?", userID).Delete(&Post{}).Error; err != nil {
            return err
        }
        return tx.Delete(&User{}, userID).Error
    })
}
```

---

## 5. Fetch Types (Lazy/Eager)

### Hibernate/JPA
```java
// Lazy (default for collections) - loads on access
@OneToMany(fetch = FetchType.LAZY)
private List<Post> posts;

// Eager - loads immediately
@ManyToOne(fetch = FetchType.EAGER)
private User author;

// Force loading with EntityGraph
@EntityGraph(attributePaths = {"posts", "posts.comments"})
Optional<User> findById(Long id);
```

### GORM
```go
// GORM is LAZY by default - relations are NOT loaded

// Eager loading with Preload
db.Preload("Posts").Find(&users)

// Nested preload
db.Preload("Posts.Comments").Find(&users)

// Conditional preload
db.Preload("Posts", "published = ?", true).Find(&users)

// Preload with custom query
db.Preload("Posts", func(db *gorm.DB) *gorm.DB {
    return db.Order("created_at DESC").Limit(5)
}).Find(&users)

// Multiple preloads
db.Preload("Posts").Preload("Profile").Find(&users)

// Preload All (use carefully!)
db.Preload(clause.Associations).Find(&users)

// Joins (like JPA JOIN FETCH) - single query
db.Joins("Profile").Find(&users)
db.Joins("LEFT JOIN profiles ON profiles.user_id = users.id").Find(&users)
```

---

## 6. Queries

### 6.1 Basic Queries

#### Hibernate/JPA
```java
// Find by ID
User user = repository.findById(1L).orElseThrow();

// Find all
List<User> users = repository.findAll();

// Find with condition
List<User> users = repository.findByName("John");

// JPQL
@Query("SELECT u FROM User u WHERE u.email LIKE %:domain")
List<User> findByEmailDomain(@Param("domain") String domain);

// Native query
@Query(value = "SELECT * FROM users WHERE age > ?1", nativeQuery = true)
List<User> findOlderThan(int age);
```

#### GORM
```go
// Find by ID
var user User
db.First(&user, 1)
db.First(&user, "id = ?", 1)

// Find all
var users []User
db.Find(&users)

// Find with condition
db.Where("name = ?", "John").Find(&users)

// Like query
db.Where("email LIKE ?", "%@gmail.com").Find(&users)

// Multiple conditions
db.Where("name = ? AND age > ?", "John", 18).Find(&users)

// Struct condition
db.Where(&User{Name: "John", Age: 20}).Find(&users)

// Map condition
db.Where(map[string]interface{}{"name": "John", "age": 20}).Find(&users)

// Raw SQL
db.Raw("SELECT * FROM users WHERE age > ?", 18).Scan(&users)

// First or Not Found
result := db.First(&user, 1)
if errors.Is(result.Error, gorm.ErrRecordNotFound) {
    // Handle not found
}
```

### 6.2 Advanced Queries

#### Hibernate/JPA
```java
// Specification pattern
Specification<User> spec = (root, query, cb) -> 
    cb.and(
        cb.equal(root.get("status"), "ACTIVE"),
        cb.greaterThan(root.get("age"), 18)
    );
repository.findAll(spec);

// Projections
interface UserSummary {
    String getName();
    String getEmail();
}
List<UserSummary> findAllProjectedBy();

// Pagination
Page<User> findAll(Pageable pageable);
```

#### GORM
```go
// Complex conditions
db.Where("status = ?", "ACTIVE").
   Where("age > ?", 18).
   Or("role = ?", "admin").
   Find(&users)

// Subquery
subQuery := db.Model(&Order{}).Select("user_id").Where("amount > ?", 100)
db.Where("id IN (?)", subQuery).Find(&users)

// Select specific fields (projection)
type UserSummary struct {
    Name  string
    Email string
}
var summaries []UserSummary
db.Model(&User{}).Select("name", "email").Scan(&summaries)

// Pagination
var users []User
var total int64

db.Model(&User{}).Count(&total)
db.Offset(0).Limit(10).Find(&users)

// Pagination helper
func Paginate(page, pageSize int) func(db *gorm.DB) *gorm.DB {
    return func(db *gorm.DB) *gorm.DB {
        offset := (page - 1) * pageSize
        return db.Offset(offset).Limit(pageSize)
    }
}
db.Scopes(Paginate(1, 10)).Find(&users)

// Group By
type Result struct {
    Status string
    Count  int
}
var results []Result
db.Model(&User{}).Select("status, count(*) as count").Group("status").Scan(&results)

// Having
db.Model(&User{}).Select("status, count(*) as count").
   Group("status").
   Having("count > ?", 5).
   Scan(&results)

// Distinct
db.Distinct("name", "email").Find(&users)

// Locking
db.Clauses(clause.Locking{Strength: "UPDATE"}).Find(&users)
```

---

## 7. Transactions

### Hibernate/JPA
```java
@Service
public class UserService {
    
    @Transactional
    public void transferMoney(Long fromId, Long toId, BigDecimal amount) {
        Account from = accountRepository.findById(fromId).orElseThrow();
        Account to = accountRepository.findById(toId).orElseThrow();
        
        from.debit(amount);
        to.credit(amount);
        
        accountRepository.save(from);
        accountRepository.save(to);
    }
    
    @Transactional(readOnly = true)
    public User getUser(Long id) {
        return userRepository.findById(id).orElseThrow();
    }
    
    @Transactional(propagation = Propagation.REQUIRES_NEW)
    public void logAction(String action) {
        // New transaction
    }
}
```

### GORM
```go
// Basic transaction
err := db.Transaction(func(tx *gorm.DB) error {
    var from, to Account
    
    if err := tx.First(&from, fromID).Error; err != nil {
        return err  // Rollback
    }
    if err := tx.First(&to, toID).Error; err != nil {
        return err  // Rollback
    }
    
    from.Balance -= amount
    to.Balance += amount
    
    if err := tx.Save(&from).Error; err != nil {
        return err  // Rollback
    }
    if err := tx.Save(&to).Error; err != nil {
        return err  // Rollback
    }
    
    return nil  // Commit
})

// Manual transaction control
tx := db.Begin()
defer func() {
    if r := recover(); r != nil {
        tx.Rollback()
    }
}()

if err := tx.Create(&user).Error; err != nil {
    tx.Rollback()
    return err
}

if err := tx.Create(&profile).Error; err != nil {
    tx.Rollback()
    return err
}

tx.Commit()

// Nested transactions (savepoints)
db.Transaction(func(tx *gorm.DB) error {
    tx.Create(&user1)
    
    tx.Transaction(func(tx2 *gorm.DB) error {
        tx2.Create(&user2)
        return errors.New("rollback user2 only")  // Only user2 rolled back
    })
    
    return nil  // user1 committed
})
```

---

## 8. Lifecycle Hooks

### Hibernate/JPA
```java
@Entity
public class User {
    @PrePersist
    public void prePersist() {
        this.createdAt = LocalDateTime.now();
    }
    
    @PreUpdate
    public void preUpdate() {
        this.updatedAt = LocalDateTime.now();
    }
    
    @PostLoad
    public void postLoad() {
        // After entity loaded
    }
    
    @PreRemove
    public void preRemove() {
        // Before delete
    }
}
```

### GORM
```go
type User struct {
    gorm.Model
    Name string
    UUID string
}

// Before create
func (u *User) BeforeCreate(tx *gorm.DB) error {
    u.UUID = uuid.New().String()
    return nil
}

// After create
func (u *User) AfterCreate(tx *gorm.DB) error {
    // Send welcome email, etc.
    return nil
}

// Before update
func (u *User) BeforeUpdate(tx *gorm.DB) error {
    // Validation, etc.
    return nil
}

// After update
func (u *User) AfterUpdate(tx *gorm.DB) error {
    // Audit log, etc.
    return nil
}

// Before delete
func (u *User) BeforeDelete(tx *gorm.DB) error {
    // Check constraints, cleanup
    return nil
}

// After delete
func (u *User) AfterDelete(tx *gorm.DB) error {
    // Cleanup related data
    return nil
}

// After find (like @PostLoad)
func (u *User) AfterFind(tx *gorm.DB) error {
    // Decrypt sensitive fields, etc.
    return nil
}

// Available hooks:
// BeforeSave, AfterSave (covers both create and update)
// BeforeCreate, AfterCreate
// BeforeUpdate, AfterUpdate
// BeforeDelete, AfterDelete
// AfterFind
```

---

## 9. Migrations

### Hibernate/JPA (with Flyway)
```sql
-- V1__Create_users_table.sql
CREATE TABLE users (
    id BIGINT AUTO_INCREMENT PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    email VARCHAR(100) UNIQUE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- V2__Add_age_column.sql
ALTER TABLE users ADD COLUMN age INT DEFAULT 0;
```

### GORM
```go
// Auto-migrate (development)
db.AutoMigrate(&User{}, &Post{}, &Comment{})

// Check migration needed
if !db.Migrator().HasTable(&User{}) {
    db.Migrator().CreateTable(&User{})
}

// Add column
if !db.Migrator().HasColumn(&User{}, "Age") {
    db.Migrator().AddColumn(&User{}, "Age")
}

// Modify column
db.Migrator().AlterColumn(&User{}, "Name")

// Drop column
db.Migrator().DropColumn(&User{}, "TempField")

// Create index
db.Migrator().CreateIndex(&User{}, "Email")
db.Migrator().CreateIndex(&User{}, "idx_user_email")

// Drop table
db.Migrator().DropTable(&User{})

// For production, use migration tools:
// - golang-migrate: https://github.com/golang-migrate/migrate
// - goose: https://github.com/pressly/goose
// - atlas: https://atlasgo.io/
```

---

## 10. Common Patterns

### 10.1 Repository Pattern

```go
// Generic repository
type Repository[T any] interface {
    Create(entity *T) error
    FindByID(id uint) (*T, error)
    FindAll() ([]T, error)
    Update(entity *T) error
    Delete(id uint) error
}

type GormRepository[T any] struct {
    db *gorm.DB
}

func NewRepository[T any](db *gorm.DB) Repository[T] {
    return &GormRepository[T]{db: db}
}

func (r *GormRepository[T]) Create(entity *T) error {
    return r.db.Create(entity).Error
}

func (r *GormRepository[T]) FindByID(id uint) (*T, error) {
    var entity T
    err := r.db.First(&entity, id).Error
    return &entity, err
}

// Usage
userRepo := NewRepository[User](db)
postRepo := NewRepository[Post](db)
```

### 10.2 Specification Pattern (like JPA Specifications)

```go
type Specification func(*gorm.DB) *gorm.DB

func WithStatus(status string) Specification {
    return func(db *gorm.DB) *gorm.DB {
        return db.Where("status = ?", status)
    }
}

func WithAgeGreaterThan(age int) Specification {
    return func(db *gorm.DB) *gorm.DB {
        return db.Where("age > ?", age)
    }
}

func WithNameLike(name string) Specification {
    return func(db *gorm.DB) *gorm.DB {
        return db.Where("name LIKE ?", "%"+name+"%")
    }
}

// Usage
func (r *UserRepository) FindAll(specs ...Specification) ([]User, error) {
    var users []User
    query := r.db
    for _, spec := range specs {
        query = spec(query)
    }
    err := query.Find(&users).Error
    return users, err
}

users, _ := repo.FindAll(
    WithStatus("active"),
    WithAgeGreaterThan(18),
)
```

### 10.3 Soft Delete

```go
// GORM automatically handles soft delete if DeletedAt field exists
type User struct {
    gorm.Model  // Includes DeletedAt
    Name string
}

// Delete (soft)
db.Delete(&user, 1)  // Sets DeletedAt, doesn't remove row

// Query excludes soft deleted by default
db.Find(&users)  // Only non-deleted users

// Include soft deleted
db.Unscoped().Find(&users)

// Permanently delete
db.Unscoped().Delete(&user, 1)

// Restore soft deleted
db.Model(&User{}).Unscoped().Where("id = ?", 1).Update("deleted_at", nil)
```

### 10.4 Audit Trail

```go
type Auditable struct {
    CreatedAt time.Time
    UpdatedAt time.Time
    CreatedBy uint
    UpdatedBy uint
}

type User struct {
    gorm.Model
    Auditable
    Name string
}

func (a *Auditable) BeforeCreate(tx *gorm.DB) error {
    userID := getCurrentUserID()  // From context
    a.CreatedBy = userID
    a.UpdatedBy = userID
    return nil
}

func (a *Auditable) BeforeUpdate(tx *gorm.DB) error {
    a.UpdatedBy = getCurrentUserID()
    return nil
}
```

---

## Quick Reference Card

```
┌─────────────────────────────────────────────────────────────┐
│                    GORM QUICK REFERENCE                     │
├─────────────────────────────────────────────────────────────┤
│ CRUD                                                        │
│   db.Create(&user)              // INSERT                   │
│   db.First(&user, 1)            // SELECT ... WHERE id=1    │
│   db.Find(&users)               // SELECT *                 │
│   db.Save(&user)                // UPDATE (all fields)      │
│   db.Delete(&user, 1)           // Soft DELETE              │
├─────────────────────────────────────────────────────────────┤
│ WHERE                                                       │
│   db.Where("name = ?", "john")                              │
│   db.Where(&User{Name: "john"})                             │
│   db.Where(map[string]any{"name": "john"})                  │
├─────────────────────────────────────────────────────────────┤
│ RELATIONS                                                   │
│   db.Preload("Posts").Find(&users)      // Eager load       │
│   db.Joins("Profile").Find(&users)      // JOIN             │
│   db.Association("Posts").Append(&post) // Add relation     │
├─────────────────────────────────────────────────────────────┤
│ TRANSACTION                                                 │
│   db.Transaction(func(tx *gorm.DB) error {                  │
│       tx.Create(&user)                                      │
│       return nil  // commit, or error to rollback           │
│   })                                                        │
├─────────────────────────────────────────────────────────────┤
│ COMMON TAGS                                                 │
│   gorm:"primaryKey"       gorm:"not null"                   │
│   gorm:"unique"           gorm:"size:100"                   │
│   gorm:"default:0"        gorm:"index"                      │
│   gorm:"-"                gorm:"foreignKey:UserID"          │
└─────────────────────────────────────────────────────────────┘
```


