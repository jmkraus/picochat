package console

import (
	"fmt"
	"golang.org/x/term"
	"os"
	"strings"
	"syscall"
)

func ReadMultilineInput() (string, bool) {
	// Switch terminal to raw mode
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
			// Catch ESC sequence
			next := make([]byte, 2)
			n, _ := in.Read(next)
			if n == 0 {
				// ESC only -> cancel
				return "", false
			}
			if next[0] == '[' {
				// Detected cursor key -> ignore
				switch next[1] {
				case 'A', 'B', 'C', 'D':
					// Arrow Up/Down/Right/Left -> skip
					continue
				}
			}
			// If not cursor key, but "true ESC key": cancel
			return "", false
		case 13, 10: // Enter
			line := strings.TrimSpace(string(currentLine))

			if len(lines) == 0 && strings.HasPrefix(line, "/") {
				// First row is a command
				fmt.Print("\r\n")
				return line, true
			}

			if enterCount == 1 && line == "" {
				// Double Enter -> done
				lines = append(lines, line)
				return strings.Join(lines, "\n"), false
			}

			enterCount++
			lines = append(lines, line)
			currentLine = []rune{}
			fmt.Print("\r\n")
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
