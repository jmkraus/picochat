package config

import (
	"math"
	"picochat/limits"
	"strings"
	"testing"
)

func f64ptr(v float64) *float64 { return &v }

func TestClampInt(t *testing.T) {
	tests := []struct {
		name      string
		in        int
		min       int
		max       int
		wantValue int
		wantChg   bool
	}{
		{name: "within range", in: 20, min: limits.MinContext, max: limits.MaxContext, wantValue: 20, wantChg: false},
		{name: "below min", in: 1, min: limits.MinContext, max: limits.MaxContext, wantValue: limits.MinContext, wantChg: true},
		{name: "above max", in: 999, min: limits.MinContext, max: limits.MaxContext, wantValue: limits.MaxContext, wantChg: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, changed := clampInt("context", tt.in, tt.min, tt.max)
			if got != tt.wantValue {
				t.Fatalf("value = %d, want %d", got, tt.wantValue)
			}
			if changed != tt.wantChg {
				t.Fatalf("changed = %v, want %v", changed, tt.wantChg)
			}
		})
	}
}

func TestClampFloat(t *testing.T) {
	tests := []struct {
		name      string
		in        float64
		min       float64
		max       float64
		wantValue float64
		wantChg   bool
	}{
		{name: "within range", in: 0.7, min: limits.MinTemperature, max: limits.MaxTemperature, wantValue: 0.7, wantChg: false},
		{name: "below min", in: -0.1, min: limits.MinTemperature, max: limits.MaxTemperature, wantValue: limits.MinTemperature, wantChg: true},
		{name: "above max", in: 3.5, min: limits.MinTemperature, max: limits.MaxTemperature, wantValue: limits.MaxTemperature, wantChg: true},
		{name: "nan", in: math.NaN(), min: limits.MinTopP, max: limits.MaxTopP, wantValue: limits.MinTopP, wantChg: true},
		{name: "negative inf", in: math.Inf(-1), min: limits.MinTopP, max: limits.MaxTopP, wantValue: limits.MinTopP, wantChg: true},
		{name: "positive inf", in: math.Inf(1), min: limits.MinTopP, max: limits.MaxTopP, wantValue: limits.MaxTopP, wantChg: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, changed := clampFloat("float", tt.in, tt.min, tt.max)
			if got != tt.wantValue {
				t.Fatalf("value = %v, want %v", got, tt.wantValue)
			}
			if changed != tt.wantChg {
				t.Fatalf("changed = %v, want %v", changed, tt.wantChg)
			}
		})
	}
}

func TestNormalizeConfig_Nil(t *testing.T) {
	var cfg *Config
	warnings := cfg.NormalizeConfig()
	if warnings != nil {
		t.Fatalf("warnings = %v, want nil", warnings)
	}
}

func TestNormalizeConfig_NoChanges(t *testing.T) {
	cfg := Config{
		Backend:     "ollama",
		Context:     20,
		Temperature: f64ptr(0.7),
		Top_p:       f64ptr(0.9),
		Effort:      "medium",
	}

	warnings := cfg.NormalizeConfig()
	if len(warnings) != 0 {
		t.Fatalf("warnings = %v, want empty", warnings)
	}
	if cfg.Backend != "ollama" || cfg.Context != 20 || cfg.Temperature == nil || *cfg.Temperature != 0.7 || cfg.Top_p == nil || *cfg.Top_p != 0.9 || cfg.Effort != "medium" {
		t.Fatalf("config unexpectedly changed: %+v", cfg)
	}
}

func TestNormalizeConfig_ClampsAndWarns(t *testing.T) {
	cfg := Config{
		Backend:     "invalid",
		Context:     0,
		Temperature: f64ptr(3.2),
		Top_p:       f64ptr(-0.5),
		Effort:      "invalid",
	}

	warnings := cfg.NormalizeConfig()

	if cfg.Context != limits.MinContext {
		t.Fatalf("context = %d, want %d", cfg.Context, limits.MinContext)
	}
	if cfg.Temperature == nil || *cfg.Temperature != limits.MaxTemperature {
		t.Fatalf("temperature = %v, want %v", cfg.Temperature, limits.MaxTemperature)
	}
	if cfg.Top_p == nil || *cfg.Top_p != limits.MinTopP {
		t.Fatalf("top_p = %v, want %v", cfg.Top_p, limits.MinTopP)
	}
	if cfg.Effort != "medium" {
		t.Fatalf("effort = %q, want %q", cfg.Effort, "medium")
	}
	if cfg.Backend != "ollama" {
		t.Fatalf("backend = %q, want %q", cfg.Backend, "ollama")
	}

	if len(warnings) != 5 {
		t.Fatalf("warnings count = %d, want 5", len(warnings))
	}

	joined := strings.Join(warnings, " | ")
	for _, field := range []string{"context", "temperature", "top_p", "effort", "backend"} {
		if !strings.Contains(joined, field) {
			t.Fatalf("warnings %q do not contain field %q", joined, field)
		}
	}
}

func TestNormalizeBackend(t *testing.T) {
	tests := []struct {
		name      string
		in        string
		wantValue string
		wantWarn  bool
	}{
		{name: "ollama", in: "ollama", wantValue: "ollama", wantWarn: false},
		{name: "openai", in: "openai", wantValue: "openai", wantWarn: false},
		{name: "responses", in: "responses", wantValue: "responses", wantWarn: false},
		{name: "case-insensitive", in: "OpenAI", wantValue: "openai", wantWarn: false},
		{name: "empty fallback", in: "", wantValue: "ollama", wantWarn: true},
		{name: "invalid fallback", in: "foo", wantValue: "ollama", wantWarn: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, warn := normalizeBackend(tt.in)
			if got != tt.wantValue {
				t.Fatalf("value = %q, want %q", got, tt.wantValue)
			}
			if warn != tt.wantWarn {
				t.Fatalf("warn = %v, want %v", warn, tt.wantWarn)
			}
		})
	}
}

func TestHasSchema(t *testing.T) {
	t.Run("nil receiver", func(t *testing.T) {
		var cfg *Config
		if got := cfg.HasSchema(); got {
			t.Fatalf("HasSchema() = %v, want false", got)
		}
	})

	t.Run("nil schema map", func(t *testing.T) {
		cfg := &Config{}
		if got := cfg.HasSchema(); got {
			t.Fatalf("HasSchema() = %v, want false", got)
		}
	})

	t.Run("empty schema map", func(t *testing.T) {
		cfg := &Config{SchemaFmt: map[string]any{}}
		if got := cfg.HasSchema(); got {
			t.Fatalf("HasSchema() = %v, want false", got)
		}
	})

	t.Run("non-empty schema map", func(t *testing.T) {
		cfg := &Config{SchemaFmt: map[string]any{"type": "object"}}
		if got := cfg.HasSchema(); !got {
			t.Fatalf("HasSchema() = %v, want true", got)
		}
	})
}
