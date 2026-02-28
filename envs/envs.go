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

	PICOCHAT_URL         EnvVar = "PICOCHAT_URL"
	PICOCHAT_MODEL       EnvVar = "PICOCHAT_MODEL"
	PICOCHAT_CONTEXT     EnvVar = "PICOCHAT_CONTEXT"
	PICOCHAT_TEMPERATURE EnvVar = "PICOCHAT_TEMPERATURE"
	PICOCHAT_TOP_P       EnvVar = "PICOCHAT_TOP_P"
	PICOCHAT_REASONING   EnvVar = "PICOCHAT_REASONING"
	PICOCHAT_QUIET       EnvVar = "PICOCHAT_QUIET"
)

var PicochatEnvVars = []EnvVar{
	PICOCHAT_URL,
	PICOCHAT_MODEL,
	PICOCHAT_CONTEXT,
	PICOCHAT_TEMPERATURE,
	PICOCHAT_TOP_P,
	PICOCHAT_REASONING,
	PICOCHAT_QUIET,
}

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

// GetSetEnvVars returns a list of set variables.
//
// Parameters:
//
//	vars ([]EnvVar) - slice of envvars
//
// Returns:
//
//	[]EnvVar - slice of set envvars
func GetSetEnvVars(vars []EnvVar) []EnvVar {
	out := make([]EnvVar, 0, len(vars))
	for _, k := range vars {
		if GetEnv(k) != "" {
			out = append(out, k)
		}
	}
	return out
}
