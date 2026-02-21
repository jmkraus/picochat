package paths

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"picochat/args"
	"picochat/envs"
	"strings"
)

const (
	cfgDefaultName   = "config.toml"
	cfgDefaultSuffix = ".toml"
)

// GetConfigPath returns the path to the configuration file.
//
// Parameters:
//
//	none
//
// Returns:
//
//	string - the configuration file path
//	error - error if any
func GetConfigPath() (string, error) {
	cfgDir, err := getConfigDir()
	if err != nil {
		return "", err
	}

	if *args.ConfigPath != "" {
		name, found := strings.CutPrefix(*args.ConfigPath, "@")
		if !found {
			return *args.ConfigPath, nil
		}
		name = EnsureSuffix(name, cfgDefaultSuffix)
		return filepath.Join(cfgDir, name), nil
	}

	return filepath.Join(cfgDir, cfgDefaultName), nil
}

// EnsureSuffix ensures that the filename ends with the given suffix.
//
// Parameters:
//
//	filename string - the original filename
//	suffix string - the suffix to ensure
//
// Returns:
//
//	string - the filename with the suffix ensured
func EnsureSuffix(filename string, suffix string) string {
	if !strings.HasSuffix(filename, suffix) {
		return filename + suffix
	}
	return filename
}

// getConfigDir determines the configuration directory using various fallbacks.
//
// Parameters:
//
//	none
//
// Returns:
//
//	string - the configuration directory path
//	error - error if any
func getConfigDir() (string, error) {
	// Fallback 1: Executable path
	exc, err := fallbackToExecutableDir()
	if err != nil {
		return "", err
	}
	if exc != "" {
		found := FileExists(filepath.Join(exc, cfgDefaultName))
		if found {
			return exc, nil
		}
	}

	// Fallback 2: $CONFIG_PATH
	if env := envs.GetEnv(envs.CONFIG_PATH); env != "" {
		return env, nil
	}

	// Fallback 3: $XDG_CONFIG_HOME or ~/.config
	xdg, err := fallbackToXDGOrHome()
	if err != nil {
		return "", err
	}
	if xdg != "" {
		return xdg, nil
	}

	return "", fmt.Errorf("no valid config path found.")
}

var overrideHistoryPath string // for Unit Tests

// OverrideHistoryPath sets a custom history path for testing purposes.
//
// Parameters:
//
//	path string - the custom history path
//
// Returns:
//
//	restore func() - restores the previous history path override
func OverrideHistoryPath(path string) (restore func()) {
	prev := overrideHistoryPath
	overrideHistoryPath = path
	return func() {
		overrideHistoryPath = prev
	}
}

// GetHistoryPath returns the path to the history directory.
//
// Parameters:
//
//	none
//
// Returns:
//
//	string - the history directory path
//	error - error if any
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

// fallbackToXDGOrHome returns the XDG config directory or the home config directory.
//
// Parameters:
//
//	none
//
// Returns:
//
//	string - the config directory path
//	error - error if any
func fallbackToXDGOrHome() (string, error) {
	if env := envs.GetEnv(envs.XDG_CONFIG_HOME); env != "" {
		return filepath.Join(env, "picochat"), nil
	}

	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}

	return filepath.Join(homeDir, ".config", "picochat"), nil
}

// fallbackToExecutableDir returns the directory of the executable.
//
// Parameters:
//
//	none
//
// Returns:
//
//	string - the executable directory path
//	error - error if any
func fallbackToExecutableDir() (string, error) {
	exc, err := os.Executable()
	if err != nil {
		return "", err
	}

	return filepath.Dir(exc), nil
}

// ExpandHomeDir checks the given path and expands its user home
//
// Parameters:
//
//	path (string) - the path with tilde
//
// Returns:
//
//	string - the expanded path
//	error  - error if any
func ExpandHomeDir(path string) (string, error) {
	if strings.HasPrefix(path, "~/") {
		homeDir, err := os.UserHomeDir()
		if err != nil {
			return "", err
		}
		path = homeDir + path[1:]
	}
	return path, nil
}

// FileExists checks if the file of a given path actually exists
//
// Parameters:
//
//	path (string) - the full file path
//
// Returns:
//
//	bool - file exists: true or false
func FileExists(path string) bool {
	info, err := os.Stat(path)
	if err == nil {
		return !info.IsDir()
	}
	return !errors.Is(err, os.ErrNotExist)
}
