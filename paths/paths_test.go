package paths

import (
	"testing"
)

func TestEnsureSuffix(t *testing.T) {
	if EnsureSuffix("foo", ".chat") != "foo.chat" {
		t.Errorf("Expected suffix to be appended")
	}
	if EnsureSuffix("bar.chat", ".chat") != "bar.chat" {
		t.Errorf("Suffix should not be duplicated")
	}
}
