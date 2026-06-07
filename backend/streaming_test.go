package backend

import (
	"fmt"
	"strings"
	"testing"
)

func TestConsumeSSEStream(t *testing.T) {
	parse := func(data string) (thinking, content string, done bool, err error) {
		switch data {
		case "ev1":
			return "t1", "c1", false, nil
		case "ev2":
			return "", "c2", true, nil
		default:
			return "", "", false, fmt.Errorf("unexpected event: %s", data)
		}
	}

	stream := strings.NewReader(strings.Join([]string{
		"event: message",
		"data: ev1",
		"",
		"data: ev2",
		"data: [DONE]",
		"",
	}, "\n"))

	var chunks []ChatChunk
	final, err := consumeSSEStream(stream, parse, func(c ChatChunk) error {
		chunks = append(chunks, c)
		return nil
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if final.Reasoning != "t1" {
		t.Fatalf("final reasoning = %q, want %q", final.Reasoning, "t1")
	}
	if final.Content != "c1c2" {
		t.Fatalf("final content = %q, want %q", final.Content, "c1c2")
	}
	if len(chunks) != 3 {
		t.Fatalf("len(chunks) = %d, want 3", len(chunks))
	}
	if !chunks[2].Done {
		t.Fatalf("expected last chunk done=true, got %+v", chunks[2])
	}
}

func TestConsumeSSEStream_ParseError(t *testing.T) {
	parse := func(data string) (thinking, content string, done bool, err error) {
		return "", "", false, fmt.Errorf("parse failed")
	}
	stream := strings.NewReader("data: broken\n")

	_, err := consumeSSEStream(stream, parse, nil)
	if err == nil || !strings.Contains(err.Error(), "parse failed") {
		t.Fatalf("expected parse error, got %v", err)
	}
}

func TestConsumeSSEStream_OnChunkError(t *testing.T) {
	parse := func(data string) (thinking, content string, done bool, err error) {
		return "", "x", false, nil
	}
	stream := strings.NewReader("data: ev\n")

	_, err := consumeSSEStream(stream, parse, func(c ChatChunk) error {
		return fmt.Errorf("chunk callback failed")
	})
	if err == nil || !strings.Contains(err.Error(), "chunk callback failed") {
		t.Fatalf("expected onChunk error, got %v", err)
	}
}

func TestConsumeSSEStream_DoneChunkCallbackError(t *testing.T) {
	parse := func(data string) (thinking, content string, done bool, err error) {
		return "", "", false, nil
	}
	stream := strings.NewReader("data: [DONE]\n")

	_, err := consumeSSEStream(stream, parse, func(c ChatChunk) error {
		if c.Done {
			return fmt.Errorf("done callback failed")
		}
		return nil
	})
	if err == nil || !strings.Contains(err.Error(), "done callback failed") {
		t.Fatalf("expected done callback error, got %v", err)
	}
}
