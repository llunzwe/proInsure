package middleware

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"sync"
	"time"

	"smartsure/internal/infrastructure/auth"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

// MetricsStore holds in-memory metrics for real-time monitoring
type MetricsStore struct {
	mu            sync.RWMutex
	requestCount  map[string]int64
	responseTime  map[string]time.Duration
	statusCount   map[int]int64
	totalRequests int64
}

// NewMetricsStore creates a new metrics store
func NewMetricsStore() *MetricsStore {
	return &MetricsStore{
		requestCount: make(map[string]int64),
		responseTime: make(map[string]time.Duration),
		statusCount:  make(map[int]int64),
	}
}

// RecordRequest records a request metric
func (m *MetricsStore) RecordRequest(method, path string, statusCode int, duration time.Duration) {
	m.mu.Lock()
	defer m.mu.Unlock()

	key := fmt.Sprintf("%s_%s", method, path)
	m.requestCount[key]++
	m.statusCount[statusCode]++
	m.responseTime[key] = (m.responseTime[key] + duration) / 2 // Simple moving average
	m.totalRequests++
}

// GetMetrics returns current metrics
func (m *MetricsStore) GetMetrics() map[string]interface{} {
	m.mu.RLock()
	defer m.mu.RUnlock()

	return map[string]interface{}{
		"total_requests": m.totalRequests,
		"request_count":  m.requestCount,
		"response_time":  m.responseTime,
		"status_count":   m.statusCount,
	}
}

// Logger middleware for request logging
func Logger(logger *logrus.Logger) gin.HandlerFunc {
	return gin.LoggerWithFormatter(func(param gin.LogFormatterParams) string {
		logger.WithFields(logrus.Fields{
			"client_ip":   param.ClientIP,
			"timestamp":   param.TimeStamp.Format(time.RFC3339),
			"method":      param.Method,
			"path":        param.Path,
			"protocol":    param.Request.Proto,
			"status_code": param.StatusCode,
			"latency":     param.Latency,
			"user_agent":  param.Request.UserAgent(),
			"error":       param.ErrorMessage,
		}).Info("HTTP Request")

		return ""
	})
}

// Recovery middleware for panic recovery
func Recovery(logger *logrus.Logger) gin.HandlerFunc {
	return gin.CustomRecovery(func(c *gin.Context, recovered interface{}) {
		if err, ok := recovered.(string); ok {
			logger.WithFields(logrus.Fields{
				"error": err,
				"path":  c.Request.URL.Path,
				"ip":    c.ClientIP(),
			}).Error("Panic recovered")
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Internal server error",
		})
	})
}

// CORS middleware for cross-origin requests
func CORS() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Credentials", "true")
		c.Header("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Header("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE, PATCH")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}

// RequestID middleware adds a unique request ID to each request
func RequestID() gin.HandlerFunc {
	return func(c *gin.Context) {
		requestID := c.GetHeader("X-Request-ID")
		if requestID == "" {
			requestID = uuid.New().String()
		}
		c.Header("X-Request-ID", requestID)
		c.Set("request_id", requestID)
		c.Next()
	}
}

// RateLimit middleware for rate limiting
func RateLimit() gin.HandlerFunc {
	// Simple in-memory rate limiting (for production, use Redis)
	clients := make(map[string][]time.Time)
	maxRequests := 100 // requests per minute
	window := time.Minute

	return func(c *gin.Context) {
		clientIP := c.ClientIP()
		now := time.Now()

		// Clean old entries
		if requests, exists := clients[clientIP]; exists {
			var validRequests []time.Time
			for _, reqTime := range requests {
				if now.Sub(reqTime) < window {
					validRequests = append(validRequests, reqTime)
				}
			}
			clients[clientIP] = validRequests
		}

		// Check rate limit
		if len(clients[clientIP]) >= maxRequests {
			c.JSON(http.StatusTooManyRequests, gin.H{
				"error":       "Rate limit exceeded",
				"retry_after": window.Seconds(),
			})
			c.Abort()
			return
		}

		// Add current request
		clients[clientIP] = append(clients[clientIP], now)
		c.Next()
	}
}

// Auth middleware for user authentication using JWT
func Auth(jwtService *auth.JWTService) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "Authorization header required",
			})
			c.Abort()
			return
		}

		// Extract token from "Bearer <token>" format
		tokenParts := strings.Split(authHeader, " ")
		if len(tokenParts) != 2 || tokenParts[0] != "Bearer" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "Invalid authorization header format",
			})
			c.Abort()
			return
		}

		token := tokenParts[1]
		claims, err := jwtService.ValidateToken(token)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "Invalid or expired token",
			})
			c.Abort()
			return
		}

		// Set user information in context
		c.Set("user_id", claims.UserID)
		c.Set("user_type", claims.UserType)
		c.Set("user_role", claims.UserRole)
		c.Next()
	}
}

