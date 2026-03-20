package chat

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"picochat/config"
	"picochat/messages"
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

// dummyHandleChat calls HandleChat and instantly closes the Stop-Channel.
func dummyHandleChat(cfg *config.Config, history *messages.ChatHistory) (*ChatResult, error) {
	stop := make(chan struct{})
	return HandleChat(cfg, history, stop)
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

	history := messages.NewHistory(cfg.Prompt, 10)
	history.AddUser("Say Hello", "")

	// Simulate HandleChat
	_, err := dummyHandleChat(cfg, history)
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

func TestHandleChat_InvalidURL(t *testing.T) {
	cfg := &config.Config{
		URL:   "://invalid-url",
		Model: "test-model",
	}

	history := messages.NewHistory(cfg.Prompt, 10)
	history.AddUser("test", "")

	_, err := dummyHandleChat(cfg, history)
	if err == nil || !strings.Contains(err.Error(), "http") {
		t.Errorf("expected http error, got: %v", err)
	}
}

func brokenStreamingHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	fmt.Fprintln(w, `{"message":{"content":"OK"}`) // invalid JSON
}

func TestHandleChat_BrokenStream(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(brokenStreamingHandler))
	defer server.Close()

	cfg := &config.Config{
		URL:   server.URL,
		Model: "test-model",
	}
	history := messages.NewHistory(cfg.Prompt, 10)
	history.AddUser("test", "")

	_, err := dummyHandleChat(cfg, history)
	if err == nil || !strings.Contains(err.Error(), "decode response") {
		t.Errorf("expected decode response failed, got: %v", err)
	}
}

func prematureEOFHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	fmt.Fprintln(w, `{"message":{"content":"Partial answer"}}`) // no done==True
	// then EOF
}

func TestHandleChat_EOFWithoutDone(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(prematureEOFHandler))
	defer server.Close()

	cfg := &config.Config{
		URL:   server.URL,
		Model: "test-model",
	}
	history := messages.NewHistory(cfg.Prompt, 10)
	history.AddUser("test", "")

	_, err := dummyHandleChat(cfg, history)
	if err != nil {
		t.Errorf("expected no error despite missing done=true, got: %v", err)
	}

	// Optional: Test, if partial answer has ben stored nevertheless
	last := history.Get()[len(history.Get())-1]
	if !strings.Contains(last.Content, "Partial answer") {
		t.Errorf("expected partial content, got: %v", last.Content)
	}
}
