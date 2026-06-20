package config

import (
	"reflect"
	"testing"
)

func TestGetTemplate(t *testing.T) {
	prev := templates
	t.Cleanup(func() {
		templates = prev
	})

	setTemplates(map[string]string{
		"sum": "line 1\nline 2",
		"eng": "Translate to English",
	})

	if got := GetTemplate("eng"); got != "Translate to English" {
		t.Fatalf("GetTemplate(eng) = %q, want %q", got, "Translate to English")
	}
	if got := GetTemplate("sum"); got != "line 1\nline 2" {
		t.Fatalf("GetTemplate(sum) = %q, want multiline value", got)
	}
	if got := GetTemplate("missing"); got != "" {
		t.Fatalf("GetTemplate(missing) = %q, want empty string", got)
	}
}

func TestListTemplates(t *testing.T) {
	prev := templates
	t.Cleanup(func() {
		templates = prev
	})

	setTemplates(map[string]string{
		"ger": "Translate to German",
		"eng": "Translate to English",
		"sum": "Summary",
	})

	got := ListTemplates()
	want := []string{"eng", "ger", "sum"}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("ListTemplates() = %#v, want %#v", got, want)
	}
}

func TestSetTemplates_EmptyClears(t *testing.T) {
	prev := templates
	t.Cleanup(func() {
		templates = prev
	})

	setTemplates(map[string]string{"x": "y"})
	setTemplates(nil)

	if got := ListTemplates(); len(got) != 0 {
		t.Fatalf("expected empty template list, got %#v", got)
	}
}
