package paths

import (
	"os"
	"path/filepath"
	"picochat/args"
)

var configPath string

func GetConfigPath() string {
	if *args.ConfigPath != "" {
		return *args.ConfigPath
	}

	// Fallback 1: $CONFIG_PATH
	if env := os.Getenv("CONFIG_PATH"); env != "" {
		return env
	}

	// Fallback 2: $XDG_CONFIG_HOME or ~/.config
	if xdg := os.Getenv("XDG_CONFIG_HOME"); xdg != "" {
		return filepath.Join(xdg, "picochat", "config.toml")
	}

	if home, err := os.UserHomeDir(); err == nil {
		return filepath.Join(home, ".config", "picochat", "config.toml")
	}

	// Fallback 3: Executable path
	if ex, err := os.Executable(); err == nil {
		return filepath.Join(filepath.Dir(ex), "config.toml")
	}

	return "config.toml"
}

func GetHistoryDir() string {
	configDir := filepath.Dir(GetConfigPath())
	historyDir := filepath.Join(configDir, "history")
	err := os.MkdirAll(historyDir, 0755)
	if err != nil {
		panic(err)
	}
	return historyDir
}

// fallbackToXDGOrHome returns the config path using XDG_CONFIG_HOME or home directory
func fallbackToXDGOrHome() string {
	xdgConfigHome := os.Getenv("XDG_CONFIG_HOME")
	if xdgConfigHome != "" {
		return filepath.Join(xdgConfigHome, "picochat", "config.toml")
	}

	homeDir, err := os.UserHomeDir()
	if err != nil {
		panic(err)
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
