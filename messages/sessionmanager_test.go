package messages

import "testing"

func TestSessionManagerAdd(t *testing.T) {
	manager := NewSessionManager()
	history := NewHistory("loaded system prompt", 10)

	index, err := manager.Add(history)
	if err != nil {
		t.Fatalf("add session failed: %v", err)
	}
	if index != 0 {
		t.Fatalf("added session index = %d, want 0", index)
	}
	if manager.ActiveIndex() != index {
		t.Fatalf("active index = %d, want %d", manager.ActiveIndex(), index)
	}
	active, err := manager.Active()
	if err != nil {
		t.Fatalf("get active session failed: %v", err)
	}
	if active != history {
		t.Fatal("expected added history to be active")
	}
}

func TestSessionManagerAddRejectsNil(t *testing.T) {
	manager := NewSessionManager()

	if _, err := manager.Add(nil); err == nil {
		t.Fatal("expected adding nil history to fail")
	}
	if got := len(manager.List()); got != 0 {
		t.Fatalf("session count = %d, want 0", got)
	}
}

func TestSessionManagerCreateAndSwitch(t *testing.T) {
	manager := NewSessionManager()

	firstIndex, first, err := manager.Create("system", 10)
	if err != nil {
		t.Fatalf("create first session failed: %v", err)
	}
	if firstIndex != 0 {
		t.Fatalf("first session index = %d, want 0", firstIndex)
	}
	if manager.ActiveIndex() != 0 {
		t.Fatalf("active index = %d, want 0", manager.ActiveIndex())
	}

	secondIndex, second, err := manager.Create("system", 10)
	if err != nil {
		t.Fatalf("create second session failed: %v", err)
	}
	if secondIndex != 1 {
		t.Fatalf("second session index = %d, want 1", secondIndex)
	}
	if second == first {
		t.Fatal("expected separate histories")
	}

	if err := manager.Switch(0); err != nil {
		t.Fatalf("switch failed: %v", err)
	}
	active, err := manager.Active()
	if err != nil {
		t.Fatalf("get active session failed: %v", err)
	}
	if active != first {
		t.Fatal("expected first session to be active")
	}

	if got := len(manager.List()); got != 2 {
		t.Fatalf("session count = %d, want 2", got)
	}
}

func TestSessionManagerCopyCompleteHistory(t *testing.T) {
	manager := NewSessionManager()
	_, source, err := manager.Create("system", 10)
	if err != nil {
		t.Fatalf("create session failed: %v", err)
	}

	if err := source.AddUser("question", ""); err != nil {
		t.Fatalf("add user message failed: %v", err)
	}
	if err := source.AddAssistant("", "answer"); err != nil {
		t.Fatalf("add assistant message failed: %v", err)
	}

	index, copied, err := manager.Copy(nil)
	if err != nil {
		t.Fatalf("copy failed: %v", err)
	}
	if index != 1 {
		t.Fatalf("copied session index = %d, want 1", index)
	}
	if copied == source {
		t.Fatal("copy returned the source history")
	}
	if copied.Len() != source.Len() {
		t.Fatalf("copied history length = %d, want %d", copied.Len(), source.Len())
	}

	copied.Messages[1].Content = "changed in copy"
	if source.Messages[1].Content == copied.Messages[1].Content {
		t.Fatal("changing copied history changed source history")
	}
}

func TestSessionManagerCopyPrefix(t *testing.T) {
	manager := NewSessionManager()
	_, source, err := manager.Create("system", 10)
	if err != nil {
		t.Fatalf("create session failed: %v", err)
	}
	for _, content := range []string{"one", "two", "three"} {
		if err := source.AddUser(content, ""); err != nil {
			t.Fatalf("add message failed: %v", err)
		}
	}

	messageIndex := 2
	_, copied, err := manager.Copy(&messageIndex)
	if err != nil {
		t.Fatalf("copy prefix failed: %v", err)
	}
	if copied.Len() != 3 {
		t.Fatalf("copied prefix length = %d, want 3", copied.Len())
	}
	if copied.GetLast().Content != "two" {
		t.Fatalf("copied prefix last message = %q, want %q", copied.GetLast().Content, "two")
	}
}

func TestSessionManagerRejectsInvalidIndexes(t *testing.T) {
	manager := NewSessionManager()
	if _, _, err := manager.Create("system", 10); err != nil {
		t.Fatalf("create session failed: %v", err)
	}

	if err := manager.Switch(-1); err == nil {
		t.Fatal("expected negative switch index to fail")
	}
	if err := manager.Switch(1); err == nil {
		t.Fatal("expected out-of-range switch index to fail")
	}

	messageIndex := 1
	if _, _, err := manager.Copy(&messageIndex); err == nil {
		t.Fatal("expected out-of-range copy index to fail")
	}
}

func TestSessionManagerLimitsSessions(t *testing.T) {
	manager := NewSessionManager()
	for i := 0; i < MaxSessions; i++ {
		if _, _, err := manager.Create("system", 10); err != nil {
			t.Fatalf("create session %d failed: %v", i, err)
		}
	}

	if _, _, err := manager.Create("system", 10); err == nil {
		t.Fatal("expected session limit error")
	}
	if got := len(manager.List()); got != MaxSessions {
		t.Fatalf("session count = %d, want %d", got, MaxSessions)
	}
}
