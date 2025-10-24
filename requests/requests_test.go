package requests_test

import (
	"testing"

	"picochat/requests"
)

func TestCleanUrl(t *testing.T) {
	tests := []struct {
		name      string
		baseURL   string
		endpoint  string
		want      string
		shouldErr bool
	}{
		{
			name:     "base URL without slash",
			baseURL:  "http://localhost:11434",
			endpoint: "models",
			want:     "http://localhost:11434/api/models",
		},
		{
			name:     "base URL with slash",
			baseURL:  "http://localhost:11434/",
			endpoint: "models",
			want:     "http://localhost:11434/api/models",
		},
		{
			name:     "base URL with /api without slash",
			baseURL:  "http://localhost:11434/api",
			endpoint: "models",
			want:     "http://localhost:11434/api/models",
		},
		{
			name:     "base URL with /api/",
			baseURL:  "http://localhost:11434/api/",
			endpoint: "models",
			want:     "http://localhost:11434/api/models",
		},
		{
			name:      "invalid URL",
			baseURL:   "://bad_url",
			endpoint:  "models",
			want:      "",
			shouldErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := requests.CleanUrl(tt.baseURL, tt.endpoint)

			if tt.shouldErr {
				if err == nil {
					t.Errorf("expected error for input %q, got none", tt.baseURL)
				}
				return
			}

			if err != nil {
				t.Errorf("unexpected error for input %q: %v", tt.baseURL, err)
			}

			if got != tt.want {
				t.Errorf("CleanUrl(%q, %q) = %q; want %q", tt.baseURL, tt.endpoint, got, tt.want)
			}
		})
	}
}
