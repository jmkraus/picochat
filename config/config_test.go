package config

import "testing"

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
