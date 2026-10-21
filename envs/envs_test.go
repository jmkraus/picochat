package envs

import (
	"os"
	"testing"
)

func TestAllowedRuntimeField_FollowsConfigRuntimeFlag(t *testing.T) {
	for _, cfg := range ConfigEnvVars {
		got := AllowedRuntimeField(cfg.JsonField)
		if got != cfg.Runtime {
			t.Errorf("field %q allowed=%v, want %v", cfg.Field, got, cfg.Runtime)
		}
	}
}

func TestAllowedRuntimeField_InvalidFieldsAreRejected(t *testing.T) {
	tests := []string{"", "unknown", "URL", " model", "temperature "}

	for _, field := range tests {
		if AllowedRuntimeField(field) {
			t.Errorf("expected field %q to be rejected", field)
		}
	}
}

func TestConfigEnvVars_JSONFieldsAreNonEmptyAndUnique(t *testing.T) {
	seen := make(map[string]struct{}, len(ConfigEnvVars))

	for _, spec := range ConfigEnvVars {
		if spec.JsonField == "" {
			t.Errorf("environment variable %q has an empty JSON field", spec.Env)
			continue
		}
		if _, exists := seen[spec.JsonField]; exists {
			t.Errorf("JSON field %q is configured more than once", spec.JsonField)
		}
		seen[spec.JsonField] = struct{}{}
	}
}

func TestUnsetApiKeyEnv_RemovesApiKeyInAllStates(t *testing.T) {
	tests := []struct {
		name  string
		value *string
	}{
		{
			name:  "unset",
			value: nil,
		},
		{
			name:  "set but empty",
			value: func() *string { value := ""; return &value }(),
		},
		{
			name:  "set with secret token",
			value: func() *string { value := "sk-secret-value"; return &value }(),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Set a value first so t.Setenv can restore the original state after
			// the test, even when this case needs the variable to be unset.
			t.Setenv("PICOCHAT_API_KEY", "original-value")
			if tt.value == nil {
				if err := os.Unsetenv("PICOCHAT_API_KEY"); err != nil {
					t.Fatalf("prepare unset API key: %v", err)
				}
			} else {
				t.Setenv("PICOCHAT_API_KEY", *tt.value)
			}

			if err := UnsetApiKeyEnv(); err != nil {
				t.Fatalf("UnsetApiKeyEnv() error = %v", err)
			}
			if _, ok := os.LookupEnv("PICOCHAT_API_KEY"); ok {
				t.Fatalf("PICOCHAT_API_KEY is still set")
			}
		})
	}
}
