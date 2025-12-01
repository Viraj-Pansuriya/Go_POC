# Exercise 07: Dependency Injection in Go

## 🎯 Goal
Learn **proper dependency injection** in Go - no magic, pure constructor injection!

---

## 📚 Spring Boot vs Go DI

### Spring Boot (Magic DI)
```java
@Service
public class UserService {
    @Autowired  // Magic! Spring finds and injects this
    private UserRepository userRepo;
    
    @Autowired
    private EmailService emailService;
}

// Usage - Spring creates everything
@RestController
public class UserController {
    @Autowired
    private UserService userService;  // Magic injection!
}
```

### Go (Explicit Constructor DI)
```go
type UserService struct {
    repo         UserRepository  // Interface, not implementation!
    emailService EmailService
}

// Constructor - YOU wire dependencies
func NewUserService(repo UserRepository, email EmailService) *UserService {
    return &UserService{
        repo:         repo,
        emailService: email,
    }
}

// Usage - YOU create and wire everything in main()
func main() {
    repo := NewPostgresUserRepository(db)
    emailSvc := NewSMTPEmailService(config)
    userSvc := NewUserService(repo, emailSvc)  // Explicit wiring!
    
    handler := NewUserHandler(userSvc)
    // ...
}
```

---

## 🤔 Why Go Chose Explicit DI?

| Spring (Magic) | Go (Explicit) |
|----------------|---------------|
| `@Autowired` finds beans | You write `New*()` functions |
| Hidden dependency graph | Dependencies visible in `main()` |
| Runtime errors if bean missing | Compile-time errors |
| Needs annotations/reflection | Just regular Go code |
| Complex to debug | Easy to trace |

**Go Philosophy**: "Clear is better than clever"

---

## 🏋️ Your Task: Refactor Exercise 06 with Proper DI

### Current Problem in Exercise 06:
```go
// controller/Webcontroller.go
var (
    repo = Repository.GetUserRepositoryInstance()  // ❌ Global variable!
)

func GetUser(c *gin.Context) {
    resp, _ := repo.FindById(id)  // ❌ Hard to test, tightly coupled
}
```

### Target Architecture:
```
main.go
    │
    ├── Creates: Repository (interface)
    ├── Creates: Service (depends on Repository)
    ├── Creates: Handler (depends on Service)
    │
    └── Wires: Handler → Gin routes
```

---

## 📁 Files to Create

```
07-dependency-injection/
├── main.go              # Wiring happens here!
├── models/
│   └── user.go
├── repository/
│   ├── user_repository.go      # Interface
│   └── user_repository_impl.go # Implementation
├── service/
│   ├── user_service.go         # Interface
│   └── user_service_impl.go    # Implementation
└── handler/
    └── user_handler.go         # HTTP handlers
```

---

## 💻 Implementation Guide

### 1. Repository Layer (Data Access)

```go
// repository/user_repository.go
package repository

type UserRepository interface {
    FindByID(id string) (*models.User, error)
    FindAll() ([]*models.User, error)
    Create(user *models.User) error
    Delete(id string) error
}
```

```go
// repository/user_repository_impl.go
package repository

type InMemoryUserRepository struct {
    users map[string]*models.User
    mu    sync.RWMutex  // Thread-safe!
}

func NewInMemoryUserRepository() *InMemoryUserRepository {
    return &InMemoryUserRepository{
        users: make(map[string]*models.User),
    }
}

func (r *InMemoryUserRepository) FindByID(id string) (*models.User, error) {
    r.mu.RLock()
    defer r.mu.RUnlock()
    
    user, ok := r.users[id]
    if !ok {
        return nil, ErrUserNotFound
    }
    return user, nil
}
// ... implement other methods
```

### 2. Service Layer (Business Logic)

```go
// service/user_service.go
package service

type UserService interface {
    GetUser(id string) (*models.User, error)
    GetAllUsers() ([]*models.User, error)
    CreateUser(name, email string) (*models.User, error)
    DeleteUser(id string) error
}
```

```go
// service/user_service_impl.go
package service

type userServiceImpl struct {
    repo repository.UserRepository  // Depends on interface!
}

// Constructor with dependency injection
func NewUserService(repo repository.UserRepository) UserService {
    return &userServiceImpl{repo: repo}
}

func (s *userServiceImpl) GetUser(id string) (*models.User, error) {
    return s.repo.FindByID(id)
}

func (s *userServiceImpl) CreateUser(name, email string) (*models.User, error) {
    user := &models.User{
        ID:    uuid.New().String(),
        Name:  name,
        Email: email,
    }
    
    if err := s.repo.Create(user); err != nil {
        return nil, err
    }
    return user, nil
}
```

