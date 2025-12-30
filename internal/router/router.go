package router

import (
	"github.com/gin-gonic/gin"

	"github.com/yourcompany/saas-platform/internal/config"
	"github.com/yourcompany/saas-platform/internal/handlers"
	"github.com/yourcompany/saas-platform/internal/middleware"
	authModule "github.com/yourcompany/saas-platform/internal/modules/auth"
	restaurantsModule "github.com/yourcompany/saas-platform/internal/modules/restaurants"
)

func SetupRouter(
	cfg *config.Config,
	healthHandler *handlers.HealthHandler,
	authHandler *authModule.Handler,
	restaurantsHandler *restaurantsModule.Handler,
) *gin.Engine {
	// Set Gin mode based on environment
	if cfg.Server.Environment == "development" {
		gin.SetMode(gin.DebugMode)
	} else {
		gin.SetMode(gin.ReleaseMode)
	}

	r := gin.New()

	// Middleware
	r.Use(gin.Logger())
	r.Use(gin.Recovery())
	r.Use(corsMiddleware())
	r.Use(requestIDMiddleware())

	// Health check endpoint
	r.GET("/health", healthHandler.HealthCheck)

	// Public routes
	api := r.Group("/api/v1")
	{
		// Auth routes (public)
		auth := api.Group("/auth")
		{
			auth.POST("/register", authHandler.Register)
			auth.POST("/login", authHandler.Login)
			auth.POST("/refresh", authHandler.RefreshToken)
		}

		// Protected routes
		protected := api.Group("")
		protected.Use(middleware.AuthMiddleware(cfg.JWT.AccessSecret))
		{
			// User routes
			protected.GET("/me", authHandler.GetMe)

			// Restaurant routes (only superadmin)
			restaurants := protected.Group("/restaurants")
			restaurants.Use(middleware.RequireSuperAdmin())
			{
				restaurants.GET("", restaurantsHandler.GetAll)
				restaurants.GET("/:id", restaurantsHandler.GetByID)
				restaurants.POST("", restaurantsHandler.Create)
				restaurants.PUT("/:id", restaurantsHandler.Update)
				restaurants.DELETE("/:id", restaurantsHandler.Delete)
			}
		}
	}

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
