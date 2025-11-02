package config

import (
	"fmt"
	"picochat/paths"
	"reflect"
	"sync"

	"github.com/BurntSushi/toml"
)

var (
	instance *Config
	once     sync.Once
	cfgName  string
	loadErr  error
)

var allowedKeys = map[string]string{
	"temperature": "Temperature",
	"top_p":       "TopP",
	"context":     "Context",
	"model":       "Model",
	"quiet":       "Quiet",
}

// Load reads and caches the configuration once.
//
// Parameters:
//   none
//
// Returns:
//   string: the path to the configuration file
//   error: an error if loading fails
func Load() (string, error) {
	once.Do(func() {
		path, err := paths.GetConfigPath()
		if err != nil {
			loadErr = err
			return
		}

		cfgName = path

		var cfg = Config{
			Temperature: 0.7,
			TopP:        0.9,
			Context:     20,
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

		if cfg.Context != 0 && (cfg.Context < 5 || cfg.Context > 100) {
			loadErr = fmt.Errorf("context size must be between 5 and 100")
			return
		}

		instance = &cfg
	})

	return cfgName, loadErr
}

// Get returns the instance of the loaded configuration.
//
// Parameters:
//   none
//
// Returns:
//   *Config: pointer to the loaded configuration
func Get() *Config {
	return instance
}

// ApplyToConfig allows changing a specific parameter after loading.
//
// Parameters:
//   key string: the configuration key to modify
//   value any: the new value for the key
//
// Returns:
//   error: an error if the key is unsupported or the value cannot be set
func ApplyToConfig(key string, value any) error {
	fieldName, ok := allowedKeys[key]
	if !ok {
		return fmt.Errorf("unsupported config key '%s'", key)
	}

	cfg := Get()
	v := reflect.ValueOf(cfg).Elem()  // dereference pointer to Config struct
	field := v.FieldByName(fieldName) // find struct field

	if !field.IsValid() {
		return fmt.Errorf("unsupported config key '%s'", fieldName)
	}
	if !field.CanSet() {
		return fmt.Errorf("cannot set config key '%s'", fieldName)
	}

	valValue := reflect.ValueOf(value)
	if valValue.Type().ConvertibleTo(field.Type()) {
		field.Set(valValue.Convert(field.Type()))
		return nil
	}

	return fmt.Errorf("cannot assign value of type %T to config key '%s'", value, fieldName)
}
