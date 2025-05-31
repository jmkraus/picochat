package config

type Config struct {
	URL     string
	Model   string
	Context int
	Prompt  string

	Temperature float64
	TopP        float64
}
