package chat

type ChatResult struct {
	Output   string  `json:"output" yaml:"output"`
	Elapsed  string  `json:"elapsed" yaml:"elapsed"`
	TokensPS float64 `json:"tokens_per_sec" yaml:"tokens_per_sec"`
	// optional fields: Model, Config etc.
	// Model     string  `json:"model,omitempty" yaml:"model,omitempty"`
}
