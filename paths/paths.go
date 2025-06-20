package paths

import (
	"fmt"
	"os"
	"path/filepath"
	"picochat/args"
	"strings"
)

func GetConfigPath() (string, error) {
	configDir, err := getConfigDir()
	if err != nil {
		return "", err
	}

	if *args.ConfigPath != "" {
		name, found := strings.CutPrefix(*args.ConfigPath, "@")
		if !found {
			return *args.ConfigPath, nil
		}

		var suffix = ".toml"
		if !HasSuffix(name, suffix) {
			name += suffix
		}
		return filepath.Join(configDir, name), nil
	}

	return filepath.Join(configDir, "config.toml"), nil
}

func HasSuffix(filename string, suffix string) bool {
	return strings.HasSuffix(filename, suffix)
}

func EnsureSuffix(filename string, suffix string) string {
	if !HasSuffix(filename, suffix) {
		return filename + suffix
	}
	return filename
}

func getConfigDir() (string, error) {
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

	return "", fmt.Errorf("no valid config path found.")
}

var overrideHistoryPath string // for Unit Tests
func OverrideHistoryPath(path string) {
	overrideHistoryPath = path
}

func GetHistoryPath() (string, error) {
	if overrideHistoryPath != "" {
		return overrideHistoryPath, nil
	}

	configDir, err := getConfigDir()
	if err != nil {
		return "", err
	}
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
		return filepath.Join(xdgConfigHome, "picochat"), nil
	}

	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}

	return filepath.Join(homeDir, ".config", "picochat"), nil
}

func fallbackToExecutableDir() (string, error) {
	ex, err := os.Executable()
	if err != nil {
		return "", err
	}

	return filepath.Dir(ex), nil
}
