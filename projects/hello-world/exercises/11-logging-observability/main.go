package main

import (
	"fmt"
	"net/http"

	"11-logging-observability/config"
	"11-logging-observability/handler"
	"11-logging-observability/logger"
	"11-logging-observability/middleware"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func main() {
	// Load configuration
	cfg, err := config.LoadConfig("./configs")
	if err != nil {
		panic("failed to load config: " + err.Error())
	}

	// Initialize logger
	log, err := logger.NewLogger(cfg.Log)
	if err != nil {
		panic("failed to init logger: " + err.Error())
	}
	defer log.Sync() // Flush buffered logs on exit

	log.Info("application starting",
		zap.String("app", cfg.App.Name),
		zap.String("version", cfg.App.Version),
		zap.String("log_level", cfg.Log.Level),
	)

	// Set Gin mode (no default logging)
	if cfg.Log.Development {
		gin.SetMode(gin.DebugMode)
	} else {
		gin.SetMode(gin.ReleaseMode)
	}

	// Create router without default middleware
	r := gin.New()

	// Add our custom middleware
	r.Use(gin.Recovery())                // Panic recovery
	r.Use(middleware.RequestID(log))     // Add request ID
	r.Use(middleware.RequestLogger(log)) // Log all requests

	// Create handlers
	userHandler := handler.NewUserHandler(log)

	// Health check
	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status":  "UP",
			"app":     cfg.App.Name,
			"version": cfg.App.Version,
		})
	})

	// User routes
	users := r.Group("/users")
	{
		users.GET("", userHandler.ListUsers)
		users.GET("/:id", userHandler.GetUser)
		users.POST("", userHandler.CreateUser)
		users.DELETE("/:id", userHandler.DeleteUser)
	}

	// Start server
	addr := fmt.Sprintf("%s:%d", cfg.Server.Host, cfg.Server.Port)
	log.Info("server starting",
		zap.String("address", addr),
	)

	if err := r.Run(addr); err != nil {
		log.Fatal("server failed to start", zap.Error(err))
	}
}
