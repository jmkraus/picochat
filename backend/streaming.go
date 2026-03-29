package backend

import (
	"bufio"
	"fmt"
	"io"
	"strings"
)

type parseEventFn func(data string) (thinking, content string, done bool, err error)

type streamAccum struct {
	Reasoning strings.Builder
	Content   strings.Builder
}

// consumeSSEStream consumes an SSE stream body, parses each data event
// with the provided parser callback, and accumulates reasoning/content.
//
// Parameters:
//
//	body (io.Reader)                 - response stream body
//	parse (parseEventFn)             - event parser callback
//	onChunk (func(ChatChunk) error)  - optional chunk callback
//
// Returns:
//
//	ChatFinal - accumulated reasoning and content
//	error     - error if stream read/parsing/callback fails
func consumeSSEStream(
	body io.Reader,
	parse parseEventFn,
	onChunk func(ChatChunk) error,
) (ChatFinal, error) {
	reader := bufio.NewReader(body)
	var acc streamAccum

	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				break
			}
			return ChatFinal{}, fmt.Errorf("read stream failed: %w", err)
		}

		line = strings.TrimSpace(line)
		if line == "" || !strings.HasPrefix(line, "data:") {
			continue
		}

		data := strings.TrimSpace(strings.TrimPrefix(line, "data:"))
		if data == "[DONE]" {
			if onChunk != nil {
				if err := onChunk(ChatChunk{Done: true}); err != nil {
					return ChatFinal{}, err
				}
			}
			break
		}

		thinking, content, done, err := parse(data)
		if err != nil {
			return ChatFinal{}, err
		}

		if thinking != "" {
			acc.Reasoning.WriteString(thinking)
		}
		if content != "" {
			acc.Content.WriteString(content)
		}

		if onChunk != nil {
			if err := onChunk(ChatChunk{
				Thinking: thinking,
				Content:  content,
				Done:     done,
			}); err != nil {
				return ChatFinal{}, err
			}
		}
	}

	return ChatFinal{
		Reasoning: acc.Reasoning.String(),
		Content:   acc.Content.String(),
	}, nil
}
