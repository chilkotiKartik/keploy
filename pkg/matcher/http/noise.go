package http

import (
	"strings"
)

// MergeNoiseConfigs deeply combines global noise definitions with per-testcase noise
func MergeNoiseConfigs(globalNoise, testNoise map[string][]string) map[string][]string {
	result := make(map[string][]string)

	for k, list := range globalNoise {
		result[k] = append([]string{}, list...)
	}

	for k, list := range testNoise {
		existing := result[k]
		for _, item := range list {
			if !contains(existing, item) {
				existing = append(existing, item)
			}
		}
		result[k] = existing
	}

	return result
}

// IsFieldIgnored checks exact path match or wildcard rule (e.g. response.body.*)
func IsFieldIgnored(path string, noiseList []string) bool {
	for _, pattern := range noiseList {
		if pattern == path {
			return true
		}
		if strings.HasSuffix(pattern, ".*") {
			prefix := strings.TrimSuffix(pattern, ".*")
			if strings.HasPrefix(path, prefix+".") {
				return true
			}
		}
		if pattern == "**" || pattern == "*" {
			return true
		}
	}
	return false
}

func contains(slice []string, val string) bool {
	for _, item := range slice {
		if item == val {
			return true
		}
	}
	return false
}
