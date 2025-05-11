package paths

import (
	"os"
	"path/filepath"
)

func GetConfigPath() string {
	// 1. Env variable
	if p := os.Getenv("CONFIG_PATH"); p != "" {
		return p
	}

	// 2. XDG or fallback to ~/.config
	if xdg := os.Getenv("XDG_CONFIG_HOME"); xdg != "" {
		return filepath.Join(xdg, "picochat", "config.toml")
	}

	if home, err := os.UserHomeDir(); err == nil {
		return filepath.Join(home, ".config", "picochat", "config.toml")
	}

	// 3. Executable dir fallback
	if ex, err := os.Executable(); err == nil {
		return filepath.Join(filepath.Dir(ex), "config.toml")
	}

	// Final fallback
	return "."
}
