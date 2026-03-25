package kafka

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/Shopify/sarama"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

// Producer represents a Kafka producer
type Producer struct {
	Producer sarama.SyncProducer
	Logger   *logrus.Logger
	Config   *Config
}

// Config represents Kafka producer configuration
type Config struct {
	Brokers         []string
	TopicPrefix     string
	MaxMessageBytes int
	RequiredAcks    sarama.RequiredAcks
	Timeout         time.Duration
	Compression     sarama.CompressionCodec
}

// Event represents a domain event
type Event struct {
	ID            string                 `json:"id"`
	Type          string                 `json:"type"`
	Source        string                 `json:"source"`
	Data          map[string]interface{} `json:"data"`
	Timestamp     time.Time              `json:"timestamp"`
	Version       string                 `json:"version"`
	CorrelationID string                 `json:"correlation_id,omitempty"`
}

// NewProducer creates a new Kafka producer
func NewProducer() (*Producer, error) {
	config := loadConfig()

	// Create Sarama configuration
	saramaConfig := sarama.NewConfig()
	saramaConfig.Producer.RequiredAcks = config.RequiredAcks
	saramaConfig.Producer.Timeout = config.Timeout
	saramaConfig.Producer.Compression = config.Compression
	saramaConfig.Producer.MaxMessageBytes = config.MaxMessageBytes
	saramaConfig.Producer.Return.Successes = true
	saramaConfig.Producer.Return.Errors = true

	// Create sync producer
	producer, err := sarama.NewSyncProducer(config.Brokers, saramaConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to create Kafka producer: %w", err)
	}

	logger := logrus.New()
	logger.Info("Successfully initialized Kafka producer")

	return &Producer{
		Producer: producer,
		Logger:   logger,
		Config:   config,
	}, nil
}

// Close closes the Kafka producer
func (c *Producer) Close() error {
	return c.Producer.Close()
}

// PublishEvent publishes an event to Kafka
func (c *Producer) PublishEvent(ctx context.Context, event *Event) error {
	// Marshal event to JSON
	eventBytes, err := json.Marshal(event)
	if err != nil {
		return fmt.Errorf("failed to marshal event: %w", err)
	}

	// Create message
	message := &sarama.ProducerMessage{
		Topic:     c.Config.TopicPrefix + "." + event.Type,
		Value:     sarama.ByteEncoder(eventBytes),
		Timestamp: event.Timestamp,
		Headers: []sarama.RecordHeader{
			{Key: []byte("event_id"), Value: []byte(event.ID)},
			{Key: []byte("event_type"), Value: []byte(event.Type)},
			{Key: []byte("source"), Value: []byte(event.Source)},
			{Key: []byte("version"), Value: []byte(event.Version)},
		},
	}

	// Add correlation ID if present
	if event.CorrelationID != "" {
		message.Headers = append(message.Headers, sarama.RecordHeader{
			Key: []byte("correlation_id"), Value: []byte(event.CorrelationID),
		})
	}

	// Send message
	partition, offset, err := c.Producer.SendMessage(message)
	if err != nil {
		return fmt.Errorf("failed to send message: %w", err)
	}

	c.Logger.WithFields(logrus.Fields{
		"event_id":   event.ID,
		"event_type": event.Type,
		"topic":      message.Topic,
		"partition":  partition,
		"offset":     offset,
	}).Info("Event published successfully")

	return nil
}

// PublishUserEvent publishes a user-related event
func (c *Producer) PublishUserEvent(ctx context.Context, eventType string, userID string, data map[string]interface{}) error {
	event := &Event{
		ID:        generateEventID(),
		Type:      eventType,
		Source:    "smartsure-api",
		Data:      data,
		Timestamp: time.Now().UTC(),
		Version:   "1.0",
	}

	return c.PublishEvent(ctx, event)
}

