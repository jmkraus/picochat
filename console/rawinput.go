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
}

func ReadMultilineInput() InputResult {
	oldState, err := term.MakeRaw(int(syscall.Stdin))
	if err != nil {
		panic(err)
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

		case 4: // Ctrl+D (EOF) → Eingabe fertig
			fmt.Print("\r\n")
			lines = append(lines, string(currentLine))
			return InputResult{Text: strings.Join(lines, "\n")}

		case 27: // ESC oder Escape-Sequenzen
			// Setze Stdin kurz auf non-blocking
			fd := int(in.Fd())
			oldMode := syscall.SetNonblock(fd, true)

			buf := make([]byte, 2)
			n, _ := in.Read(buf)

			// Setze blocking wieder zurück
			syscall.SetNonblock(fd, false)
			_ = oldMode // Rückgabewert uninteressant hier

			if n == 0 {
				// Reines ESC → sofort abbrechen
				return InputResult{Aborted: true}
			}

			if buf[0] == '[' && n > 1 {
				// Pfeiltasten (z. B. ESC[A, ESC[B, ESC[C, ESC[D])
				switch buf[1] {
				case 'A', 'B', 'C', 'D':
					continue
				}
			}

			// Irgendeine andere Escape-Sequenz → ignorieren oder abbrechen
			return InputResult{Aborted: true}

		case 127: // Backspace
			if len(currentLine) > 0 {
				currentLine = currentLine[:len(currentLine)-1]
				fmt.Print("\b \b")
			}

		case 13, 10: // Enter
			line := string(currentLine)

			// Sofortiger Command
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
