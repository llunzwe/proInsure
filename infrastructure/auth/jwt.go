package auth

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

var (
	ErrInvalidToken        = errors.New("invalid token")
	ErrExpiredToken        = errors.New("token has expired")
	ErrMissingToken        = errors.New("authorization token required")
	ErrInvalidRefreshToken = errors.New("invalid refresh token")
	ErrRefreshTokenExpired = errors.New("refresh token has expired")
	ErrTokenRevoked        = errors.New("token has been revoked")
)

// Claims represents the JWT claims
type Claims struct {
	UserID    uuid.UUID `json:"user_id"`
	Email     string    `json:"email"`
	UserType  string    `json:"user_type"`
	UserRole  string    `json:"user_role"`
	TokenType string    `json:"token_type"` // "access" or "refresh"
	jwt.RegisteredClaims
}

// TokenPair represents access and refresh token pair
type TokenPair struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresIn    int64  `json:"expires_in"` // seconds
	TokenType    string `json:"token_type"`
}

// RefreshTokenData represents stored refresh token data
type RefreshTokenData struct {
	UserID    uuid.UUID `json:"user_id"`
	TokenID   string    `json:"token_id"`
	CreatedAt time.Time `json:"created_at"`
	ExpiresAt time.Time `json:"expires_at"`
	IsRevoked bool      `json:"is_revoked"`
}

// JWTService handles JWT operations with refresh token support
type JWTService struct {
	secretKey           []byte
	accessTokenDuration time.Duration
	refreshTokenDuration time.Duration
	tokenStore          TokenStore // Interface for storing refresh tokens
}

// TokenStore interface for refresh token storage
type TokenStore interface {
	StoreRefreshToken(tokenData *RefreshTokenData) error
	GetRefreshToken(tokenID string) (*RefreshTokenData, error)
	RevokeRefreshToken(tokenID string) error
	RevokeUserTokens(userID uuid.UUID) error
	CleanupExpiredTokens() error
}

// NewJWTService creates a new JWT service with refresh token support
func NewJWTService(secretKey string, accessTokenDuration, refreshTokenDuration time.Duration, tokenStore TokenStore) *JWTService {
	return &JWTService{
		secretKey:             []byte(secretKey),
		accessTokenDuration:   accessTokenDuration,
		refreshTokenDuration:  refreshTokenDuration,
		tokenStore:           tokenStore,
	}
}

// NewJWTServiceLegacy creates a JWT service without refresh tokens (backward compatibility)
func NewJWTServiceLegacy(secretKey string, tokenDuration time.Duration) *JWTService {
	return &JWTService{
		secretKey:           []byte(secretKey),
		accessTokenDuration: tokenDuration,
		tokenStore:         nil, // No token store for legacy mode
	}
}

// GenerateTokenPair generates both access and refresh tokens for a user
func (j *JWTService) GenerateTokenPair(userID uuid.UUID, email, userType, userRole string) (*TokenPair, error) {
	// Generate unique token ID for refresh token
	tokenID, err := j.generateTokenID()
	if err != nil {
		return nil, err
	}

	// Create access token
	accessToken, err := j.generateAccessToken(userID, email, userType, userRole)
	if err != nil {
		return nil, err
	}

	// Create refresh token
	refreshToken, err := j.generateRefreshToken(userID, tokenID)
	if err != nil {
		return nil, err
	}

	// Store refresh token data if token store is available
	if j.tokenStore != nil {
		tokenData := &RefreshTokenData{
			UserID:    userID,
			TokenID:   tokenID,
			CreatedAt: time.Now(),
			ExpiresAt: time.Now().Add(j.refreshTokenDuration),
			IsRevoked: false,
		}

		if err := j.tokenStore.StoreRefreshToken(tokenData); err != nil {
			return nil, err
		}
	}

	return &TokenPair{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresIn:    int64(j.accessTokenDuration.Seconds()),
		TokenType:    "Bearer",
	}, nil
}

// generateAccessToken creates an access token
func (j *JWTService) generateAccessToken(userID uuid.UUID, email, userType, userRole string) (string, error) {
	claims := &Claims{
		UserID:    userID,
		Email:     email,
		UserType:  userType,
		UserRole:  userRole,
		TokenType: "access",
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(j.accessTokenDuration)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			Issuer:    "smartsure-api",
			Subject:   userID.String(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(j.secretKey)
}

// generateRefreshToken creates a refresh token
func (j *JWTService) generateRefreshToken(userID uuid.UUID, tokenID string) (string, error) {
	claims := &Claims{
		UserID:    userID,
		TokenType: "refresh",
		RegisteredClaims: jwt.RegisteredClaims{
			ID:        tokenID,
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(j.refreshTokenDuration)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			Issuer:    "smartsure-api",
			Subject:   userID.String(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(j.secretKey)
}

// generateTokenID generates a unique token identifier
func (j *JWTService) generateTokenID() (string, error) {
	bytes := make([]byte, 32)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}

// GenerateToken generates a new JWT token for a user (legacy method for backward compatibility)
func (j *JWTService) GenerateToken(userID uuid.UUID, email, userType, userRole string) (string, error) {
	return j.generateAccessToken(userID, email, userType, userRole)
}

// ValidateToken validates a JWT token and returns the claims
func (j *JWTService) ValidateToken(tokenString string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, ErrInvalidToken
		}
		return j.secretKey, nil
	})

	if err != nil {
		if errors.Is(err, jwt.ErrTokenExpired) {
			return nil, ErrExpiredToken
		}
		return nil, ErrInvalidToken
	}

	claims, ok := token.Claims.(*Claims)
	if !ok || !token.Valid {
		return nil, ErrInvalidToken
	}

	return claims, nil
}

// RefreshTokenPair refreshes both access and refresh tokens using a valid refresh token
func (j *JWTService) RefreshTokenPair(refreshTokenString string) (*TokenPair, error) {
	// Validate refresh token
	claims, err := j.ValidateRefreshToken(refreshTokenString)
	if err != nil {
		return nil, err
	}

	// Check if refresh token is stored and valid
	if j.tokenStore != nil {
		tokenData, err := j.tokenStore.GetRefreshToken(claims.ID)
		if err != nil {
			return nil, ErrInvalidRefreshToken
		}

		if tokenData.IsRevoked || time.Now().After(tokenData.ExpiresAt) {
			return nil, ErrRefreshTokenExpired
		}
	}

	// Revoke old refresh token for security (token rotation)
	if j.tokenStore != nil {
		j.tokenStore.RevokeRefreshToken(claims.ID)
	}

	// Generate new token pair
	return j.GenerateTokenPair(claims.UserID, claims.Email, claims.UserType, claims.UserRole)
}

// ValidateRefreshToken validates a refresh token specifically
func (j *JWTService) ValidateRefreshToken(tokenString string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, ErrInvalidToken
		}
		return j.secretKey, nil
	})

	if err != nil {
		if errors.Is(err, jwt.ErrTokenExpired) {
			return nil, ErrRefreshTokenExpired
		}
		return nil, ErrInvalidRefreshToken
	}

	claims, ok := token.Claims.(*Claims)
	if !ok || !token.Valid {
		return nil, ErrInvalidRefreshToken
	}

	// Ensure this is actually a refresh token
	if claims.TokenType != "refresh" {
		return nil, ErrInvalidRefreshToken
	}

	return claims, nil
}

