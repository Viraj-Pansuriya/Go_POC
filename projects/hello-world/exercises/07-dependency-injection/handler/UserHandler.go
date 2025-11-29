package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/viraj/go-mono-repo/projects/hello-world/exercises/07-dependency-injection/Model"
	"github.com/viraj/go-mono-repo/projects/hello-world/exercises/07-dependency-injection/service"
)

type userHandler struct {
	service service.UserService // Depends on interface!
}

func NewUserHandler(svc service.UserService) *userHandler {
	return &userHandler{service: svc}
}

func (h *userHandler) DeleteUser(c *gin.Context) {
	id := c.Param("id")
	_ = h.service.DeleteUser(id)
	c.JSON(200, gin.H{"id": id})
}

func (h *userHandler) GetUser(c *gin.Context) {
	id := c.Param("id")
	user, err := h.service.FindById(id)
	if err != nil {
		c.JSON(404, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, user)
}

func (h *userHandler) CreateUser(c *gin.Context) {
	var req Model.CreateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	user, err := h.service.AddUser(req.ID, req.Name, req.Email)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(201, user)
}

// RegisterRoutes adds routes to Gin engine
func (h *userHandler) RegisterRoutes(r *gin.Engine) {
	api := r.Group("/api")
	// NOTE : below braces are just a block scope only nothing fancy;
	{
		api.GET("/users/:id", h.GetUser)
		api.POST("/users", h.CreateUser)
		api.DELETE("/users/:id", h.DeleteUser)
	}
}
