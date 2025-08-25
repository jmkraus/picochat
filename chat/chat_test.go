package chat

import (
	"fmt"
	"math"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

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
	history.Add("user", "Say Hello")

	// Simulate HandleChat
	_, err := HandleChat(cfg, history)
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
	history.Add("user", "test")

	_, err := HandleChat(cfg, history)
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
	history.Add("user", "test")

	_, err := HandleChat(cfg, history)
	if err == nil || !strings.Contains(err.Error(), "stream") {
		t.Errorf("expected stream decoding error, got: %v", err)
	}
}

func prematureEOFHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	fmt.Fprintln(w, `{"message":{"content":"Teilantwort"}}`) // no done==True
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
	history.Add("user", "test")

	_, err := HandleChat(cfg, history)
	if err != nil {
		t.Errorf("expected no error despite missing done=true, got: %v", err)
	}

	// Optional: Test, ob Antwort trotzdem gespeichert wurde
	last := history.Get()[len(history.Get())-1]
	if !strings.Contains(last.Content, "Teilantwort") {
		t.Errorf("expected partial content, got: %v", last.Content)
	}
}

func TestElapsedTime(t *testing.T) {
	tests := []struct {
		name    string
		start   time.Time
		wantSec int
		wantStr string
	}{
		{
			name:    "Zero elapsed time",
			start:   time.Now(),
			wantSec: 0,
			wantStr: "00:00",
		},
		{
			name:    "Exactly one minute",
			start:   time.Now().Add(-1 * time.Minute),
			wantSec: 60,
			wantStr: "01:00",
		},
		{
			name:    "More than one minute",
			start:   time.Now().Add(-90 * time.Second),
			wantSec: 90,
			wantStr: "01:30",
		},
		{
			name:    "Exactly one hour",
			start:   time.Now().Add(-1 * time.Hour),
			wantSec: 3600,
			wantStr: "60:00",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotSec, gotStr := elapsedTime(tt.start)
			if gotSec != tt.wantSec {
				t.Errorf("elapsedTime() got = %d, want %d", gotSec, tt.wantSec)
			}
			if gotStr != tt.wantStr {
				t.Errorf("elapsedTime() got = %s, want %s", gotStr, tt.wantStr)
			}
		})
	}
}

func TestTokenSpeed(t *testing.T) {
	tests := []struct {
		name string
		t    int
		s    string
		want float64
	}{
		{
			name: "Zero time",
			t:    0,
			s:    "example text",
			want: 0,
		},
		{
			name: "Zero tokens",
			t:    1,
			s:    "",
			want: 0,
		},
		{
			name: "Small number of tokens",
			t:    2,
			s:    "hello world",
			want: 1.3,
		},
		{
			name: "Large number of tokens",
			t:    10,
			s:    "this is a longer text with multiple words",
			want: 1.0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tokenSpeed(tt.t, tt.s)
			if !equalFloat64(got, tt.want, 0.0001) {
				t.Errorf("tokenSpeed(%s) got = %v, want %v", tt.s, got, tt.want)
			}
		})
	}
}

// Helper function to compare float64 values with a tolerance
func equalFloat64(a, b, tolerance float64) bool {
	return math.Abs(a-b) <= tolerance
}
