package command_test

import (
	"picochat/command"
	"picochat/messages"
	"strings"
	"testing"
)

func TestHandleClear(t *testing.T) {
	h := messages.NewHistory("initial system prompt", 50)
	h.Add("user", "hello", "")

	result := command.HandleCommand("/clear", h, strings.NewReader(""))
	if !h.IsEmpty() {
		t.Error("expected history to be cleared except system prompt")
	}
	if !strings.Contains(result.Output, "cleared") {
		t.Errorf("unexpected output: %s", result.Output)
	}
}

func TestHandleHelp(t *testing.T) {
	h := messages.NewHistory("prompt", 50)
	result := command.HandleCommand("/help", h, strings.NewReader(""))
	if !strings.Contains(result.Output, "/save") {
		t.Errorf("expected help to contain /save, got: %s", result.Output)
	}
}

func TestHandleLoad_WithFilename(t *testing.T) {
	h := messages.NewHistory("system prompt", 50)
	input := strings.NewReader("dummy.chat\n")

	result := command.HandleCommand("/load", h, input)

	if !strings.Contains(result.Output, "failed") && !strings.Contains(result.Output, "success") {
		t.Errorf("Unexpected load result: %s", result.Output)
	}
}
