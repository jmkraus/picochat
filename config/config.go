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
	path := paths.GetConfigPath()
	log.Printf("Configuration file used: %s", path)

	if _, err := toml.DecodeFile(path, &cfg); err != nil {
		return types.Config{}, fmt.Errorf("failed to decode TOML file %q: %w", path, err)
	}

	// Basic validation
	if cfg.URL == "" || cfg.Model == "" || cfg.Prompt == "" {
		return types.Config{}, fmt.Errorf("required fields URL, Model, or Prompt are missing in config")
	}

	return cfg, nil
}
