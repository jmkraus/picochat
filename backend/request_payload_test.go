package backend

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"picochat/messages"
	"testing"
)

func fptr(v float64) *float64 { return &v }

func TestOllamaChatStream_RequestPayload(t *testing.T) {
	tests := []struct {
		name        string
		temperature *float64
		topP        *float64
		wantOptions bool
		wantTemp    bool
		wantTopP    bool
	}{
		{
			name:        "omit options when both unset",
			wantOptions: false,
		},
		{
			name:        "include only temperature",
			temperature: fptr(0.2),
			wantOptions: true,
			wantTemp:    true,
			wantTopP:    false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var gotPath string
			var gotBody map[string]any

			srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				gotPath = r.URL.Path
				raw, err := io.ReadAll(r.Body)
				if err != nil {
					t.Fatalf("read request body failed: %v", err)
				}
				if err := json.Unmarshal(raw, &gotBody); err != nil {
					t.Fatalf("decode request body failed: %v", err)
				}
				w.Header().Set("Content-Type", "application/json")
				_, _ = w.Write([]byte("{\"message\":{\"content\":\"ok\"},\"done\":true}\n"))
			}))
			defer srv.Close()

			c := newOllamaClient(srv.URL)
			_, err := c.ChatStream(ChatInput{
				Model:       "m",
				Messages:    []messages.Message{{Role: messages.RoleUser, Content: "hi"}},
				Temperature: tt.temperature,
				TopP:        tt.topP,
			}, nil)
			if err != nil {
				t.Fatalf("ChatStream failed: %v", err)
			}

			if gotPath != "/api/chat" {
				t.Fatalf("path = %q, want %q", gotPath, "/api/chat")
			}

			options, hasOptions := gotBody["options"].(map[string]any)
			if hasOptions != tt.wantOptions {
				t.Fatalf("options present = %v, want %v", hasOptions, tt.wantOptions)
			}
			if !tt.wantOptions {
				return
			}

			_, hasTemp := options["temperature"]
			if hasTemp != tt.wantTemp {
				t.Fatalf("temperature present = %v, want %v", hasTemp, tt.wantTemp)
			}

			_, hasTopP := options["top_p"]
			if hasTopP != tt.wantTopP {
				t.Fatalf("top_p present = %v, want %v", hasTopP, tt.wantTopP)
			}
		})
	}
}

func TestOpenAIChatCompletions_RequestPayload(t *testing.T) {
	tests := []struct {
		name        string
		temperature *float64
		topP        *float64
		wantTemp    bool
		wantTopP    bool
	}{
		{
			name:     "omit sampling when unset",
			wantTemp: false,
			wantTopP: false,
		},
		{
			name:     "include only top_p",
			topP:     fptr(0.8),
			wantTemp: false,
			wantTopP: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var gotPath string
			var gotAuth string
			var gotBody map[string]any

			srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				gotPath = r.URL.Path
				gotAuth = r.Header.Get("Authorization")
				raw, err := io.ReadAll(r.Body)
				if err != nil {
					t.Fatalf("read request body failed: %v", err)
				}
				if err := json.Unmarshal(raw, &gotBody); err != nil {
					t.Fatalf("decode request body failed: %v", err)
				}
				w.Header().Set("Content-Type", "text/event-stream")
				_, _ = w.Write([]byte("data: {\"choices\":[{\"delta\":{\"content\":\"ok\"},\"finish_reason\":null}]}\n"))
				_, _ = w.Write([]byte("data: [DONE]\n"))
			}))
			defer srv.Close()

			c := newOpenAIClient(srv.URL, "sk-test")
			_, err := c.ChatStream(ChatInput{
				Model:       "m",
				Messages:    []messages.Message{{Role: messages.RoleUser, Content: "hi"}},
				Temperature: tt.temperature,
				TopP:        tt.topP,
			}, nil)
			if err != nil {
				t.Fatalf("ChatStream failed: %v", err)
			}

			if gotPath != "/v1/chat/completions" {
				t.Fatalf("path = %q, want %q", gotPath, "/v1/chat/completions")
			}
			if gotAuth != "Bearer sk-test" {
				t.Fatalf("auth = %q, want %q", gotAuth, "Bearer sk-test")
			}

			_, hasTemp := gotBody["temperature"]
			if hasTemp != tt.wantTemp {
				t.Fatalf("temperature present = %v, want %v", hasTemp, tt.wantTemp)
			}
			_, hasTopP := gotBody["top_p"]
			if hasTopP != tt.wantTopP {
				t.Fatalf("top_p present = %v, want %v", hasTopP, tt.wantTopP)
			}
		})
	}
}

func TestResponsesAPI_RequestPayload(t *testing.T) {
	tests := []struct {
		name           string
		format         map[string]any
		wantFormatType string
		wantSchema     bool
	}{
		{
			name:           "plain text when no schema",
			format:         nil,
			wantFormatType: "text",
			wantSchema:     false,
		},
		{
			name: "json schema when provided",
			format: map[string]any{
				"type":                 "object",
				"properties":           map[string]any{"x": map[string]any{"type": "string"}},
				"required":             []any{"x"},
				"additionalProperties": false,
			},
			wantFormatType: "json_schema",
			wantSchema:     true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var gotPath string
			var gotAuth string
			var gotBody map[string]any

			srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				gotPath = r.URL.Path
				gotAuth = r.Header.Get("Authorization")
				raw, err := io.ReadAll(r.Body)
				if err != nil {
					t.Fatalf("read request body failed: %v", err)
				}
				if err := json.Unmarshal(raw, &gotBody); err != nil {
					t.Fatalf("decode request body failed: %v", err)
				}
				w.Header().Set("Content-Type", "text/event-stream")
				_, _ = w.Write([]byte("data: {\"type\":\"response.output_text.delta\",\"delta\":\"ok\"}\n"))
				_, _ = w.Write([]byte("data: [DONE]\n"))
			}))
			defer srv.Close()

			c := newOpenAIResponsesClient(srv.URL, "sk-test")
			_, err := c.ChatStream(ChatInput{
				Model:    "m",
				Messages: []messages.Message{{Role: messages.RoleUser, Content: "hi"}},
				Format:   tt.format,
			}, nil)
			if err != nil {
				t.Fatalf("ChatStream failed: %v", err)
			}

			if gotPath != "/v1/responses" {
				t.Fatalf("path = %q, want %q", gotPath, "/v1/responses")
			}
			if gotAuth != "Bearer sk-test" {
				t.Fatalf("auth = %q, want %q", gotAuth, "Bearer sk-test")
			}

			text, ok := gotBody["text"].(map[string]any)
			if !ok {
				t.Fatalf("text block missing in payload")
			}
			format, ok := text["format"].(map[string]any)
			if !ok {
				t.Fatalf("text.format block missing in payload")
			}
			if format["type"] != tt.wantFormatType {
				t.Fatalf("format type = %v, want %q", format["type"], tt.wantFormatType)
			}
			_, hasSchema := format["schema"]
			if hasSchema != tt.wantSchema {
				t.Fatalf("schema present = %v, want %v", hasSchema, tt.wantSchema)
			}
		})
	}
}
