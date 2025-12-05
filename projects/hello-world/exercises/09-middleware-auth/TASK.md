# Exercise 09: Middleware & Authentication (JWT)

## 🎯 Goal
Learn how middleware works in Go/Gin and implement JWT-based authentication - just like Spring Security!

---

## 📚 Spring Security vs Go Middleware

### Spring Security
```java
@Configuration
@EnableWebSecurity
public class SecurityConfig {
    @Bean
    public SecurityFilterChain filterChain(HttpSecurity http) {
        http
            .authorizeRequests()
            .antMatchers("/api/public/**").permitAll()
            .antMatchers("/api/admin/**").hasRole("ADMIN")
            .anyRequest().authenticated()
            .and()
            .addFilterBefore(jwtFilter, UsernamePasswordAuthenticationFilter.class);
        return http.build();
    }
}

@Component
public class JwtFilter extends OncePerRequestFilter {
    @Override
    protected void doFilterInternal(HttpServletRequest request, 
                                     HttpServletResponse response, 
                                     FilterChain chain) {
        String token = extractToken(request);
        if (isValid(token)) {
            SecurityContextHolder.getContext().setAuthentication(auth);
        }
        chain.doFilter(request, response);
    }
}
```

### Go/Gin Middleware
```go
// Much simpler! Just a function
func AuthMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        token := extractToken(c)
        if !isValid(token) {
            c.AbortWithStatusJSON(401, gin.H{"error": "unauthorized"})
            return
        }
        c.Set("userID", claims.UserID)  // Like SecurityContext
        c.Next()  // Continue to next handler
    }
}

// Apply to routes
r.GET("/api/public", publicHandler)                    // No auth
r.GET("/api/private", AuthMiddleware(), privateHandler) // With auth
```

---

## 🏗️ What is Middleware?

Middleware = Functions that run **before** or **after** your handler

```
Request → [Logger] → [Auth] → [RateLimit] → Handler → [ResponseTime] → Response
              ↑         ↑          ↑                         ↑
           Middleware  Middleware  Middleware             Middleware
```

### Gin Middleware Signature
```go
func MyMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        // Before handler
        start := time.Now()
        
        c.Next()  // Call next handler/middleware
        
        // After handler
        duration := time.Since(start)
        log.Printf("Request took %v", duration)
    }
}
```

---

## 🚀 Your Tasks

### Task 1: Basic Logging Middleware
Create a middleware that logs every request.

```go
// middleware/logger.go
func Logger() gin.HandlerFunc {
    return func(c *gin.Context) {
        // TODO: Log request start time, method, path
        // TODO: Call c.Next()
        // TODO: Log status code, duration
    }
}
```

### Task 2: JWT Authentication
Implement JWT token generation and validation.

**Install JWT library:**
```bash
go get -u github.com/golang-jwt/jwt/v5
```

```go
// auth/jwt.go
type Claims struct {
    UserID uint   `json:"user_id"`
    Email  string `json:"email"`
    Role   string `json:"role"`
    jwt.RegisteredClaims
}

func GenerateToken(userID uint, email, role string) (string, error) {
    // TODO: Create claims with expiration
    // TODO: Sign with secret key
    // TODO: Return token string
}

func ValidateToken(tokenString string) (*Claims, error) {
    // TODO: Parse and validate token
    // TODO: Return claims or error
}
```

### Task 3: Auth Middleware
Create middleware that validates JWT from Authorization header.

```go
// middleware/auth.go
func AuthRequired() gin.HandlerFunc {
    return func(c *gin.Context) {
        // TODO: Extract token from "Authorization: Bearer <token>"
        // TODO: Validate token
        // TODO: Set user info in context
        // TODO: c.AbortWithStatusJSON(401, ...) if invalid
        // TODO: c.Next() if valid
    }
}
```

### Task 4: Role-Based Authorization
Create middleware that checks user roles.

```go
// middleware/roles.go
func RequireRole(roles ...string) gin.HandlerFunc {
    return func(c *gin.Context) {
        // TODO: Get user role from context
        // TODO: Check if user has required role
        // TODO: c.AbortWithStatusJSON(403, ...) if not authorized
    }
}
```

### Task 5: Wire Everything Together

