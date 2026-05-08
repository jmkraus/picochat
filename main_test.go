package main

import (
	"os"
	"path/filepath"
	"picochat/args"
	"picochat/config"
	"picochat/messages"
	"testing"
)

func newTestSession() *Session {
	cfg := &config.Config{
		URL:      "://invalid-url",
		Model:    "test-model",
		Prompt:   "system prompt",
		OutputFmt: "plain",
	}
	history := messages.NewHistory(cfg.Prompt, 10)
	return &Session{
		Config:  cfg,
		History: history,
		Quiet:   true,
	}
}

func TestSendPrompt_AppendsUserAndClearsImagePath(t *testing.T) {
	session := newTestSession()

	tmpDir := t.TempDir()
	imagePath := filepath.Join(tmpDir, "image.jpg")
	if err := os.WriteFile(imagePath, []byte("dummy-image"), 0644); err != nil {
		t.Fatalf("write image file failed: %v", err)
	}
	session.Config.ImagePath = imagePath

	sendPrompt(session, "hello")

	if session.Config.ImagePath != "" {
		t.Fatalf("expected image path to be cleared, got %q", session.Config.ImagePath)
	}
	if session.History.Len() != 2 {
		t.Fatalf("expected history length 2, got %d", session.History.Len())
	}

	last := session.History.GetLast()
	if last.Role != messages.RoleUser {
		t.Fatalf("expected last role %q, got %q", messages.RoleUser, last.Role)
	}
	if last.Content != "hello" {
		t.Fatalf("expected last content %q, got %q", "hello", last.Content)
	}
	if len(last.Images) != 1 {
		t.Fatalf("expected one image payload, got %d", len(last.Images))
	}
}

func TestSendPrompt_InvalidImageDoesNotAppend(t *testing.T) {
	session := newTestSession()
	session.Config.ImagePath = "/path/does/not/exist.jpg"

	sendPrompt(session, "hello")

	if session.History.Len() != 1 {
		t.Fatalf("expected history length 1, got %d", session.History.Len())
	}
	if session.Config.ImagePath == "" {
		t.Fatal("expected image path to remain set after add-user failure")
	}
}

func TestRunChat_InvalidURLDoesNotAppendAssistant(t *testing.T) {
	session := newTestSession()
	if err := session.History.AddUser("hello", ""); err != nil {
		t.Fatalf("add user failed: %v", err)
	}
	before := session.History.Len()

	runChat(session)

	after := session.History.Len()
	if after != before {
		t.Fatalf("history length changed on failed runChat: before=%d after=%d", before, after)
	}
}

func TestRepeatPrompt_InvalidURLDoesNotAppendAssistant(t *testing.T) {
	session := newTestSession()
	if err := session.History.AddUser("hello", ""); err != nil {
		t.Fatalf("add user failed: %v", err)
	}
	before := session.History.Len()

	repeatPrompt(session)

	after := session.History.Len()
	if after != before {
		t.Fatalf("history length changed on failed repeatPrompt: before=%d after=%d", before, after)
	}
}

func TestInitSessionFromArgs_ShowVersionShortCircuit(t *testing.T) {
	prevShowVersion := *args.ShowVersion
	t.Cleanup(func() {
		*args.ShowVersion = prevShowVersion
	})

	*args.ShowVersion = true

	showVersion, session, warnings, err := initSessionFromArgs()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !showVersion {
		t.Fatal("expected showVersion to be true")
	}
	if session != nil {
		t.Fatalf("expected nil session, got %+v", session)
	}
	if warnings != nil {
		t.Fatalf("expected nil warnings, got %v", warnings)
	}
}
