package backend

import (
	"strings"
	"testing"

	"picochat/messages"
)

func TestParseOpenAIChatCompletionEvent(t *testing.T) {
	t.Run("content delta", func(t *testing.T) {
		data := `{"choices":[{"delta":{"content":"hello"},"finish_reason":null}]}`
		thinking, content, done, err := parseOpenAIChatCompletionEvent(data)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if thinking != "" || content != "hello" || done {
			t.Fatalf("unexpected parse result: thinking=%q content=%q done=%v", thinking, content, done)
		}
	})

	t.Run("reasoning content preferred over reasoning", func(t *testing.T) {
		data := `{"choices":[{"delta":{"reasoning_content":"r1","reasoning":"r2"},"finish_reason":null}]}`
		thinking, content, done, err := parseOpenAIChatCompletionEvent(data)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if thinking != "r1" || content != "" || done {
			t.Fatalf("unexpected parse result: thinking=%q content=%q done=%v", thinking, content, done)
		}
	})

	t.Run("done when finish reason present", func(t *testing.T) {
		data := `{"choices":[{"delta":{"content":"x"},"finish_reason":"stop"}]}`
		_, _, done, err := parseOpenAIChatCompletionEvent(data)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if !done {
			t.Fatal("expected done=true")
		}
	})

	t.Run("empty choices", func(t *testing.T) {
		data := `{"choices":[]}`
		thinking, content, done, err := parseOpenAIChatCompletionEvent(data)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if thinking != "" || content != "" || done {
			t.Fatalf("unexpected parse result: thinking=%q content=%q done=%v", thinking, content, done)
		}
	})

	t.Run("invalid json", func(t *testing.T) {
		_, _, _, err := parseOpenAIChatCompletionEvent("{")
		if err == nil {
			t.Fatal("expected error for invalid json")
		}
	})
}

