package utils

import (
	"testing"
)

func TestFormatList_WithBullets(t *testing.T) {
	items := []string{"first.chat", "second.chat"}
	expected := "History files:\n - first.chat\n - second.chat"

	result := FormatList(items, "history files", false)

	if result != expected {
		t.Errorf("Expected:\n%q\nGot:\n%q", expected, result)
	}
}

func TestFormatList_WithNumbers(t *testing.T) {
	items := []string{"model-a", "model-b", "model-c"}
	expected := "Language models:\n(01) model-a\n(02) model-b\n(03) model-c"

	result := FormatList(items, "language models", true)

	if result != expected {
		t.Errorf("Expected:\n%q\nGot:\n%q", expected, result)
	}
}

func TestFormatList_Empty(t *testing.T) {
	items := []string{}
	expected := "no items found."

	result := FormatList(items, "items", false)

	if result != expected {
		t.Errorf("Expected:\n%q\nGot:\n%q", expected, result)
	}
}
