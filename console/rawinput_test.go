package console

import (
	"io"
	"os"
	"strings"
	"testing"
)

// --- insertCharAt --------------------------------------------------

func TestInsertCharAt(t *testing.T) {
	tests := []struct {
		name      string
		input     string
		cursorPos int
		char      rune
		wantLine  string
		wantPos   int
	}{
		{"insert middle", "Helo", 2, 'l', "Hello", 3},
		{"insert start", "ello", 0, 'H', "Hello", 1},
		{"insert end", "Hell", 4, 'o', "Hello", 5},
		{"insert beyond length", "Hi", 99, '!', "Hi!", 3},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			line := []rune(tt.input)
			gotLine, gotPos := insertCharAt(line, tt.cursorPos, tt.char)
			if string(gotLine) != tt.wantLine {
				t.Errorf("insertCharAt() = %q, want %q", string(gotLine), tt.wantLine)
			}
			if gotPos != tt.wantPos {
				t.Errorf("cursorPos = %d, want %d", gotPos, tt.wantPos)
			}
		})
	}
}

// --- deleteCharAt --------------------------------------------------

func TestDeleteCharAt(t *testing.T) {
	tests := []struct {
		name      string
		input     string
		cursorPos int
		wantLine  string
		wantPos   int
	}{
		{"delete middle", "Hello", 3, "Helo", 2},
		{"delete end", "Hello", 5, "Hell", 4},
		{"delete start", "Hello", 1, "ello", 0},
		{"delete empty", "", 0, "", 0},
		{"delete pos 0", "Hello", 0, "Hello", 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			line := []rune(tt.input)
			gotLine, gotPos := deleteCharAt(line, tt.cursorPos)
			if string(gotLine) != tt.wantLine {
				t.Errorf("deleteCharAt() = %q, want %q", string(gotLine), tt.wantLine)
			}
			if gotPos != tt.wantPos {
				t.Errorf("cursorPos = %d, want %d", gotPos, tt.wantPos)
			}
		})
	}
}

// --- updateCurrentLine ---------------------------------------------

func captureOutput(f func()) string {
	old := os.Stdout // save original
	r, w, _ := os.Pipe()
	os.Stdout = w

	f() // execute the function we want to test

	_ = w.Close()
	os.Stdout = old // restore

	var buf strings.Builder
	_, _ = io.Copy(&buf, r)
	_ = r.Close()

	return buf.String()
}

func TestUpdateCurrentLine(t *testing.T) {
	out := captureOutput(func() {
		updateCurrentLine([]rune("Hello"), true, 5)
	})
	if !strings.Contains(out, ">>> Hello") {
		t.Errorf("expected prompt + text, got %q", out)
	}

	out = captureOutput(func() {
		updateCurrentLine([]rune("World"), false, 3)
	})
	if !strings.Contains(out, "World") {
		t.Errorf("expected text without prompt, got %q", out)
	}
}
