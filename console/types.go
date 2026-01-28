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

const (
	// Cursor controls
	EscClearLine       string = "\r\033[K"
	EscCursorBack      string = "\033[D"
	EscCursorForward   string = "\033[C"
	EscCursorToColumn  string = "\033[%dG"
	EscDisableLineWrap string = "\033[?7l"
	EscEnableLineWrap  string = "\033[?7h"
	EscDisableCursor   string = "\033[?25l"
	EscEnableCursor    string = "\033[?25h"
)

const (
	// Foreground colors
	EscBlack   string = "\033[30m"
	EscRed     string = "\033[31m"
	EscGreen   string = "\033[32m"
	EscYellow  string = "\033[33m"
	EscBlue    string = "\033[34m"
	EscMagenta string = "\033[35m"
	EscCyan    string = "\033[36m"
	EscWhite   string = "\033[37m"

	// Bright foreground colors
	EscBrightBlack   string = "\033[90m"
	EscBrightRed     string = "\033[91m"
	EscBrightGreen   string = "\033[92m"
	EscBrightYellow  string = "\033[93m"
	EscBrightBlue    string = "\033[94m"
	EscBrightMagenta string = "\033[95m"
	EscBrightCyan    string = "\033[96m"
	EscBrightWhite   string = "\033[97m"

	// Background colors
	EscBgBlack   string = "\033[40m"
	EscBgRed     string = "\033[41m"
	EscBgGreen   string = "\033[42m"
	EscBgYellow  string = "\033[43m"
	EscBgBlue    string = "\033[44m"
	EscBgMagenta string = "\033[45m"
	EscBgCyan    string = "\033[46m"
	EscBgWhite   string = "\033[47m"

	// Bright background colors
	EscBgBrightBlack   string = "\033[100m"
	EscBgBrightRed     string = "\033[101m"
	EscBgBrightGreen   string = "\033[102m"
	EscBgBrightYellow  string = "\033[103m"
	EscBgBrightBlue    string = "\033[104m"
	EscBgBrightMagenta string = "\033[105m"
	EscBgBrightCyan    string = "\033[106m"
	EscBgBrightWhite   string = "\033[107m"

	// Reset
	EscColorReset string = "\033[0m"
)
