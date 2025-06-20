package console

import (
	"fmt"
	"os"
)

const err_prefix string = "× "  // or error: "
const warn_prefix string = "! " // or "warning: "

// Error prints an error message to stderr, prefixed with "error:"
func Error(msg string) {
	fmt.Fprintf(os.Stderr, "%s%s\n", err_prefix, msg)
}

// Errorf prints a custom error message to stderr
func Errorf(msg string, err error) {
	fmt.Fprintf(os.Stderr, "%s%s: %v\n", err_prefix, msg, err)
}

// Warn prints a warning message to stderr, prefixed with "warning:"
func Warn(msg string) {
	fmt.Fprintf(os.Stderr, "%s%s\n", warn_prefix, msg)
}

// Info prints a message to stdout (normal output)
func Info(msg string) {
	fmt.Fprintln(os.Stdout, msg)
}
