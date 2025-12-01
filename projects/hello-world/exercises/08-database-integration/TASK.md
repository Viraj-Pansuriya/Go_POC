# Exercise 08: Database Integration with GORM

## 🎯 Goal
Connect to a database using **GORM** - Go's most popular ORM (like Hibernate for Java)!

---

## 📚 Spring/Hibernate vs GORM

### Spring + JPA/Hibernate
```java
@Entity
@Table(name = "users")
public class User {
    @Id
    @GeneratedValue(strategy = GenerationType.IDENTITY)
    private Long id;
    
    @Column(nullable = false)
    private String name;
    
    @Column(unique = true)
    private String email;
}

@Repository
public interface UserRepository extends JpaRepository<User, Long> {
    Optional<User> findByEmail(String email);
}
```

### Go + GORM
```go
type User struct {
    ID    uint   `gorm:"primaryKey"`
    Name  string `gorm:"not null"`
    Email string `gorm:"unique"`
}

// No interface needed! GORM provides methods directly
db.Create(&user)
db.First(&user, id)
db.Where("email = ?", email).First(&user)
```

---

## 🚀 Setup

### 1. Create module
```bash
cd exercises/08-database-integration
go mod init github.com/viraj/go-mono-repo/projects/hello-world/exercises/08-database-integration
```

### 2. Install GORM + SQLite driver
```bash
go get -u gorm.io/gorm
go get -u gorm.io/driver/sqlite
```

SQLite = no database server needed! File-based DB, perfect for learning.

---

## 🏋️ Your Task: User CRUD with Real Database

### 1. Create User model with GORM tags

```go
// models/user.go
package models

import "gorm.io/gorm"

type User struct {
    gorm.Model         // Adds ID, CreatedAt, UpdatedAt, DeletedAt
    Name  string `gorm:"size:100;not null"`
    Email string `gorm:"size:100;uniqueIndex;not null"`
    Age   int    `gorm:"default:0"`
}
```

`gorm.Model` gives you:
```go
type Model struct {
    ID        uint           `gorm:"primaryKey"`
    CreatedAt time.Time
    UpdatedAt time.Time
    DeletedAt gorm.DeletedAt `gorm:"index"`  // Soft delete!
}
```

### 2. Create Database connection

```go
// database/database.go
package database

import (
    "gorm.io/driver/sqlite"
    "gorm.io/gorm"
)

var DB *gorm.DB

func Connect() error {
    var err error
    DB, err = gorm.Open(sqlite.Open("app.db"), &gorm.Config{})
    if err != nil {
        return err
    }
    
    // Auto-migrate (creates/updates tables)
    DB.AutoMigrate(&models.User{})
    
    return nil
}
```

### 3. Create Repository using GORM

```go
// repository/user_repository.go
package repository

type UserRepository interface {
    Create(user *models.User) error
    FindByID(id uint) (*models.User, error)
    FindByEmail(email string) (*models.User, error)
    FindAll() ([]models.User, error)
    Update(user *models.User) error
    Delete(id uint) error
}

type userRepository struct {
    db *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepository {
    return &userRepository{db: db}
}

func (r *userRepository) Create(user *models.User) error {
    return r.db.Create(user).Error
}

func (r *userRepository) FindByID(id uint) (*models.User, error) {
    var user models.User
    err := r.db.First(&user, id).Error
    if err != nil {
        return nil, err
    }
    return &user, nil
}

func (r *userRepository) FindByEmail(email string) (*models.User, error) {
    var user models.User
    err := r.db.Where("email = ?", email).First(&user).Error
    return &user, err
}

func (r *userRepository) FindAll() ([]models.User, error) {
    var users []models.User
    err := r.db.Find(&users).Error
    return users, err
}

func (r *userRepository) Update(user *models.User) error {
    return r.db.Save(user).Error
}

func (r *userRepository) Delete(id uint) error {
    return r.db.Delete(&models.User{}, id).Error  // Soft delete!
}
```

### 4. Wire everything in main.go

