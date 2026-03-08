package envs

import "testing"

func TestAllowedField_ConfigFieldsAreAllowed(t *testing.T) {
	for _, cfg := range ConfigEnvVars {
		if !AllowedField(cfg.Field) {
			t.Errorf("expected field %q to be allowed", cfg.Field)
		}
	}
}

func TestAllowedField_InvalidFieldsAreRejected(t *testing.T) {
	tests := []string{"", "unknown", "URL", " model", "temperature "}

	for _, field := range tests {
		if AllowedField(field) {
			t.Errorf("expected field %q to be rejected", field)
		}
	}
}
