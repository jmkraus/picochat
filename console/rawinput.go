package console

import (
	"fmt"
	"io"
	"os"
	"strings"

	"golang.org/x/term"
)

// ReadMultilineInput reads multiline input from stdin. It handles raw mode, escape sequences,
// command detection, and returns an InputResult containing the entered text, flags for EOF,
// Aborted, IsCommand, and any error.
// Parameters:
//
//	none
//
// Returns:
//
//	InputResult - A structure containing the entered text and flags indicating specific states
func ReadMultilineInput() InputResult {
	in := os.Stdin
	fd := int(in.Fd())
	if !term.IsTerminal(fd) {
		// stdin is not an interactive console → pipe or file
		data, err := io.ReadAll(os.Stdin)
		if err != nil {
			return InputResult{Error: fmt.Errorf("error reading stdin: %v", err)}
		}
		text := strings.TrimRight(string(data), "\n")
		return InputResult{Text: text, EOF: true}
	} else {
		oldState, err := term.MakeRaw(fd)
		if err != nil {
			return InputResult{Error: fmt.Errorf("error enabling raw input mode: %v", err)}
		}
		defer term.Restore(fd, oldState)
	}

	var lines []string
	var currentLine []rune

	for {
		b := make([]byte, 1)
		_, err := in.Read(b)
		if err != nil {
			break
		}

		switch b[0] {
		case 3: // Ctrl+C
			fmt.Print("\r\n")
			return InputResult{Aborted: true}

		case 4: // Ctrl+D (EOF) → input finished
			fmt.Print("\r\n")
			if len(currentLine) > 0 {
				lines = append(lines, string(currentLine))
			}
			return InputResult{Text: strings.Join(lines, "\n"), EOF: false}

		case 27: // ESC or escape sequences
			// Temporarily set Stdin to non-blocking
			err = setNonblock(fd, true)
			if err != nil {
				return InputResult{Error: fmt.Errorf("could not set stdin to non-blocking mode: %w", err)}
			}

			buf := make([]byte, 3)
			n, _ := in.Read(buf)

			// Revert back to Stdin blocking mode
			err = setNonblock(fd, false)
			if err != nil {
				return InputResult{Error: fmt.Errorf("could not unset stdin to blocking mode: %w", err)}
			}

			if n == 0 {
				// Plain ESC → abort immediately
				return InputResult{Aborted: true}
			}

			if buf[0] == '[' && n > 1 {
				// Arrow keys (e.g. ESC[A, ESC[B, ESC[C, ESC[D])
				switch buf[1] {
				case 'A': // Up
					recalled := PrevCommand()
					if recalled != "" {
						fmt.Print("\r\033[K")
						fmt.Printf(">>> %s", recalled)
						currentLine = []rune(recalled)
					}
					continue
				case 'B': // Down
					recalled := NextCommand()
					fmt.Print("\r\033[K")
					fmt.Printf(">>> %s", recalled)
					currentLine = []rune(recalled)
					continue
				case 'C', 'D':
					continue
				}
			}

			// Any other escape sequence → ignore or abort
			return InputResult{Aborted: true}

		case 127: // Backspace
			if len(currentLine) > 0 {
				currentLine = currentLine[:len(currentLine)-1]
				fmt.Print("\b \b")
			}

		case 13, 10: // Enter
			line := string(currentLine)
			trimLine := strings.TrimSpace(line)

			// Input is a command
			if len(lines) == 0 && strings.HasPrefix(trimLine, "/") {
				return InputResult{Text: trimLine, IsCommand: true}
			}

			lines = append(lines, line)
			currentLine = []rune{}
			fmt.Print("\r\n")

		default:
			currentLine = append(currentLine, rune(b[0]))
			fmt.Print(string(b))
		}
	}
	return InputResult{Text: strings.Join(lines, "\n")}
}