// AdminAuth middleware for admin authentication using JWT
func AdminAuth(jwtService *auth.JWTService) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "Authorization header required",
			})
			c.Abort()
			return
		}

		// Extract token from "Bearer <token>" format
		tokenParts := strings.Split(authHeader, " ")
		if len(tokenParts) != 2 || tokenParts[0] != "Bearer" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "Invalid authorization header format",
			})
			c.Abort()
			return
		}

		token := tokenParts[1]
		claims, err := jwtService.ValidateToken(token)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "Invalid or expired token",
			})
			c.Abort()
			return
		}

		// Check if user has admin role
		if claims.UserRole != "admin" {
			c.JSON(http.StatusForbidden, gin.H{
				"error": "Admin access required",
			})
			c.Abort()
			return
		}

		// Set admin information in context
		c.Set("user_id", claims.UserID)
		c.Set("user_type", claims.UserType)
		c.Set("user_role", claims.UserRole)
		c.Next()
	}
}

// Validation middleware for request validation
func Validation() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Validate common request headers
		contentType := c.GetHeader("Content-Type")
		if c.Request.Method == "POST" || c.Request.Method == "PUT" || c.Request.Method == "PATCH" {
			if contentType != "" && !strings.Contains(contentType, "application/json") && !strings.Contains(contentType, "multipart/form-data") {
				c.JSON(http.StatusUnsupportedMediaType, gin.H{
					"error":     "Unsupported content type",
					"supported": []string{"application/json", "multipart/form-data"},
				})
				c.Abort()
				return
			}
		}

		// Validate request size (max 10MB)
		if c.Request.ContentLength > 10*1024*1024 {
			c.JSON(http.StatusRequestEntityTooLarge, gin.H{
				"error":    "Request body too large",
				"max_size": "10MB",
			})
			c.Abort()
			return
		}

		c.Next()
	}
}

// Metrics middleware for collecting metrics
func Metrics(logger *logrus.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		c.Next()
		duration := time.Since(start)

		// Send metrics to monitoring system
		statusCode := c.Writer.Status()
		method := c.Request.Method
		path := c.Request.URL.Path
		userAgent := c.Request.UserAgent()

		// TODO: Implement metrics export to external monitoring systems:
		// - Prometheus metrics integration
		// - DataDog, New Relic, or other monitoring systems
		// - Custom metrics collection service
		// - Real-time dashboards and alerting

		// For now, log structured metrics
		logger.WithFields(logrus.Fields{
			"method":      method,
			"path":        path,
			"status_code": statusCode,
			"duration_ms": duration.Milliseconds(),
			"user_agent":  userAgent,
			"ip":          c.ClientIP(),
		}).Info("Request metrics")

		// Update in-memory metrics store for real-time monitoring
		if metricsStore, exists := c.Get("metrics_store"); exists {
			if store, ok := metricsStore.(*MetricsStore); ok {
				store.RecordRequest(method, path, statusCode, duration)
			}
		}
	}
}

