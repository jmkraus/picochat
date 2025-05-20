package config

import (
	"fmt"
	"log"
	"picochat/paths"
	"picochat/types"

	"github.com/BurntSushi/toml"
)

// Load reads the configuration from the specified file.
func Load() (types.Config, error) {
	var cfg types.Config
	path, err := paths.GetConfigPath()
	if err != nil {
		return types.Config{}, err
	}
	log.Printf("Configuration file used: %s", path)

	if _, err := toml.DecodeFile(path, &cfg); err != nil {
		return types.Config{}, fmt.Errorf("failed to decode TOML file %q: %w", path, err)
	}

	// Basic validation
	if cfg.URL == "" || cfg.Model == "" || cfg.Prompt == "" {
		return types.Config{}, fmt.Errorf("required fields URL, Model, or Prompt are missing in config")
	}

	if cfg.Context != 0 && (cfg.Context < 5 || cfg.Context > 100) {
		return types.Config{}, fmt.Errorf("Context size must be between 5 and 100")
	}

	return cfg, nil
}
