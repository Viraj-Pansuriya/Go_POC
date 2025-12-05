package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"10-configuration-management/config"

	"github.com/gin-gonic/gin"
)

func main() {
	// Get profile from environment (like Spring's SPRING_PROFILES_ACTIVE)
	profile := os.Getenv("APP_PROFILE")
	if profile == "" {
		profile = "dev" // Default to dev
	}

	fmt.Printf("🔧 Loading configuration with profile: %s\n", profile)

	// Load configuration
	cfg, err := config.LoadConfig("./configs", profile)
	if err != nil {
		log.Fatalf("❌ Failed to load config: %v", err)
	}

	// Validate configuration (fail fast!)
	if err := cfg.Validate(); err != nil {
		log.Fatalf("❌ Invalid configuration: %v", err)
	}

	// Print loaded config (useful for debugging)
	printConfig(cfg)

	// Set Gin mode based on debug flag
	if cfg.App.Debug {
		gin.SetMode(gin.DebugMode)
	} else {
		gin.SetMode(gin.ReleaseMode)
	}

	// Create router
	r := gin.Default()

	// Health check endpoint
	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status":  "UP",
			"app":     cfg.App.Name,
			"version": cfg.App.Version,
			"profile": profile,
		})
	})

	// Config info endpoint (don't expose in prod!)
	r.GET("/config", func(c *gin.Context) {
		if !cfg.App.Debug {
			c.JSON(http.StatusForbidden, gin.H{"error": "not available in production"})
			return
		}
		// Mask sensitive values
		c.JSON(http.StatusOK, gin.H{
			"server": gin.H{
				"host": cfg.Server.Host,
				"port": cfg.Server.Port,
			},
			"database": gin.H{
				"host":   cfg.Database.Host,
				"port":   cfg.Database.Port,
				"dbname": cfg.Database.DBName,
			},
			"app": gin.H{
				"name":    cfg.App.Name,
				"version": cfg.App.Version,
				"debug":   cfg.App.Debug,
			},
		})
	})

	// Start server
	addr := fmt.Sprintf("%s:%d", cfg.Server.Host, cfg.Server.Port)
	fmt.Printf("🚀 Starting %s v%s on %s\n", cfg.App.Name, cfg.App.Version, addr)

	if err := r.Run(addr); err != nil {
		log.Fatalf("❌ Failed to start server: %v", err)
	}
}

func printConfig(cfg *config.Config) {
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	fmt.Printf("📦 App: %s v%s\n", cfg.App.Name, cfg.App.Version)
	fmt.Printf("🖥️  Server: %s:%d\n", cfg.Server.Host, cfg.Server.Port)
	fmt.Printf("🗄️  Database: %s@%s:%d/%s\n",
		cfg.Database.User, cfg.Database.Host, cfg.Database.Port, cfg.Database.DBName)
	fmt.Printf("🐛 Debug: %v\n", cfg.App.Debug)
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
}

