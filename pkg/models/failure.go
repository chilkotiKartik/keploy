package models

type RiskLevel string

const (
	RiskNone RiskLevel = "NONE"
	RiskLow  RiskLevel = "LOW"
	RiskHigh RiskLevel = "HIGH"
)

type FailureInfo struct {
	Risk               RiskLevel `json:"risk" yaml:"risk"`
	AssessmentCategory string    `json:"assessmentCategory" yaml:"assessmentCategory"`
	Message            string    `json:"message" yaml:"message"`
}

// IsAutoReplayHighRisk returns whether a test failure is considered high risk and ineligible for auto-pass
func IsAutoReplayHighRisk(f *FailureInfo) bool {
	if f == nil {
		return false
	}
	if f.Risk == RiskHigh {
		return true
	}
	switch f.AssessmentCategory {
	case "SchemaAdded", "SchemaBroken", "SchemaModified", "StatusCodeMismatch", "HeaderMismatch":
		return true
	default:
		return false
	}
}
