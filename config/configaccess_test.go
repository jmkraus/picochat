package config

import (
	"math"
	"os"
	"picochat/envs"
	"picochat/vartypes"
	"strings"
	"testing"
)

func setEnv(t *testing.T, key, value string) func() {
	t.Helper()

	prev, hadPrev := os.LookupEnv(key)
	if err := os.Setenv(key, value); err != nil {
		t.Fatalf("Setenv(%q) failed: %v", key, err)
	}

	return func() {
		if hadPrev {
			_ = os.Setenv(key, prev)
			return
		}
		_ = os.Unsetenv(key)
	}
}

func TestConfig_applyConfigValue(t *testing.T) {
	t.Run("sets string field", func(t *testing.T) {
		cfg := Config{Backend: "ollama"}
		if err := cfg.applyConfigValue("backend", "openai"); err != nil {
			t.Fatalf("applyConfigValue returned error: %v", err)
		}
		if cfg.Backend != "openai" {
			t.Fatalf("Backend = %q, want %q", cfg.Backend, "openai")
		}
	})

	t.Run("sets int field", func(t *testing.T) {
		cfg := Config{Context: 20}
		if err := cfg.applyConfigValue("context", 42); err != nil {
			t.Fatalf("applyConfigValue returned error: %v", err)
		}
		if cfg.Context != 42 {
			t.Fatalf("Context = %d, want %d", cfg.Context, 42)
		}
	})

	t.Run("sets float pointer field", func(t *testing.T) {
		cfg := Config{Temperature: nil}
		if err := cfg.applyConfigValue("temperature", 0.7); err != nil {
			t.Fatalf("applyConfigValue returned error: %v", err)
		}
		if cfg.Temperature == nil || *cfg.Temperature != 0.7 {
			t.Fatalf("Temperature = %v, want %v", cfg.Temperature, 0.7)
		}
	})

	t.Run("unknown field is ignored", func(t *testing.T) {
		cfg := Config{Backend: "ollama"}
		if err := cfg.applyConfigValue("does_not_exist", "x"); err != nil {
			t.Fatalf("applyConfigValue returned error: %v", err)
		}
		if cfg.Backend != "ollama" {
			t.Fatalf("Backend = %q, want %q", cfg.Backend, "ollama")
		}
	})
}

func TestConfig_applyEnvValues(t *testing.T) {
	t.Run("applies env vars and skips empty values", func(t *testing.T) {
		restoreCtx := setEnv(t, "PICOCHAT_CONTEXT", "50")
		defer restoreCtx()
		restoreReasoning := setEnv(t, "PICOCHAT_REASONING", "yes")
		defer restoreReasoning()
		restoreTemp := setEnv(t, "PICOCHAT_TEMPERATURE", "0.6")
		defer restoreTemp()
		restoreModelEmpty := setEnv(t, "PICOCHAT_MODEL", "")
		defer restoreModelEmpty()

		cfg := defaultConfig()
		cfg.Model = "keepme"
		if err := cfg.applyEnvValues(); err != nil {
			t.Fatalf("applyEnvValues returned error: %v", err)
		}

		if cfg.Context != 50 {
			t.Fatalf("Context = %d, want %d", cfg.Context, 50)
		}
		if cfg.Reasoning != true {
			t.Fatalf("Reasoning = %v, want %v", cfg.Reasoning, true)
		}
		if cfg.Temperature == nil || *cfg.Temperature != 0.6 {
			t.Fatalf("Temperature = %v, want %v", cfg.Temperature, 0.6)
		}
		if cfg.Model != "keepme" {
			t.Fatalf("Model = %q, want %q (empty env must be skipped)", cfg.Model, "keepme")
		}
	})

	t.Run("returns convert error", func(t *testing.T) {
		restore := setEnv(t, "PICOCHAT_CONTEXT", "not-an-int")
		defer restore()

		cfg := defaultConfig()
		err := cfg.applyEnvValues()
		if err == nil {
			t.Fatalf("applyEnvValues expected error, got nil")
		}
		if !strings.Contains(err.Error(), "convert type for env PICOCHAT_CONTEXT failed") {
			t.Fatalf("error = %q, want to contain %q", err.Error(), "convert type for env PICOCHAT_CONTEXT failed")
		}
	})

	t.Run("returns apply config error", func(t *testing.T) {
		orig := envs.ConfigEnvVars
		envs.ConfigEnvVars = append(envs.ConfigEnvVars, envs.EnvSpec{
			Env:   "PICOCHAT_BAD_CONTEXT",
			Type:  vartypes.VarString,
			Field: "context",
		})
		t.Cleanup(func() {
			envs.ConfigEnvVars = orig
		})

		restore := setEnv(t, "PICOCHAT_BAD_CONTEXT", "not-a-number")
		defer restore()

		cfg := defaultConfig()
		err := cfg.applyEnvValues()
		if err == nil {
			t.Fatalf("applyEnvValues expected error, got nil")
		}
		if !strings.Contains(err.Error(), "apply config value for env PICOCHAT_BAD_CONTEXT failed") {
			t.Fatalf("error = %q, want to contain %q", err.Error(), "apply config value for env PICOCHAT_BAD_CONTEXT failed")
		}
	})
}

