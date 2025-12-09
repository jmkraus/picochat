package config

type Config struct {
	URL         string
	Model       string
	Context     int
	Prompt      string
	Temperature float64
	TopP        float64
	Reasoning   bool
	Quiet       bool

	FilePath string `toml:"-"`
	////IMAGES
	ImagePath string `toml:"-"`
}
