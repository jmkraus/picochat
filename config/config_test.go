package config

import (
	"strings"
	"testing"
)

func TestDefaultConfig(t *testing.T) {
	cfg := DefaultConfig()

	if cfg.URL != "http://localhost:11434/api" {
		t.Errorf("URL default = %q, want %q", cfg.URL, "http://localhost:11434/api")
	}
	if cfg.Model != "gpt-oss:latest" {
		t.Errorf("Model default = %q, want %q", cfg.Model, "gpt-oss:latest")
	}
	if cfg.Prompt != "You are a Large Language Model. Answer as concisely as possible. Your answers should be informative, helpful and engaging." {
		t.Errorf("Prompt default = %q, want expected prompt text", cfg.Prompt)
	}
	if cfg.Context != 20 {
		t.Errorf("Context default = %d, want %d", cfg.Context, 20)
	}
	if cfg.Temperature != 0.7 {
		t.Errorf("Temperature default = %v, want %v", cfg.Temperature, 0.7)
	}
	if cfg.Top_p != 0.9 {
		t.Errorf("Top_p default = %v, want %v", cfg.Top_p, 0.9)
	}
	if cfg.Reasoning {
		t.Errorf("Reasoning default = %v, want %v", cfg.Reasoning, false)
	}
	if cfg.Quiet {
		t.Errorf("Quiet default = %v, want %v", cfg.Quiet, false)
	}

	if cfg.ConfigPath != "" {
		t.Errorf("ConfigPath default = %q, want empty string", cfg.ConfigPath)
	}
	if cfg.ImagePath != "" {
		t.Errorf("ImagePath default = %q, want empty string", cfg.ImagePath)
	}
	if cfg.OutputFmt != "" {
		t.Errorf("OutputFmt default = %q, want empty string", cfg.OutputFmt)
	}
	if cfg.SchemaFmt != nil {
		t.Errorf("SchemaFmt default = %v, want nil", cfg.SchemaFmt)
	}
}

func TestApplyEnvValues_OverridesConfigFields(t *testing.T) {
	cfg := DefaultConfig()

	t.Setenv("PICOCHAT_URL", "https://example.org/api")
	t.Setenv("PICOCHAT_MODEL", "gpt-4.1-mini")
	t.Setenv("PICOCHAT_CONTEXT", "42")
	t.Setenv("PICOCHAT_TEMPERATURE", "0.2")
	t.Setenv("PICOCHAT_TOP_P", "0.6")
	t.Setenv("PICOCHAT_REASONING", "true")
	t.Setenv("PICOCHAT_QUIET", "true")

	if err := applyEnvValues(&cfg); err != nil {
		t.Fatalf("applyEnvValues returned error: %v", err)
	}

	if cfg.URL != "https://example.org/api" {
		t.Errorf("URL = %q, want %q", cfg.URL, "https://example.org/api")
	}
	if cfg.Model != "gpt-4.1-mini" {
		t.Errorf("Model = %q, want %q", cfg.Model, "gpt-4.1-mini")
	}
	if cfg.Context != 42 {
		t.Errorf("Context = %d, want %d", cfg.Context, 42)
	}
	if cfg.Temperature != 0.2 {
		t.Errorf("Temperature = %v, want %v", cfg.Temperature, 0.2)
	}
	if cfg.Top_p != 0.6 {
		t.Errorf("Top_p = %v, want %v", cfg.Top_p, 0.6)
	}
	if !cfg.Reasoning {
		t.Errorf("Reasoning = %v, want %v", cfg.Reasoning, true)
	}
	if !cfg.Quiet {
		t.Errorf("Quiet = %v, want %v", cfg.Quiet, true)
	}
}

func TestApplyEnvValues_InvalidValueReturnsError(t *testing.T) {
	cfg := DefaultConfig()

	t.Setenv("PICOCHAT_CONTEXT", "not-an-int")

	err := applyEnvValues(&cfg)
	if err == nil {
		t.Fatal("applyEnvValues returned nil error, want error")
	}
	if !strings.Contains(err.Error(), "PICOCHAT_CONTEXT") {
		t.Errorf("error = %q, want to contain %q", err.Error(), "PICOCHAT_CONTEXT")
	}
}

func TestApplyEnvValues_EmptyValueIsIgnored(t *testing.T) {
	cfg := DefaultConfig()

	t.Setenv("PICOCHAT_MODEL", "")
	t.Setenv("PICOCHAT_QUIET", "")

	if err := applyEnvValues(&cfg); err != nil {
		t.Fatalf("applyEnvValues returned error: %v", err)
	}

	if cfg.Model != "gpt-oss:latest" {
		t.Errorf("Model = %q, want %q", cfg.Model, "gpt-oss:latest")
	}
	if cfg.Quiet {
		t.Errorf("Quiet = %v, want %v", cfg.Quiet, false)
	}
}
