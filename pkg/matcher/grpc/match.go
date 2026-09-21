package grpc

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"

	"go.keploy.io/server/v2/pkg/models"
	httpMatcher "go.keploy.io/server/v2/pkg/matcher/http"
	"go.uber.org/zap"
)

type Matcher struct {
	logger *zap.Logger
}

func NewMatcher(logger *zap.Logger) *Matcher {
	return &Matcher{logger: logger}
}

func (m *Matcher) Match(ctx context.Context, tc *models.TestCase, actualResp *models.GrpcResp) (bool, error) {
	if len(tc.Assertions) > 0 {
		pass, err := m.EvaluateAssertions(ctx, tc.Assertions, actualResp)
		if err != nil || !pass {
			return false, err
		}
	}
	return true, nil
}

func (m *Matcher) EvaluateAssertions(ctx context.Context, assertions map[string]interface{}, actualResp *models.GrpcResp) (bool, error) {
	for key, expected := range assertions {
		switch key {
		case "status_code", "grpc_status":
			actualStatus := strconv.Itoa(int(actualResp.Err.Code))
			if fmt.Sprintf("%v", expected) != actualStatus {
				m.logger.Warn("grpc status assertion failed", zap.String("expected", fmt.Sprintf("%v", expected)), zap.String("actual", actualStatus))
				return false, nil
			}
		case "json_contains", "body_contains":
			var actualBody interface{}
			if err := json.Unmarshal([]byte(actualResp.Body), &actualBody); err != nil {
				return false, fmt.Errorf("failed to unmarshal grpc response body for assertion: %w", err)
			}
			if !httpMatcher.CompareJSONContains(actualBody, expected) {
				m.logger.Warn("grpc body assertion failed", zap.Any("expected", expected), zap.Any("actual", actualBody))
				return false, nil
			}
		}
	}
	return true, nil
}
