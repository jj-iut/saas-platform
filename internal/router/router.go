package router

import (
	"github.com/gin-gonic/gin"

	"github.com/yourcompany/saas-platform/internal/handlers"
)

func SetupRouter(healthHandler *handlers.HealthHandler) *gin.Engine {
	// Set Gin mode based on environment
	gin.SetMode(gin.ReleaseMode)

	r := gin.New()

	// Middleware
	r.Use(gin.Logger())
	r.Use(gin.Recovery())
	r.Use(corsMiddleware())
	r.Use(requestIDMiddleware())

	// Health check endpoint
	r.GET("/health", healthHandler.HealthCheck)

	// API routes
	// v1 := r.Group("/api/v1")
	// {
	// 	// Add your API routes here
	// 	// Example: v1.GET("/users", userHandler.GetUsers)
	// }

	return r
}

func corsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE, PATCH")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}

func requestIDMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		requestID := c.GetHeader("X-Request-ID")
		if requestID == "" {
			requestID = c.GetString("X-Request-ID")
		}
		c.Header("X-Request-ID", requestID)
		c.Next()
	}
}
