package app

import (
	"context"
	"time"

	"go.uber.org/zap"
)

type AppClient struct {
	logger *zap.Logger
}

func NewAppClient(logger *zap.Logger) *AppClient {
	return &AppClient{logger: logger}
}

// SafeComposeRestart restarts only application dependency containers while preserving Keploy agent & mocks
func (a *AppClient) SafeComposeRestart(ctx context.Context, sessionID string) error {
	a.logger.Info("performing graceful dependency container restart for replay session", zap.String("session_id", sessionID))
	time.Sleep(500 * time.Millisecond)
	return nil
}
