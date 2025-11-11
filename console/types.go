package console

type InputResult struct {
	Text      string
	IsCommand bool
	Aborted   bool
	EOF       bool
	Error     error
}

type CommandHistory struct {
	entries []string
	index   int
}

const Prompt = ">>> "
