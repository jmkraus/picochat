package main

import (
	"encoding/json"
	"os"
	"path/filepath"
	"picochat/args"
	"picochat/command"
	"picochat/config"
	"picochat/messages"
	"picochat/paths"
	"strings"
	"testing"
)

func newTestInstance() *Instance {
	cfg := &config.Config{
		URL:       "://invalid-url",
		Model:     "test-model",
		Prompt:    "system prompt",
		OutputFmt: "plain",
	}
	history := messages.NewHistory(cfg.Prompt, 10)
	manager := messages.NewSessionManager()
	if _, err := manager.Add(history); err != nil {
		panic(err)
	}
	return &Instance{
		Config:   cfg,
		History:  history,
		Sessions: manager,
		Quiet:    true,
	}
}

func TestSendPrompt_AppendsUserAndClearsImagePath(t *testing.T) {
	instance := newTestInstance()

	tmpDir := t.TempDir()
	imagePath := filepath.Join(tmpDir, "image.jpg")
	if err := os.WriteFile(imagePath, []byte("dummy-image"), 0644); err != nil {
		t.Fatalf("write image file failed: %v", err)
	}
	instance.Config.ImagePath = imagePath

	sendPrompt(instance, "hello")

	if instance.Config.ImagePath != "" {
		t.Fatalf("expected image path to be cleared, got %q", instance.Config.ImagePath)
	}
	if instance.History.Len() != 2 {
		t.Fatalf("expected history length 2, got %d", instance.History.Len())
	}

	last := instance.History.GetLast()
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
	instance := newTestInstance()
	instance.Config.ImagePath = "/path/does/not/exist.jpg"

	sendPrompt(instance, "hello")

	if instance.History.Len() != 1 {
		t.Fatalf("expected history length 1, got %d", instance.History.Len())
	}
	if instance.Config.ImagePath == "" {
		t.Fatal("expected image path to remain set after add-user failure")
	}
}

func TestRunChat_InvalidURLDoesNotAppendAssistant(t *testing.T) {
	instance := newTestInstance()
	if err := instance.History.AddUser("hello", ""); err != nil {
		t.Fatalf("add user failed: %v", err)
	}
	before := instance.History.Len()

	runChat(instance)

	after := instance.History.Len()
	if after != before {
		t.Fatalf("history length changed on failed runChat: before=%d after=%d", before, after)
	}
}

func TestRetryPrompt_InvalidURLDoesNotAppendAssistant(t *testing.T) {
	instance := newTestInstance()
	if err := instance.History.AddUser("hello", ""); err != nil {
		t.Fatalf("add user failed: %v", err)
	}
	before := instance.History.Len()

	retryPrompt(instance)

	after := instance.History.Len()
	if after != before {
		t.Fatalf("history length changed on failed retryPrompt: before=%d after=%d", before, after)
	}
}

func TestSyncActiveHistory(t *testing.T) {
	instance := newTestInstance()
	original := instance.History

	_, active, err := instance.Sessions.Create("new system prompt", 10)
	if err != nil {
		t.Fatalf("create session failed: %v", err)
	}
	if instance.History != original {
		t.Fatal("creating a session unexpectedly changed instance history")
	}

	if err := getActiveHistory(instance); err != nil {
		t.Fatalf("sync active history failed: %v", err)
	}
	if instance.History != active {
		t.Fatal("instance history was not updated to active session")
	}
}

func TestSyncActiveHistoryWithoutManager(t *testing.T) {
	instance := &Instance{}
	if err := getActiveHistory(instance); err == nil {
		t.Fatal("expected missing session manager error")
	}
}

