package paths

import (
	"fmt"
	"os"
	"path/filepath"
	"picochat/args"
)

func GetConfigPath() (string, error) {
	if *args.ConfigPath != "" {
		return *args.ConfigPath, nil
	}

	// Fallback 1: $CONFIG_PATH
	if env := os.Getenv("CONFIG_PATH"); env != "" {
		return env, nil
	}

	// Fallback 2: $XDG_CONFIG_HOME or ~/.config
	xdg, err := fallbackToXDGOrHome()
	if err != nil {
		return "", err
	}
	if xdg != "" {
		return xdg, nil
	}

	// Fallback 3: Executable path
	ex, err := fallbackToExecutableDir()
	if err != nil {
		return "", err
	}
	if ex != "" {
		return ex, nil
	}

	return "", fmt.Errorf("No valid config path found.")
}

var overrideHistoryDir string // for Unit Tests
func OverrideHistoryDir(path string) {
	overrideHistoryDir = path
}

func GetHistoryDir() (string, error) {
	if overrideHistoryDir != "" {
		return overrideHistoryDir, nil
	}

	configPath, err := GetConfigPath()
	if err != nil {
		return "", err
	}
	configDir := filepath.Dir(configPath)
	historyDir := filepath.Join(configDir, "history")
	err = os.MkdirAll(historyDir, 0755)
	if err != nil {
		return "", err
	}
	return historyDir, nil
}

func fallbackToXDGOrHome() (string, error) {
	xdgConfigHome := os.Getenv("XDG_CONFIG_HOME")
	if xdgConfigHome != "" {
		return filepath.Join(xdgConfigHome, "picochat", "config.toml"), nil
	}

	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}

	return filepath.Join(homeDir, ".config", "picochat", "config.toml"), nil
}

func fallbackToExecutableDir() (string, error) {
	ex, err := os.Executable()
	if err != nil {
		return "", err
	}

	return filepath.Join(filepath.Dir(ex), "config.toml"), nil
}