```go
// main.go
func main() {
    r := gin.Default()
    
    // Global middleware (applies to all routes)
    r.Use(middleware.Logger())
    r.Use(middleware.CORS())
    
    // Public routes (no auth)
    public := r.Group("/api/v1")
    {
        public.POST("/register", handler.Register)
        public.POST("/login", handler.Login)
    }
    
    // Protected routes (auth required)
    protected := r.Group("/api/v1")
    protected.Use(middleware.AuthRequired())
    {
        protected.GET("/profile", handler.GetProfile)
        protected.PUT("/profile", handler.UpdateProfile)
    }
    
    // Admin routes (auth + admin role required)
    admin := r.Group("/api/v1/admin")
    admin.Use(middleware.AuthRequired())
    admin.Use(middleware.RequireRole("admin"))
    {
        admin.GET("/users", handler.ListAllUsers)
        admin.DELETE("/users/:id", handler.DeleteUser)
    }
    
    r.Run(":8080")
}
```

---

## 📁 Files to Create

```
09-middleware-auth/
├── main.go
├── go.mod
├── auth/
│   └── jwt.go           # Token generation/validation
├── middleware/
│   ├── logger.go        # Request logging
│   ├── auth.go          # JWT validation middleware
│   ├── roles.go         # Role-based authorization
│   └── cors.go          # CORS handling
├── handler/
│   ├── auth_handler.go  # Login, Register endpoints
│   └── user_handler.go  # Protected user endpoints
├── model/
│   └── user.go          # User model
└── service/
    └── auth_service.go  # Business logic
```

---

## 💡 Key Patterns

### 1. c.Set() / c.Get() - Request-scoped storage
```go
// In middleware
c.Set("userID", claims.UserID)

// In handler
userID, exists := c.Get("userID")
if !exists {
    // Handle error
}
```

### 2. c.Next() vs c.Abort()
```go
c.Next()   // Continue to next middleware/handler
c.Abort()  // Stop chain, but still runs deferred
c.AbortWithStatusJSON(401, gin.H{"error": "unauthorized"})
```

### 3. Middleware Order Matters!
```go
r.Use(Logger())      // First: logs everything
r.Use(Recovery())    // Second: catches panics
r.Use(AuthRequired()) // Third: checks auth
// Handler runs last
```

---

## 🆚 Spring Security Comparison

| Spring Security | Go/Gin | Notes |
|-----------------|--------|-------|
| `SecurityFilterChain` | `r.Use(middleware)` | Global middleware |
| `@PreAuthorize("hasRole")` | `RequireRole("admin")` | Role check |
| `SecurityContextHolder` | `c.Get("userID")` | Request context |
| `OncePerRequestFilter` | `gin.HandlerFunc` | Middleware function |
| `UserDetailsService` | Custom service | Load user details |
| `@EnableWebSecurity` | Not needed | No annotations in Go |

---

## ✅ Expected Behavior

```bash
# Register
curl -X POST http://localhost:8080/api/v1/register \
  -H "Content-Type: application/json" \
  -d '{"email":"user@test.com","password":"secret123"}'
# {"message": "registered successfully"}

# Login
curl -X POST http://localhost:8080/api/v1/login \
  -H "Content-Type: application/json" \
  -d '{"email":"user@test.com","password":"secret123"}'
# {"token": "eyJhbGciOiJIUzI1NiIs..."}

# Access protected route (with token)
curl http://localhost:8080/api/v1/profile \
  -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIs..."
# {"id": 1, "email": "user@test.com"}

# Access without token
curl http://localhost:8080/api/v1/profile
# {"error": "authorization header required"} (401)

# Access admin route without admin role
curl http://localhost:8080/api/v1/admin/users \
  -H "Authorization: Bearer <user_token>"
# {"error": "forbidden"} (403)
```

---

## 🎓 What You'll Learn

1. **Middleware pattern** - Request/response pipeline
2. **JWT basics** - Token structure, signing, validation
3. **Authentication** - Verifying identity
4. **Authorization** - Checking permissions
5. **Context passing** - Sharing data across middleware
6. **Error responses** - Proper HTTP status codes

---

## ⏱️ Estimated Time: 30-40 minutes

This is where your app becomes secure! 🔐

