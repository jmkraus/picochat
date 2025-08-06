package utils

import (
	"testing"
)

func TestFormatList_WithBullets(t *testing.T) {
	items := []string{"first.chat", "second.chat"}
	expected := "Available history files:\n - first.chat\n - second.chat"

	result := FormatList(items, "history files", false)

	if result != expected {
		t.Errorf("Expected:\n%q\nGot:\n%q", expected, result)
	}
}

func TestFormatList_WithNumbers(t *testing.T) {
	items := []string{"model-a", "model-b", "model-c"}
	expected := "Available models:\n(01) model-a\n(02) model-b\n(03) model-c"

	result := FormatList(items, "models", true)

	if result != expected {
		t.Errorf("Expected:\n%q\nGot:\n%q", expected, result)
	}
}

func TestFormatList_Empty(t *testing.T) {
	items := []string{}
	expected := "no items found."

	result := FormatList(items, "items", false)

	if result != expected {
		t.Errorf("Expected:\n%q\nGot:\n%q", expected, result)
	}
}

func TestStripReasoning(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"Answer here <think>reasoning inside</think> end.", "Answer here  end."},
		{"Answer here </think> reasoning inside end.", " reasoning inside end."},
		{"Answer here <think>reasoning inside", "Answer here <think>reasoning inside"},
		{"Answer here", "Answer here"},
	}

	for _, tt := range tests {
		result := StripReasoning(tt.input)
		if result != tt.expected {
			t.Errorf("Expected %q, got %q for input %q", tt.expected, result, tt.input)
		}
	}
}

func TestTrimEmptyLines(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{
			"\n\nLine 1\nLine 2\n\n\n",
			"Line 1\nLine 2",
		},
		{
			"\n\n  \nLine 1\n\n  \n",
			"Line 1",
		},
		{
			"Line 1\nLine 2",
			"Line 1\nLine 2",
		},
		{
			"\n\n\n\n",
			"",
		},
		{
			"",
			"",
		},
		{
			"\n  \nHello\n\nWorld\n\n",
			"Hello\n\nWorld",
		},
	}

	for _, tt := range tests {
		result := TrimEmptyLines(tt.input)
		if result != tt.expected {
			t.Errorf("Expected %q, got %q", tt.expected, result)
		}
	}
}

func TestExtractCodeBlock(t *testing.T) {
	text := "Some explanation.\n```go\nfmt.Println(\"hi\")\n```"
	code := ExtractCodeBlock(text)
	if code != "fmt.Println(\"hi\")\n" {
		t.Errorf("Unexpected extracted code block: %q", code)
	}
}

func TestExtractCodeBlock_Empty(t *testing.T) {
	text := "No code block here"
	code := ExtractCodeBlock(text)
	if code != "" {
		t.Errorf("Expected empty string, got %q", code)
	}
}
