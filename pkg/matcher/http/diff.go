package http

import (
	"fmt"
	"reflect"
)

type RiskLevel string
type AssessmentType string

const (
	RiskNone RiskLevel = "NONE"
	RiskLow  RiskLevel = "LOW"
	RiskHigh RiskLevel = "HIGH"

	NoChange       AssessmentType = "NO_CHANGE"
	SchemaAdded    AssessmentType = "SCHEMA_ADDED"
	SchemaModified AssessmentType = "SCHEMA_MODIFIED"
	SchemaDeleted  AssessmentType = "SCHEMA_DELETED"
)

// CollectJSONWithIndices traverses JSON payloads preserving array indices in path keys
func CollectJSONWithIndices(prefix string, val interface{}, target map[string]interface{}) {
	if val == nil {
		target[prefix] = nil
		return
	}

	switch v := val.(type) {
	case map[string]interface{}:
		for k, child := range v {
			newKey := k
			if prefix != "" {
				newKey = prefix + "." + k
			}
			CollectJSONWithIndices(newKey, child, target)
		}
	case []interface{}:
		for i, child := range v {
			newKey := fmt.Sprintf("%s[%d]", prefix, i)
			CollectJSONWithIndices(newKey, child, target)
		}
	default:
		target[prefix] = v
	}
}

// GradeFailureAssessment ensures modified items flag as SCHEMA_MODIFIED rather than SCHEMA_ADDED
func GradeFailureAssessment(expectedFlat, actualFlat map[string]interface{}) (RiskLevel, AssessmentType) {
	hasModifications := false
	hasAdditions := false

	for k, expVal := range expectedFlat {
		actVal, exists := actualFlat[k]
		if !exists {
			return RiskHigh, SchemaDeleted
		}
		if !reflect.DeepEqual(expVal, actVal) {
			hasModifications = true
		}
	}

	for k := range actualFlat {
		if _, exists := expectedFlat[k]; !exists {
			hasAdditions = true
		}
	}

	if hasModifications {
		return RiskHigh, SchemaModified
	}
	if hasAdditions {
		return RiskLow, SchemaAdded
	}
	return RiskNone, NoChange
}
