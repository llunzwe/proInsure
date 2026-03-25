package auth

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/google/uuid"
	"golang.org/x/net/context"
)

// RedisTokenStore implements TokenStore interface using Redis
type RedisTokenStore struct {
	client redis.Cmdable
	prefix string
	ctx    context.Context
}

// NewRedisTokenStore creates a new Redis-based token store
func NewRedisTokenStore(client redis.Cmdable, prefix string) *RedisTokenStore {
	return &RedisTokenStore{
		client: client,
		prefix: prefix,
		ctx:    context.Background(),
	}
}

// StoreRefreshToken stores a refresh token in Redis
func (r *RedisTokenStore) StoreRefreshToken(tokenData *RefreshTokenData) error {
	key := r.tokenKey(tokenData.TokenID)

	// Serialize token data to JSON
	data, err := json.Marshal(tokenData)
	if err != nil {
		return fmt.Errorf("failed to marshal token data: %w", err)
	}

	// Store with expiration
	expiration := time.Until(tokenData.ExpiresAt)
	if expiration <= 0 {
		return fmt.Errorf("token already expired")
	}

	return r.client.Set(r.ctx, key, data, expiration).Err()
}

// GetRefreshToken retrieves a refresh token from Redis
func (r *RedisTokenStore) GetRefreshToken(tokenID string) (*RefreshTokenData, error) {
	key := r.tokenKey(tokenID)

	data, err := r.client.Get(r.ctx, key).Result()
	if err != nil {
		if err == redis.Nil {
			return nil, fmt.Errorf("refresh token not found")
		}
		return nil, fmt.Errorf("failed to get refresh token: %w", err)
	}

	var tokenData RefreshTokenData
	if err := json.Unmarshal([]byte(data), &tokenData); err != nil {
		return nil, fmt.Errorf("failed to unmarshal token data: %w", err)
	}

	return &tokenData, nil
}

// RevokeRefreshToken revokes a specific refresh token
func (r *RedisTokenStore) RevokeRefreshToken(tokenID string) error {
	key := r.tokenKey(tokenID)
	return r.client.Del(r.ctx, key).Err()
}

// RevokeUserTokens revokes all refresh tokens for a user
func (r *RedisTokenStore) RevokeUserTokens(userID uuid.UUID) error {
	// Find all tokens for this user
	pattern := r.userTokenPattern(userID)
	keys, err := r.client.Keys(r.ctx, pattern).Result()
	if err != nil {
		return fmt.Errorf("failed to find user tokens: %w", err)
	}

	if len(keys) == 0 {
		return nil
	}

	// Delete all found tokens
	return r.client.Del(r.ctx, keys...).Err()
}

// CleanupExpiredTokens removes expired tokens (Redis handles this automatically with TTL)
func (r *RedisTokenStore) CleanupExpiredTokens() error {
	// Redis automatically removes expired keys, so this is a no-op
	// We could implement additional cleanup logic if needed
	return nil
}

// tokenKey generates the Redis key for a token
func (r *RedisTokenStore) tokenKey(tokenID string) string {
	return fmt.Sprintf("%s:token:%s", r.prefix, tokenID)
}

// userTokenPattern generates the pattern to find all tokens for a user
func (r *RedisTokenStore) userTokenPattern(userID uuid.UUID) string {
	return fmt.Sprintf("%s:token:*", r.prefix)
}

// InMemoryTokenStore implements TokenStore interface using in-memory storage
// This is useful for development/testing or when Redis is not available
type InMemoryTokenStore struct {
	tokens map[string]*RefreshTokenData
}

// NewInMemoryTokenStore creates a new in-memory token store
func NewInMemoryTokenStore() *InMemoryTokenStore {
	return &InMemoryTokenStore{
		tokens: make(map[string]*RefreshTokenData),
	}
}

// StoreRefreshToken stores a refresh token in memory
func (m *InMemoryTokenStore) StoreRefreshToken(tokenData *RefreshTokenData) error {
	m.tokens[tokenData.TokenID] = tokenData
	return nil
}

// GetRefreshToken retrieves a refresh token from memory
func (m *InMemoryTokenStore) GetRefreshToken(tokenID string) (*RefreshTokenData, error) {
	tokenData, exists := m.tokens[tokenID]
	if !exists {
		return nil, fmt.Errorf("refresh token not found")
	}

	// Check if token is expired
	if time.Now().After(tokenData.ExpiresAt) {
		delete(m.tokens, tokenID) // Clean up expired token
		return nil, fmt.Errorf("refresh token expired")
	}

	return tokenData, nil
}

// RevokeRefreshToken revokes a specific refresh token
func (m *InMemoryTokenStore) RevokeRefreshToken(tokenID string) error {
	delete(m.tokens, tokenID)
	return nil
}

// RevokeUserTokens revokes all refresh tokens for a user
func (m *InMemoryTokenStore) RevokeUserTokens(userID uuid.UUID) error {
	for tokenID, tokenData := range m.tokens {
		if tokenData.UserID == userID {
			delete(m.tokens, tokenID)
		}
	}
	return nil
}

// CleanupExpiredTokens removes expired tokens from memory
func (m *InMemoryTokenStore) CleanupExpiredTokens() error {
	now := time.Now()
	for tokenID, tokenData := range m.tokens {
		if now.After(tokenData.ExpiresAt) {
			delete(m.tokens, tokenID)
		}
	}
	return nil
}
