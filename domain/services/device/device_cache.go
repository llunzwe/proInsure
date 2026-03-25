package device

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"
	"time"

	"github.com/google/uuid"

	"smartsure/internal/domain/models"
)

// CacheKey types for different cache entries
type CacheKeyType string

const (
	CacheKeyDevice         CacheKeyType = "device"
	CacheKeyRiskScore      CacheKeyType = "risk_score"
	CacheKeyPremium        CacheKeyType = "premium"
	CacheKeyTradeInValue   CacheKeyType = "trade_in"
	CacheKeyInsurableValue CacheKeyType = "insurable_value"
	CacheKeyEligibility    CacheKeyType = "eligibility"
	CacheKeyBehaviorScore  CacheKeyType = "behavior_score"
	CacheKeyFailureRisk    CacheKeyType = "failure_risk"
)

// CacheEntry represents a cached value with expiration
type CacheEntry struct {
	Value     interface{}
	ExpiresAt time.Time
	Version   int
}

// IsExpired checks if cache entry is expired
func (e *CacheEntry) IsExpired() bool {
	return time.Now().After(e.ExpiresAt)
}

// DeviceCache provides caching for frequently accessed device data
type DeviceCache struct {
	store      map[string]*CacheEntry
	mu         sync.RWMutex
	defaultTTL time.Duration
	maxEntries int
}

// NewDeviceCache creates a new device cache
func NewDeviceCache(defaultTTL time.Duration, maxEntries int) *DeviceCache {
	if maxEntries <= 0 {
		maxEntries = 10000 // Default max entries
	}
	if defaultTTL <= 0 {
		defaultTTL = 5 * time.Minute // Default TTL
	}

	cache := &DeviceCache{
		store:      make(map[string]*CacheEntry),
		defaultTTL: defaultTTL,
		maxEntries: maxEntries,
	}

	// Start cleanup goroutine
	go cache.cleanupExpired()

	return cache
}

// BuildKey builds a cache key from components
func (c *DeviceCache) BuildKey(keyType CacheKeyType, deviceID uuid.UUID, suffix ...string) string {
	key := fmt.Sprintf("%s:%s", keyType, deviceID.String())
	for _, s := range suffix {
		key += ":" + s
	}
	return key
}

// Get retrieves a value from cache
func (c *DeviceCache) Get(key string) (interface{}, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	entry, exists := c.store[key]
	if !exists {
		return nil, false
	}

	if entry.IsExpired() {
		return nil, false
	}

	return entry.Value, true
}

// GetDevice retrieves a cached device
func (c *DeviceCache) GetDevice(deviceID uuid.UUID) (*models.Device, bool) {
	key := c.BuildKey(CacheKeyDevice, deviceID)
	value, ok := c.Get(key)
	if !ok {
		return nil, false
	}

	device, ok := value.(*models.Device)
	return device, ok
}

// GetRiskScore retrieves cached risk score
func (c *DeviceCache) GetRiskScore(deviceID uuid.UUID) (float64, bool) {
	key := c.BuildKey(CacheKeyRiskScore, deviceID)
	value, ok := c.Get(key)
	if !ok {
		return 0, false
	}

	score, ok := value.(float64)
	return score, ok
}

// GetPremium retrieves cached premium
func (c *DeviceCache) GetPremium(deviceID uuid.UUID) (float64, bool) {
	key := c.BuildKey(CacheKeyPremium, deviceID)
	value, ok := c.Get(key)
	if !ok {
		return 0, false
	}

	premium, ok := value.(float64)
	return premium, ok
}

// GetTradeInValue retrieves cached trade-in value
func (c *DeviceCache) GetTradeInValue(deviceID uuid.UUID) (float64, bool) {
	key := c.BuildKey(CacheKeyTradeInValue, deviceID)
	value, ok := c.Get(key)
	if !ok {
		return 0, false
	}

	tradeIn, ok := value.(float64)
	return tradeIn, ok
}