func TestConfig_NormalizeConfig(t *testing.T) {
	t.Run("nil receiver returns nil", func(t *testing.T) {
		var cfg *Config
		if got := cfg.NormalizeConfig(); got != nil {
			t.Fatalf("warnings = %v, want nil", got)
		}
	})

	t.Run("clamps and normalizes with warnings", func(t *testing.T) {
		temp := -0.1
		topP := 2.0
		cfg := Config{
			Context:     1,
			Temperature: &temp,
			Top_p:       &topP,
			Effort:      "MID",
			Backend:     "nope",
		}
		warnings := cfg.NormalizeConfig()

		if cfg.Context != MinContext {
			t.Fatalf("Context = %d, want %d", cfg.Context, MinContext)
		}
		if cfg.Temperature == nil || *cfg.Temperature != MinTemperature {
			t.Fatalf("Temperature = %v, want %v", cfg.Temperature, MinTemperature)
		}
		if cfg.Top_p == nil || *cfg.Top_p != MaxTopP {
			t.Fatalf("Top_p = %v, want %v", cfg.Top_p, MaxTopP)
		}
		if cfg.Effort != "medium" {
			t.Fatalf("Effort = %q, want %q", cfg.Effort, "medium")
		}
		if cfg.Backend != "ollama" {
			t.Fatalf("Backend = %q, want %q", cfg.Backend, "ollama")
		}
		if len(warnings) != 5 {
			t.Fatalf("warnings count = %d, want %d; warnings=%v", len(warnings), 5, warnings)
		}
	})

	t.Run("normalizes backend case without warning", func(t *testing.T) {
		cfg := Config{Backend: "Responses", Context: 20, Effort: "high"}
		warnings := cfg.NormalizeConfig()
		if cfg.Backend != "responses" {
			t.Fatalf("Backend = %q, want %q", cfg.Backend, "responses")
		}
		for _, w := range warnings {
			if strings.Contains(w, "backend") {
				t.Fatalf("unexpected backend warning: %q", w)
			}
		}
	})

	t.Run("clamps NaN temperature", func(t *testing.T) {
		temp := math.NaN()
		cfg := Config{Context: 20, Temperature: &temp}
		warnings := cfg.NormalizeConfig()
		if cfg.Temperature == nil || *cfg.Temperature != MinTemperature {
			t.Fatalf("Temperature = %v, want %v", cfg.Temperature, MinTemperature)
		}
		found := false
		for _, w := range warnings {
			if strings.Contains(w, "temperature") {
				found = true
				break
			}
		}
		if !found {
			t.Fatalf("expected temperature warning, got %v", warnings)
		}
	})

	t.Run("no changes yields no warnings", func(t *testing.T) {
		temp := 0.7
		topP := 0.9
		cfg := Config{
			Backend:     "ollama",
			Context:     20,
			Temperature: &temp,
			Top_p:       &topP,
			Effort:      "high",
		}
		warnings := cfg.NormalizeConfig()
		if len(warnings) != 0 {
			t.Fatalf("warnings = %v, want empty", warnings)
		}
	})
}

func TestConfig_HasSchema(t *testing.T) {
	t.Run("nil receiver", func(t *testing.T) {
		var cfg *Config
		if cfg.HasSchema() {
			t.Fatalf("HasSchema() = true, want false")
		}
	})

	t.Run("empty schema", func(t *testing.T) {
		cfg := Config{SchemaFmt: map[string]any{}}
		if cfg.HasSchema() {
			t.Fatalf("HasSchema() = true, want false")
		}
	})

	t.Run("non-empty schema", func(t *testing.T) {
		cfg := Config{SchemaFmt: map[string]any{"type": "object"}}
		if !cfg.HasSchema() {
			t.Fatalf("HasSchema() = false, want true")
		}
	})
}
