package handler

import (
	"net/http"
	"strconv"
	"time"

	"11-logging-observability/middleware"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// User represents a simple user model
type User struct {
	ID        uint      `json:"id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"created_at"`
}

// UserHandler handles user-related HTTP requests
type UserHandler struct {
	logger *zap.Logger
	// In real app: userService UserService
}

// NewUserHandler creates a new UserHandler
func NewUserHandler(logger *zap.Logger) *UserHandler {
	return &UserHandler{
		logger: logger.With(zap.String("handler", "user")),
	}
}

// GetUser retrieves a user by ID
// Demonstrates request-scoped logging
func (h *UserHandler) GetUser(c *gin.Context) {
	// Get request-scoped logger (includes request_id)
	logger := middleware.GetLogger(c)

	// Parse user ID
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		logger.Warn("invalid user id",
			zap.String("id_param", idStr),
			zap.Error(err),
		)
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user id"})
		return
	}

	logger.Info("fetching user", zap.Uint64("user_id", id))

	// Simulate database lookup
	user := &User{
		ID:        uint(id),
		Name:      "John Doe",
		Email:     "john@example.com",
		CreatedAt: time.Now().Add(-24 * time.Hour),
	}

	logger.Debug("user found",
		zap.Uint("user_id", user.ID),
		zap.String("email", user.Email),
	)

	c.JSON(http.StatusOK, user)
}

// ListUsers lists all users
func (h *UserHandler) ListUsers(c *gin.Context) {
	logger := middleware.GetLogger(c)

	logger.Info("listing users")

	// Simulate slow operation
	start := time.Now()
	time.Sleep(50 * time.Millisecond)

	users := []User{
		{ID: 1, Name: "John Doe", Email: "john@example.com"},
		{ID: 2, Name: "Jane Smith", Email: "jane@example.com"},
	}

	logger.Info("users retrieved",
		zap.Int("count", len(users)),
		zap.Duration("db_latency", time.Since(start)),
	)

	c.JSON(http.StatusOK, users)
}

// CreateUser creates a new user
// Demonstrates error logging
func (h *UserHandler) CreateUser(c *gin.Context) {
	logger := middleware.GetLogger(c)

	var user User
	if err := c.ShouldBindJSON(&user); err != nil {
		logger.Error("failed to parse user",
			zap.Error(err),
		)
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	logger.Info("creating user",
		zap.String("name", user.Name),
		zap.String("email", user.Email),
	)

	// Simulate creation
	user.ID = 123
	user.CreatedAt = time.Now()

	logger.Info("user created",
		zap.Uint("user_id", user.ID),
	)

	c.JSON(http.StatusCreated, user)
}

// DeleteUser demonstrates error logging
func (h *UserHandler) DeleteUser(c *gin.Context) {
	logger := middleware.GetLogger(c)

	idStr := c.Param("id")
	id, _ := strconv.ParseUint(idStr, 10, 32)

	// Simulate error condition
	if id == 999 {
		logger.Error("failed to delete user",
			zap.Uint64("user_id", id),
			zap.String("reason", "user not found"),
		)
		c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
		return
	}

	logger.Info("user deleted", zap.Uint64("user_id", id))
	c.JSON(http.StatusOK, gin.H{"message": "user deleted"})
}

