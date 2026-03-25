package external

import (
	"context"

	"smartsure/internal/domain/services/fraud_detection"
	"smartsure/pkg/external"

	"github.com/google/uuid"
)

// OpenAIAdapter adapts OpenAI service to AIService interface
type OpenAIAdapter struct {
	openaiService *external.OpenAIFraudService
}

// NewOpenAIAdapter creates a new OpenAI adapter
func NewOpenAIAdapter(openaiService *external.OpenAIFraudService) *OpenAIAdapter {
	return &OpenAIAdapter{
		openaiService: openaiService,
	}
}

// PredictFraud implements AIService interface
func (a *OpenAIAdapter) PredictFraud(ctx context.Context, features map[string]interface{}) (*fraud_detection.AIPrediction, error) {
	return a.openaiService.PredictFraud(ctx, features)
}

// TrainModel implements AIService interface
func (a *OpenAIAdapter) TrainModel(ctx context.Context, trainingData []fraud_detection.TrainingData) error {
	return a.openaiService.TrainModel(ctx, trainingData)
}

// UpdateModel implements AIService interface
func (a *OpenAIAdapter) UpdateModel(ctx context.Context, modelID uuid.UUID, newData []fraud_detection.TrainingData) error {
	return a.openaiService.UpdateModel(ctx, modelID, newData)
}

// GetModelMetrics implements AIService interface
func (a *OpenAIAdapter) GetModelMetrics(ctx context.Context, modelID uuid.UUID) (*fraud_detection.ModelMetrics, error) {
	return a.openaiService.GetModelMetrics(ctx, modelID)
}