func TestInitInstanceFromArgs_InitializesSessionManager(t *testing.T) {
	previous := struct {
		configPath  string
		historyFile string
		quiet       bool
		model       string
		showVersion bool
		image       string
		output      string
		schema      string
	}{
		configPath:  *args.ConfigPath,
		historyFile: *args.HistoryFile,
		quiet:       *args.Quiet,
		model:       *args.Model,
		showVersion: *args.ShowVersion,
		image:       *args.Image,
		output:      *args.Output,
		schema:      *args.Schema,
	}
	t.Cleanup(func() {
		*args.ConfigPath = previous.configPath
		*args.HistoryFile = previous.historyFile
		*args.Quiet = previous.quiet
		*args.Model = previous.model
		*args.ShowVersion = previous.showVersion
		*args.Image = previous.image
		*args.Output = previous.output
		*args.Schema = previous.schema
	})

	*args.ConfigPath = ""
	*args.HistoryFile = ""
	*args.Quiet = true
	*args.Model = ""
	*args.ShowVersion = false
	*args.Image = ""
	*args.Output = ""
	*args.Schema = ""

	showVersion, instance, _, err := initInstanceFromArgs()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if showVersion {
		t.Fatal("expected normal startup, got version response")
	}
	if instance == nil {
		t.Fatal("expected initialized instance")
	}
	if instance.Sessions == nil {
		t.Fatal("expected session manager")
	}
	if got := len(instance.Sessions.List()); got != 1 {
		t.Fatalf("session count = %d, want 1", got)
	}
	if instance.Sessions.ActiveIndex() != 0 {
		t.Fatalf("active session index = %d, want 0", instance.Sessions.ActiveIndex())
	}
	active, err := instance.Sessions.Active()
	if err != nil {
		t.Fatalf("get active session failed: %v", err)
	}
	if active != instance.History {
		t.Fatal("instance history is not the active registered session")
	}
}

func TestInitInstanceFromArgs_ShowVersionShortCircuit(t *testing.T) {
	prevShowVersion := *args.ShowVersion
	t.Cleanup(func() {
		*args.ShowVersion = prevShowVersion
	})

	*args.ShowVersion = true

	showVersion, instance, warnings, err := initInstanceFromArgs()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !showVersion {
		t.Fatal("expected showVersion to be true")
	}
	if instance != nil {
		t.Fatalf("expected nil instance, got %+v", instance)
	}
	if warnings != nil {
		t.Fatalf("expected nil warnings, got %v", warnings)
	}
}

func TestInitInstanceFromArgs_LoadedHistoryContextAndChatNew(t *testing.T) {
	historyDir := t.TempDir()
	restoreHistoryPath := paths.OverrideHistoryPath(historyDir)
	t.Cleanup(restoreHistoryPath)

	loadedMessages := []messages.Message{{Role: messages.RoleSystem, Content: "system prompt"}}
	for i := 0; i < 24; i++ {
		loadedMessages = append(loadedMessages, messages.Message{
			Role:    messages.RoleUser,
			Content: "user message",
		})
	}
	data, err := json.Marshal(loadedMessages)
	if err != nil {
		t.Fatalf("marshal history failed: %v", err)
	}
	if err := os.WriteFile(filepath.Join(historyDir, "loaded.chat"), data, 0644); err != nil {
		t.Fatalf("write history failed: %v", err)
	}

	previousHistoryFile := *args.HistoryFile
	previousShowVersion := *args.ShowVersion
	previousQuiet := *args.Quiet
	t.Cleanup(func() {
		*args.HistoryFile = previousHistoryFile
		*args.ShowVersion = previousShowVersion
		*args.Quiet = previousQuiet
	})

	*args.HistoryFile = "loaded"
	*args.ShowVersion = false
	*args.Quiet = true

	_, instance, _, err := initInstanceFromArgs()
	if err != nil {
		t.Fatalf("initialize from loaded history failed: %v", err)
	}

	cfg, _, err := config.Get()
	if err != nil {
		t.Fatalf("get config failed: %v", err)
	}
	wantContext := len(loadedMessages) + 1
	if cfg.Context > wantContext {
		wantContext = cfg.Context
	}
	if got := instance.History.MaxCtx(); got != wantContext {
		t.Fatalf("loaded history context = %d, want %d", got, wantContext)
	}

	result := command.HandleCommand("/chat new", instance.History, instance.Sessions, strings.NewReader(""))
	if result.Error != nil {
		t.Fatalf("/chat new failed: %v", result.Error)
	}
	if !result.SessionChanged {
		t.Fatal("/chat new did not report a session change")
	}

	active, err := instance.Sessions.Active()
	if err != nil {
		t.Fatalf("get new active session failed: %v", err)
	}
	if got := active.MaxCtx(); got != wantContext {
		t.Fatalf("new session context = %d, want %d", got, wantContext)
	}
}
