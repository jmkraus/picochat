package backend

import "testing"

func TestBuildOllamaURL(t *testing.T) {
	tests := []struct {
		name     string
		baseURL  string
		endpoint string
		want     string
	}{
		{
			name:     "no path adds api",
			baseURL:  "http://localhost:11434",
			endpoint: "chat",
			want:     "http://localhost:11434/api/chat",
		},
		{
			name:     "existing api path stays",
			baseURL:  "http://localhost:11434/api",
			endpoint: "tags",
			want:     "http://localhost:11434/api/tags",
		},
		{
			name:     "custom path is trusted",
			baseURL:  "http://localhost:11434/ollama",
			endpoint: "version",
			want:     "http://localhost:11434/ollama/version",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := buildOllamaURL(tt.baseURL, tt.endpoint)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if got != tt.want {
				t.Fatalf("got %q, want %q", got, tt.want)
			}
		})
	}
}

func TestBuildOpenAIURL(t *testing.T) {
	tests := []struct {
		name     string
		baseURL  string
		endpoint string
		want     string
	}{
		{
			name:     "no path adds v1",
			baseURL:  "https://api.openai.com",
			endpoint: "models",
			want:     "https://api.openai.com/v1/models",
		},
		{
			name:     "existing v1 path stays",
			baseURL:  "https://api.openai.com/v1",
			endpoint: "chat/completions",
			want:     "https://api.openai.com/v1/chat/completions",
		},
		{
			name:     "custom path is trusted",
			baseURL:  "http://localhost:8080/openai",
			endpoint: "chat/completions",
			want:     "http://localhost:8080/openai/chat/completions",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := buildOpenAIURL(tt.baseURL, tt.endpoint)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if got != tt.want {
				t.Fatalf("got %q, want %q", got, tt.want)
			}
		})
	}
}

func TestBuildOpenAIURL_ResponsesEndpoint(t *testing.T) {
	tests := []struct {
		name     string
		baseURL  string
		endpoint string
		want     string
	}{
		{
			name:     "no path adds v1 for responses",
			baseURL:  "https://api.openai.com",
			endpoint: "responses",
			want:     "https://api.openai.com/v1/responses",
		},
		{
			name:     "existing v1 path keeps responses endpoint",
			baseURL:  "https://api.openai.com/v1",
			endpoint: "responses",
			want:     "https://api.openai.com/v1/responses",
		},
		{
			name:     "custom proxy path is trusted for responses",
			baseURL:  "http://localhost:8080/openai",
			endpoint: "responses",
			want:     "http://localhost:8080/openai/responses",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := buildOpenAIURL(tt.baseURL, tt.endpoint)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if got != tt.want {
				t.Fatalf("got %q, want %q", got, tt.want)
			}
		})
	}
}
