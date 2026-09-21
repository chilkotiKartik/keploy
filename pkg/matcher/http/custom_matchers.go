package http

import (
	"fmt"
	"regexp"
	"time"

	"github.com/google/uuid"
)

type CustomMatcherFunc func(actual interface{}) bool

var DefaultCustomMatchers = map[string]CustomMatcherFunc{
	"@isUUID": func(actual interface{}) bool {
		str, ok := actual.(string)
		if !ok {
			return false
		}
		_, err := uuid.Parse(str)
		return err == nil
	},
	"@isDateTime": func(actual interface{}) bool {
		str, ok := actual.(string)
		if !ok {
			return false
		}
		for _, layout := range []string{time.RFC3339, time.RFC3339Nano, "2006-01-02 15:04:05"} {
			if _, err := time.Parse(layout, str); err == nil {
				return true
			}
		}
		return false
	},
	"@isNumber": func(actual interface{}) bool {
		_, ok := toFloat64(actual)
		return ok
	},
}

// EvaluateCustomMatcher checks whether an expected assertion tag matches the actual value
func EvaluateCustomMatcher(expectedTag string, actual interface{}) (bool, bool) {
	if fn, exists := DefaultCustomMatchers[expectedTag]; exists {
		return fn(actual), true
	}
	if r := regexp.MustCompile(`^@regex\((.+)\)$`); r.MatchString(expectedTag) {
		match := r.FindStringSubmatch(expectedTag)
		re, err := regexp.Compile(match[1])
		if err != nil {
			return false, true
		}
		return re.MatchString(fmt.Sprintf("%v", actual)), true
	}
	return false, false
}
