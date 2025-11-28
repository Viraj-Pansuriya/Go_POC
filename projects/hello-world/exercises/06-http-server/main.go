package main

import (
	"github.com/gin-gonic/gin"
	"github.com/viraj/go-mono-repo/projects/hello-world/exercises/06-http-server/controller"
)

func main() {
	r := gin.Default()
	r.Use(gin.Logger())
	r.Use(gin.Recovery())
	// Health check
	r.GET("/health", controller.HealthCheck)
	r.GET("/api/user/:id", authHandler(), controller.GetUser)
	r.POST("/api/user/add", controller.AddUser)
	r.DELETE("/api/user/delete/:id", controller.DeleteUser)
	r.Run(":8080")
}

// just to to learn stuff for middleware , there is no meaning of adding;
func authHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		qr := c.Query("token")
		if qr == "" {
			c.AbortWithStatusJSON(401, gin.H{"error": "unauthorized"})
			return
		}
		c.Next()
	}
}
