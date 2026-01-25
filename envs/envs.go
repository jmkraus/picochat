package envs

import (
	"os"
)

// EnvVar represents the valid environment variables in this package.
type EnvVar string

const (
	CONFIG_PATH     EnvVar = "CONFIG_PATH"
	TMUX            EnvVar = "TMUX"
	XDG_CONFIG_HOME EnvVar = "XDG_CONFIG_HOME"
)

// GetEnv encapsulates reading environment variables
// and ensures with the use of constants proper naming.
//
// Parameters:
//
//	envvar (EnvVar) - the name of the environment variable
//
// Returns:
//
//	string - the value of the environment variable
func GetEnv(envvar EnvVar) string {
	return os.Getenv(string(envvar))
}
