package command_test

import (
	"picochat/command"
	"picochat/types"
	"strings"
	"testing"
)

func TestHandleClear(t *testing.T) {
	h := types.NewHistory("initial system prompt", 50)
	h.Add("user", "hello")

	result := command.Handle("/clear", h, strings.NewReader(""))
	if !h.IsEmpty() {
		t.Error("expected history to be cleared except system prompt")
	}
	if !strings.Contains(result.Output, "cleared") {
		t.Errorf("unexpected output: %s", result.Output)
	}
}

func TestHandleHelp(t *testing.T) {
	h := types.NewHistory("prompt", 50)
	result := command.Handle("/help", h, strings.NewReader(""))
	if !strings.Contains(result.Output, "/save") {
		t.Errorf("expected help to contain /save, got: %s", result.Output)
	}
}

func TestHandleLoad_WithFilename(t *testing.T) {
	h := types.NewHistory("system prompt", 50)
	input := strings.NewReader("dummy.chat\n")

	result := command.Handle("/load", h, input)

	// hier kannst du prüfen, ob die Meldung korrekt ist — oder Fehler bei ungültiger Datei
	if !strings.Contains(result.Output, "failed") && !strings.Contains(result.Output, "success") {
		t.Errorf("Unexpected load result: %s", result.Output)
	}
}
