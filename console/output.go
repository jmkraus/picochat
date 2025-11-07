package console

import (
	"fmt"
	"os"
)

const err_prefix string = "×"  // or "error:"
const warn_prefix string = "!" // or "warning:"

// Error is a simplified wrapper for Errorf
// Parameters:
//
//	err (error) - error to be printed
//
// Returns:
//
//	none
func Error(err error) {
	Errorf("%v", err)
}

// Errorf prints a custom error message to stderr
// Parameters:
//
//	format (string) - format string
//	args (any)      - arguments for the format string
//
// Returns:
//
//	none
func Errorf(format string, args ...any) {
	msg := fmt.Sprintf(format, args...)
	fmt.Fprintf(os.Stderr, "%s %s\n", err_prefix, msg)
}

// Warn prints a warning message to stderr, prefixed with "warning:"
// Parameters:
//
//	msg (string) - warning message to be printed
//
// Returns:
//
//	none
func Warn(msg string) {
	fmt.Fprintf(os.Stderr, "%s %s\n", warn_prefix, msg)
}

// Info prints a message to stdout (normal output)
// Parameters:
//
//	msg (string) - message to be printed
//
// Returns:
//
//	none
func Info(msg string) {
	fmt.Fprintln(os.Stdout, msg)
}
