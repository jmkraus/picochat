package format

import (
	"testing"
)

// TestAllowedKeys tests the AllowedKeys function with various inputs.
func TestAllowedKeys(t *testing.T) {
	tests := []struct {
		input    string
		expected string
		ok       bool
	}{
		{"plain", "plain", true},
		{"Plain", "plain", true},
		{"json", "json", true},
		{"JSON", "json", true},
		{"yaml", "yaml", true},
		{"YAML", "yaml", true},
		{"xml", "xml", false},
		{"XML", "xml", false},
		{"", "", false},
		{" ", " ", false},
	}

	for _, test := range tests {
		normalizedInput, ok := AllowedKeys(test.input)
		if normalizedInput != test.expected || ok != test.ok {
			t.Errorf("AllowedKeys(%q) = (%q, %v); expected (%q, %v)", test.input, normalizedInput, ok, test.expected, test.ok)
		}
	}
}
