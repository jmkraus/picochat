package config

import (
	"math"
	"strings"
)

const (
	MinContext     = 3
	MaxContext     = 100
	MinTemperature = 0.0
	MaxTemperature = 2.0
	MinTopP        = 0.0
	MaxTopP        = 1.0
)

// clampInt clamps an integer value to the given inclusive range.
//
// Parameters:
//
//	name (string) - logical field name (for warnings/logging context)
//	v    (int)    - input value to clamp
//	min  (int)    - minimum allowed value
//	max  (int)    - maximum allowed value
//
// Returns:
//
//	int  - clamped or original value
//	bool - true if value was changed
func clampInt(name string, v, min, max int) (int, bool) {
	_ = name
	if v < min {
		return min, true
	}
	if v > max {
		return max, true
	}
	return v, false
}

// clampFloat clamps a float64 value to the given inclusive range.
// NaN and Inf values are treated as invalid and clamped to boundaries.
//
// Parameters:
//
//	name (string)  - logical field name (for warnings/logging context)
//	v    (float64) - input value to clamp
//	min  (float64) - minimum allowed value
//	max  (float64) - maximum allowed value
//
// Returns:
//
//	float64 - clamped or original value
//	bool    - true if value was changed
func clampFloat(name string, v, min, max float64) (float64, bool) {
	_ = name
	if math.IsNaN(v) || math.IsInf(v, -1) {
		return min, true
	}
	if math.IsInf(v, 1) {
		return max, true
	}
	if v < min {
		return min, true
	}
	if v > max {
		return max, true
	}
	return v, false
}

// normalizeEffort validates and normalizes the effort value.
//
// Parameters:
//
//	raw (string) - input effort value
//
// Returns:
//
//	string - normalized effort value
//	bool   - true if fallback/alias handling was applied
func normalizeEffort(raw string) (string, bool) {
	value := strings.ToLower(strings.TrimSpace(raw))

	switch value {
	case "none", "low", "medium", "high":
		return value, false
	case "mid":
		return "medium", true
	case "":
		return "medium", true
	default:
		return "medium", true
	}
}

// normalizeBackend validates and normalizes the backend value.
//
// Parameters:
//
//	raw (string) - input backend value
//
// Returns:
//
//	string - normalized backend value
//	bool   - true if fallback/alias handling was applied
func normalizeBackend(raw string) (string, bool) {
	value := strings.ToLower(strings.TrimSpace(raw))

	switch value {
	case "ollama", "openai", "responses":
		return value, false
	default:
		return "ollama", true
	}
}
