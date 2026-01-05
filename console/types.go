package console

// struct types
type InputResult struct {
	Text      string
	IsCommand bool
	Aborted   bool
	EOF       bool
	Error     error
}

type commandHistory struct {
	entries []string
	index   int
}

// Constants
const Prompt = ">>> "
const crlf = "\r\n"

// Esc Sequences
const escClearLine = "\r\033[K"
const escCursorBack = "\033[D"
const escCursorForward = "\033[C"
const escCursorToColumn = "\033[%dG"
const escDisableLineWrap = "\033[?7l"
const escEnableLineWrap = "\033[?7h"
const escDisableCursor = "\033[?25l"
const escEnableCursor = "\033[?25h"
