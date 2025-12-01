package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/viraj/go-mono-repo/projects/hello-world/exercises/08-database-integration/service"
)

// =============================================================================
// USER HANDLER - HTTP Layer (Controller)
// =============================================================================
// In Java/Spring: @RestController with @RequestMapping
// Go/Gin: Struct with methods that handle HTTP requests

// UserHandler handles HTTP requests for user operations
type UserHandler struct {
	service service.UserService
}

// NewUserHandler creates a UserHandler with injected service
// Java equivalent: @RestController with @Autowired service
func NewUserHandler(service service.UserService) *UserHandler {
	return &UserHandler{service: service}
}

// RegisterRoutes sets up routes for user endpoints
// Java equivalent: @RequestMapping annotations on methods
func (h *UserHandler) RegisterRoutes(r *gin.Engine) {
	// Group routes under /api/users
	// Java: @RequestMapping("/api/users") on class
	users := r.Group("/api/users")
	{
		users.POST("", h.Register)           // POST /api/users
		users.GET("", h.GetAll)              // GET /api/users
		users.GET("/:id", h.GetByID)         // GET /api/users/:id
		users.PUT("/:id", h.Update)          // PUT /api/users/:id
		users.DELETE("/:id", h.Delete)       // DELETE /api/users/:id
		users.GET("/search", h.Search)       // GET /api/users/search?q=xxx
		users.PUT("/:id/profile", h.UpdateProfile) // PUT /api/users/:id/profile
	}
}

// =============================================================================
// REQUEST/RESPONSE DTOs
// =============================================================================
// In Java: Separate DTO classes or records
// Go: Structs with json tags

// RegisterRequest is the DTO for user registration
// Java: public record RegisterRequest(String name, String email, Integer age) {}
type RegisterRequest struct {
	Name  string `json:"name" binding:"required"`
	Email string `json:"email" binding:"required,email"`
	Age   int    `json:"age" binding:"gte=0"`
}

// UpdateProfileRequest is the DTO for profile updates
type UpdateProfileRequest struct {
	Bio       string `json:"bio"`
	AvatarURL string `json:"avatar_url"`
	Website   string `json:"website"`
}

// UserResponse is the DTO for user responses
// Separate from model to control what's exposed
type UserResponse struct {
	ID    uint   `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
	Age   int    `json:"age"`
}

// PaginatedResponse wraps paginated results
type PaginatedResponse struct {
	Data       interface{} `json:"data"`
	Page       int         `json:"page"`
	PageSize   int         `json:"page_size"`
	TotalCount int64       `json:"total_count"`
	TotalPages int         `json:"total_pages"`
}

// ErrorResponse is the standard error response
type ErrorResponse struct {
	Error string `json:"error"`
}

// =============================================================================
// HANDLER METHODS
// =============================================================================

// Register handles POST /api/users
// Java: @PostMapping public ResponseEntity<User> register(@RequestBody @Valid RegisterRequest request)
func (h *UserHandler) Register(c *gin.Context) {
	var req RegisterRequest

	// Bind and validate JSON body
	// Like @RequestBody @Valid in Spring
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
		return
	}

	// Call service
	user, err := h.service.Register(req.Name, req.Email, req.Age)
	if err != nil {
		// Determine status code based on error type
		// In Spring, you'd use @ExceptionHandler or throw specific exceptions
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
		return
	}

	// Return created user
	// 201 Created with the new resource
	c.JSON(http.StatusCreated, user)
}

// GetAll handles GET /api/users
// Java: @GetMapping public ResponseEntity<List<User>> getAll(@RequestParam Optional<Integer> page, ...)
func (h *UserHandler) GetAll(c *gin.Context) {
	// Check for pagination parameters
	pageStr := c.Query("page")
	pageSizeStr := c.Query("page_size")

	if pageStr != "" || pageSizeStr != "" {
		// Paginated response
		page, _ := strconv.Atoi(pageStr)
		pageSize, _ := strconv.Atoi(pageSizeStr)

		if page < 1 {
			page = 1
		}
		if pageSize < 1 {
			pageSize = 10
		}

		users, total, err := h.service.GetAllPaginated(page, pageSize)
		if err != nil {
			c.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
			return
		}

		totalPages := int(total) / pageSize
		if int(total)%pageSize > 0 {
			totalPages++
		}

		c.JSON(http.StatusOK, PaginatedResponse{
			Data:       users,
			Page:       page,
			PageSize:   pageSize,
			TotalCount: total,
			TotalPages: totalPages,
		})
		return
	}

	// Non-paginated response
	users, err := h.service.GetAll()
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, users)
}

// GetByID handles GET /api/users/:id
// Java: @GetMapping("/{id}") public ResponseEntity<User> getById(@PathVariable Long id)
func (h *UserHandler) GetByID(c *gin.Context) {
	// Parse path parameter
	// In Spring: @PathVariable Long id (auto-parsed)
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "invalid user id"})
		return
	}

	// Check for profile inclusion
	// GET /api/users/1?include=profile
	include := c.Query("include")

	var user interface{}
	if include == "profile" {
		user, err = h.service.GetByIDWithProfile(uint(id))
	} else {
		user, err = h.service.GetByID(uint(id))
	}

	if err != nil {
		c.JSON(http.StatusNotFound, ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, user)
}

// Update handles PUT /api/users/:id
// Java: @PutMapping("/{id}") public ResponseEntity<User> update(@PathVariable Long id, @RequestBody UpdateRequest request)
func (h *UserHandler) Update(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "invalid user id"})
		return
	}

	// Get existing user
	user, err := h.service.GetByID(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, ErrorResponse{Error: err.Error()})
		return
	}

	// Bind update data
	var req RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
		return
	}

	// Update fields
	user.Name = req.Name
	user.Email = req.Email
	user.Age = req.Age

	c.JSON(http.StatusOK, user)
}

// Delete handles DELETE /api/users/:id
// Java: @DeleteMapping("/{id}") public ResponseEntity<Void> delete(@PathVariable Long id)
func (h *UserHandler) Delete(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "invalid user id"})
		return
	}

	if err := h.service.Delete(uint(id)); err != nil {
		c.JSON(http.StatusNotFound, ErrorResponse{Error: err.Error()})
		return
	}

	// 204 No Content for successful delete
	c.Status(http.StatusNoContent)
}

// Search handles GET /api/users/search?q=xxx
// Java: @GetMapping("/search") public List<User> search(@RequestParam String q)
func (h *UserHandler) Search(c *gin.Context) {
	query := c.Query("q")

	users, err := h.service.SearchUsers(query)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, users)
}

// UpdateProfile handles PUT /api/users/:id/profile
func (h *UserHandler) UpdateProfile(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "invalid user id"})
		return
	}

	var req UpdateProfileRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
		return
	}

	if err := h.service.UpdateProfile(uint(id), req.Bio, req.AvatarURL, req.Website); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "profile updated"})
}

