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
	ClearLine       string = "\r\033[K"
	CursorBack      string = "\033[D"
	CursorForward   string = "\033[C"
	CursorToColumn  string = "\033[%dG"
	DisableLineWrap string = "\033[?7l"
	EnableLineWrap  string = "\033[?7h"
	DisableCursor   string = "\033[?25l"
	EnableCursor    string = "\033[?25h"
)

const (
	// Foreground colors
	Black   string = "\033[30m"
	Red     string = "\033[31m"
	Green   string = "\033[32m"
	Yellow  string = "\033[33m"
	Blue    string = "\033[34m"
	Magenta string = "\033[35m"
	Cyan    string = "\033[36m"
	White   string = "\033[37m"

	// Bright foreground colors
	BrightBlack   string = "\033[90m"
	BrightRed     string = "\033[91m"
	BrightGreen   string = "\033[92m"
	BrightYellow  string = "\033[93m"
	BrightBlue    string = "\033[94m"
	BrightMagenta string = "\033[95m"
	BrightCyan    string = "\033[96m"
	BrightWhite   string = "\033[97m"

	// Background colors
	BgBlack   string = "\033[40m"
	BgRed     string = "\033[41m"
	BgGreen   string = "\033[42m"
	BgYellow  string = "\033[43m"
	BgBlue    string = "\033[44m"
	BgMagenta string = "\033[45m"
	BgCyan    string = "\033[46m"
	BgWhite   string = "\033[47m"

	// Bright background colors
	BgBrightBlack   string = "\033[100m"
	BgBrightRed     string = "\033[101m"
	BgBrightGreen   string = "\033[102m"
	BgBrightYellow  string = "\033[103m"
	BgBrightBlue    string = "\033[104m"
	BgBrightMagenta string = "\033[105m"
	BgBrightCyan    string = "\033[106m"
	BgBrightWhite   string = "\033[107m"

	// Reset
	ColorReset string = "\033[0m"
)
