package command

import (
	"fmt"
	"picochat/envs"
	"picochat/messages"
	"picochat/vartypes"
	"reflect"
	"strings"
	"testing"
)

func TestExtractCodeBlock(t *testing.T) {
	text := "Some explanation.\n```go\nfmt.Println(\"hi\")\n```"
	code, found := extractCodeBlock(text)
	if code != "fmt.Println(\"hi\")\n" {
		t.Errorf("Unexpected extracted code block: %q", code)
	}
	if code == "fmt.Println(\"hi\")\n" && found == false {
		t.Errorf("'Found' flag for ExtractCode reported False, but should be True")
	}
}

func TestExtractCodeBlock_Empty(t *testing.T) {
	text := "No code block here"
	code, found := extractCodeBlock(text)
	if code != "" {
		t.Errorf("Expected empty string, got %q", code)
	}
	if code == "" && found == true {
		t.Errorf("'Found' flag for ExtractCode reported True, but should be False")
	}
}

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
		key, val, err := parseKeyVal(tt.input)
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
		_, _, err := parseKeyVal(input)
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
		_, _, err := parseKeyVal(input)

		if spec.Runtime && err != nil {
			t.Errorf("expected runtime field %q to parse, got error: %v", spec.Field, err)
		}
		if !spec.Runtime && err == nil {
			t.Errorf("expected non-runtime field %q to be rejected", spec.Field)
		}
	}
}

func TestParseIndex_Valid(t *testing.T) {
	index, err := parseIndex("7")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if index != 7 {
		t.Fatalf("expected 7, got %d", index)
	}
}

func TestParseIndex_Invalid(t *testing.T) {
	_, err := parseIndex("abc")
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if got, want := err.Error(), "value not an integer"; got != want {
		t.Fatalf("expected %q, got %q", want, got)
	}
}

func TestResolveCopyPayload_DefaultAssistant(t *testing.T) {
	h := messages.NewHistory("sys", 10)
	if err := h.AddAssistant("", "assistant answer"); err != nil {
		t.Fatalf("failed to add assistant message: %v", err)
	}

	payload, err := resolveCopyPayload("", h)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if got, want := payload.Text, "assistant answer"; got != want {
		t.Fatalf("expected %q, got %q", want, got)
	}
	if got, want := payload.Info, "Last assistant prompt written to clipboard."; got != want {
		t.Fatalf("expected %q, got %q", want, got)
	}
}

func TestResolveCopyPayload_ByIndex(t *testing.T) {
	h := messages.NewHistory("sys", 10)
	if err := h.AddUser("hello", ""); err != nil {
		t.Fatalf("failed to add user message: %v", err)
	}

	payload, err := resolveCopyPayload("#1", h)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if got, want := payload.Info, "Message #1 written to clipboard"; got != want {
		t.Fatalf("expected %q, got %q", want, got)
	}
	if !strings.Contains(payload.Text, "(1:user)") || !strings.Contains(payload.Text, "hello") {
		t.Fatalf("unexpected payload text: %q", payload.Text)
	}
}

func TestResolveCopyPayload_UnknownArg(t *testing.T) {
	h := messages.NewHistory("sys", 10)

	_, err := resolveCopyPayload("invalid", h)
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if got, want := err.Error(), "unknown copy argument"; got != want {
		t.Fatalf("expected %q, got %q", want, got)
	}
}
