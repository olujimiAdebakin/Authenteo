package router

import (
	"net/http"

	"authentio/internal/handler"
	"authentio/internal/middleware"
	"authentio/pkg/logger"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// SetupRouter configures the Gin engine with routes, middlewares, and health checks.
func SetupRouter(h *handler.Handler) *gin.Engine {
	// Initialize the Gin engine
	r := gin.New()

	// --- Global Middlewares ---
	r.Use(gin.Recovery())                         // Recover from panics
	r.Use(middleware.RequestLogger())              // Custom structured request logger
	r.Use(middleware.CORSMiddleware())             // CORS handler
	r.Use(middleware.RateLimiterMiddlewareRedis()) // Optional: switch to in-mem for local dev

	// --- Health check ---
	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	// --- API v1 routes ---
	api := r.Group("/api/v1")
	{
		auth := api.Group("/auth")
		{
			auth.POST("/register", h.Auth.Register)
			auth.POST("/login", h.Auth.Login)
			auth.POST("/refresh", h.Auth.RefreshToken)

			// Two-factor routes
			auth.POST("/2fa/send", h.TwoFA.SendOTP)
			auth.POST("/2fa/verify", h.TwoFA.VerifyOTP)
		}

		user := api.Group("/user")
		user.Use(middleware.AuthRequired())
		{
			user.GET("/me", h.User.GetProfile)
			user.PUT("/update", h.User.UpdateProfile)
		}
	}

	// --- 404 handler ---
	r.NoRoute(func(c *gin.Context) {
		logger.Logger.Warn("Route not found", zap.String("path", c.Request.URL.Path))
		c.JSON(http.StatusNotFound, gin.H{"error": "endpoint not found"})
	})

	return r
}
