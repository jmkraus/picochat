package main

import (
	"bufio"
	"fmt"
	"golang.org/x/term"
	"os"
	"picochat/args"
	"picochat/chat"
	"picochat/command"
	"picochat/config"
	"picochat/console"
	"picochat/messages"
	"picochat/version"
	"strings"
	"syscall"
)

func sendPrompt(prompt string, cfg *config.Config, history *messages.ChatHistory) {
	history.Add(messages.RoleUser, prompt)

	msg, err := chat.HandleChat(cfg, history)
	if err != nil {
		console.Error(err)
	} else {
		console.Info(msg)
	}
}

func repeatPrompt(cfg *config.Config, history *messages.ChatHistory) {
	if history.Len() < 2 {
		console.Warn("chat history is empty")
		return
	}

	lastEntry := history.GetLast()
	if lastEntry.Role != messages.RoleUser {
		console.Warn("last entry in history is not a user prompt")
		return
	}

	msg, err := chat.HandleChat(cfg, history)
	if err != nil {
		console.Error(err)
	} else {
		console.Info(msg)
	}
}

func readMultilineInput() (string, bool) {
	// Terminal in Raw Mode schalten
	oldState, err := term.MakeRaw(int(syscall.Stdin))
	if err != nil {
		panic(err)
	}
	defer term.Restore(int(syscall.Stdin), oldState)

	in := os.Stdin
	var lines []string
	var currentLine []rune

	enterCount := 0

	for {
		b := make([]byte, 1)
		_, err := in.Read(b)
		if err != nil {
			break
		}

		switch b[0] {
		case 27: // ESC
			return "", false
		case 13: // Enter (CR)
			if enterCount == 1 {
				// Double Enter -> done
				lines = append(lines, string(currentLine))
				return strings.Join(lines, "\n"), true
			}
			enterCount++
			lines = append(lines, string(currentLine))
			currentLine = []rune{}
			fmt.Print("\n")
		case 10: // LF (Unix Enter)
			// sometimes just \n -> treat like  CR
			if enterCount == 1 {
				lines = append(lines, string(currentLine))
				return strings.Join(lines, "\n"), true
			}
			enterCount++
			lines = append(lines, string(currentLine))
			currentLine = []rune{}
			fmt.Print("\n")
		case 127: // Backspace
			if len(currentLine) > 0 {
				currentLine = currentLine[:len(currentLine)-1]
				fmt.Print("\b \b")
			}
			enterCount = 0
		default:
			currentLine = append(currentLine, rune(b[0]))
			fmt.Print(string(b))
			enterCount = 0
		}
	}

	return strings.Join(lines, "\n"), true
}

func oldReadMultilineInput() (string, bool) {
	scanner := bufio.NewScanner(os.Stdin)
	var lines []string
	firstLine := true

	for scanner.Scan() {
		line := scanner.Text()
		trimmed := strings.TrimSpace(line)

		if firstLine && strings.HasPrefix(trimmed, "/") {
			return trimmed, true // input, isCommand
		}

		switch trimmed {
		case "/cancel":
			return "", false
		case "/done", "///":
			return strings.Join(lines, "\n"), false
		default:
			lines = append(lines, line)
			firstLine = false
			result, suffix := strings.CutSuffix(trimmed, "///")
			if suffix {
				return result, false
			}
		}
	}

	if err := scanner.Err(); err != nil {
		console.Errorf("reading standard input: %v", err)
		return "", false
	}

	return strings.Join(lines, "\n"), false
}

func main() {
	args.Parse()

	if *args.ShowVersion {
		console.Info(fmt.Sprintf("picochat version is %s", version.Version))
		os.Exit(0)
	}

	err := config.Load()
	if err != nil {
		console.Errorf("load configuration failed: %v", err)
		os.Exit(1)
	}
	cfg := config.Get()

	var history *messages.ChatHistory
	if *args.HistoryFile != "" {
		history, err = messages.LoadHistoryFromFile(*args.HistoryFile)
		if err != nil {
			console.Errorf("load history failed: %v", err)
			os.Exit(1)
		}
	} else {
		history = messages.NewHistory(cfg.Prompt, cfg.Context)
	}

	console.Info("PicoChat started. Help with '/?'")

	for {
		fmt.Print("\n>>> ")

		input, isCommand := readMultilineInput()

		if input == "" && !isCommand {
			console.Info("Input canceled.")
			continue
		}

		if isCommand {
			result := command.HandleCommand(input, history, os.Stdin)
			if result.Output != "" {
				console.Info(result.Output)
			}
			if result.Error != nil {
				console.Errorf("command handler failed: %v", result.Error)
			}
			if result.Quit {
				break
			}
			if result.Repeat {
				repeatPrompt(cfg, history)
			} else if result.Prompt != "" {
				sendPrompt(result.Prompt, cfg, history)
			}
			continue
		}

		sendPrompt(input, cfg, history)
	}
}
