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
	{Env: "PICOCHAT_URL", Type: "string", Field: "url"},
	{Env: "PICOCHAT_MODEL", Type: "string", Field: "model"},
	{Env: "PICOCHAT_CONTEXT", Type: "int", Field: "context"},
	{Env: "PICOCHAT_TEMPERATURE", Type: "float", Field: "temperature"},
	{Env: "PICOCHAT_TOP_P", Type: "float", Field: "top_p"},
	{Env: "PICOCHAT_REASONING", Type: "bool", Field: "reasoning"},
	{Env: "PICOCHAT_QUIET", Type: "bool", Field: "quiet"},
}

var allowedFields map[string]bool

func init() {
	allowedFields = make(map[string]bool)
	for _, envVar := range ConfigEnvVars {
		allowedFields[envVar.Field] = true
	}
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

// AllowedField checks if the given field name is valid.
//
// Parameters:
//
//	field (string) - the field name to be checked
//
// Returns:
//
//	bool - field name is valid: true or false
func AllowedField(field string) bool {
	return allowedFields[field]
}
