# Exercise 06: HTTP Server with Gin Framework

## 🎯 Goal
Build REST APIs with **Gin** - Go's most popular web framework!

---

## 📚 Spring Boot vs Gin

### Spring Boot
```java
@RestController
@RequestMapping("/api")
public class UserController {
    
    @GetMapping("/users/{id}")
    public User getUser(@PathVariable Long id) {
        return userService.findById(id);
    }
    
    @PostMapping("/users")
    public User createUser(@RequestBody User user) {
        return userService.save(user);
    }
}
```

### Gin (Very Similar!)
```go
func main() {
    r := gin.Default()
    
    r.GET("/api/users/:id", getUser)
    r.POST("/api/users", createUser)
    
    r.Run(":8080")
}

func getUser(c *gin.Context) {
    id := c.Param("id")  // Like @PathVariable
    c.JSON(200, user)    // Auto JSON!
}

func createUser(c *gin.Context) {
    var user User
    c.BindJSON(&user)    // Like @RequestBody
    c.JSON(201, user)
}
```

---

## 🚀 Setup

### 1. Initialize module (if not done)
```bash
cd exercises/06-http-server
go mod init github.com/viraj/go-mono-repo/projects/hello-world/exercises/06-http-server
```

### 2. Install Gin
```bash
go get -u github.com/gin-gonic/gin
```

---

## 🏋️ Your Task: Build a User API with Gin

### 1. Create endpoints:

| Method | Path | Handler | Description |
|--------|------|---------|-------------|
| GET | `/health` | `healthCheck` | Return status |
| GET | `/api/users` | `getUsers` | List all users |
| GET | `/api/users/:id` | `getUserByID` | Get one user |
| POST | `/api/users` | `createUser` | Create user |
| PUT | `/api/users/:id` | `updateUser` | Update user |
| DELETE | `/api/users/:id` | `deleteUser` | Delete user |

### 2. User struct
```go
type User struct {
    ID    int    `json:"id"`
    Name  string `json:"name" binding:"required"`
    Email string `json:"email" binding:"required,email"`
}
```

### 3. In-memory storage
```go
var users = []User{
    {ID: 1, Name: "John", Email: "john@example.com"},
    {ID: 2, Name: "Jane", Email: "jane@example.com"},
}
var nextID = 3
```

---

## 📁 Files to Create

```
06-http-server/
├── go.mod
├── main.go
├── handlers/
│   └── user_handler.go
└── models/
    └── user.go
```

---

## 💡 Gin Cheat Sheet

### Basic Server
```go
package main

import "github.com/gin-gonic/gin"

func main() {
    r := gin.Default()  // Includes logger & recovery middleware
    
    r.GET("/ping", func(c *gin.Context) {
        c.JSON(200, gin.H{"message": "pong"})
    })
    
    r.Run(":8080")  // Listen on port 8080
}
```

### Path Parameters (like @PathVariable)
```go
r.GET("/users/:id", func(c *gin.Context) {
    id := c.Param("id")  // Get path parameter
    c.JSON(200, gin.H{"id": id})
})
```

### Query Parameters (like @RequestParam)
```go
// GET /search?name=john&age=25
r.GET("/search", func(c *gin.Context) {
    name := c.Query("name")           // "john"
    age := c.DefaultQuery("age", "0") // "25" or default "0"
})
```

### Request Body (like @RequestBody)
```go
r.POST("/users", func(c *gin.Context) {
    var user User
    
    if err := c.ShouldBindJSON(&user); err != nil {
        c.JSON(400, gin.H{"error": err.Error()})
        return
    }
    
    c.JSON(201, user)
})
```

### Response Methods
```go
c.JSON(200, user)                    // JSON response
c.String(200, "Hello %s", name)      // Plain text
c.HTML(200, "index.html", data)      // HTML template
c.Redirect(302, "/new-url")          // Redirect
c.AbortWithStatus(401)               // Stop & return status
```

### Route Groups (like @RequestMapping on class)
```go
api := r.Group("/api")
{
    api.GET("/users", getUsers)
    api.POST("/users", createUser)
    
    // Nested group
    v1 := api.Group("/v1")
    {
        v1.GET("/users", getUsersV1)
    }
}
```

