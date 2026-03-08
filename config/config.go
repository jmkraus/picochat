package config

import (
	"encoding/json"
	"fmt"
	"picochat/envs"
	"picochat/paths"
	"sync"

	"github.com/BurntSushi/toml"
)

type Config struct {
	URL         string  `json:"url"`
	Model       string  `json:"model"`
	Prompt      string  `json:"prompt"`
	Context     int     `json:"context"`
	Temperature float64 `json:"temperature"`
	Top_p       float64 `json:"top_p"`
	Reasoning   bool    `json:"reasoning"`
	Quiet       bool    `json:"quiet"`

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
	cfg := DefaultConfig()

	// 2. Config file
	if paths.FileExists(path) {
		if _, err := toml.DecodeFile(path, &cfg); err != nil {
			loadErr = fmt.Errorf("decode toml file %q failed: %w", path, err)
			return
		}
	} else {
		path = "No config.toml found - fallback to internal defaults"
	}

	if cfg.Context < MinCtx || cfg.Context > MaxCtx {
		loadErr = fmt.Errorf("context size must be between %d and %d", MinCtx, MaxCtx)
		return
	}

	cfg.ConfigPath = path
	instance = &cfg
}

// DefaultConfig defines the default values before loading the config file or evaluation env vars.
//
// Parameters:
//
//	none
//
// Returns:
//
//	Config - a filled Config struct
func DefaultConfig() Config {
	return Config{
		URL:         "http://localhost:11434/api",
		Model:       "gpt-oss:latest",
		Prompt:      "You are a Large Language Model. Answer as concisely as possible. Your answers should be informative, helpful and engaging.",
		Context:     20,
		Temperature: 0.7,
		Top_p:       0.9,
		Reasoning:   false,
		Quiet:       false,
	}
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

	if !envs.AllowedRuntimeField(key) {
		return fmt.Errorf("unsupported config key '%s'", key)
	}

	next := *cfg // work on copy to avoid compromised config
	if err := applyConfig(&next, key, value); err != nil {
		return fmt.Errorf("apply config failed: %w", err)
	}

	if key == "context" && (next.Context < MinCtx || next.Context > MaxCtx) {
		return fmt.Errorf("context size must be between %d and %d", MinCtx, MaxCtx)
	}

	*cfg = next
	return nil
}

// applyConfig alters a specific config element.
//
// Parameters:
//
//	cfg (*Config) - the instance of the configuration struct
//	key (string)  - the configuration key to modify
//	value (any)   - the new value for the key
//
// Returns:
//
//	error - error if any
func applyConfig(cfg *Config, key string, val any) error {
	patch := map[string]any{key: val}
	b, err := json.Marshal(patch)
	if err != nil {
		return err
	}
	return json.Unmarshal(b, cfg)
}
