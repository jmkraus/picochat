package format

import (
	"bytes"
	"strings"
	"testing"

	"picochat/chat"
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
		{"json-pretty", "json-pretty", true},
		{"JSON-PrettY", "json-pretty", true},
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

func testResult() *chat.ChatResult {
	return &chat.ChatResult{
		Output:   "Hello World",
		Elapsed:  "00:10",
		TokensPS: 12.3,
		// Model:    "test-model",
	}
}

func TestRenderResult_JSON(t *testing.T) {
	var buf bytes.Buffer

	err := RenderResult(&buf, testResult(), "json", false)
	if err != nil {
		t.Fatalf("RenderResult returned error: %v", err)
	}

	out := buf.String()

	if !strings.Contains(out, `"output":"Hello World"`) {
		t.Errorf("JSON output missing field 'output': %s", out)
	}

	if !strings.Contains(out, `"tokens_per_sec":12.3`) {
		t.Errorf("JSON output missing field 'tokens_per_sec': %s", out)
	}
}

func TestRenderResult_JSONPretty(t *testing.T) {
	var buf bytes.Buffer

	err := RenderResult(&buf, testResult(), "json-pretty", false)
	if err != nil {
		t.Fatalf("RenderResult returned error: %v", err)
	}

	out := buf.String()

	if !strings.Contains(out, "\n  \"output\": \"Hello World\"") {
		t.Errorf("Pretty JSON not indented as expected:\n%s", out)
	}
}

func TestRenderResult_YAML(t *testing.T) {
	var buf bytes.Buffer

	err := RenderResult(&buf, testResult(), "yaml", false)
	if err != nil {
		t.Fatalf("RenderResult returned error: %v", err)
	}

	out := buf.String()

	if !strings.Contains(out, "output: Hello World") {
		t.Errorf("YAML output missing 'output': %s", out)
	}

	if !strings.Contains(out, "tokens_per_sec: 12.3") {
		t.Errorf("YAML output missing 'tokens_per_sec': %s", out)
	}
}

func TestRenderResult_UnknownFormat(t *testing.T) {
	var buf bytes.Buffer

	err := RenderResult(&buf, testResult(), "unknown", false)
	if err == nil {
		t.Fatal("expected error for unknown format, got nil")
	}
}
