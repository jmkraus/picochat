package config

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"picochat/types"

	"github.com/BurntSushi/toml"
)

var configPath string

// getConfigPath determines the configuration file path using the following priority:
// 1. Command line argument (-config)
// 2. Environment variable CONFIG_PATH
// 3. XDG_CONFIG_HOME (or ~/.config if not set)
// 4. Executable directory
func getConfigPath() string {
	// 1. Command line argument
	flag.StringVar(&configPath, "config", "", "Path to configuration file")

	// 2. Environment variable
	if configPath == "" {
		configPath = os.Getenv("CONFIG_PATH")
	}

	// 3. Fallback to XDG_CONFIG_HOME or home directory
	if configPath == "" {
		configPath = fallbackToXDGOrHome()
	}

	// 4. Fallback to executable directory
	if configPath == "" {
		configPath = fallbackToExecutableDir()
	}

	log.Printf("Configuration file used: %s", configPath)
	return configPath
}

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
	path := getConfigPath()

	if _, err := toml.DecodeFile(path, &cfg); err != nil {
		return types.Config{}, fmt.Errorf("failed to decode TOML file %q: %w", path, err)
	}

	// Basic validation
	if cfg.URL == "" || cfg.Model == "" || cfg.Prompt == "" {
		return types.Config{}, fmt.Errorf("required fields URL, Model, or Prompt are missing in config")
	}

	return cfg, nil
}