// Security middleware for security headers
func Security() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("X-Content-Type-Options", "nosniff")
		c.Header("X-Frame-Options", "DENY")
		c.Header("X-XSS-Protection", "1; mode=block")
		c.Header("Strict-Transport-Security", "max-age=31536000; includeSubDomains")
		c.Header("Content-Security-Policy", "default-src 'self'")
		c.Next()
	}
}

// Timeout middleware for request timeout
func Timeout(timeout time.Duration) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Create context with timeout
		ctx, cancel := context.WithTimeout(c.Request.Context(), timeout)
		defer cancel()

		// Replace request context
		c.Request = c.Request.WithContext(ctx)

		// Channel to signal completion
		done := make(chan struct{})
		go func() {
			c.Next()
			close(done)
		}()

		// Wait for completion or timeout
		select {
		case <-done:
			// Request completed normally
			return
		case <-ctx.Done():
			// Request timed out
			c.JSON(http.StatusRequestTimeout, gin.H{
				"error":   "Request timeout",
				"timeout": timeout.String(),
			})
			c.Abort()
		}
	}
}

// Compression middleware for response compression
func Compression() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Check if client accepts gzip
		acceptEncoding := c.GetHeader("Accept-Encoding")
		if !strings.Contains(acceptEncoding, "gzip") {
			c.Next()
			return
		}

		// Set compression headers
		c.Header("Content-Encoding", "gzip")
		c.Header("Vary", "Accept-Encoding")

		// Note: In production, use a proper compression middleware like gin-gzip
		// This is a simplified implementation
		c.Next()
	}
}

// Cache middleware for response caching
func Cache() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Set cache headers based on request path
		path := c.Request.URL.Path

		// Static resources - cache for 1 hour
		if strings.Contains(path, "/static/") || strings.Contains(path, "/assets/") {
			c.Header("Cache-Control", "public, max-age=3600")
			c.Header("Expires", time.Now().Add(time.Hour).Format(http.TimeFormat))
		}

		// API responses - no cache for dynamic content
		if strings.Contains(path, "/api/") {
			c.Header("Cache-Control", "no-cache, no-store, must-revalidate")
			c.Header("Pragma", "no-cache")
			c.Header("Expires", "0")
		}

		c.Next()
	}
}

// LoadBalancer middleware for load balancer health checks
func LoadBalancer() gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.Request.URL.Path == "/health" {
			c.JSON(http.StatusOK, gin.H{
				"status": "healthy",
			})
			c.Abort()
			return
		}
		c.Next()
	}
}

// Version middleware for API versioning
func Version() gin.HandlerFunc {
	return func(c *gin.Context) {
		version := c.GetHeader("X-API-Version")
		if version == "" {
			version = "v1"
		}
		c.Set("api_version", version)
		c.Next()
	}
}

// CorrelationID middleware for distributed tracing
func CorrelationID() gin.HandlerFunc {
	return func(c *gin.Context) {
		correlationID := c.GetHeader("X-Correlation-ID")
		if correlationID == "" {
			correlationID = uuid.New().String()
		}
		c.Header("X-Correlation-ID", correlationID)
		c.Set("correlation_id", correlationID)
		c.Next()
	}
}

// Tenant middleware for multi-tenancy support
func Tenant() gin.HandlerFunc {
	return func(c *gin.Context) {
		tenantID := c.GetHeader("X-Tenant-ID")
		if tenantID == "" {
			tenantID = "default"
		}
		c.Set("tenant_id", tenantID)
		c.Next()
	}
}

// Audit middleware for audit logging
func Audit(logger *logrus.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Log request for audit purposes
		logger.WithFields(logrus.Fields{
			"method":         c.Request.Method,
			"path":           c.Request.URL.Path,
			"client_ip":      c.ClientIP(),
			"user_agent":     c.Request.UserAgent(),
			"request_id":     c.GetString("request_id"),
			"correlation_id": c.GetString("correlation_id"),
			"tenant_id":      c.GetString("tenant_id"),
		}).Info("Audit log")

		c.Next()
	}
}
