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
	var cursorPos int

	for {
		b := make([]byte, 1)
		_, err := in.Read(b)
		if err != nil {
			break
		}

		firstLine := len(lines) == 0

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
						fmt.Printf("%s%s", Prompt, recalled)
						currentLine = []rune(recalled)
					}
					continue
				case 'B': // Down
					recalled := NextCommand()
					fmt.Print("\r\033[K")
					fmt.Printf("%s%s", Prompt, recalled)
					currentLine = []rune(recalled)
					continue
				case 'C': // Right
					if cursorPos < len(currentLine) {
						cursorPos++
						fmt.Print("\033[C") // move cursor right
					}
					continue
				case 'D': // Left
					if cursorPos > 0 {
						cursorPos--
						fmt.Print("\033[D") // move cursor left
					}
					continue
				}
			}

			// Any other escape sequence → ignore or abort
			return InputResult{Aborted: true}

		case 127: // Backspace
			if cursorPos > 0 {
				currentLine, cursorPos = deleteCharAt(currentLine, cursorPos)
				updateCurrentLine(currentLine, firstLine, cursorPos)
			}
		case 13, 10: // Enter
			line := string(currentLine)
			trimLine := strings.TrimSpace(line)

			// Input is a command
			if firstLine && strings.HasPrefix(trimLine, "/") {
				return InputResult{Text: trimLine, IsCommand: true}
			}

			lines = append(lines, line)
			currentLine = []rune{}
			cursorPos = 0
			fmt.Print("\r\n")

		default:
			currentLine, cursorPos = insertCharAt(currentLine, cursorPos, rune(b[0]))
			updateCurrentLine(currentLine, firstLine, cursorPos)
		}
	}
	return InputResult{Text: strings.Join(lines, "\n")}
}

// deleteCharAt deletes the character at the cursor position in the current line.
// If the cursor is at the beginning (pos 0), no deletion occurs.
// Parameters:
//
//	line ([]rune)   - the current line to be edited
//	cursorPos (int) - the actual cursor position in the line
//
// Returns:
//
//	[]rune - the modified line
//	int    - new cursor position
func deleteCharAt(line []rune, cursorPos int) ([]rune, int) {
	if cursorPos == 0 || len(line) == 0 {
		return line, cursorPos
	}

	// Split at cursor position
	before := line[:cursorPos]
	after := line[cursorPos:]

	// Remove last character from before part
	if len(before) > 0 {
		before = before[:len(before)-1]
	}

	// Merge back together
	newLine := append(before, after...)
	newCursorPos := cursorPos - 1

	return newLine, newCursorPos
}

// insertCharAt inserts a character at the cursor position in the current line.
// Parameters:
//
//	line ([]rune)   - the current line to be edited
//	cursorPos (int) - the actual cursor position in the line
//	char (rune)     - the single character to be inserted
//
// Returns:
//
//	[]rune - the modified line
//	int    - new cursor position
func insertCharAt(line []rune, cursorPos int, char rune) ([]rune, int) {
	if cursorPos > len(line) {
		cursorPos = len(line)
	}

	// Split at cursor position
	before := line[:cursorPos]
	after := line[cursorPos:]

	// Create new slice with capacity for one more character
	newLine := make([]rune, 0, len(line)+1)
	newLine = append(newLine, before...)
	newLine = append(newLine, char)
	newLine = append(newLine, after...)

	newCursorPos := cursorPos + 1

	return newLine, newCursorPos
}

// updateCurrentLine redraws the current line after an edit
// Parameters:
//
//	line ([]rune)    - the current line to be drawn
//	firstLine (bool) - is the line the first line (true / false)
//	cursorPos        - the cursor position for the redraw
//
// Returns:
//
//	none
func updateCurrentLine(line []rune, firstLine bool, cursorPos int) {
	// draw line
	fmt.Print("\r\033[K") // cursor to beginning

	if firstLine {
		// first line: with prompt
		fmt.Printf("%s%s", Prompt, string(line))
		fmt.Printf("\033[%dG", cursorPos+len(Prompt)+1)
	} else {
		// subsequent lines: without prompt
		fmt.Printf("%s", string(line))
		fmt.Printf("\033[%dG", cursorPos+1)
	}
}
