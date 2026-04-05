package config

import (
	"fmt"
	"math"
)

const (
	MinContext     = 3
	MaxContext     = 100
	MinTemperature = 0.0
	MaxTemperature = 2.0
	MinTopP        = 0.0
	MaxTopP        = 1.0
)

// ClampInt clamps an integer value to the given inclusive range.
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
func ClampInt(name string, v, min, max int) (int, bool) {
	_ = name
	if v < min {
		return min, true
	}
	if v > max {
		return max, true
	}
	return v, false
}

// ClampFloat clamps a float64 value to the given inclusive range.
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
func ClampFloat(name string, v, min, max float64) (float64, bool) {
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

// NormalizeConfig clamps numeric config values to their valid ranges.
// It mutates cfg and returns warning messages for changed values.
//
// Parameters:
//
//	cfg (*Config) - target config to normalize
//
// Returns:
//
//	[]string - warnings for each clamped field
func NormalizeConfig(cfg *Config) []string {
	if cfg == nil {
		return nil
	}

	var warnings []string

	origCtx := cfg.Context
	if v, changed := ClampInt("context", cfg.Context, MinContext, MaxContext); changed {
		cfg.Context = v
		warnings = append(warnings, fmt.Sprintf("config value 'context' (%d) out of range [%d..%d], clamped to %d", origCtx, MinCtx, MaxCtx, v))
	}

	origTemp := cfg.Temperature
	if v, changed := ClampFloat("temperature", cfg.Temperature, MinTemperature, MaxTemperature); changed {
		cfg.Temperature = v
		warnings = append(warnings, fmt.Sprintf("config value 'temperature' (%g) out of range [%g..%g], clamped to %g", origTemp, MinTemperature, MaxTemperature, v))
	}

	origTopP := cfg.Top_p
	if v, changed := ClampFloat("top_p", cfg.Top_p, MinTopP, MaxTopP); changed {
		cfg.Top_p = v
		warnings = append(warnings, fmt.Sprintf("config value 'top_p' (%g) out of range [%g..%g], clamped to %g", origTopP, MinTopP, MaxTopP, v))
	}

	return warnings
}
