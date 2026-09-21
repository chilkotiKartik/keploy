package http

import (
	"encoding/json"
	"fmt"
	"reflect"

	"go.keploy.io/server/v2/pkg/models"
)

// NormalizeAssertionValue converts assertion maps/slices into standard interface representations
func NormalizeAssertionValue(val interface{}) interface{} {
	if val == nil {
		return nil
	}
	switch v := val.(type) {
	case map[models.AssertionType]interface{}:
		res := make(map[string]interface{}, len(v))
		for k, val := range v {
			res[string(k)] = NormalizeAssertionValue(val)
		}
		return res
	case map[interface{}]interface{}:
		res := make(map[string]interface{}, len(v))
		for k, val := range v {
			res[fmt.Sprintf("%v", k)] = NormalizeAssertionValue(val)
		}
		return res
	case map[string]interface{}:
		res := make(map[string]interface{}, len(v))
		for k, val := range v {
			res[k] = NormalizeAssertionValue(val)
		}
		return res
	case []interface{}:
		res := make([]interface{}, len(v))
		for i, item := range v {
			res[i] = NormalizeAssertionValue(item)
		}
		return res
	default:
		return v
	}
}

// CompareJSONContains recursively checks if expected subset exists within actual JSON structure
func CompareJSONContains(actual, expected interface{}) bool {
	if expected == nil {
		return actual == nil
	}
	if actual == nil {
		return false
	}

	actualNorm := NormalizeAssertionValue(actual)
	expectedNorm := NormalizeAssertionValue(expected)

	switch exp := expectedNorm.(type) {
	case map[string]interface{}:
		actMap, ok := actualNorm.(map[string]interface{})
		if !ok {
			return false
		}
		for k, expVal := range exp {
			actVal, exists := actMap[k]
			if !exists {
				return false
			}
			if !CompareJSONContains(actVal, expVal) {
				return false
			}
		}
		return true

	case []interface{}:
		actSlice, ok := actualNorm.([]interface{})
		if !ok {
			return false
		}
		for _, expItem := range exp {
			matched := false
			for _, actItem := range actSlice {
				if CompareJSONContains(actItem, expItem) {
					matched = true
					break
				}
			}
			if !matched {
				return false
			}
		}
		return true

	default:
		return compareScalarValues(actualNorm, expectedNorm)
	}
}

func compareScalarValues(a, b interface{}) bool {
	if reflect.DeepEqual(a, b) {
		return true
	}
	valA, okA := toFloat64(a)
	valB, okB := toFloat64(b)
	if okA && okB {
		return valA == valB
	}
	return fmt.Sprintf("%v", a) == fmt.Sprintf("%v", b)
}

func toFloat64(v interface{}) (float64, bool) {
	switch n := v.(type) {
	case float64:
		return n, true
	case float32:
		return float64(n), true
	case int:
		return float64(n), true
	case int32:
		return float64(n), true
	case int64:
		return float64(n), true
	case uint:
		return float64(n), true
	case uint64:
		return float64(n), true
	case json.Number:
		f, err := n.Float64()
		return f, err == nil
	default:
		return 0, false
	}
}
