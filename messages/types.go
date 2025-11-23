package messages

type ChatRequest struct {
	Model     string       `json:"model"`
	Messages  []Message    `json:"messages"`
	Reasoning *Reasoning   `json:"reasoning,omitempty"`
	Options   *ChatOptions `json:"options,omitempty"`
	Stream    bool         `json:"stream"`
	Think     bool         `json:"think,omitempty"`
}

type Reasoning struct {
	Effort string `json:"effort,omitempty"`
}

type ChatOptions struct {
	Temperature   float64 `json:"temperature,omitempty"`
	TopP          float64 `json:"top_p,omitempty"`
	TopK          int     `json:"top_k,omitempty"`
	RepeatPenalty float64 `json:"repeat_penalty,omitempty"`
	NumCtx        int     `json:"num_ctx,omitempty"`
}

type StreamResponse struct {
	Message         Message `json:"message"`
	Done            bool    `json:"done"`
	PromptEvalCount int     `json:"prompt_eval_count"`
	EvalCount       int     `json:"eval_count"`
}
