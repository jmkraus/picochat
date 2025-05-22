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
	"picochat/args"
	"picochat/command"
	"picochat/config"
	"picochat/types"
	"picochat/version"
	"strings"
)

func readMultilineInput() (string, bool, bool) {
	scanner := bufio.NewScanner(os.Stdin)
	var lines []string
	firstLine := true

	for scanner.Scan() {
		line := scanner.Text()
		trimmed := strings.TrimSpace(line)

		if firstLine && strings.HasPrefix(trimmed, "/") {
			return trimmed, true, false // input, isCommand, quit=false
		}

		if trimmed == "/done" {
			break
		}

		lines = append(lines, line)
		firstLine = false
	}

	return strings.Join(lines, "\n"), false, false
}

func main() {
	args.Parse()

	if *args.ShowVersion {
		fmt.Println("picochat version is", version.Version)
		os.Exit(0)
	}

	err := config.Load()
	if err != nil {
		log.Fatalf("Error while loading configuration: %v", err)
	}
	cfg := config.Get()

	var history *types.ChatHistory
	if *args.HistoryFile != "" {
		history, err = types.LoadHistoryFromFile(*args.HistoryFile)
		if err != nil {
			log.Fatalf("Could not load history: %v", err)
		}
	} else {
		history = types.NewHistory(cfg.Prompt, cfg.Context)
	}

	log.Println("Chat with PicoAI started. Help with '/?'.")

	for {
		fmt.Print("\n>>> ")

		input, isCommand, quit := readMultilineInput()
		if quit {
			break
		}

		if isCommand {
			result := command.Handle(input, history, os.Stdin)
			if result.Output != "" {
				fmt.Println(result.Output)
			}
			if result.Quit {
				log.Println("Chat has ended.")
				break
			}
			continue
		}

		history.Add("user", input)

		reqBody := types.ChatRequest{
			Model:    cfg.Model,
			Messages: history.Messages,
			Stream:   true,
		}

		jsonData, err := json.Marshal(reqBody)
		if err != nil {
			fmt.Println("Error while Marshaling:", err)
			continue
		}

		resp, err := http.Post(cfg.URL+"/chat", "application/json", bytes.NewBuffer(jsonData))
		if err != nil {
			fmt.Println("Error while HTTP-Request:", err)
			continue
		}
		defer resp.Body.Close()

		decoder := json.NewDecoder(resp.Body)
		var fullReply strings.Builder
		var promptEvalCount, evalCount int
		fmt.Println("Assistant: ")

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
				promptEvalCount = res.PromptEvalCount
				evalCount = res.EvalCount
				if promptEvalCount != 0 && evalCount != 0 {
					// Ignore token output if zero (PicoAI doesn't have them)
					fmt.Println()
					log.Printf("Token stats: prompt_eval_count=%d, eval_count=%d", promptEvalCount, evalCount)
				}
				break
			}
		}

		fmt.Println()
		history.Add("assistant", fullReply.String())
	}
}
