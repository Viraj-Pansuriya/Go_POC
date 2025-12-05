package middleware

import (
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/viraj/go-mono-repo/projects/hello-world/exercises/09-middleware-auth/auth"
)

// Much simpler! Just a function
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// accessToken
		token := ExtractToken(c)
		claims, err := auth.ValidateToken(token)
		if err != nil {
			c.AbortWithStatusJSON(401, gin.H{"error": "unauthorized"})
			return
		}
		c.Set("userId", claims.UserID)
		c.Next() // Continue to next handler
	}
}
func ExtractToken(c *gin.Context) string {
	brr := c.Request.Header.Get("Authorization")
	str, ok := strings.CutPrefix(brr, "Bearer ")
	if ok == false {
		return brr
	}
	return str
}
func RoleRequired(s string) gin.HandlerFunc {
	return func(c *gin.Context) {
		token := ExtractToken(c)
		claims, err := auth.ValidateToken(token)
		if err != nil {
			c.JSON(401, gin.H{"error": "unauthorized"})
			return
		}
		if s != claims.Role {
			c.JSON(401, gin.H{"error": "unauthorized due to role mismatch"})
			return
		}
		c.Next()
	}
}
