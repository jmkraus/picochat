package chat

import "picochat/messages"

type ChatRequest struct {
	Model     string             `json:"model"`
	Messages  []messages.Message `json:"messages"`
	Reasoning *Reasoning         `json:"reasoning,omitempty"`
	Options   *ChatOptions       `json:"options,omitempty"`
	Stream    bool               `json:"stream"`
	Think     bool               `json:"think,omitempty"`
	Format    any                `json:"format,omitempty"`
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
	Message         messages.Message `json:"message"`
	Done            bool             `json:"done"`
	PromptEvalCount int              `json:"prompt_eval_count"`
	EvalCount       int              `json:"eval_count"`
}

type ChatResult struct {
	Output   string  `json:"output" yaml:"output"`
	Elapsed  string  `json:"elapsed" yaml:"elapsed"`
	TokensPS float64 `json:"tokens_per_sec" yaml:"tokens_per_sec"`
	// optional fields: Model, Config etc.
	// Model     string  `json:"model,omitempty" yaml:"model,omitempty"`
}
