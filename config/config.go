package config

import (
	"fmt"
	"picochat/paths"
	"sync"

	"github.com/BurntSushi/toml"
)

var (
	instance *Config
	once     sync.Once
	loadErr  error
)

// Load reads and caches the configuration.
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
		Temperature: 0.7,
		TopP:        0.9,
		Context:     20,
		Reasoning:   false,
		Quiet:       false,
	}

	if _, err := toml.DecodeFile(path, &cfg); err != nil {
		loadErr = fmt.Errorf("failed to decode TOML file %q: %w", path, err)
		return
	}

	if cfg.URL == "" || cfg.Model == "" || cfg.Prompt == "" {
		loadErr = fmt.Errorf("required fields URL, Model, or Prompt are missing in config")
		return
	}

	if cfg.Context != 0 && (cfg.Context < 3 || cfg.Context > 100) {
		loadErr = fmt.Errorf("context size must be between 3 and 100")
		return
	}

	cfg.ConfigPath = path
	instance = &cfg
}

// Get loads the configuration once and returns the instance of the loaded configuration.
//
// Parameters:
//
//	none
//
// Returns:
//
//	*Config - pointer to the loaded configuration
//	 error  - an error if the loading of the config file failed
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

	case "image":
		strVal, ok := value.(string)
		if !ok {
			return fmt.Errorf("value for key '%s' must be a string", key)
		}
		cfg.ImagePath = strVal
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
		return fmt.Errorf("unsupported config key '%s'", key)
	}
}
