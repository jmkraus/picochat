package config

import (
	"encoding/json"
	"fmt"
	"picochat/envs"
	"picochat/vartypes"
)

// applyConfig updates a specific config element.
//
// Parameters:
//
//	key (string)  - the configuration key to modify
//	value (any)   - the new value for the key
//
// Returns:
//
//	error - error if any
func (c *Config) applyConfigValue(key string, val any) error {
	patch := map[string]any{key: val}
	b, err := json.Marshal(patch)
	if err != nil {
		return err
	}
	return json.Unmarshal(b, c)
}

// applyEnvValues updates config fields according to
// set environment variables.
//
// Parameters:
//
//	none
//
// Returns:
//
//	error - error if any
func (c *Config) applyEnvValues() error {
	for _, spec := range envs.ConfigEnvVars {
		envVal, lookup := envs.GetEnv(spec.Env)
		if !lookup || envVal == "" {
			continue // Skip if not set or empty
		}

		v, err := vartypes.Convert(spec.Type, envVal)
		if err != nil {
			return fmt.Errorf("convert type for env %s failed: %w", spec.Env, err)
		}
		if err := c.applyConfigValue(spec.Field, v); err != nil {
			return fmt.Errorf("apply config value for env %s failed: %w", spec.Env, err)
		}
	}
	return nil
}

// NormalizeConfig clamps numeric config values to their valid ranges and
// normalizes effort/backend values to supported sets.
// It mutates cfg and returns warning messages for changed values.
//
// Parameters:
//
//	none
//
// Returns:
//
//	[]string - warnings for each clamped field
func (c *Config) NormalizeConfig() []string {
	if c == nil {
		return nil
	}

	var warnings []string

	origCtx := c.Context
	if v, changed := clampInt("context", c.Context, MinContext, MaxContext); changed {
		c.Context = v
		warnings = append(warnings, fmt.Sprintf("config value 'context' (%d) out of range [%d..%d], clamped to %d", origCtx, MinContext, MaxContext, v))
	}

	if c.Temperature != nil {
		origTemp := *c.Temperature
		if v, changed := clampFloat("temperature", *c.Temperature, MinTemperature, MaxTemperature); changed {
			c.Temperature = &v
			warnings = append(warnings, fmt.Sprintf("config value 'temperature' (%g) out of range [%g..%g], clamped to %g", origTemp, MinTemperature, MaxTemperature, v))
		}
	}

	if c.Top_p != nil {
		origTopP := *c.Top_p
		if v, changed := clampFloat("top_p", *c.Top_p, MinTopP, MaxTopP); changed {
			c.Top_p = &v
			warnings = append(warnings, fmt.Sprintf("config value 'top_p' (%g) out of range [%g..%g], clamped to %g", origTopP, MinTopP, MaxTopP, v))
		}
	}

	origEffort := c.Effort
	if v, warn := normalizeEffort(c.Effort); v != c.Effort {
		c.Effort = v
		if warn {
			warnings = append(warnings, fmt.Sprintf("config value 'effort' (%q) invalid, normalized to %q", origEffort, v))
		}
	}

	origBackend := c.Backend
	if v, warn := normalizeBackend(c.Backend); v != c.Backend {
		c.Backend = v
		if warn {
			warnings = append(warnings, fmt.Sprintf("config value 'backend' (%q) invalid, normalized to %q", origBackend, v))
		}
	}

	return warnings
}

// HasSchema checks if a JSON schema for structured output was loaded.
//
// Parameters:
//
//	none
//
// Returns:
//
//	bool - true or false
func (c *Config) HasSchema() bool {
	if c == nil {
		return false
	}
	return len(c.SchemaFmt) > 0
}