func TestParseResponsesEvent(t *testing.T) {
	tests := []struct {
		name         string
		data         string
		wantThinking string
		wantContent  string
		wantDone     bool
	}{
		{
			name:        "text delta",
			data:        `{"type":"response.output_text.delta","delta":"hello"}`,
			wantContent: "hello",
		},
		{
			name:         "reasoning delta",
			data:         `{"type":"response.reasoning.delta","delta":"think"}`,
			wantThinking: "think",
		},
		{
			name:         "reasoning summary delta",
			data:         `{"type":"response.reasoning_summary_text.delta","delta":"summary"}`,
			wantThinking: "summary",
		},
		{
			name:     "completed",
			data:     `{"type":"response.completed"}`,
			wantDone: true,
		},
		{
			name: "unknown event",
			data: `{"type":"other.event","delta":"x"}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			thinking, content, done, err := parseResponsesEvent(tt.data)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if thinking != tt.wantThinking || content != tt.wantContent || done != tt.wantDone {
				t.Fatalf("unexpected parse result: thinking=%q content=%q done=%v", thinking, content, done)
			}
		})
	}

	t.Run("invalid json", func(t *testing.T) {
		_, _, _, err := parseResponsesEvent("{")
		if err == nil {
			t.Fatal("expected error for invalid json")
		}
	})
}

func TestBuildResponsesText(t *testing.T) {
	t.Run("empty schema uses plain text", func(t *testing.T) {
		got := buildResponsesText(nil)
		if got == nil {
			t.Fatal("expected non-nil response text")
		}
		if got.Format.Type != "text" {
			t.Fatalf("format type = %q, want %q", got.Format.Type, "text")
		}
	})

	t.Run("schema uses json_schema format", func(t *testing.T) {
		schema := map[string]any{"type": "object"}
		got := buildResponsesText(schema)
		if got == nil {
			t.Fatal("expected non-nil response text")
		}
		if got.Format.Type != "json_schema" {
			t.Fatalf("format type = %q, want %q", got.Format.Type, "json_schema")
		}
		if got.Format.Name != "user" {
			t.Fatalf("format name = %q, want %q", got.Format.Name, "user")
		}
		if got.Format.Schema["type"] != "object" {
			t.Fatalf("unexpected schema: %+v", got.Format.Schema)
		}
	})
}

func TestMapMessagesToOpenAIChatMessages(t *testing.T) {
	in := []messages.Message{
		{Role: messages.RoleSystem, Content: "sys"},
		{Role: messages.RoleUser, Content: "hello"},
		{
			Role:    messages.RoleUser,
			Content: "caption",
			Images:  []string{"data:image/png;base64,AAA", " ", "data:image/jpeg;base64,BBB"},
		},
	}

	out := mapMessagesToOpenAIChatMessages(in)
	if len(out) != 3 {
		t.Fatalf("len(out) = %d, want 3", len(out))
	}

	if c, ok := out[0].Content.(string); !ok || c != "sys" {
		t.Fatalf("unexpected system content mapping: %#v", out[0].Content)
	}
	if c, ok := out[1].Content.(string); !ok || c != "hello" {
		t.Fatalf("unexpected user content mapping: %#v", out[1].Content)
	}

	parts, ok := out[2].Content.([]map[string]any)
	if !ok {
		t.Fatalf("expected multimodal content parts, got %T", out[2].Content)
	}
	if len(parts) != 3 {
		t.Fatalf("len(parts) = %d, want 3", len(parts))
	}
	if parts[0]["type"] != "text" {
		t.Fatalf("first part type = %v, want text", parts[0]["type"])
	}
	if parts[1]["type"] != "image_url" || parts[2]["type"] != "image_url" {
		t.Fatalf("image parts not mapped correctly: %+v", parts)
	}
}

func TestMapMessagesToResponsesInput(t *testing.T) {
	in := []messages.Message{
		{
			Role:    messages.RoleUser,
			Content: "prompt",
			Images:  []string{"data:image/png;base64,AAA", ""},
		},
		{
			Role:    messages.RoleAssistant,
			Content: "answer",
		},
	}

	out := mapMessagesToResponsesInput(in)
	if len(out) != 2 {
		t.Fatalf("len(out) = %d, want 2", len(out))
	}
	if out[0].Role != messages.RoleUser {
		t.Fatalf("first role = %q, want %q", out[0].Role, messages.RoleUser)
	}
	if len(out[0].Content) != 2 {
		t.Fatalf("first content parts len = %d, want 2", len(out[0].Content))
	}
	if out[0].Content[0].Type != "input_text" || out[0].Content[0].Text != "prompt" {
		t.Fatalf("unexpected first text part: %+v", out[0].Content[0])
	}
	if out[0].Content[1].Type != "input_image" || !strings.HasPrefix(out[0].Content[1].ImageURL, "data:image") {
		t.Fatalf("unexpected image part: %+v", out[0].Content[1])
	}
	if out[0].Content[1].Detail != "auto" {
		t.Fatalf("image detail = %q, want auto", out[0].Content[1].Detail)
	}

	if len(out[1].Content) != 1 || out[1].Content[0].Type != "input_text" {
		t.Fatalf("unexpected assistant mapping: %+v", out[1].Content)
	}
}

func TestNormalizeOllamaImages(t *testing.T) {
	in := []messages.Message{
		{
			Role:   messages.RoleUser,
			Images: []string{"data:image/png;base64,ABC", "rawB64"},
		},
	}

	out := normalizeOllamaImages(in)
	if len(out) != 1 || len(out[0].Images) != 2 {
		t.Fatalf("unexpected output shape: %+v", out)
	}
	if out[0].Images[0] != "ABC" {
		t.Fatalf("first image = %q, want %q", out[0].Images[0], "ABC")
	}
	if out[0].Images[1] != "rawB64" {
		t.Fatalf("second image = %q, want rawB64", out[0].Images[1])
	}

	// Ensure input was not modified in place.
	if in[0].Images[0] != "data:image/png;base64,ABC" {
		t.Fatalf("input was modified: %+v", in)
	}
}
