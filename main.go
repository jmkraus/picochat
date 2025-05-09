package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"picochat/command"
	"picochat/config"
	"picochat/types"
	"strings"
)

func readMultilineInput() (string, bool, bool) {
	scanner := bufio.NewScanner(os.Stdin)
	var lines []string
	first := true

	for scanner.Scan() {
		line := scanner.Text()
		trimmed := strings.TrimSpace(line)

		if first && strings.HasPrefix(trimmed, "/") {
			return trimmed, true, false // input, isCommand, quit=false
		}

		if trimmed == "/done" {
			break
		}

		lines = append(lines, line)
		first = false
	}

	return strings.Join(lines, "\n"), false, false
}

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Error while loading configuration: %v", err)
	}

	history := []types.Message{
		{Role: "system", Content: cfg.Prompt},
	}

	log.Println("Chat with PicoAI started. Exit with '/bye'.")

	for {
		fmt.Print("\n>>> ")

		input, isCommand, quit := readMultilineInput()
		if quit {
			break
		}

		if isCommand {
			output, quit := command.Handle(input, &history)
			if output != "" {
				fmt.Println(output)
			}
			if quit {
				log.Println("Chat has ended.")
				break
			}
			continue
		}

		history = append(history, types.Message{Role: "user", Content: input})

		reqBody := types.ChatRequest{
			Model:    cfg.Model,
			Messages: history,
			Stream:   true,
		}

		jsonData, err := json.Marshal(reqBody)
		if err != nil {
			fmt.Println("Error while Marshaling:", err)
			continue
		}

		resp, err := http.Post(cfg.URL, "application/json", bytes.NewBuffer(jsonData))
		if err != nil {
			fmt.Println("Error while HTTP-Request:", err)
			continue
		}
		defer resp.Body.Close()

		decoder := json.NewDecoder(resp.Body)
		var fullReply strings.Builder
		// var promptEvalCount, evalCount int
		fmt.Print("Assistant: ")

		for {
			var res types.StreamResponse
			if err := decoder.Decode(&res); err != nil {
				if err == io.EOF {
					break
				}
				fmt.Println("\nError while Stream-Decoding:", err)
				break
			}

			fmt.Print(res.Message.Content)
			fullReply.WriteString(res.Message.Content)

			if res.Done {
				// promptEvalCount = res.PromptEvalCount
				// evalCount = res.EvalCount
				// fmt.Printf("\nToken stats: prompt_eval_count=%d, eval_count=%d\n", promptEvalCount, evalCount)
				raw, err := json.MarshalIndent(res, "", "  ")
				if err != nil {
					fmt.Println("Error marshaling final response:", err)
				} else {
					fmt.Printf("\n[i] Final response (raw):\n%s\n", string(raw))
				}
				break
			}
		}

		fmt.Println()
		history = append(history, types.Message{Role: "assistant", Content: fullReply.String()})
	}
}
