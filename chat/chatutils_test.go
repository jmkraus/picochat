package chat

import (
	"math"
	"testing"
	"time"
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

func TestElapsedTime(t *testing.T) {
	tests := []struct {
		name    string
		start   time.Time
		wantSec int
		wantStr string
	}{
		{
			name:    "Zero elapsed time",
			start:   time.Now(),
			wantSec: 0,
			wantStr: "0s",
		},
		{
			name:    "Exactly one minute",
			start:   time.Now().Add(-1 * time.Minute),
			wantSec: 60,
			wantStr: "1m 0s",
		},
		{
			name:    "More than one minute",
			start:   time.Now().Add(-90 * time.Second),
			wantSec: 90,
			wantStr: "1m 30s",
		},
		{
			name:    "Exactly one hour",
			start:   time.Now().Add(-1 * time.Hour),
			wantSec: 3600,
			wantStr: "60m 0s",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotSec, gotStr := elapsedTime(tt.start)
			if gotSec != tt.wantSec {
				t.Errorf("elapsedTime() got = %d, want %d", gotSec, tt.wantSec)
			}
			if gotStr != tt.wantStr {
				t.Errorf("elapsedTime() got = %s, want %s", gotStr, tt.wantStr)
			}
		})
	}
}

func TestTokenSpeed(t *testing.T) {
	tests := []struct {
		name string
		t    int
		s    string
		want float64
	}{
		{
			name: "Zero time",
			t:    0,
			s:    "example text",
			want: 0,
		},
		{
			name: "Zero tokens",
			t:    1,
			s:    "",
			want: 0,
		},
		{
			name: "Small number of tokens",
			t:    2,
			s:    "hello world",
			want: 1.3,
		},
		{
			name: "Large number of tokens",
			t:    10,
			s:    "this is a longer text with multiple words",
			want: 1.0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tokenSpeed(tt.t, tt.s)
			if !equalFloat64(got, tt.want, 0.0001) {
				t.Errorf("tokenSpeed(%s) got = %v, want %v", tt.s, got, tt.want)
			}
		})
	}
}

// Helper function to compare float64 values with a tolerance
func equalFloat64(a, b, tolerance float64) bool {
	return math.Abs(a-b) <= tolerance
}
