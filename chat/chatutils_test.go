package chat

import (
	"testing"
)

func TestSplitReasoning(t *testing.T) {
	tests := []struct {
		input             string
		expectedReasoning string
		expectedContent   string
	}{
		{
			input:             "Answer here <think>reasoning inside</think> end.",
			expectedReasoning: "reasoning inside",
			expectedContent:   "Answer here  end.",
		},
		{
			input:             "Answer here </think> reasoning inside end.",
			expectedReasoning: "Answer here ",
			expectedContent:   " reasoning inside end.",
		},
		{
			input:             "Answer here <think>reasoning inside",
			expectedReasoning: "",
			expectedContent:   "Answer here <think>reasoning inside",
		},
		{
			input:             "Answer here",
			expectedReasoning: "",
			expectedContent:   "Answer here",
		},
	}

	for _, tt := range tests {
		reasoning, content := splitReasoning(tt.input)

		if reasoning != tt.expectedReasoning || content != tt.expectedContent {
			t.Errorf(
				"For input %q expected reasoning=%q, content=%q; got reasoning=%q, content=%q",
				tt.input,
				tt.expectedReasoning,
				tt.expectedContent,
				reasoning,
				content,
			)
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
