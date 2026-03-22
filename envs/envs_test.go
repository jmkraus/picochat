package envs

import (
	"strings"
	"testing"
)

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

func TestConfigEnvVarsTable_HidesSensitiveValues(t *testing.T) {
	t.Setenv("PICOCHAT_API_KEY", "sk-secret-value")

	table := ConfigEnvVarsTable()
	if strings.Contains(table, "sk-secret-value") {
		t.Fatalf("table leaked sensitive value")
	}
	if !strings.Contains(table, "PICOCHAT_API_KEY") {
		t.Fatalf("table does not contain API key row")
	}
	if !strings.Contains(table, "[hidden]") {
		t.Fatalf("table does not mark sensitive value as hidden")
	}
}

func TestConfigEnvVarsTable_ShowsNonSensitiveValues(t *testing.T) {
	t.Setenv("PICOCHAT_MODEL", "test-model")

	table := ConfigEnvVarsTable()
	if !strings.Contains(table, "PICOCHAT_MODEL") {
		t.Fatalf("table does not contain model row")
	}
	if !strings.Contains(table, "test-model") {
		t.Fatalf("non-sensitive value should be visible")
	}
}
