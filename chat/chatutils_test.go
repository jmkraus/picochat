package chat

import (
	"testing"
)

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
		result := stripReasoning(tt.input)
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
		result := trimEmptyLines(tt.input)
		if result != tt.expected {
			t.Errorf("Expected %q, got %q", tt.expected, result)
		}
	}
}
