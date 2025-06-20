package paths

import (
	"testing"
)

func TestHasSuffix(t *testing.T) {
	if !HasSuffix("test.chat", ".chat") {
		t.Error("Expected HasSuffix to return true for 'test.chat'")
	}
	if HasSuffix("test.txt", ".chat") {
		t.Error("Expected HasSuffix to return false for 'test.txt'")
	}
}

func TestEnsureSuffix(t *testing.T) {
	if EnsureSuffix("foo", ".chat") != "foo.chat" {
		t.Errorf("Expected suffix to be appended")
	}
	if EnsureSuffix("bar.chat", ".chat") != "bar.chat" {
		t.Errorf("Suffix should not be duplicated")
	}
}
