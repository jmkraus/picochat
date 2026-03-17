package config

import (
	"encoding/json"
	"fmt"
	"picochat/envs"
	"picochat/paths"
	"picochat/vartypes"
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
	instance       *Config
	once           sync.Once
	loadErr        error
	initConfigPath string
)

// Init sets startup config path override.
//
// Parameters:
//
//	configPath (string) - the path to the config file
//
// Returns:
//
//	none
func Init(configPath string) {
	initConfigPath = configPath
}

// load reads and caches the configuration.
//
// Parameters:
//
//	none
//
// Returns:
//
//	none
func load(configPathArg string) {
	path, err := paths.GetConfigPath(configPathArg)
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
		path = "none"
	}

	// 3. Environment variables
	err = applyEnvValues(&cfg)
	if err != nil {
		loadErr = fmt.Errorf("apply env var values failed: %w", err)
		return
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
	once.Do(func() {
		load(initConfigPath) // load takes string arg
	})
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
	if err := applyConfigValue(&next, key, value); err != nil {
		return fmt.Errorf("apply config value failed: %w", err)
	}

	if key == "context" && (next.Context < MinCtx || next.Context > MaxCtx) {
		return fmt.Errorf("context size must be between %d and %d", MinCtx, MaxCtx)
	}

	*cfg = next
	return nil
}

// applyConfig updates a specific config element.
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
func applyConfigValue(cfg *Config, key string, val any) error {
	patch := map[string]any{key: val}
	b, err := json.Marshal(patch)
	if err != nil {
		return err
	}
	return json.Unmarshal(b, cfg)
}

// applyEnvValues updates config fields according to
// set environment variables.
//
// Parameters:
//
//	cfg (*Config) - the instance of the configuration struct
//
// Returns:
//
//	error - error if any
func applyEnvValues(cfg *Config) error {
	for _, spec := range envs.ConfigEnvVars {
		envVal, lookup := envs.GetEnv(spec.Env)
		if !lookup || envVal == "" {
			continue // Skip if not set or empty
		}

		v, err := vartypes.Convert(spec.Type, envVal)
		if err != nil {
			return fmt.Errorf("convert type for env %s failed: %w", spec.Env, err)
		}
		applyConfigValue(cfg, spec.Field, v)
	}
	return nil
}