### Middleware
```go
// Global middleware
r.Use(gin.Logger())
r.Use(gin.Recovery())

// Route-specific middleware
r.GET("/admin", authMiddleware(), adminHandler)

// Custom middleware
func authMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        token := c.GetHeader("Authorization")
        if token == "" {
            c.AbortWithStatusJSON(401, gin.H{"error": "unauthorized"})
            return
        }
        c.Next()  // Continue to handler
    }
}
```

---

## ✅ Expected Behavior

```bash
# Health check
curl http://localhost:8080/health
# {"status":"ok"}

# Get all users
curl http://localhost:8080/api/users
# [{"id":1,"name":"John","email":"john@example.com"},...]

# Get user by ID
curl http://localhost:8080/api/users/1
# {"id":1,"name":"John","email":"john@example.com"}

# Create user
curl -X POST http://localhost:8080/api/users \
  -H "Content-Type: application/json" \
  -d '{"name":"Bob","email":"bob@example.com"}'
# {"id":3,"name":"Bob","email":"bob@example.com"}

# Update user
curl -X PUT http://localhost:8080/api/users/1 \
  -H "Content-Type: application/json" \
  -d '{"name":"John Updated","email":"john.updated@example.com"}'
# {"id":1,"name":"John Updated","email":"john.updated@example.com"}

# Delete user
curl -X DELETE http://localhost:8080/api/users/1
# {"message":"user deleted"}

# Validation error
curl -X POST http://localhost:8080/api/users \
  -d '{"name":"Bob"}'
# {"error":"Key: 'User.Email' Error:Field validation for 'Email' failed..."}
```

---

## 🎓 What You'll Learn

1. **Gin router** - Clean routing like Spring
2. **Path parameters** - `:id` syntax
3. **JSON binding** - Auto-parse request body
4. **Validation** - Built-in with `binding` tags
5. **Route groups** - Organize routes
6. **gin.H** - Quick JSON maps

---

## 🆚 Spring Boot vs Gin Comparison

| Feature | Spring Boot | Gin |
|---------|-------------|-----|
| Routing | `@GetMapping` | `r.GET()` |
| Path param | `@PathVariable` | `c.Param()` |
| Query param | `@RequestParam` | `c.Query()` |
| Body | `@RequestBody` | `c.BindJSON()` |
| Validation | `@Valid` | `binding:"required"` |
| Groups | `@RequestMapping` | `r.Group()` |
| Middleware | `@Component` Filter | `r.Use()` |
| JSON | Auto Jackson | `c.JSON()` |

---

## 📝 Sample main.go

```go
package main

import (
    "github.com/gin-gonic/gin"
)

type User struct {
    ID    int    `json:"id"`
    Name  string `json:"name" binding:"required"`
    Email string `json:"email" binding:"required,email"`
}

var users = []User{
    {ID: 1, Name: "John", Email: "john@example.com"},
    {ID: 2, Name: "Jane", Email: "jane@example.com"},
}

func main() {
    r := gin.Default()
    
    // Health check
    r.GET("/health", func(c *gin.Context) {
        c.JSON(200, gin.H{"status": "ok"})
    })
    
    // API routes
    api := r.Group("/api")
    {
        api.GET("/users", getUsers)
        api.GET("/users/:id", getUserByID)
        api.POST("/users", createUser)
        api.PUT("/users/:id", updateUser)
        api.DELETE("/users/:id", deleteUser)
    }
    
    r.Run(":8080")
}

// TODO: Implement handlers
func getUsers(c *gin.Context) {
    c.JSON(200, users)
}

func getUserByID(c *gin.Context) {
    // Your implementation
}

func createUser(c *gin.Context) {
    // Your implementation
}

func updateUser(c *gin.Context) {
    // Your implementation
}

func deleteUser(c *gin.Context) {
    // Your implementation
}
```

---

## ⏱️ Estimated Time: 20-25 minutes

Much cleaner than raw net/http! 🚀