// Set stores a value in cache with default TTL
func (c *DeviceCache) Set(key string, value interface{}) {
	c.SetWithTTL(key, value, c.defaultTTL)
}

// SetWithTTL stores a value in cache with custom TTL
func (c *DeviceCache) SetWithTTL(key string, value interface{}, ttl time.Duration) {
	c.mu.Lock()
	defer c.mu.Unlock()

	// Check max entries limit
	if len(c.store) >= c.maxEntries {
		c.evictOldest()
	}

	version := 1
	if existing, ok := c.store[key]; ok {
		version = existing.Version + 1
	}

	c.store[key] = &CacheEntry{
		Value:     value,
		ExpiresAt: time.Now().Add(ttl),
		Version:   version,
	}
}

// SetDevice caches a device
func (c *DeviceCache) SetDevice(device *models.Device) {
	key := c.BuildKey(CacheKeyDevice, device.ID)
	c.SetWithTTL(key, device, c.defaultTTL*2) // Devices cached longer
}

// SetRiskScore caches risk score
func (c *DeviceCache) SetRiskScore(deviceID uuid.UUID, score float64) {
	key := c.BuildKey(CacheKeyRiskScore, deviceID)
	c.SetWithTTL(key, score, c.defaultTTL)
}

// SetPremium caches premium
func (c *DeviceCache) SetPremium(deviceID uuid.UUID, premium float64) {
	key := c.BuildKey(CacheKeyPremium, deviceID)
	c.SetWithTTL(key, premium, 30*time.Minute) // Premium cached for 30 minutes
}

// SetTradeInValue caches trade-in value
func (c *DeviceCache) SetTradeInValue(deviceID uuid.UUID, value float64) {
	key := c.BuildKey(CacheKeyTradeInValue, deviceID)
	c.SetWithTTL(key, value, 24*time.Hour) // Trade-in values cached for 24 hours
}

// Delete removes a value from cache
func (c *DeviceCache) Delete(key string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	delete(c.store, key)
}

// InvalidateDevice removes all cached data for a device
func (c *DeviceCache) InvalidateDevice(deviceID uuid.UUID) {
	c.mu.Lock()
	defer c.mu.Unlock()

	// Build all possible keys for this device
	keysToDelete := []string{
		c.BuildKey(CacheKeyDevice, deviceID),
		c.BuildKey(CacheKeyRiskScore, deviceID),
		c.BuildKey(CacheKeyPremium, deviceID),
		c.BuildKey(CacheKeyTradeInValue, deviceID),
		c.BuildKey(CacheKeyInsurableValue, deviceID),
		c.BuildKey(CacheKeyEligibility, deviceID),
		c.BuildKey(CacheKeyBehaviorScore, deviceID),
		c.BuildKey(CacheKeyFailureRisk, deviceID),
	}

	for _, key := range keysToDelete {
		delete(c.store, key)
	}
}

// Clear removes all entries from cache
func (c *DeviceCache) Clear() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.store = make(map[string]*CacheEntry)
}

// Size returns the number of entries in cache
func (c *DeviceCache) Size() int {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return len(c.store)
}

// GetStats returns cache statistics
func (c *DeviceCache) GetStats() map[string]interface{} {
	c.mu.RLock()
	defer c.mu.RUnlock()

	expired := 0
	for _, entry := range c.store {
		if entry.IsExpired() {
			expired++
		}
	}

	return map[string]interface{}{
		"total_entries":   len(c.store),
		"expired_entries": expired,
		"active_entries":  len(c.store) - expired,
		"max_entries":     c.maxEntries,
		"default_ttl":     c.defaultTTL.String(),
	}
}

// evictOldest removes the oldest entry from cache
func (c *DeviceCache) evictOldest() {
	var oldestKey string
	var oldestTime time.Time

	for key, entry := range c.store {
		if oldestKey == "" || entry.ExpiresAt.Before(oldestTime) {
			oldestKey = key
			oldestTime = entry.ExpiresAt
		}
	}

	if oldestKey != "" {
		delete(c.store, oldestKey)
	}
}

