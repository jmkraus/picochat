package convert

import (
	"fmt"
	"picochat/envs"
	"picochat/vartypes"
	"reflect"
	"testing"
)

func TestParseArgs_Valid(t *testing.T) {
	tests := []struct {
		input    string
		wantKey  string
		wantType string
	}{
		{"temperature=0.7", "temperature", "float64"},
		{"top_p=0.9", "top_p", "float64"},
		{"context=42", "context", "int"},
		{"model=gpt-4", "model", "string"},
	}

	for _, tt := range tests {
		key, val, err := ParseKeyVal(tt.input)
		if err != nil {
			t.Errorf("unexpected error for input %q: %v", tt.input, err)
			continue
		}
		if key != tt.wantKey {
			t.Errorf("expected key %q, got %q", tt.wantKey, key)
		}
		if reflect.TypeOf(val).String() != tt.wantType {
			t.Errorf("expected type %q, got %T", tt.wantType, val)
		}
	}
}

func TestParseArgs_Invalid(t *testing.T) {
	tests := []string{
		"notakeyvalue",
		"=nokey",
		"novalue=",
		"temperature=abc",
		"unknown=123",
		"url=http://example.com",
		"quiet=true",
		"context=notanint",
	}

	for _, input := range tests {
		_, _, err := ParseKeyVal(input)
		if err == nil {
			t.Errorf("expected error for input %q, got nil", input)
		}
	}
}

// Ensures envs runtime schema and ParseKeyVal stay in sync.
func TestParseKeyVal_RuntimeFieldsFromEnvSpec(t *testing.T) {
	for _, spec := range envs.ConfigEnvVars {
		var sample string
		switch spec.Type {
		case vartypes.VarString:
			sample = "abc"
		case vartypes.VarInt:
			sample = "7"
		case vartypes.VarFloat:
			sample = "0.5"
		case vartypes.VarBool:
			sample = "true"
		default:
			t.Fatalf("unsupported spec type %q for field %q", spec.Type, spec.Field)
		}

		input := fmt.Sprintf("%s=%s", spec.Field, sample)
		_, _, err := ParseKeyVal(input)

		if spec.Runtime && err != nil {
			t.Errorf("expected runtime field %q to parse, got error: %v", spec.Field, err)
		}
		if !spec.Runtime && err == nil {
			t.Errorf("expected non-runtime field %q to be rejected", spec.Field)
		}
	}
}
