package backend

import "testing"

func TestBuildOllamaURL(t *testing.T) {
	tests := []struct {
		name     string
		baseURL  string
		endpoint string
		want     string
		wantErr  bool
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
			name:     "legacy v1 path is rewritten to api",
			baseURL:  "http://localhost:11434/v1",
			endpoint: "version",
			want:     "http://localhost:11434/api/version",
		},
		{
			name:     "custom path is rejected",
			baseURL:  "http://localhost:11434/ollama",
			endpoint: "version",
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := buildOllamaURL(tt.baseURL, tt.endpoint)
			if tt.wantErr {
				if err == nil {
					t.Fatalf("expected error, got nil")
				}
				return
			}
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
		wantErr  bool
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
			name:     "legacy api path is rewritten to v1",
			baseURL:  "http://localhost:8080/api",
			endpoint: "chat/completions",
			want:     "http://localhost:8080/v1/chat/completions",
		},
		{
			name:     "custom path is rejected",
			baseURL:  "http://localhost:8080/openai",
			endpoint: "chat/completions",
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := buildOpenAIURL(tt.baseURL, tt.endpoint)
			if tt.wantErr {
				if err == nil {
					t.Fatalf("expected error, got nil")
				}
				return
			}
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
		wantErr  bool
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
			name:     "legacy api path rewritten for responses",
			baseURL:  "http://localhost:8080/api",
			endpoint: "responses",
			want:     "http://localhost:8080/v1/responses",
		},
		{
			name:     "custom proxy path rejected for responses",
			baseURL:  "http://localhost:8080/openai",
			endpoint: "responses",
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := buildOpenAIURL(tt.baseURL, tt.endpoint)
			if tt.wantErr {
				if err == nil {
					t.Fatalf("expected error, got nil")
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if got != tt.want {
				t.Fatalf("got %q, want %q", got, tt.want)
			}
		})
	}
}
