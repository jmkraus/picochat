package paths

import (
	"flag"
	"os"
	"path/filepath"
)

var configPath string

func init() {
	if flag.Lookup("config") == nil {
		flag.StringVar(&configPath, "config", "", "Path to configuration file")
	}
}

// getConfigPath determines the configuration file path using the following priority:
// 1. Command line argument (-config)
// 2. Environment variable CONFIG_PATH
// 3. XDG_CONFIG_HOME (or ~/.config if not set)
// 4. Executable directory
func GetConfigPath() string {

	// 1. Command line argument

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

	return configPath
}

func GetHistoryDir() string {
	configDir := filepath.Dir(GetConfigPath())
	historyDir := filepath.Join(configDir, "history")
	os.MkdirAll(historyDir, os.ModePerm)
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