```go
// main.go
package main

func main() {
    // Connect to database
    if err := database.Connect(); err != nil {
        log.Fatal("Failed to connect to database:", err)
    }
    
    // Create dependencies
    userRepo := repository.NewUserRepository(database.DB)
    userService := service.NewUserService(userRepo)
    userHandler := handler.NewUserHandler(userService)
    
    // Setup Gin
    r := gin.Default()
    userHandler.RegisterRoutes(r)
    
    r.Run(":8080")
}
```

---

## 📁 Files to Create

```
08-database-integration/
├── main.go
├── database/
│   └── database.go
├── models/
│   └── user.go
├── repository/
│   └── user_repository.go
├── service/
│   └── user_service.go
└── handler/
    └── user_handler.go
```

---

## 💡 GORM Cheat Sheet

### Basic CRUD
```go
// Create
db.Create(&user)

// Read
db.First(&user, 1)                    // Find by primary key
db.First(&user, "id = ?", 1)          // Find by condition
db.Find(&users)                        // Find all
db.Where("name = ?", "John").First(&user)

// Update
db.Save(&user)                         // Update all fields
db.Model(&user).Update("name", "New")  // Update single field
db.Model(&user).Updates(User{Name: "New", Age: 20})  // Update multiple

// Delete
db.Delete(&user, 1)                    // Soft delete (if DeletedAt exists)
db.Unscoped().Delete(&user, 1)         // Hard delete
```

### Queries
```go
// Where conditions
db.Where("name = ?", "John").Find(&users)
db.Where("name LIKE ?", "%jo%").Find(&users)
db.Where("age > ?", 18).Find(&users)
db.Where("name IN ?", []string{"John", "Jane"}).Find(&users)

// Order, Limit, Offset
db.Order("created_at desc").Find(&users)
db.Limit(10).Offset(0).Find(&users)

// Count
var count int64
db.Model(&User{}).Count(&count)

// Select specific fields
db.Select("name", "email").Find(&users)
```

### Associations (like JPA relations)
```go
type User struct {
    gorm.Model
    Name    string
    Posts   []Post    `gorm:"foreignKey:UserID"`  // Has many
    Profile Profile   `gorm:"foreignKey:UserID"`  // Has one
}

type Post struct {
    gorm.Model
    Title  string
    UserID uint
}

// Eager loading (like @EntityGraph)
db.Preload("Posts").Find(&users)
```

---

## 🆚 JPA/Hibernate vs GORM

| Feature | JPA/Hibernate | GORM |
|---------|--------------|------|
| Entity annotations | `@Entity`, `@Column` | Struct tags `gorm:"..."` |
| Repository | Interface + Spring | Manual implementation |
| Auto ID | `@GeneratedValue` | `gorm.Model` or `primaryKey` |
| Soft delete | Manual | Built-in with `DeletedAt` |
| Migrations | Flyway/Liquibase | `AutoMigrate()` |
| Transactions | `@Transactional` | `db.Transaction()` |
| Query builder | JPQL/Criteria | Method chaining |
| Lazy loading | Default | Explicit with `Preload` |

---

## ✅ Expected Behavior

```bash
# Create user
curl -X POST http://localhost:8080/api/users \
  -H "Content-Type: application/json" \
  -d '{"name":"John","email":"john@test.com","age":25}'
# {"ID":1,"CreatedAt":"...","Name":"John","Email":"john@test.com","Age":25}

# Get user
curl http://localhost:8080/api/users/1
# {"ID":1,"CreatedAt":"...","Name":"John","Email":"john@test.com","Age":25}

# Get all users
curl http://localhost:8080/api/users
# [{"ID":1,...},{"ID":2,...}]

# Update user
curl -X PUT http://localhost:8080/api/users/1 \
  -d '{"name":"John Updated","email":"john@test.com","age":26}'

# Delete user (soft delete)
curl -X DELETE http://localhost:8080/api/users/1
# User still in DB but with DeletedAt set
```

---

## 🎓 What You'll Learn

1. **GORM basics** - Go's ORM
2. **Struct tags** - `gorm:"..."` annotations
3. **AutoMigrate** - Schema management
4. **Soft deletes** - Built-in with `DeletedAt`
5. **Query building** - Method chaining
6. **Repository pattern** - With real database

---

## ⏱️ Estimated Time: 25-30 minutes

This is real database work! 🗄️


