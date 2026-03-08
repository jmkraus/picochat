package envs

import "testing"

func TestAllowedRuntimeField_FollowsConfigRuntimeFlag(t *testing.T) {
	for _, cfg := range ConfigEnvVars {
		got := AllowedRuntimeField(cfg.Field)
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
