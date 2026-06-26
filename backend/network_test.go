package backend

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestBuildProviderURL_ErrorBranches(t *testing.T) {
	t.Run("invalid base url", func(t *testing.T) {
		_, err := buildProviderURL("://bad", "v1", "models")
		if err == nil {
			t.Fatal("expected error, got nil")
		}
	})

	t.Run("empty endpoint", func(t *testing.T) {
		_, err := buildProviderURL("https://api.example.com", "v1", "")
		if err == nil {
			t.Fatal("expected error, got nil")
		}
	})

	t.Run("empty api root", func(t *testing.T) {
		_, err := buildProviderURL("https://api.example.com", "", "models")
		if err == nil {
			t.Fatal("expected error, got nil")
		}
	})
}

func TestFetchOpenAIModels(t *testing.T) {
	t.Run("missing api key", func(t *testing.T) {
		_, err := fetchOpenAIModels("https://api.example.com", "")
		if err == nil || !strings.Contains(err.Error(), "missing OpenAI API key") {
			t.Fatalf("expected missing key error, got %v", err)
		}
	})

	t.Run("non-200 response", func(t *testing.T) {
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			http.Error(w, "nope", http.StatusBadRequest)
		}))
		defer srv.Close()

		_, err := fetchOpenAIModels(srv.URL, "k")
		if err == nil || !strings.Contains(err.Error(), "non-200 response") {
			t.Fatalf("expected non-200 error, got %v", err)
		}
	})

	t.Run("invalid json", func(t *testing.T) {
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			_, _ = fmt.Fprint(w, "{")
		}))
		defer srv.Close()

		_, err := fetchOpenAIModels(srv.URL, "k")
		if err == nil || !strings.Contains(err.Error(), "decode response failed") {
			t.Fatalf("expected decode error, got %v", err)
		}
	})

	t.Run("success with filtered empty ids", func(t *testing.T) {
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if got := r.Header.Get("Authorization"); got != "Bearer test-key" {
				http.Error(w, "missing auth", http.StatusUnauthorized)
				return
			}
			w.WriteHeader(http.StatusOK)
			_, _ = fmt.Fprint(w, `{"data":[{"id":"gpt-a"},{"id":""},{"id":"gpt-b"}]}`)
		}))
		defer srv.Close()

		got, err := fetchOpenAIModels(srv.URL, "test-key")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(got) != 2 || got[0] != "gpt-a" || got[1] != "gpt-b" {
			t.Fatalf("unexpected models: %+v", got)
		}
	})
}

func TestOpenAIAndResponsesServerVersion(t *testing.T) {
	openaiV, err := (&openAIClient{}).GetServerVersion()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(openaiV, "Chat Completions") {
		t.Fatalf("unexpected version string: %q", openaiV)
	}

	responsesV, err := (&openAIResponsesClient{}).GetServerVersion()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(responsesV, "Responses API") {
		t.Fatalf("unexpected version string: %q", responsesV)
	}
}

func TestChatStream_InputValidationErrors(t *testing.T) {
	dummyInput := ChatInput{}

	t.Run("openai missing key", func(t *testing.T) {
		_, err := (&openAIClient{baseURL: "https://api.example.com"}).ChatStream(dummyInput, nil)
		if err == nil || !strings.Contains(err.Error(), "missing OpenAI API key") {
			t.Fatalf("expected missing key error, got %v", err)
		}
	})

	t.Run("openai missing base url", func(t *testing.T) {
		_, err := (&openAIClient{apiKey: "k"}).ChatStream(dummyInput, nil)
		if err == nil || !strings.Contains(err.Error(), "missing OpenAI base URL") {
			t.Fatalf("expected missing base url error, got %v", err)
		}
	})

	t.Run("responses missing key", func(t *testing.T) {
		_, err := (&openAIResponsesClient{baseURL: "https://api.example.com"}).ChatStream(dummyInput, nil)
		if err == nil || !strings.Contains(err.Error(), "missing OpenAI API key") {
			t.Fatalf("expected missing key error, got %v", err)
		}
	})

	t.Run("responses missing base url", func(t *testing.T) {
		_, err := (&openAIResponsesClient{apiKey: "k"}).ChatStream(dummyInput, nil)
		if err == nil || !strings.Contains(err.Error(), "missing OpenAI base URL") {
			t.Fatalf("expected missing base url error, got %v", err)
		}
	})
}

func TestOllamaInvalidBaseURLBranches(t *testing.T) {
	c := &ollamaClient{baseURL: "://bad"}

	if _, err := c.ChatStream(ChatInput{}, nil); err == nil {
		t.Fatal("expected error for invalid chat base url, got nil")
	}
	if _, err := c.GetAvailableModels(); err == nil {
		t.Fatal("expected error for invalid models base url, got nil")
	}
	if _, err := c.GetServerVersion(); err == nil {
		t.Fatal("expected error for invalid version base url, got nil")
	}
}