### 3. Handler Layer (HTTP)

```go
// handler/user_handler.go
package handler

type UserHandler struct {
    service service.UserService  // Depends on interface!
}

// Constructor with dependency injection
func NewUserHandler(svc service.UserService) *UserHandler {
    return &UserHandler{service: svc}
}

func (h *UserHandler) GetUser(c *gin.Context) {
    id := c.Param("id")
    
    user, err := h.service.GetUser(id)
    if err != nil {
        c.JSON(404, gin.H{"error": err.Error()})
        return
    }
    
    c.JSON(200, user)
}

func (h *UserHandler) CreateUser(c *gin.Context) {
    var req CreateUserRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(400, gin.H{"error": err.Error()})
        return
    }
    
    user, err := h.service.CreateUser(req.Name, req.Email)
    if err != nil {
        c.JSON(500, gin.H{"error": err.Error()})
        return
    }
    
    c.JSON(201, user)
}

// RegisterRoutes adds routes to Gin engine
func (h *UserHandler) RegisterRoutes(r *gin.Engine) {
    api := r.Group("/api")
    {
        api.GET("/users/:id", h.GetUser)
        api.GET("/users", h.GetAllUsers)
        api.POST("/users", h.CreateUser)
        api.DELETE("/users/:id", h.DeleteUser)
    }
}
```

### 4. Main - The Wiring Point!

```go
// main.go
package main

func main() {
    // ===== Create Dependencies (bottom-up) =====
    
    // 1. Repository (no dependencies)
    userRepo := repository.NewInMemoryUserRepository()
    
    // 2. Service (depends on repository)
    userService := service.NewUserService(userRepo)
    
    // 3. Handler (depends on service)
    userHandler := handler.NewUserHandler(userService)
    
    // ===== Setup Gin =====
    r := gin.Default()
    
    // Health check
    r.GET("/health", func(c *gin.Context) {
        c.JSON(200, gin.H{"status": "ok"})
    })
    
    // Register user routes
    userHandler.RegisterRoutes(r)
    
    // ===== Start Server =====
    log.Println("Server starting on :8080")
    r.Run(":8080")
}
```

---

## 🧪 Bonus: Easy Testing with DI!

```go
// service/user_service_test.go

// Mock repository
type MockUserRepository struct {
    users map[string]*models.User
}

func (m *MockUserRepository) FindByID(id string) (*models.User, error) {
    user, ok := m.users[id]
    if !ok {
        return nil, errors.New("not found")
    }
    return user, nil
}

func TestGetUser(t *testing.T) {
    // Create mock with test data
    mockRepo := &MockUserRepository{
        users: map[string]*models.User{
            "1": {ID: "1", Name: "Test User", Email: "test@test.com"},
        },
    }
    
    // Inject mock into service
    svc := NewUserService(mockRepo)
    
    // Test!
    user, err := svc.GetUser("1")
    if err != nil {
        t.Fatal(err)
    }
    if user.Name != "Test User" {
        t.Errorf("expected 'Test User', got '%s'", user.Name)
    }
}
```

---

## ✅ Expected Outcome

1. **No global variables** for dependencies
2. **Interfaces everywhere** - easy to mock/swap
3. **Constructor injection** - dependencies explicit
4. **main.go is the composition root** - all wiring visible
5. **Easy to test** - just inject mocks!

---

## 🎓 What You'll Learn

1. **Constructor injection** - The Go way
2. **Interfaces for dependencies** - Loose coupling
3. **Composition root pattern** - Wire in main()
4. **Testability** - Easy mocking
5. **Clean architecture** - Handler → Service → Repository

---

## 💡 Key Patterns

### Pattern 1: Accept Interfaces, Return Structs
```go
// Constructor accepts interface (flexible)
func NewUserService(repo UserRepository) UserService {
    return &userServiceImpl{repo: repo}  // Returns concrete type wrapped in interface
}
```

### Pattern 2: Private Struct, Public Interface
```go
// Public interface
type UserService interface {
    GetUser(id string) (*User, error)
}

// Private implementation (lowercase = unexported)
type userServiceImpl struct {
    repo UserRepository
}
```

### Pattern 3: Registration Methods
```go
// Handler registers its own routes
func (h *UserHandler) RegisterRoutes(r *gin.Engine) {
    r.GET("/users/:id", h.GetUser)
}

// Clean main.go
userHandler.RegisterRoutes(r)
```

---

## ⏱️ Estimated Time: 25-30 minutes

This is how production Go apps are structured! 🏗️


