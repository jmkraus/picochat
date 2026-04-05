package config

import (
	"math"
	"strings"
	"testing"
)

func TestClampInt(t *testing.T) {
	tests := []struct {
		name      string
		in        int
		min       int
		max       int
		wantValue int
		wantChg   bool
	}{
		{name: "within range", in: 20, min: MinContext, max: MaxContext, wantValue: 20, wantChg: false},
		{name: "below min", in: 1, min: MinContext, max: MaxContext, wantValue: MinContext, wantChg: true},
		{name: "above max", in: 999, min: MinContext, max: MaxContext, wantValue: MaxContext, wantChg: true},
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
		{name: "within range", in: 0.7, min: MinTemperature, max: MaxTemperature, wantValue: 0.7, wantChg: false},
		{name: "below min", in: -0.1, min: MinTemperature, max: MaxTemperature, wantValue: MinTemperature, wantChg: true},
		{name: "above max", in: 3.5, min: MinTemperature, max: MaxTemperature, wantValue: MaxTemperature, wantChg: true},
		{name: "nan", in: math.NaN(), min: MinTopP, max: MaxTopP, wantValue: MinTopP, wantChg: true},
		{name: "negative inf", in: math.Inf(-1), min: MinTopP, max: MaxTopP, wantValue: MinTopP, wantChg: true},
		{name: "positive inf", in: math.Inf(1), min: MinTopP, max: MaxTopP, wantValue: MaxTopP, wantChg: true},
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
	warnings := NormalizeConfig(nil)
	if warnings != nil {
		t.Fatalf("warnings = %v, want nil", warnings)
	}
}

func TestNormalizeConfig_NoChanges(t *testing.T) {
	cfg := Config{
		Context:     20,
		Temperature: 0.7,
		Top_p:       0.9,
	}

	warnings := NormalizeConfig(&cfg)
	if len(warnings) != 0 {
		t.Fatalf("warnings = %v, want empty", warnings)
	}
	if cfg.Context != 20 || cfg.Temperature != 0.7 || cfg.Top_p != 0.9 {
		t.Fatalf("config unexpectedly changed: %+v", cfg)
	}
}

func TestNormalizeConfig_ClampsAndWarns(t *testing.T) {
	cfg := Config{
		Context:     0,
		Temperature: 3.2,
		Top_p:       -0.5,
	}

	warnings := NormalizeConfig(&cfg)

	if cfg.Context != MinContext {
		t.Fatalf("context = %d, want %d", cfg.Context, MinContext)
	}
	if cfg.Temperature != MaxTemperature {
		t.Fatalf("temperature = %v, want %v", cfg.Temperature, MaxTemperature)
	}
	if cfg.Top_p != MinTopP {
		t.Fatalf("top_p = %v, want %v", cfg.Top_p, MinTopP)
	}

	if len(warnings) != 3 {
		t.Fatalf("warnings count = %d, want 3", len(warnings))
	}

	joined := strings.Join(warnings, " | ")
	for _, field := range []string{"context", "temperature", "top_p"} {
		if !strings.Contains(joined, field) {
			t.Fatalf("warnings %q do not contain field %q", joined, field)
		}
	}
}
