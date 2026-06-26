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

	if got, err := GetTemplate("eng"); err != nil || got != "Translate to English" {
		t.Fatalf("GetTemplate(eng) = (%q, %v), want (%q, nil)", got, err, "Translate to English")
	}
	if got, err := GetTemplate("sum"); err != nil || got != "line 1\nline 2" {
		t.Fatalf("GetTemplate(sum) = (%q, %v), want multiline value", got, err)
	}
	if got, err := GetTemplate("missing"); err == nil || got != "" {
		t.Fatalf("GetTemplate(missing) = (%q, %v), want (empty, error)", got, err)
	}
	if got, err := GetTemplate(""); err != nil || got != "" {
		t.Fatalf("GetTemplate(\"\") = (%q, %v), want (empty, nil)", got, err)
	}
	if got, err := GetTemplate("   "); err != nil || got != "" {
		t.Fatalf("GetTemplate(\"   \") = (%q, %v), want (empty, nil)", got, err)
	}
	if got, err := GetTemplate(" eng "); err != nil || got != "Translate to English" {
		t.Fatalf("GetTemplate(\" eng \") = (%q, %v), want (%q, nil)", got, err, "Translate to English")
	}
	if got := GetTemplateDescription("eng"); got != "Translate EN" {
		t.Fatalf("GetTemplateDescription(eng) = %q, want %q", got, "Translate EN")
	}
	if got := GetTemplateDescription("missing"); got != "" {
		t.Fatalf("GetTemplateDescription(missing) = %q, want empty string", got)
	}
}

func TestGetTemplate_EmptyPromptReturnsError(t *testing.T) {
	prev := templates
	t.Cleanup(func() {
		templates = prev
	})

	setTemplates(map[string]Template{
		"empty": {Prompt: "", Description: "Empty"},
	})

	got, err := GetTemplate("empty")
	if err == nil || got != "" {
		t.Fatalf("GetTemplate(empty) = (%q, %v), want (empty, error)", got, err)
	}
}

func TestGetTemplate_NilTemplatesMapReturnsNotFound(t *testing.T) {
	prev := templates
	t.Cleanup(func() {
		templates = prev
	})

	templates = nil
	got, err := GetTemplate("missing")
	if err == nil || got != "" {
		t.Fatalf("GetTemplate(missing) with nil map = (%q, %v), want (empty, error)", got, err)
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
