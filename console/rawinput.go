package console

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/mattn/go-runewidth"
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
		fmt.Print(escDisableLineWrap)
		defer func() {
			fmt.Print(escEnableLineWrap)
			term.Restore(fd, oldState)
		}()
	}

	var lines []string
	var currentLine []rune
	var cursorPos int

	reader := bufio.NewReader(in)

	for {
		r, _, err := reader.ReadRune()
		if err != nil {
			break
		}

		firstLine := len(lines) == 0

		switch r {
		case 3: // Ctrl+C
			fmt.Print(crlf)
			return InputResult{Aborted: true}

		case 4: // Ctrl+D (EOF) → input finished
			fmt.Print(crlf)
			if len(currentLine) > 0 {
				lines = append(lines, string(currentLine))
			}
			return InputResult{Text: strings.Join(lines, "\n"), EOF: false}

		case 27: // ESC or escape sequences
			// Peek to see if there are more bytes (escape sequence)
			_ = setNonblock(fd, true)
			peekBuf, err := reader.Peek(2)
			_ = setNonblock(fd, false)
			if err != nil || len(peekBuf) < 2 {
				// Plain ESC → abort immediately
				return InputResult{Aborted: true}
			}

			// Read the escape sequence
			buf := make([]byte, 2)
			n, _ := reader.Read(buf)

			if n < 2 || buf[0] != '[' {
				// Unknown escape sequence
				continue
			}

			// Arrow keys (e.g. ESC[A, ESC[B, ESC[C, ESC[D])
			switch buf[1] {
			case 'A': // Up
				recalled := PrevCommand()
				if recalled != "" {
					currentLine = []rune(recalled)
					cursorPos = len(currentLine)
					updateCurrentLine(currentLine, firstLine, cursorPos)
				}
				continue
			case 'B': // Down
				recalled := NextCommand()
				currentLine = []rune(recalled)
				cursorPos = len(currentLine)
				updateCurrentLine(currentLine, firstLine, cursorPos)
				continue
			case 'C': // Right
				if cursorPos < len(currentLine) {
					width := runewidth.RuneWidth(currentLine[cursorPos])
					cursorPos++
					// Move cursor forward by visual width of char
					for range width {
						fmt.Print(escCursorForward)
					}
				}
				continue

			case 'D': // Left
				if cursorPos > 0 {
					cursorPos--
					width := runewidth.RuneWidth(currentLine[cursorPos])
					// Move cursor back by visual width of char
					for range width {
						fmt.Print(escCursorBack)
					}
				}
				continue
			}
			continue // ignore everything else
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
			fmt.Print(crlf)
		default:
			currentLine, cursorPos = insertCharAt(currentLine, cursorPos, rune(r))
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
	fmt.Printf("\r%s", escClearLine)

	prefix := ""
	prefixWidth := 0
	if firstLine {
		prefix = Prompt
		prefixWidth = runewidth.StringWidth(Prompt)
	}

	visualPos := visualWidth(line, cursorPos)

	fmt.Printf("%s%s", prefix, string(line))
	fmt.Printf(escCursorToColumn, visualPos+prefixWidth+1)
}

// visualWidth calculates the visual display width of a rune slice up to a given position.
// This accounts for characters that occupy multiple terminal columns (e.g., CJK characters,
// emojis) or zero columns (e.g., combining characters).
// Parameters:
//
//	line ([]rune) - the line of runes to measure
//	pos (int)     - the position up to which to calculate the width
//
// Returns:
//
//	int - the total visual width in terminal columns
func visualWidth(line []rune, pos int) int {
	if pos > len(line) {
		pos = len(line)
	}
	width := 0
	for _, r := range line[:pos] {
		width += runewidth.RuneWidth(r)
	}
	return width
}