// cleanupExpired periodically removes expired entries
func (c *DeviceCache) cleanupExpired() {
	ticker := time.NewTicker(5 * time.Minute)
	defer ticker.Stop()

	for range ticker.C {
		c.mu.Lock()
		keysToDelete := []string{}

		for key, entry := range c.store {
			if entry.IsExpired() {
				keysToDelete = append(keysToDelete, key)
			}
		}

		for _, key := range keysToDelete {
			delete(c.store, key)
		}
		c.mu.Unlock()
	}
}

// === Computed Value Cache ===

// ComputedValueCache caches expensive computations
type ComputedValueCache struct {
	cache *DeviceCache
}

// NewComputedValueCache creates a new computed value cache
func NewComputedValueCache() *ComputedValueCache {
	return &ComputedValueCache{
		cache: NewDeviceCache(15*time.Minute, 5000),
	}
}

// GetOrCompute retrieves from cache or computes if missing
func (cvc *ComputedValueCache) GetOrCompute(
	ctx context.Context,
	deviceID uuid.UUID,
	keyType CacheKeyType,
	compute func() (interface{}, error),
	ttl time.Duration,
) (interface{}, error) {
	key := cvc.cache.BuildKey(keyType, deviceID)

	// Try to get from cache
	if value, ok := cvc.cache.Get(key); ok {
		return value, nil
	}

	// Compute the value
	value, err := compute()
	if err != nil {
		return nil, err
	}

	// Cache the result
	cvc.cache.SetWithTTL(key, value, ttl)

	return value, nil
}

// GetRiskScoreOrCompute gets risk score from cache or computes
func (cvc *ComputedValueCache) GetRiskScoreOrCompute(
	ctx context.Context,
	deviceID uuid.UUID,
	compute func() (float64, error),
) (float64, error) {
	value, err := cvc.GetOrCompute(ctx, deviceID, CacheKeyRiskScore, func() (interface{}, error) {
		return compute()
	}, 15*time.Minute)

	if err != nil {
		return 0, err
	}

	score, ok := value.(float64)
	if !ok {
		return 0, fmt.Errorf("invalid risk score type in cache")
	}

	return score, nil
}

// GetPremiumOrCompute gets premium from cache or computes
func (cvc *ComputedValueCache) GetPremiumOrCompute(
	ctx context.Context,
	deviceID uuid.UUID,
	compute func() (float64, error),
) (float64, error) {
	value, err := cvc.GetOrCompute(ctx, deviceID, CacheKeyPremium, func() (interface{}, error) {
		return compute()
	}, 30*time.Minute)

	if err != nil {
		return 0, err
	}

	premium, ok := value.(float64)
	if !ok {
		return 0, fmt.Errorf("invalid premium type in cache")
	}

	return premium, nil
}

// InvalidateDevice invalidates all cached computations for a device
func (cvc *ComputedValueCache) InvalidateDevice(deviceID uuid.UUID) {
	cvc.cache.InvalidateDevice(deviceID)
}

// === Relationship Cache ===

// RelationshipCache caches device relationships to avoid N+1 queries
type RelationshipCache struct {
	cache *DeviceCache
}

// NewRelationshipCache creates a new relationship cache
func NewRelationshipCache() *RelationshipCache {
	return &RelationshipCache{
		cache: NewDeviceCache(10*time.Minute, 10000),
	}
}

// GetPolicies gets cached policies for device
func (rc *RelationshipCache) GetPolicies(deviceID uuid.UUID) ([]uuid.UUID, bool) {
	key := rc.cache.BuildKey(CacheKeyDevice, deviceID, "policies")
	value, ok := rc.cache.Get(key)
	if !ok {
		return nil, false
	}

	policies, ok := value.([]uuid.UUID)
	return policies, ok
}

