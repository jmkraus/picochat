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
	expected := "No items found."

	result := FormatList(items, "items", false)

	if result != expected {
		t.Errorf("Expected:\n%q\nGot:\n%q", expected, result)
	}
}

func TestListAvailableModels_Empty(t *testing.T) {
	_, err := ListAvailableModels([]string{})
	if err == nil {
		t.Fatal("expected error for empty model list, got nil")
	}
}

func TestListAvailableModels_SingleEntry(t *testing.T) {
	models := []string{"model-x"}

	got, err := ListAvailableModels(models)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	want := "Language models:\n(01) model-x"
	if got != want {
		t.Fatalf("result = %q, want %q", got, want)
	}

	v, ok := GetModelsByIndex(1)
	if !ok || v != "model-x" {
		t.Fatalf("GetModelsByIndex(1) = (%q, %v), want (model-x, true)", v, ok)
	}
}

func TestListAvailableModels_MultipleEntries_SortsAndIndexes(t *testing.T) {
	tests := []struct {
		name   string
		input  []string
		sorted []string
	}{
		{
			name:   "already sorted",
			input:  []string{"alpha", "beta", "gamma"},
			sorted: []string{"alpha", "beta", "gamma"},
		},
		{
			name:   "unsorted input",
			input:  []string{"gamma", "alpha", "beta"},
			sorted: []string{"alpha", "beta", "gamma"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			models := append([]string(nil), tt.input...)
			got, err := ListAvailableModels(models)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			want := "Language models:\n(01) " + tt.sorted[0] + "\n(02) " + tt.sorted[1] + "\n(03) " + tt.sorted[2]
			if got != want {
				t.Fatalf("result = %q, want %q", got, want)
			}

			for i, m := range tt.sorted {
				v, ok := GetModelsByIndex(i + 1)
				if !ok || v != m {
					t.Fatalf("GetModelsByIndex(%d) = (%q, %v), want (%q, true)", i+1, v, ok, m)
				}
			}

			for i, m := range tt.sorted {
				if models[i] != m {
					t.Fatalf("models slice not sorted in place at index %d: got %q, want %q", i, models[i], m)
				}
			}
		})
	}
}

func TestCapitalize(t *testing.T) {
	tests := []struct {
		name string
		in   string
		want string
	}{
		{name: "empty", in: "", want: ""},
		{name: "all uppercase", in: "HELLO WORLD", want: "Hello world"},
		{name: "mixed case", in: "hElLo WoRLD", want: "Hello world"},
		{name: "already lowercase", in: "hello world", want: "Hello world"},
		{name: "non-letter first rune", in: "123ABC", want: "123abc"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := capitalize(tt.in)
			if got != tt.want {
				t.Fatalf("capitalize(%q) = %q, want %q", tt.in, got, tt.want)
			}
		})
	}
}