// RevokeUserTokens revokes all refresh tokens for a user
func (j *JWTService) RevokeUserTokens(userID uuid.UUID) error {
	if j.tokenStore == nil {
		return nil // No-op if no token store
	}
	return j.tokenStore.RevokeUserTokens(userID)
}

// RevokeRefreshToken revokes a specific refresh token
func (j *JWTService) RevokeRefreshToken(tokenID string) error {
	if j.tokenStore == nil {
		return nil // No-op if no token store
	}
	return j.tokenStore.RevokeRefreshToken(tokenID)
}

// RefreshToken generates a new token from an existing valid token (legacy method)
func (j *JWTService) RefreshToken(tokenString string) (string, error) {
	claims, err := j.ValidateToken(tokenString)
	if err != nil {
		return "", err
	}

	// Generate new token with same claims but updated timestamps
	return j.GenerateToken(claims.UserID, claims.Email, claims.UserType, claims.UserRole)
}

// AuthMiddleware returns a Gin middleware for JWT authentication
func (j *JWTService) AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error":   "unauthorized",
				"message": "authorization token required",
			})
			c.Abort()
			return
		}

		// Extract token from "Bearer <token>" format
		tokenParts := strings.Split(authHeader, " ")
		if len(tokenParts) != 2 || tokenParts[0] != "Bearer" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error":   "unauthorized",
				"message": "invalid authorization header format",
			})
			c.Abort()
			return
		}

		claims, err := j.ValidateToken(tokenParts[1])
		if err != nil {
			var message string
			switch err {
			case ErrExpiredToken:
				message = "token has expired"
			case ErrInvalidToken:
				message = "invalid token"
			default:
				message = "token validation failed"
			}

			c.JSON(http.StatusUnauthorized, gin.H{
				"error":   "unauthorized",
				"message": message,
			})
			c.Abort()
			return
		}

		// Set user information in context
		c.Set("user_id", claims.UserID)
		c.Set("user_email", claims.Email)
		c.Set("user_type", claims.UserType)
		c.Set("user_role", claims.UserRole)
		c.Set("claims", claims)

		c.Next()
	}
}

// OptionalAuthMiddleware returns a Gin middleware for optional JWT authentication
func (j *JWTService) OptionalAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.Next()
			return
		}

		tokenParts := strings.Split(authHeader, " ")
		if len(tokenParts) != 2 || tokenParts[0] != "Bearer" {
			c.Next()
			return
		}

		claims, err := j.ValidateToken(tokenParts[1])
		if err != nil {
			c.Next()
			return
		}

		// Set user information in context if token is valid
		c.Set("user_id", claims.UserID)
		c.Set("user_email", claims.Email)
		c.Set("user_type", claims.UserType)
		c.Set("user_role", claims.UserRole)
		c.Set("claims", claims)

		c.Next()
	}
}

// RequireRole returns a middleware that requires specific user roles
func RequireRole(allowedRoles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		userRole, exists := c.Get("user_role")
		if !exists {
			c.JSON(http.StatusForbidden, gin.H{
				"error":   "forbidden",
				"message": "insufficient permissions",
			})
			c.Abort()
			return
		}

		role := userRole.(string)
		for _, allowedRole := range allowedRoles {
			if role == allowedRole {
				c.Next()
				return
			}
		}

		c.JSON(http.StatusForbidden, gin.H{
			"error":   "forbidden",
			"message": "insufficient role permissions",
		})
		c.Abort()
	}
}

// RequireUserType returns a middleware that requires specific user types
func RequireUserType(allowedTypes ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		userType, exists := c.Get("user_type")
		if !exists {
			c.JSON(http.StatusForbidden, gin.H{
				"error":   "forbidden",
				"message": "insufficient permissions",
			})
			c.Abort()
			return
		}

		uType := userType.(string)
		for _, allowedType := range allowedTypes {
			if uType == allowedType {
				c.Next()
				return
			}
		}

		c.JSON(http.StatusForbidden, gin.H{
			"error":   "forbidden",
			"message": "insufficient type permissions",
		})
		c.Abort()
	}
}
