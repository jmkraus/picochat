package types

type ChatRequest struct {
	Model    string    `json:"model"`
	Messages []Message `json:"messages"`
	Stream   bool      `json:"stream"`
}

type StreamResponse struct {
	Message         Message `json:"message"`
	Done            bool    `json:"done"`
	PromptEvalCount int     `json:"prompt_eval_count"`
	EvalCount       int     `json:"eval_count"`
}

type CommandResult struct {
	Output     string
	Quit       bool
	NewHistory *ChatHistory
	Error      error
}

// temporary solution (it's a var not a type)
// TODO: refactoring to a dedicated caches module
var SelectionCache = struct {
	HistoryFiles []string
	Models       []string
}{}
