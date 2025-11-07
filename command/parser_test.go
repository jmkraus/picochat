package command

import (
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
		key, val, err := parseArgs(tt.input)
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
		"context=notanint",
	}

	for _, input := range tests {
		_, _, err := parseArgs(input)
		if err == nil {
			t.Errorf("expected error for input %q, got nil", input)
		}
	}
}
