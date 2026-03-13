package console

import (
	"fmt"
	"os"
)

const info_prefix string = BrightGreen + Bold + "✓" + ColorReset
const warn_prefix string = BrightYellow + Bold + "!" + ColorReset
const err_prefix string = BrightRed + Bold + "×" + ColorReset

// const warn_prefix string = BrightYellow + Bold + "!" + ColorReset

// Error prints a custom error message to stderr
//
// Parameters:
//
//	msg (string) - error message to be printed
//
// Returns:
//
//	none
func Error(msg string) {
	fmt.Print(ClearLine)
	fmt.Fprintf(os.Stderr, "%s %s\n", err_prefix, msg)
}

// Warn prints a warning message to stderr, prefixed with "warning:"
//
// Parameters:
//
//	msg (string) - warning message to be printed
//
// Returns:
//
//	none
func Warn(msg string) {
	fmt.Print(ClearLine)
	fmt.Fprintf(os.Stderr, "%s %s\n", warn_prefix, msg)
}

// Info prints a message to stdout (normal output)
//
// Parameters:
//
//	msg (string) - message to be printed
//
// Returns:
//
//	none
func Info(msg string) {
	fmt.Print(ClearLine)
	fmt.Fprintf(os.Stdout, "%s %s\n", info_prefix, msg)
}

// NewLine writes a newline if the current mode isn't "quiet".
//
// Parameters:
//
//	quiet (bool) - status of quiet mode
//
// Returns:
//
//	none
func NewLine(quiet bool) {
	if quiet {
		return
	}
	fmt.Println()
}

// SetCursorPos places the cursor into the given column
//
// Parameters:
//
//	col (int) - the column of the new cursor position
//
// Returns:
//
//	none
func SetCursorPos(col int) {
	fmt.Printf(CursorToColumn, col)

}

// Colorize is a Helper function for enclosing text in color Esc sequences.
// A reset Esc sequence is added to the end of string.
//
// Parameters:
//
//	color (string) - esc sequence for the color (use constants)
//	text (string)  - the text
//
// Returns:
//
//	string - text enclosed in esc sequences
func Colorize(color, text string) string {
	return color + text + ColorReset
}

// Style is a Helper function for enclosing text in font style Esc sequences.
// A reset Esc sequence is added to the end of string.
//
// Parameters:
//
//	fontstyle (string) - esc sequence for the font style (use constants)
//	text (string)      - the text
//
// Returns:
//
//	string - text enclosed in esc sequences
func Style(fontstyle, text string) string {
	return fontstyle + text + Regular
}

// colorPrint is a Helper function for printing colorized text
//
// Parameters:
//
//	color (string) - esc sequence for the color (use constants)
//	newline (bool) - flag if text output should be finished with newline
//	a (any)        - interface for arbitrary output data
//
// Returns:
//
//	none
func colorPrint(color string, newline bool, a ...any) {
	text := Colorize(color, fmt.Sprint(a...))
	if newline {
		fmt.Println(text)
	} else {
		fmt.Print(text)
	}
}

// ColorPrint is the interface for color text output without newline
//
// Parameters:
//
//	color (string) - esc sequence for the color (use constants)
//	a (any)        - interface for arbitrary output data
//
// Returns:
//
//	none
func ColorPrint(color string, a ...any) {
	colorPrint(color, false, a...)
}

// ColorPrintln is the interface for color text output with newline
//
// Parameters:
//
//	color (string) - esc sequence for the color (use constants)
//	a (any)        - interface for arbitrary output data
//
// Returns:
//
//	none
func ColorPrintln(color string, a ...any) {
	colorPrint(color, true, a...)
}
