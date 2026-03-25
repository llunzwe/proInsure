package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"smartsure/pkg/ratelimit"
)

// RateLimitMiddleware provides rate limiting for API endpoints
type RateLimitMiddleware struct {
	endpointLimiter *ratelimit.APIEndpointLimiter
	tenantLimiter   *ratelimit.TenantRateLimiter
}

// NewRateLimitMiddleware creates a new rate limit middleware
func NewRateLimitMiddleware(endpointLimiter *ratelimit.APIEndpointLimiter, tenantLimiter *ratelimit.TenantRateLimiter) *RateLimitMiddleware {
	return &RateLimitMiddleware{
		endpointLimiter: endpointLimiter,
		tenantLimiter:   tenantLimiter,
	}
}

// EndpointRateLimit returns Gin middleware for endpoint-based rate limiting
func (rlm *RateLimitMiddleware) EndpointRateLimit() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get client information
		clientIP := c.ClientIP()
		method := c.Request.Method
		path := c.Request.URL.Path

		// Get user ID from context if authenticated
		var userID *uuid.UUID
		if uid, exists := c.Get("user_id"); exists {
			if id, ok := uid.(uuid.UUID); ok {
				userID = &id
			}
		}

		// Check endpoint rate limit
		if !rlm.endpointLimiter.AllowEndpointRequest(method, path, clientIP, userID) {
			c.JSON(http.StatusTooManyRequests, gin.H{
				"error":   "rate_limit_exceeded",
				"message": "Too many requests. Please try again later.",
				"retry_after": 60, // seconds
			})
			c.Abort()
			return
		}

		c.Next()
	}
}

// TenantRateLimit returns Gin middleware for tenant-based rate limiting
func (rlm *RateLimitMiddleware) TenantRateLimit() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Skip if no tenant limiter configured
		if rlm.tenantLimiter == nil {
			c.Next()
			return
		}

		// Get tenant ID from header or context
		tenantIDStr := c.GetHeader("X-Tenant-ID")
		if tenantIDStr == "" {
			// Try to get from subdomain or path
			host := c.Request.Host
			if parts := strings.Split(host, "."); len(parts) >= 3 {
				tenantIDStr = parts[0]
			}
		}

		// If no tenant ID, skip tenant limiting
		if tenantIDStr == "" {
			c.Next()
			return
		}

		tenantID, err := uuid.Parse(tenantIDStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error":   "invalid_tenant_id",
				"message": "Invalid tenant ID format",
			})
			c.Abort()
			return
		}

		// Get user ID
		var userID uuid.UUID
		if uid, exists := c.Get("user_id"); exists {
			if id, ok := uid.(uuid.UUID); ok {
				userID = id
			} else {
				c.Next() // No user ID, skip tenant limiting
				return
			}
		} else {
			c.Next() // No user ID, skip tenant limiting
			return
		}

		// Check tenant rate limit
		allowed, err := rlm.tenantLimiter.AllowTenantRequest(c.Request.Context(), tenantID, userID)
		if err != nil {
			// Log error but allow request (fail open)
			c.Next()
			return
		}

		if !allowed {
			c.JSON(http.StatusTooManyRequests, gin.H{
				"error":   "tenant_rate_limit_exceeded",
				"message": "Tenant rate limit exceeded. Please try again later.",
				"retry_after": 300, // 5 minutes
			})
			c.Abort()
			return
		}

		c.Next()
	}
}

// CombinedRateLimit returns middleware that combines endpoint and tenant rate limiting
func (rlm *RateLimitMiddleware) CombinedRateLimit() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Apply endpoint rate limiting first
		rlm.EndpointRateLimit()(c)

		// If request was blocked, don't continue
		if c.IsAborted() {
			return
		}

		// Apply tenant rate limiting
		rlm.TenantRateLimit()(c)
	}
}

// IPWhitelistMiddleware provides IP whitelisting functionality
type IPWhitelistMiddleware struct {
	whitelistedIPs map[string]bool
}

// NewIPWhitelistMiddleware creates a new IP whitelist middleware
func NewIPWhitelistMiddleware(whitelistedIPs []string) *IPWhitelistMiddleware {
	whitelist := make(map[string]bool)
	for _, ip := range whitelistedIPs {
		whitelist[ip] = true
	}

	return &IPWhitelistMiddleware{
		whitelistedIPs: whitelist,
	}
}

// IPWhitelist returns middleware that restricts access to whitelisted IPs
func (iwm *IPWhitelistMiddleware) IPWhitelist() gin.HandlerFunc {
	return func(c *gin.Context) {
		clientIP := c.ClientIP()

		// Check if IP is whitelisted
		if !iwm.whitelistedIPs[clientIP] {
			c.JSON(http.StatusForbidden, gin.H{
				"error":   "ip_not_whitelisted",
				"message": "Access denied from this IP address",
			})
			c.Abort()
			return
		}

		c.Next()
	}
}

// GeoBlockMiddleware provides geographic blocking functionality
type GeoBlockMiddleware struct {
	blockedCountries map[string]bool
}

// NewGeoBlockMiddleware creates a new geographic blocking middleware
func NewGeoBlockMiddleware(blockedCountries []string) *GeoBlockMiddleware {
	blocklist := make(map[string]bool)
	for _, country := range blockedCountries {
		blocklist[strings.ToUpper(country)] = true
	}

	return &GeoBlockMiddleware{
		blockedCountries: blocklist,
	}
}

// GeoBlock returns middleware that blocks requests from specific countries
func (gbm *GeoBlockMiddleware) GeoBlock() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get country from header (set by reverse proxy or geo service)
		country := c.GetHeader("X-Country-Code")
		if country == "" {
			// If no country header, allow request (assume geo service not configured)
			c.Next()
			return
		}

		country = strings.ToUpper(country)

		// Check if country is blocked
		if gbm.blockedCountries[country] {
			c.JSON(http.StatusForbidden, gin.H{
				"error":   "country_blocked",
				"message": "Access denied from this country",
			})
			c.Abort()
			return
		}

		c.Next()
	}
}

// RateLimitStatusMiddleware provides rate limit status in response headers
type RateLimitStatusMiddleware struct {
	rateLimiter *ratelimit.RateLimiter
}

// NewRateLimitStatusMiddleware creates a new rate limit status middleware
func NewRateLimitStatusMiddleware(rateLimiter *ratelimit.RateLimiter) *RateLimitStatusMiddleware {
	return &RateLimitStatusMiddleware{
		rateLimiter: rateLimiter,
	}
}

// AddRateLimitHeaders returns middleware that adds rate limit headers to responses
func (rlsm *RateLimitStatusMiddleware) AddRateLimitHeaders() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get rate limit key (user ID or IP)
		var key string
		if userID, exists := c.Get("user_id"); exists {
			if id, ok := userID.(uuid.UUID); ok {
				key = id.String()
			}
		} else {
			key = "ip:" + c.ClientIP()
		}

		// Note: In a real implementation, you'd need to expose rate limiter statistics
		// For now, we'll just add basic headers
		c.Header("X-RateLimit-Key", key)
		c.Header("X-RateLimit-Remaining", "unlimited") // Placeholder
		c.Header("X-RateLimit-Reset", "3600")         // Placeholder: 1 hour

		c.Next()
	}
}
