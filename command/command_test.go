package command

import (
	"os"
	"path/filepath"

	"fmt"
	"picochat/messages"
	"picochat/paths"
	"strings"
	"testing"
)

func TestHandleClear(t *testing.T) {
	h := messages.NewHistory("initial system prompt", 50)
	h.AddUser("hello", "")

	result := HandleCommand("/clear", h, strings.NewReader(""))
	if !h.IsEmpty() {
		t.Error("expected history to be cleared except system prompt")
	}
	if !strings.Contains(result.Info, "cleared") {
		t.Errorf("unexpected output: %s", result.Info)
	}
}

func TestHandleHelp(t *testing.T) {
	h := messages.NewHistory("prompt", 50)
	result := HandleCommand("/help", h, strings.NewReader(""))
	if !strings.Contains(result.Output, "/save") {
		t.Errorf("expected help to contain /save, got: %s", result.Info)
	}
}

func TestHandleLoad_WithFilename(t *testing.T) {
	const dummyChat = `[
  {
    "role": "system",
    "content": "You are a LLM"
  },
  {
    "role": "user",
    "content": "Hello, are you there?"
  },
  {
    "role": "assistant",
    "content": "How can I assist you today?"
  }
]`

	tmpDir := t.TempDir()
	restore := paths.OverrideHistoryPath(tmpDir)
	t.Cleanup(restore)

	err := os.WriteFile(filepath.Join(tmpDir, "dummy.chat"), []byte(dummyChat), 0644)
	if err != nil {
		t.Fatalf("failed to write dummy.chat: %v", err)
	}

	h := messages.NewHistory("system prompt", 50)
	input := strings.NewReader("dummy.chat\n")

	result := HandleCommand("/load", h, input)

	if result.Error != nil {
		t.Fatalf("expected no error, got: %v", result.Error)
	}
	if result.Info != "History file \"dummy.chat\" loaded." {
		t.Fatalf("unexpected info: %q", result.Info)
	}

	loaded := h.Get()
	if len(loaded) != 3 {
		t.Fatalf("expected 3 messages, got %d", len(loaded))
	}

	if loaded[0].Role != messages.RoleSystem || loaded[0].Content != "You are a LLM" {
		t.Fatalf("unexpected message[0]: %#v", loaded[0])
	}
	if loaded[1].Role != messages.RoleUser || loaded[1].Content != "Hello, are you there?" {
		t.Fatalf("unexpected message[1]: %#v", loaded[1])
	}
	if loaded[2].Role != messages.RoleAssistant || loaded[2].Content != "How can I assist you today?" {
		t.Fatalf("unexpected message[2]: %#v", loaded[2])
	}
}

func TestHandleCommand_Paste_UsesRuneCount(t *testing.T) {
	prevReadClipboard := readClipboard
	t.Cleanup(func() {
		readClipboard = prevReadClipboard
	})

	readClipboard = func() (string, error) {
		return "Hello, 世界", nil
	}

	history := messages.NewHistory("sys", 10)
	result := HandleCommand("/paste", history, strings.NewReader(""))

	if result.Error != nil {
		t.Fatalf("expected no error, got %v", result.Error)
	}
	if result.Pasted != "Hello, 世界" {
		t.Fatalf("unexpected pasted payload: %q", result.Pasted)
	}
	wantInfo := "Pasted 9 characters from clipboard."
	if result.Info != wantInfo {
		t.Fatalf("info = %q, want %q", result.Info, wantInfo)
	}
}

func TestHandleCommand_Paste_UnknownTemplateError(t *testing.T) {
	prevReadClipboard := readClipboard
	t.Cleanup(func() {
		readClipboard = prevReadClipboard
	})

	readClipboard = func() (string, error) { return "Hello", nil }

	history := messages.NewHistory("sys", 10)
	result := HandleCommand("/paste unknown", history, strings.NewReader(""))

	if result.Error == nil {
		t.Fatal("expected error, got nil")
	}
	if !strings.Contains(result.Error.Error(), `template key "unknown" not found`) {
		t.Fatalf("unexpected error: %v", result.Error)
	}
}

