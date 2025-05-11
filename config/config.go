package config

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"picochat/paths"
	"picochat/types"

	"github.com/BurntSushi/toml"
)

var configPath string

// fallbackToXDGOrHome returns the config path using XDG_CONFIG_HOME or home directory
func fallbackToXDGOrHome() string {
	xdgConfigHome := os.Getenv("XDG_CONFIG_HOME")
	if xdgConfigHome != "" {
		return filepath.Join(xdgConfigHome, "picochat", "config.toml")
	}

	homeDir, err := os.UserHomeDir()
	if err != nil {
		log.Printf("Failed to get home directory: %v", err)
		return ""
	}

	return filepath.Join(homeDir, ".config", "picochat", "config.toml")
}

// fallbackToExecutableDir returns the config path using the executable directory
func fallbackToExecutableDir() string {
	ex, err := os.Executable()
	if err != nil {
		panic(err)
	}

	return filepath.Join(filepath.Dir(ex), "config.toml")
}

// Load reads the configuration from the specified file.
func Load() (types.Config, error) {
	var cfg types.Config
	path := paths.GetConfigPath()

	if _, err := toml.DecodeFile(path, &cfg); err != nil {
		return types.Config{}, fmt.Errorf("failed to decode TOML file %q: %w", path, err)
	}

	// Basic validation
	if cfg.URL == "" || cfg.Model == "" || cfg.Prompt == "" {
		return types.Config{}, fmt.Errorf("required fields URL, Model, or Prompt are missing in config")
	}

	return cfg, nil
}