// SetPolicies caches policy IDs for device
func (rc *RelationshipCache) SetPolicies(deviceID uuid.UUID, policyIDs []uuid.UUID) {
	key := rc.cache.BuildKey(CacheKeyDevice, deviceID, "policies")
	rc.cache.SetWithTTL(key, policyIDs, 10*time.Minute)
}

// GetClaims gets cached claims for device
func (rc *RelationshipCache) GetClaims(deviceID uuid.UUID) ([]uuid.UUID, bool) {
	key := rc.cache.BuildKey(CacheKeyDevice, deviceID, "claims")
	value, ok := rc.cache.Get(key)
	if !ok {
		return nil, false
	}

	claims, ok := value.([]uuid.UUID)
	return claims, ok
}

// SetClaims caches claim IDs for device
func (rc *RelationshipCache) SetClaims(deviceID uuid.UUID, claimIDs []uuid.UUID) {
	key := rc.cache.BuildKey(CacheKeyDevice, deviceID, "claims")
	rc.cache.SetWithTTL(key, claimIDs, 10*time.Minute)
}

// === Batch Cache Operations ===

// BatchCache provides batch caching operations
type BatchCache struct {
	cache *DeviceCache
}

// NewBatchCache creates a new batch cache
func NewBatchCache() *BatchCache {
	return &BatchCache{
		cache: NewDeviceCache(5*time.Minute, 20000),
	}
}

// GetMultiple retrieves multiple values from cache
func (bc *BatchCache) GetMultiple(keys []string) map[string]interface{} {
	results := make(map[string]interface{})

	for _, key := range keys {
		if value, ok := bc.cache.Get(key); ok {
			results[key] = value
		}
	}

	return results
}

// SetMultiple stores multiple values in cache
func (bc *BatchCache) SetMultiple(values map[string]interface{}, ttl time.Duration) {
	for key, value := range values {
		bc.cache.SetWithTTL(key, value, ttl)
	}
}

// GetDevices retrieves multiple devices from cache
func (bc *BatchCache) GetDevices(deviceIDs []uuid.UUID) map[uuid.UUID]*models.Device {
	results := make(map[uuid.UUID]*models.Device)

	for _, id := range deviceIDs {
		key := bc.cache.BuildKey(CacheKeyDevice, id)
		if value, ok := bc.cache.Get(key); ok {
			if device, ok := value.(*models.Device); ok {
				results[id] = device
			}
		}
	}

	return results
}

// SetDevices caches multiple devices
func (bc *BatchCache) SetDevices(devices []*models.Device) {
	for _, device := range devices {
		key := bc.cache.BuildKey(CacheKeyDevice, device.ID)
		bc.cache.SetWithTTL(key, device, 10*time.Minute)
	}
}

// === JSON Cache for API Responses ===

// JSONCache caches JSON responses
type JSONCache struct {
	cache *DeviceCache
}

// NewJSONCache creates a new JSON cache
func NewJSONCache() *JSONCache {
	return &JSONCache{
		cache: NewDeviceCache(2*time.Minute, 5000),
	}
}

// GetJSON retrieves cached JSON
func (jc *JSONCache) GetJSON(endpoint string, params map[string]string) ([]byte, bool) {
	key := jc.buildJSONKey(endpoint, params)
	value, ok := jc.cache.Get(key)
	if !ok {
		return nil, false
	}

	data, ok := value.([]byte)
	return data, ok
}

// SetJSON caches JSON response
func (jc *JSONCache) SetJSON(endpoint string, params map[string]string, data []byte) {
	key := jc.buildJSONKey(endpoint, params)
	jc.cache.SetWithTTL(key, data, 2*time.Minute)
}

func (jc *JSONCache) buildJSONKey(endpoint string, params map[string]string) string {
	paramJSON, _ := json.Marshal(params)
	return fmt.Sprintf("json:%s:%s", endpoint, string(paramJSON))
}
