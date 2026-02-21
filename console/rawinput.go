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

// PromptWidth returns the width of the predefined prompt symbols
// as correct runewidth calculation.
//
// Pramaters:
//
//	none
//
// Returns:
//
//	int - the width of the symbols
func PromptWidth() int {
	return runewidth.StringWidth(Prompt)
}

// ReadMultilineInput reads multiline input from stdin. It handles raw mode,
// ape sequences, command detection, and returns an InputResult containing
// the entered text, flags for EOF, Aborted, IsCommand, and any error.
//
// Parameters:
//
//	none
//
// Returns:
//
//	InputResult - A structure containing entered text and specific states
func ReadMultilineInput() InputResult {
	in := os.Stdin
	fd := int(in.Fd())
	if !term.IsTerminal(fd) {
		// stdin is not an interactive console → pipe or file
		data, err := io.ReadAll(in)
		if err != nil {
			return InputResult{Error: fmt.Errorf("read stdin failed: %w", err)}
		}
		text := strings.TrimRight(string(data), "\n")
		trimText := strings.TrimSpace(text)
		return InputResult{Text: trimText, EOF: true, IsCommand: strings.HasPrefix(trimText, "/")}
	} else {
		oldState, err := term.MakeRaw(fd)
		if err != nil {
			return InputResult{Error: fmt.Errorf("enable raw input mode failed: %w", err)}
		}
		fmt.Print(DisableLineWrap)
		defer func() {
			fmt.Print(EnableLineWrap)
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
			if firstLine {
				fmt.Print(ClearLine + Prompt)
			}
			return InputResult{Aborted: true}

		case 4: // Ctrl+D (EOF) → input finished
			fmt.Println()
			if len(currentLine) > 0 {
				lines = append(lines, string(currentLine))
			}
			return InputResult{Text: strings.Join(lines, "\n"), EOF: false}

		case 27: //  or Escape sequences
			// Peek to see if there are more bytes (Escape sequence)
			_ = setNonblock(fd, true)
			peekBuf, err := reader.Peek(2)
			_ = setNonblock(fd, false)
			if err != nil || len(peekBuf) < 2 {
				// Plain  → abort immediately
				if firstLine {
					fmt.Print(ClearLine + Prompt)
				}
				return InputResult{Aborted: true}
			}

			// Read the Escape sequence
			buf := make([]byte, 2)
			n, _ := reader.Read(buf)

			if n < 2 || buf[0] != '[' {
				// Unknown Escape sequence
				continue
			}

			// Arrow keys (e.g. [A, [B, [C, [D])
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
						fmt.Print(CursorForward)
					}
				}
				continue
			case 'D': // Left
				if cursorPos > 0 {
					cursorPos--
					width := runewidth.RuneWidth(currentLine[cursorPos])
					// Move cursor back by visual width of char
					for range width {
						fmt.Print(CursorBack)
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

			// first line is empty
			if firstLine && len(trimLine) == 0 {
				fmt.Print(ClearLine + Prompt)
			}

			lines = append(lines, line)
			currentLine = []rune{}
			cursorPos = 0
			fmt.Print("\r\n") // Println not sufficient here
		default:
			currentLine, cursorPos = insertCharAt(currentLine, cursorPos, rune(r))
			updateCurrentLine(currentLine, firstLine, cursorPos)
		}
	}

	return InputResult{Text: strings.Join(lines, "\n")}
}

// deleteCharAt deletes the character at the cursor position in the current line.
// If the cursor is at the beginning (pos 0), no deletion occurs.
//
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
//
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
//
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
	fmt.Print(ClearLine)

	prefix := ""
	prefixWidth := 0
	if firstLine {
		prefix = Prompt
		prefixWidth = runewidth.StringWidth(Prompt)
	}

	visualPos := visualWidth(line, cursorPos)

	if firstLine && len(line) == 0 {
		fmt.Print(prefix + Shadow)
	} else {
		fmt.Print(prefix + string(line))
	}
	fmt.Printf(CursorToColumn, visualPos+prefixWidth+1)
}

// visualWidth calculates the visual display width of a rune slice up to
// a given position. This accounts for characters that occupy multiple
// terminal columns (e.g., CJK characters, emojis) or zero columns
// (e.g., combining characters).
//
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

// getTerminalWidth returns the current width of the terminal window
//
// Parameters:
//
//	fd (int) - a file descriptor
//
// Returns:
//
//	int - width of the window
func getTerminalWidth(fd int) int {
	width, _, err := term.GetSize(fd)
	if err != nil {
		// Fallback to a default width if terminal size cannot be determined
		return 80
	}
	return width
}
