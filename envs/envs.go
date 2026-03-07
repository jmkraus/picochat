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

// Preparation for future use
type Config struct {
	Env     string
	Type    string
	Field   string
	Default any
}

var ConfigEnvVars = []Config{
	{Env: "PICOCHAT_URL", Type: "string", Field: "url", Default: "http://localhost:11434/api"},
	{Env: "PICOCHAT_MODEL", Type: "string", Field: "model", Default: "gpt-oss:latest"},
	{Env: "PICOCHAT_CONTEXT", Type: "int", Field: "context", Default: 20},
	{Env: "PICOCHAT_TEMPERATURE", Type: "float", Field: "temperature", Default: 0.70},
	{Env: "PICOCHAT_TOP_P", Type: "float", Field: "top_p", Default: 0.90},
	{Env: "PICOCHAT_REASONING", Type: "bool", Field: "reasoning", Default: false},
	{Env: "PICOCHAT_QUIET", Type: "bool", Field: "quiet", Default: false},
}

// GetEnv encapsulates reading environment variables
// and ensures with the use of constants their proper naming.
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
