package config

import (
	"log"

	"go.uber.org/zap"
)

func NewLogger() *zap.Logger {

	// Initialize zap logger
	logger, err := zap.NewProduction()
	if err != nil {
		log.Fatalf("Failed to initialize zap logger: %v", err)
	}

	// No need to defer logger.Sync() in Lambda
	return logger
}
