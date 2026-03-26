package backend

import (
	"testing"

	"picochat/config"
)

func TestNew_SelectsBackendByFlavor(t *testing.T) {
	tests := []struct {
		name    string
		backend string
		want    any
	}{
		{name: "default ollama", backend: "", want: &ollamaClient{}},
		{name: "ollama explicit", backend: "ollama", want: &ollamaClient{}},
		{name: "openai", backend: "openai", want: &openAIClient{}},
		{name: "responses", backend: "responses", want: &openAIResponsesClient{}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := &config.Config{Backend: tt.backend}
			got := New(cfg)

			switch tt.want.(type) {
			case *ollamaClient:
				if _, ok := got.(*ollamaClient); !ok {
					t.Fatalf("expected *ollamaClient, got %T", got)
				}
			case *openAIClient:
				if _, ok := got.(*openAIClient); !ok {
					t.Fatalf("expected *openAIClient, got %T", got)
				}
			case *openAIResponsesClient:
				if _, ok := got.(*openAIResponsesClient); !ok {
					t.Fatalf("expected *openAIResponsesClient, got %T", got)
				}
			}
		})
	}
}