func TestHandleCommand_Paste_ReadClipboardError(t *testing.T) {
	prevReadClipboard := readClipboard
	t.Cleanup(func() {
		readClipboard = prevReadClipboard
	})

	readClipboard = func() (string, error) {
		return "", fmt.Errorf("clipboard read failed")
	}

	history := messages.NewHistory("sys", 10)
	result := HandleCommand("/paste", history, strings.NewReader(""))

	if result.Error == nil {
		t.Fatal("expected error, got nil")
	}
	if !strings.Contains(result.Error.Error(), "clipboard read failed") {
		t.Fatalf("unexpected error: %v", result.Error)
	}
}

func TestHandleCommand_Save_ExistingFile_NoOverwrite(t *testing.T) {
	tmpDir := t.TempDir()
	restore := paths.OverrideHistoryPath(tmpDir)
	t.Cleanup(restore)

	h := messages.NewHistory("system prompt", 50)
	if err := h.AddUser("first", ""); err != nil {
		t.Fatalf("add user failed: %v", err)
	}

	existingName := "existing"
	if _, err := messages.SaveHistoryToFile(existingName, h.Get(), false); err != nil {
		t.Fatalf("initial save failed: %v", err)
	}

	result := HandleCommand("/save "+existingName, h, strings.NewReader("n\n"))
	if result.Error != nil {
		t.Fatalf("expected no error, got: %v", result.Error)
	}
	if result.Warn != "Save canceled." {
		t.Fatalf("unexpected warn: %q", result.Warn)
	}
}

func TestHandleCommand_Save_ExistingFile_Overwrite(t *testing.T) {
	tmpDir := t.TempDir()
	restore := paths.OverrideHistoryPath(tmpDir)
	t.Cleanup(restore)

	h := messages.NewHistory("system prompt", 50)
	if err := h.AddUser("first", ""); err != nil {
		t.Fatalf("add user failed: %v", err)
	}

	existingName := "existing"
	if _, err := messages.SaveHistoryToFile(existingName, h.Get(), false); err != nil {
		t.Fatalf("initial save failed: %v", err)
	}

	result := HandleCommand("/save "+existingName, h, strings.NewReader("y\n"))
	if result.Error != nil {
		t.Fatalf("expected no error, got: %v", result.Error)
	}
	if !strings.Contains(result.Info, "History saved as file") {
		t.Fatalf("unexpected info: %q", result.Info)
	}
}

func TestParseCommandArgs(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		wantCmd string
		wantArg string
	}{
		{
			name:    "empty input",
			input:   "",
			wantCmd: "",
			wantArg: "",
		},
		{
			name:    "whitespace only",
			input:   "   \t   ",
			wantCmd: "",
			wantArg: "",
		},
		{
			name:    "copy alias",
			input:   "/c",
			wantCmd: "copy",
			wantArg: "",
		},
		{
			name:    "paste alias",
			input:   "/v",
			wantCmd: "paste",
			wantArg: "",
		},
		{
			name:    "help alias",
			input:   "/?",
			wantCmd: "help",
			wantArg: "",
		},
		{
			name:    "hallo alias",
			input:   "/hallo",
			wantCmd: "hello",
			wantArg: "",
		},
		{
			name:    "normalization lowercase and trim slash",
			input:   "/MoDeLs  #2",
			wantCmd: "models",
			wantArg: "#2",
		},
		{
			name:    "unknown command passthrough normalized",
			input:   "/FoObAr keep This ARG",
			wantCmd: "foobar",
			wantArg: "keep This ARG",
		},
		{
			name:    "without slash",
			input:   "copy think",
			wantCmd: "copy",
			wantArg: "think",
		},
		{
			name:    "args keep content but normalize spacing",
			input:   "/set   temperature=0.7    now",
			wantCmd: "set",
			wantArg: "temperature=0.7 now",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd, arg := parseCommandArgs(tt.input)
			if cmd != tt.wantCmd || arg != tt.wantArg {
				t.Fatalf("parseCommandArgs(%q) = (%q, %q), want (%q, %q)", tt.input, cmd, arg, tt.wantCmd, tt.wantArg)
			}
		})
	}
}
