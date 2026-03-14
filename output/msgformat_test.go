package output

import (
	"picochat/messages"
	"regexp"
	"testing"
)

var ansiRE = regexp.MustCompile(`\x1b\[[0-9;]*m`)

func stripANSI(s string) string {
	return ansiRE.ReplaceAllString(s, "")
}

func TestFormatMessage_NoHeaderNoColor(t *testing.T) {
	msg := messages.Message{Role: messages.RoleAssistant, Content: "hello world"}

	got := FormatMessage(msg, 0, false, false)
	if got != "hello world" {
		t.Fatalf("expected %q, got %q", "hello world", got)
	}
}

func TestFormatMessage_WithHeader(t *testing.T) {
	msg := messages.Message{Role: messages.RoleUser, Content: "hello"}

	got := stripANSI(FormatMessage(msg, 2, true, false))
	want := "(2:user)\nhello"
	if got != want {
		t.Fatalf("expected %q, got %q", want, got)
	}
}

func TestFormatConversation(t *testing.T) {
	msgs := []messages.Message{
		{Role: messages.RoleSystem, Content: "sys"},
		{Role: messages.RoleUser, Content: "hi"},
		{Role: messages.RoleAssistant, Content: "hey"},
	}

	got := stripANSI(FormatConversation(msgs))
	want := "(0:system)\nsys\n\n(1:user)\nhi\n\n(2:assistant)\nhey\n\n"
	if got != want {
		t.Fatalf("expected %q, got %q", want, got)
	}
}
