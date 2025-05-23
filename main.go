package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"picochat/args"
	"picochat/chat"
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
		if err := chat.HandleChat(cfg, history); err != nil {
			log.Fatalf("Chat error: %v", err)
		}
	}
}
