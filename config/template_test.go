package config

import (
	"strings"
	"testing"
)

func TestGetTemplate(t *testing.T) {
	prev := templates
	t.Cleanup(func() {
		templates = prev
	})

	setTemplates(map[string]Template{
		"sum": {Prompt: "line 1\nline 2", Description: "Summary"},
		"eng": {Prompt: "Translate to English", Description: "Translate EN"},
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
	if got := GetTemplateDescription("eng"); got != "Translate EN" {
		t.Fatalf("GetTemplateDescription(eng) = %q, want %q", got, "Translate EN")
	}
	if got := GetTemplateDescription("missing"); got != "" {
		t.Fatalf("GetTemplateDescription(missing) = %q, want empty string", got)
	}
}

func TestListTemplates(t *testing.T) {
	prev := templates
	t.Cleanup(func() {
		templates = prev
	})

	setTemplates(map[string]Template{
		"ger": {Prompt: "Translate to German", Description: "DE"},
		"eng": {Prompt: "Translate to English", Description: "EN"},
		"sum": {Prompt: "Summary", Description: "SUM"},
	})

	got := ListTemplates()
	if !strings.Contains(got, "| Key") || !strings.Contains(got, "| Description") {
		t.Fatalf("ListTemplates() missing header columns, got:\n%s", got)
	}
	if !strings.Contains(got, "| eng") || !strings.Contains(got, "| EN") {
		t.Fatalf("ListTemplates() missing eng row, got:\n%s", got)
	}
	if !strings.Contains(got, "| ger") || !strings.Contains(got, "| DE") {
		t.Fatalf("ListTemplates() missing ger row, got:\n%s", got)
	}
	if !strings.Contains(got, "| sum") || !strings.Contains(got, "| SUM") {
		t.Fatalf("ListTemplates() missing sum row, got:\n%s", got)
	}
}

func TestSetTemplates_EmptyClears(t *testing.T) {
	prev := templates
	t.Cleanup(func() {
		templates = prev
	})

	setTemplates(map[string]Template{"x": {Prompt: "y"}})
	setTemplates(nil)

	got := ListTemplates()
	if !strings.Contains(got, "| Key") || !strings.Contains(got, "| Description") {
		t.Fatalf("expected header-only table, got:\n%s", got)
	}
}
