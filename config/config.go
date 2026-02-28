package config

import (
	"fmt"
	"picochat/envs"
	"picochat/paths"
	"strconv"
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

	// 1. Default values
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

	// 2. Config file
	if paths.FileExists(path) {
		if _, err := toml.DecodeFile(path, &cfg); err != nil {
			loadErr = fmt.Errorf("decode toml file %q failed: %w", path, err)
			return
		}
	} else {
		path = "No config.toml found - fallback to internal defaults"
	}

	// 3. Environment variables
	if err := applyEnvOverrides(&cfg); err != nil {
		loadErr = fmt.Errorf("set config with env vars failed: %w", err)
		return
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
//	error - error if any
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

// applyEnvOverrides checks all environment variables if set
// and updates the respective config entry accordingly.
//
// Parameters:
//
//	cfg (*Config) - the instance of the Config struct
//
// Returns:
//
//	error - error if any
func applyEnvOverrides(cfg *Config) error {
	if v := envs.GetEnv(envs.PICOCHAT_URL); v != "" {
		cfg.URL = v
	}
	if v := envs.GetEnv(envs.PICOCHAT_MODEL); v != "" {
		cfg.Model = v
	}
	if v := envs.GetEnv(envs.PICOCHAT_CONTEXT); v != "" {
		n, err := strconv.Atoi(v)
		if err != nil {
			return fmt.Errorf("invalid %s %q: %w", envs.PICOCHAT_CONTEXT, v, err)
		}
		cfg.Context = n
	}
	if v := envs.GetEnv(envs.PICOCHAT_TEMPERATURE); v != "" {
		f, err := strconv.ParseFloat(v, 64)
		if err != nil {
			return fmt.Errorf("invalid %s %q: %w", envs.PICOCHAT_TEMPERATURE, v, err)
		}
		cfg.Temperature = f
	}
	if v := envs.GetEnv(envs.PICOCHAT_TOP_P); v != "" {
		f, err := strconv.ParseFloat(v, 64)
		if err != nil {
			return fmt.Errorf("invalid %s %q: %w", envs.PICOCHAT_TOP_P, v, err)
		}
		cfg.TopP = f
	}
	if v := envs.GetEnv(envs.PICOCHAT_REASONING); v != "" {
		b, err := parseBool01(v)
		if err != nil {
			return fmt.Errorf("invalid %s %q: %w", envs.PICOCHAT_REASONING, v, err)
		}
		cfg.Reasoning = b
	}
	if v := envs.GetEnv(envs.PICOCHAT_QUIET); v != "" {
		b, err := parseBool01(v)
		if err != nil {
			return fmt.Errorf("invalid %s %q: %w", envs.PICOCHAT_QUIET, v, err)
		}
		cfg.Quiet = b
	}

	return nil
}

// parseBool01 is a helper function, checking for 0 or 1
// and returning a matching boolean value.
//
// Parameters:
//
//	s (string) - the value to be parsed
//
// Returns:
//
//	bool  - the parsed boolean value
//	error - error if any
func parseBool01(s string) (bool, error) {
	switch s {
	case "0":
		return false, nil
	case "1":
		return true, nil
	default:
		return false, fmt.Errorf("expected 0 or 1")
	}
}
