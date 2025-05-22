package config

import (
	"fmt"
	"log"
	"picochat/paths"
	"sync"

	"github.com/BurntSushi/toml"
)

type Config struct {
	URL     string
	Model   string
	Context int
	Prompt  string
}

var (
	instance *Config
	once     sync.Once
	loadErr  error
)

// Load reads and caches the configuration once.
func Load() error {
	once.Do(func() {
		path, err := paths.GetConfigPath()
		if err != nil {
			loadErr = err
			return
		}

		log.Printf("Configuration file used: %s", path)

		var cfg Config
		if _, err := toml.DecodeFile(path, &cfg); err != nil {
			loadErr = fmt.Errorf("failed to decode TOML file %q: %w", path, err)
			return
		}

		if cfg.URL == "" || cfg.Model == "" || cfg.Prompt == "" {
			loadErr = fmt.Errorf("required fields URL, Model, or Prompt are missing in config")
			return
		}

		if cfg.Context != 0 && (cfg.Context < 5 || cfg.Context > 100) {
			loadErr = fmt.Errorf("Context size must be between 5 and 100")
			return
		}

		instance = &cfg
	})

	return loadErr
}

func Get() *Config {
	return instance
}
