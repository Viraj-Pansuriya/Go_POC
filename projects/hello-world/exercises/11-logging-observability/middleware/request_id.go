package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

// RequestID adds a unique request ID to each request
// Similar to Spring's MDC (Mapped Diagnostic Context)
func RequestID(logger *zap.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Check for existing request ID (from upstream service)
		requestID := c.GetHeader("X-Request-ID")
		if requestID == "" {
			requestID = uuid.New().String()
		}

		// Store in context for other handlers
		c.Set("request_id", requestID)

		// Add to response header for client reference
		c.Header("X-Request-ID", requestID)

		// Create child logger with request context
		// All logs using this logger will include request_id
		reqLogger := logger.With(
			zap.String("request_id", requestID),
		)
		c.Set("logger", reqLogger)

		c.Next()
	}
}

// GetLogger retrieves the request-scoped logger from context
func GetLogger(c *gin.Context) *zap.Logger {
	if logger, exists := c.Get("logger"); exists {
		return logger.(*zap.Logger)
	}
	// Fallback to no-op logger
	return zap.NewNop()
}
