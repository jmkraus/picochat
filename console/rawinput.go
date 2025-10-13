package console

import (
	"fmt"
	"os"
	"strings"
	"syscall"

	"golang.org/x/term"
)

type InputResult struct {
	Text      string
	IsCommand bool
	Aborted   bool
	Error     error
}

func ReadMultilineInput() InputResult {
	oldState, err := term.MakeRaw(int(syscall.Stdin))
	if err != nil {
		return InputResult{Error: fmt.Errorf("Switching to terminal raw mode failed.")}
	}
	defer term.Restore(int(syscall.Stdin), oldState)

	in := os.Stdin
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
			fmt.Print("^C")
			return InputResult{Aborted: true}

		case 4: // Ctrl+D (EOF) → input finished
			fmt.Print("\r\n")
			lines = append(lines, string(currentLine))
			return InputResult{Text: strings.Join(lines, "\n")}

		case 27: // ESC or escape sequences
			// Temporarily set Stdin to non-blocking
			fd := int(in.Fd())
			oldMode := syscall.SetNonblock(fd, true)

			buf := make([]byte, 2)
			n, _ := in.Read(buf)

			// Restore blocking mode
			syscall.SetNonblock(fd, false)
			_ = oldMode // return value not relevant here

			if n == 0 {
				// Plain ESC → abort immediately
				return InputResult{Aborted: true}
			}

			if buf[0] == '[' && n > 1 {
				// Arrow keys (e.g. ESC[A, ESC[B, ESC[C, ESC[D])
				switch buf[1] {
				case 'A', 'B', 'C', 'D':
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

			// Input is a command
			if len(lines) == 0 && strings.HasPrefix(line, "/") {
				return InputResult{Text: line, IsCommand: true}
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
