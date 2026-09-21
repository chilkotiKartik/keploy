package http

import (
	"testing"
)

func TestCompareJSONContains(t *testing.T) {
	actual := map[string]interface{}{
		"id":       123.0,
		"name":     "Keploy",
		"active":   true,
		"tags":     []interface{}{"test", "mock", "replay"},
		"metadata": map[string]interface{}{"env": "prod", "version": 2.0},
	}

	expected := map[string]interface{}{
		"id":       123,
		"name":     "Keploy",
		"metadata": map[string]interface{}{"env": "prod"},
	}

	if !CompareJSONContains(actual, expected) {
		t.Errorf("expected CompareJSONContains to return true for subset match, got false")
	}
}
