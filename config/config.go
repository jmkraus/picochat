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
	cfgName  string
	loadErr  error
)

// Load reads and caches the configuration once.
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

func Get() *Config {
	return instance
}
