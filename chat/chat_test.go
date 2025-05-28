package chat_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"picochat/chat"
	"picochat/config"
	"picochat/types"
)

// Simulated Streaming-Handler (PicoAI)
func streamingHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// Simulated Streaming-Body (Chunked JSON)
	chunks := []string{
		`{"message":{"content":"Hallo"}}`,
		`{"message":{"content":" Welt"},"done":true,"prompt_eval_count":5,"eval_count":10}`,
	}

	for _, chunk := range chunks {
		fmt.Fprintln(w, chunk)
	}
}

func TestHandleChat(t *testing.T) {
	// Starting fake server
	server := httptest.NewServer(http.HandlerFunc(streamingHandler))
	defer server.Close()

	cfg := &config.Config{
		URL:    server.URL,
		Model:  "test-model",
		Prompt: "You are a test bot",
	}

	history := types.NewHistory(cfg.Prompt, 10)
	history.Add("user", "Sag Hallo")

	// Simulate HandleChat
	err := chat.HandleChat(cfg, history)
	if err != nil {
		t.Fatalf("HandleChat returned error: %v", err)
	}

	// Check if bot reply was stored
	messages := history.Get()
	if len(messages) != 3 {
		t.Errorf("expected 3 messages (system, user, assistant), got %d", len(messages))
	}

	last := messages[len(messages)-1]
	if !strings.Contains(last.Content, "Hallo Welt") {
		t.Errorf("unexpected assistant response: %q", last.Content)
	}
}
