package console

import (
	"fmt"
	"os"
)

// Error prints an error message to stderr, prefixed with "error:"
func Error(msg string) {
	fmt.Fprintf(os.Stderr, "error: %s\n", msg)
}

// Errorf prints a custom error message to stderr
func Errorf(msg string, err error) {
	fmt.Fprintf(os.Stderr, "%s: %v\n", msg, err)
}

// Warn prints a warning message to stderr, prefixed with "warning:"
func Warn(msg string) {
	fmt.Fprintf(os.Stderr, "warning: %s\n", msg)
}

// Info prints a message to stdout (normal output)
func Info(msg string) {
	fmt.Fprintln(os.Stdout, msg)
}

// Success prints a success message to stdout, prefixed with "✓"
func Success(msg string) {
	fmt.Fprintf(os.Stdout, "✓ %s\n", msg)
}
