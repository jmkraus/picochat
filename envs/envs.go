package envs

import (
	"os"
	"picochat/utils"
	"picochat/vartypes"
	"strconv"
)

// EnvVar represents the valid environment variables in this package.
type EnvVar string

type EnvSpec struct {
	Env       EnvVar
	Type      vartypes.VarType
	Field     string
	Runtime   bool
	Sensitive bool
}

var ConfigEnvVars = []EnvSpec{
	{Env: "PICOCHAT_BACKEND", Type: vartypes.VarString, Field: "backend"},
	{Env: "PICOCHAT_URL", Type: vartypes.VarString, Field: "url"},
	{Env: "PICOCHAT_API_KEY", Type: vartypes.VarString, Field: "api_key", Sensitive: true},
	{Env: "PICOCHAT_MODEL", Type: vartypes.VarString, Field: "model"},
	{Env: "PICOCHAT_CONTEXT", Type: vartypes.VarInt, Field: "context", Runtime: true},
	{Env: "PICOCHAT_TEMPERATURE", Type: vartypes.VarFloat, Field: "temperature", Runtime: true},
	{Env: "PICOCHAT_TOP_P", Type: vartypes.VarFloat, Field: "top_p", Runtime: true},
	{Env: "PICOCHAT_REASONING", Type: vartypes.VarBool, Field: "reasoning", Runtime: true},
	{Env: "PICOCHAT_EFFORT", Type: vartypes.VarString, Field: "effort", Runtime: true},
	{Env: "PICOCHAT_QUIET", Type: vartypes.VarBool, Field: "quiet"},
}

var configByField map[string]EnvSpec

func init() {
	configByField = make(map[string]EnvSpec, len(ConfigEnvVars))
	for _, v := range ConfigEnvVars {
		configByField[v.Field] = v
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
//	bool   - environment variable is actually set (but can be empty)
func GetEnv(envvar EnvVar) (string, bool) {
	return os.LookupEnv(string(envvar))
}

// AllowedRuntimeField checks if the given field can be set at runtime.
//
// Parameters:
//
//	field (string) - the field name to be checked
//
// Returns:
//
//	bool - field can be set at runtime: true or false
func AllowedRuntimeField(field string) bool {
	cfg, ok := configByField[field]
	return ok && cfg.Runtime
}

// ConfigByField returns the config metadata for a field.
//
// Parameters:
//
//	field (string) - the config field name
//
// Returns:
//
//	EnvSpec - metadata for the field
//	bool    - true if field exists
func ConfigByField(field string) (EnvSpec, bool) {
	cfg, ok := configByField[field]
	return cfg, ok
}

// ConfigEnvVarsTable builds a table from env var state and values.
//
// Parameters:
//
//	none
//
// Returns:
//
//	string - the full markdown table
func ConfigEnvVarsTable() string {
	tableData := make([][]string, 0, len(ConfigEnvVars)+1)
	tableData = append(tableData, []string{"Env", "Set", "Value"})

	for _, spec := range ConfigEnvVars {
		val, lookup := GetEnv(spec.Env)
		set := strconv.FormatBool(lookup)

		if lookup && val == "" {
			val = "[empty]"
		}
		if lookup && spec.Sensitive && val != "" {
			val = "[hidden]"
		}
		tableData = append(tableData, []string{string(spec.Env), set, val})
	}

	return utils.MarkdownTable(tableData)
}
