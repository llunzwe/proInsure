package routes

import (
	"time"

	"github.com/gin-gonic/gin"
)

// Handlers represents all HTTP handlers
type Handlers struct{}

// NewHandlers creates new HTTP handlers
func NewHandlers() *Handlers {
	return &Handlers{}
}

// SetupRoutes configures all API routes
func SetupRoutes(router *gin.Engine, handlers *Handlers) {
	// Root endpoint
	router.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"service":     "SmartSure Go Backend",
			"version":     "1.0.0",
			"status":      "operational",
			"description": "Core Insurance Logic & Business Operations",
			"endpoints": gin.H{
				"health": "/health",
			},
			"timestamp": time.Now().UTC(),
		})
	})

	// API version 1
	v1 := router.Group("/api/v1")
	{
		// Health check
		v1.GET("/health", func(c *gin.Context) {
			c.JSON(200, gin.H{
				"status":    "healthy",
				"service":   "smartsure-api",
				"version":   "1.0.0",
				"timestamp": time.Now().UTC(),
			})
		})
	}
}