// PublishPolicyEvent publishes a policy-related event
func (c *Producer) PublishPolicyEvent(ctx context.Context, eventType string, policyID string, data map[string]interface{}) error {
	event := &Event{
		ID:        generateEventID(),
		Type:      eventType,
		Source:    "smartsure-api",
		Data:      data,
		Timestamp: time.Now().UTC(),
		Version:   "1.0",
	}

	return c.PublishEvent(ctx, event)
}

// PublishClaimEvent publishes a claim-related event
func (c *Producer) PublishClaimEvent(ctx context.Context, eventType string, claimID string, data map[string]interface{}) error {
	event := &Event{
		ID:        generateEventID(),
		Type:      eventType,
		Source:    "smartsure-api",
		Data:      data,
		Timestamp: time.Now().UTC(),
		Version:   "1.0",
	}

	return c.PublishEvent(ctx, event)
}

// PublishDeviceEvent publishes a device-related event
func (c *Producer) PublishDeviceEvent(ctx context.Context, eventType string, deviceID string, data map[string]interface{}) error {
	event := &Event{
		ID:        generateEventID(),
		Type:      eventType,
		Source:    "smartsure-api",
		Data:      data,
		Timestamp: time.Now().UTC(),
		Version:   "1.0",
	}

	return c.PublishEvent(ctx, event)
}

// PublishIoTEvent publishes an IoT-related event
func (c *Producer) PublishIoTEvent(ctx context.Context, eventType string, deviceID string, data map[string]interface{}) error {
	event := &Event{
		ID:        generateEventID(),
		Type:      eventType,
		Source:    "smartsure-iot",
		Data:      data,
		Timestamp: time.Now().UTC(),
		Version:   "1.0",
	}

	return c.PublishEvent(ctx, event)
}

// PublishFraudDetectionEvent publishes a fraud detection event
func (c *Producer) PublishFraudDetectionEvent(ctx context.Context, eventType string, entityID string, data map[string]interface{}) error {
	event := &Event{
		ID:        generateEventID(),
		Type:      eventType,
		Source:    "smartsure-fraud-detection",
		Data:      data,
		Timestamp: time.Now().UTC(),
		Version:   "1.0",
	}

	return c.PublishEvent(ctx, event)
}

// HealthCheck performs a health check on the Kafka producer
func (c *Producer) HealthCheck(ctx context.Context) error {
	// Try to send a test message to check connectivity
	testEvent := &Event{
		ID:        "health-check",
		Type:      "health_check",
		Source:    "smartsure-api",
		Data:      map[string]interface{}{"status": "healthy"},
		Timestamp: time.Now().UTC(),
		Version:   "1.0",
	}

	return c.PublishEvent(ctx, testEvent)
}

// loadConfig loads Kafka producer configuration from environment
func loadConfig() *Config {
	return &Config{
		Brokers:         viper.GetStringSlice("kafka.brokers"),
		TopicPrefix:     viper.GetString("kafka.topic_prefix"),
		MaxMessageBytes: viper.GetInt("kafka.max_message_bytes"),
		RequiredAcks:    getRequiredAcks(viper.GetString("kafka.required_acks")),
		Timeout:         viper.GetDuration("kafka.timeout"),
		Compression:     getCompression(viper.GetString("kafka.compression")),
	}
}

// getRequiredAcks converts string to sarama.RequiredAcks
func getRequiredAcks(acks string) sarama.RequiredAcks {
	switch acks {
	case "none":
		return sarama.NoResponse
	case "one":
		return sarama.WaitForLocal
	case "all":
		return sarama.WaitForAll
	default:
		return sarama.WaitForLocal
	}
}

// getCompression converts string to sarama.CompressionCodec
func getCompression(compression string) sarama.CompressionCodec {
	switch compression {
	case "gzip":
		return sarama.CompressionGZIP
	case "snappy":
		return sarama.CompressionSnappy
	case "lz4":
		return sarama.CompressionLZ4
	case "zstd":
		return sarama.CompressionZSTD
	default:
		return sarama.CompressionNone
	}
}

// generateEventID generates a unique event ID
func generateEventID() string {
	return fmt.Sprintf("evt_%d", time.Now().UnixNano())
}
