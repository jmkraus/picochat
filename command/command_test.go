package command_test

import (
	"picochat/command"
	"picochat/messages"
	"strings"
	"testing"
)

func TestHandleClear(t *testing.T) {
	h := messages.NewHistory("initial system prompt", 50)
	h.AddUser("hello", "")

	result := command.HandleCommand("/clear", h, strings.NewReader(""))
	if !h.IsEmpty() {
		t.Error("expected history to be cleared except system prompt")
	}
	if !strings.Contains(result.Info, "cleared") {
		t.Errorf("unexpected output: %s", result.Info)
	}
}

func TestHandleHelp(t *testing.T) {
	h := messages.NewHistory("prompt", 50)
	result := command.HandleCommand("/help", h, strings.NewReader(""))
	if !strings.Contains(result.Info, "/save") {
		t.Errorf("expected help to contain /save, got: %s", result.Info)
	}
}

func TestHandleLoad_WithFilename(t *testing.T) {
	h := messages.NewHistory("system prompt", 50)
	input := strings.NewReader("dummy.chat\n")

	result := command.HandleCommand("/load", h, input)

	if !strings.Contains(result.Info, "failed") && !strings.Contains(result.Info, "success") {
		t.Errorf("Unexpected load result: %s", result.Info)
	}
}
