package command

import (
	"strings"
	"testing"

	"picochat/messages"
)

func testChatSessions(t *testing.T) (*messages.SessionManager, *messages.ChatHistory) {
	t.Helper()

	sessions := messages.NewSessionManager()
	_, history, err := sessions.Create("system prompt", 10)
	if err != nil {
		t.Fatalf("create session failed: %v", err)
	}
	return sessions, history
}

func TestHandleChatCommand_List(t *testing.T) {
	sessions, history := testChatSessions(t)
	result := handleChatCommand([]string{""}, history, sessions)

	if result.Error != nil {
		t.Fatalf("unexpected error: %v", result.Error)
	}
	if !strings.Contains(result.Output, "* 1: 1 messages") {
		t.Fatalf("unexpected chat list: %q", result.Output)
	}
	if result.SessionChanged {
		t.Fatal("listing chats should not change the active session")
	}
}

func TestHandleChatCommand_NewAndSwitch(t *testing.T) {
	sessions, history := testChatSessions(t)

	result := handleChatCommand([]string{"new"}, history, sessions)
	if result.Error != nil {
		t.Fatalf("create chat failed: %v", result.Error)
	}
	if !result.SessionChanged {
		t.Fatal("creating a chat should mark the session as changed")
	}
	if sessions.ActiveIndex() != 1 {
		t.Fatalf("active index = %d, want 1", sessions.ActiveIndex())
	}

	result = handleChatCommand([]string{"1"}, history, sessions)
	if result.Error != nil {
		t.Fatalf("switch chat failed: %v", result.Error)
	}
	if sessions.ActiveIndex() != 0 {
		t.Fatalf("active index = %d, want 0", sessions.ActiveIndex())
	}
}

func TestHandleChatCommand_CopyPrefix(t *testing.T) {
	sessions, history := testChatSessions(t)
	if err := history.AddUser("one", ""); err != nil {
		t.Fatalf("add first message failed: %v", err)
	}
	if err := history.AddUser("two", ""); err != nil {
		t.Fatalf("add second message failed: %v", err)
	}

	result := handleChatCommand([]string{"copy", "1"}, history, sessions)
	if result.Error != nil {
		t.Fatalf("copy chat failed: %v", result.Error)
	}
	if !result.SessionChanged {
		t.Fatal("copying a chat should mark the session as changed")
	}

	copied, err := sessions.Active()
	if err != nil {
		t.Fatalf("get copied chat failed: %v", err)
	}
	if copied.Len() != 2 {
		t.Fatalf("copied history length = %d, want 2", copied.Len())
	}
	if copied.GetLast().Content != "one" {
		t.Fatalf("copied last message = %q, want %q", copied.GetLast().Content, "one")
	}
}

func TestHandleChatCommand_RejectsInvalidChat(t *testing.T) {
	sessions, history := testChatSessions(t)

	result := handleChatCommand([]string{"0"}, history, sessions)
	if result.Error == nil {
		t.Fatal("expected invalid chat number to fail")
	}
	if sessions.ActiveIndex() != 0 {
		t.Fatalf("active index = %d, want 0", sessions.ActiveIndex())
	}
}

func testCopyClipboard(t *testing.T) {
	t.Helper()
	previous := writeClipboard
	writeClipboard = func(string) error { return nil }
	t.Cleanup(func() { writeClipboard = previous })
}

func TestHandleCopyCommand_DefaultAssistant(t *testing.T) {
	testCopyClipboard(t)
	h := messages.NewHistory("sys", 10)
	if err := h.AddAssistant("", "assistant answer"); err != nil {
		t.Fatalf("failed to add assistant message: %v", err)
	}

	result := handleCopyCommand("", h)
	if result.Error != nil {
		t.Fatalf("expected no error, got %v", result.Error)
	}
	if result.Warn != "" {
		t.Fatalf("unexpected warning: %q", result.Warn)
	}
	if got, want := result.Info, "Last assistant prompt copied to clipboard."; got != want {
		t.Fatalf("expected %q, got %q", want, got)
	}
}

func TestHandleCopyCommand_ByIndex(t *testing.T) {
	testCopyClipboard(t)
	h := messages.NewHistory("sys", 10)
	if err := h.AddUser("hello", ""); err != nil {
		t.Fatalf("failed to add user message: %v", err)
	}

	result := handleCopyCommand("#1", h)
	if result.Error != nil {
		t.Fatalf("expected no error, got %v", result.Error)
	}
	if got, want := result.Info, "Message #1 copied to clipboard."; got != want {
		t.Fatalf("expected %q, got %q", want, got)
	}
}

func TestHandleCopyCommand_All(t *testing.T) {
	testCopyClipboard(t)
	h := messages.NewHistory("system prompt", 10)
	if err := h.AddUser("user message", ""); err != nil {
		t.Fatalf("failed to add user message: %v", err)
	}
	if err := h.AddAssistant("internal reasoning", "assistant response"); err != nil {
		t.Fatalf("failed to add assistant message: %v", err)
	}

	result := handleCopyCommand("all", h)
	if result.Error != nil {
		t.Fatalf("expected no error, got %v", result.Error)
	}

	if got, want := result.Info, "Full conversation copied to clipboard."; got != want {
		t.Errorf("payload info = %q, want %q", got, want)
	}
}

func TestHandleCopyCommand_UnknownArg(t *testing.T) {
	h := messages.NewHistory("sys", 10)

	result := handleCopyCommand("invalid", h)
	if result.Error == nil {
		t.Fatal("expected error, got nil")
	}
	if got, want := result.Error.Error(), "unknown copy argument"; got != want {
		t.Fatalf("expected %q, got %q", want, got)
	}
}
