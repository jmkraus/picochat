package config

import (
	"fmt"
	"picochat/paths"
	"sync"

	"github.com/BurntSushi/toml"
)

type Config struct {
	URL         string
	Model       string
	Prompt      string
	Context     int
	Temperature float64
	TopP        float64
	Reasoning   bool
	Quiet       bool

	ConfigPath string `toml:"-"`
	ImagePath  string `toml:"-"` ////IMAGES
	OutputFmt  string `toml:"-"`
	SchemaFmt  any    `toml:"-"`
}

const (
	MinCtx = 3
	MaxCtx = 100
)

var (
	instance *Config
	once     sync.Once
	loadErr  error
)

// load reads and caches the configuration.
//
// Parameters:
//
//	none
//
// Returns:
//
//	none
func load() {
	path, err := paths.GetConfigPath()
	if err != nil {
		loadErr = err
		return
	}

	var cfg = Config{
		URL:         "http://localhost:11434/api",
		Model:       "gpt-oss:latest",
		Prompt:      "You are a Large Language Model. Answer as concisely as possible. Your answers should be informative, helpful and engaging.",
		Context:     20,
		Temperature: 0.7,
		TopP:        0.9,
		Reasoning:   false,
		Quiet:       false,
	}

	if paths.FileExists(path) {
		if _, err := toml.DecodeFile(path, &cfg); err != nil {
			loadErr = fmt.Errorf("failed to decode TOML file %q: %w", path, err)
			return
		}
	} else {
		path = "No contig.toml found - fallback to internal defaults"
	}

	if cfg.Context < MinCtx || cfg.Context > MaxCtx {
		loadErr = fmt.Errorf("context size must be between %d and %d", MinCtx, MaxCtx)
		return
	}

	cfg.ConfigPath = path
	instance = &cfg
}

// Get loads the configuration once and returns the instance of the
// loaded configuration.
//
// Parameters:
//
//	none
//
// Returns:
//
//	*Config - pointer to the loaded configuration
//	error   - error if any
func Get() (*Config, error) {
	once.Do(load)
	return instance, loadErr
}

// Set allows changing a specific parameter after loading.
//
// Parameters:
//
//	key (string) - the configuration key to modify
//	value (any)  - the new value for the key
//
// Returns:
//
//	error - an error if the key is unsupported or the value has the wrong type
func Set(key string, value any) error {
	cfg, err := Get()
	if err != nil {
		return fmt.Errorf("cannot apply config change: %w", err)
	}

	switch key {
	case "context":
		intVal, ok := value.(int)
		if !ok {
			return fmt.Errorf("value for key '%s' must be an integer", key)
		}
		cfg.Context = intVal
		return nil

	case "model":
		strVal, ok := value.(string)
		if !ok {
			return fmt.Errorf("value for key '%s' must be a string", key)
		}
		cfg.Model = strVal
		return nil

	case "reasoning":
		boolVal, ok := value.(bool)
		if !ok {
			return fmt.Errorf("value for key '%s' must be a boolean", key)
		}
		cfg.Reasoning = boolVal
		return nil

	case "temperature":
		floatVal, ok := value.(float64)
		if !ok {
			return fmt.Errorf("value for key '%s' must be a float", key)
		}
		cfg.Temperature = floatVal
		return nil

	case "top_p":
		floatVal, ok := value.(float64)
		if !ok {
			return fmt.Errorf("value for key '%s' must be a float", key)
		}
		cfg.TopP = floatVal
		return nil

	default:
		// Don't forget to update command/parser.go --> validateAndConvert()
		return fmt.Errorf("unsupported config key '%s'", key)
	}
}
