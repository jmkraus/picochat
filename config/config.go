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
	Backend     string  `json:"backend"`
	URL         string  `json:"url"`
	APIKey      string  `json:"api_key"`
	Model       string  `json:"model"`
	Prompt      string  `json:"prompt"`
	Context     int     `json:"context"`
	Temperature float64 `json:"temperature"`
	Top_p       float64 `json:"top_p"`
	Reasoning   bool    `json:"reasoning"`
	Effort      string  `json:"effort"`
	Quiet       bool    `json:"quiet"`

	ConfigPath string            `toml:"-"`
	ImagePath  string            `toml:"-"` ////IMAGES
	OutputFmt  string            `toml:"-"`
	SchemaFmt  map[string]any    `toml:"-"`
	Templates  map[string]string `json:"-" toml:"Templates"`
}

var (
	instance       *Config
	once           sync.Once
	loadWarn       []string
	loadError      error
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
//	configPathArg (string) - the path to the config file
//
// Returns:
//
//	none
func load(configPathArg string) {
	path, err := paths.GetConfigPath(configPathArg)
	if err != nil {
		loadError = err
		return
	}

	// 1. Default values
	cfg := defaultConfig()

	// 2. Config file
	if paths.FileExists(path) {
		if _, err := toml.DecodeFile(path, &cfg); err != nil {
			loadError = fmt.Errorf("decode toml file %q failed: %w", path, err)
			return
		}
	} else {
		if len(initConfigPath) > 0 {
			// Show a warning only if the configuration file is explicitly set via the -config argument.
			loadWarn = append(loadWarn, "config file not found - fallback to default or env vars")
		}
		path = "none"
	}

	// 3. Environment variables
	err = applyEnvValues(&cfg)
	if err != nil {
		loadError = fmt.Errorf("apply env var values failed: %w", err)
		return
	}

	// 4. Check value contraints
	loadWarn = append(loadWarn, NormalizeConfig(&cfg)...)

	// 5. Load templates
	setTemplates(cfg.Templates)

	cfg.ConfigPath = path
	instance = &cfg
}

// defaultConfig defines the default values before loading the config file or evaluation env vars.
//
// Parameters:
//
//	none
//
// Returns:
//
//	Config - a filled Config struct
func defaultConfig() Config {
	return Config{
		URL:         "http://localhost:11434",
		Backend:     "ollama",
		APIKey:      "ollama",
		Model:       "gpt-oss:latest",
		Prompt:      "You are a Large Language Model. Answer as concisely as possible. Your answers should be informative, helpful and engaging.",
		Context:     20,
		Temperature: 0.7,
		Top_p:       0.9,
		Reasoning:   false,
		Effort:      "medium",
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
//	*Config  - pointer to the loaded configuration
//	[]string - one or more warnings if any
//	error    - error if any
func Get() (*Config, []string, error) {
	once.Do(func() {
		load(initConfigPath) // load takes string arg
	})
	return instance, loadWarn, loadError
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
//	[]strings - one or more warnings if any
//	error     - error if any
func Set(key string, value any) ([]string, error) {
	cfg, _, err := Get()
	if err != nil {
		return nil, fmt.Errorf("get config data failed: %w", err)
	}

	if !envs.AllowedRuntimeField(key) {
		return nil, fmt.Errorf("unsupported config key %q", key)
	}

	next := *cfg // work on copy to avoid compromised config
	if err := applyConfigValue(&next, key, value); err != nil {
		return nil, fmt.Errorf("apply config value failed: %w", err)
	}

	warnings := NormalizeConfig(&next)

	*cfg = next
	return warnings, nil
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
