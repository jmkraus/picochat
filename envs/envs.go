package envs

import (
	"os"
	"picochat/utils"
	"picochat/vartypes"
)

// EnvVar represents the valid environment variables in this package.
type EnvVar string

type EnvSpec struct {
	Env     EnvVar
	Type    vartypes.VarType
	Field   string
	Runtime bool
}

var ConfigEnvVars = []EnvSpec{
	{Env: "PICOCHAT_URL", Type: vartypes.VarString, Field: "url", Runtime: false},
	{Env: "PICOCHAT_MODEL", Type: vartypes.VarString, Field: "model", Runtime: true},
	{Env: "PICOCHAT_CONTEXT", Type: vartypes.VarInt, Field: "context", Runtime: true},
	{Env: "PICOCHAT_TEMPERATURE", Type: vartypes.VarFloat, Field: "temperature", Runtime: true},
	{Env: "PICOCHAT_TOP_P", Type: vartypes.VarFloat, Field: "top_p", Runtime: true},
	{Env: "PICOCHAT_REASONING", Type: vartypes.VarBool, Field: "reasoning", Runtime: true},
	{Env: "PICOCHAT_QUIET", Type: vartypes.VarBool, Field: "quiet", Runtime: false},
}

var allowedRuntimeFields map[string]bool
var configByField map[string]EnvSpec

func init() {
	allowedRuntimeFields = make(map[string]bool, len(ConfigEnvVars))
	configByField = make(map[string]EnvSpec, len(ConfigEnvVars))
	for _, v := range ConfigEnvVars {
		configByField[v.Field] = v
		if v.Runtime {
			allowedRuntimeFields[v.Field] = true
		}
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
//		string - the value of the environment variable
//	 bool   - environment variable is actually set (but can be empty)
func GetEnv(envvar EnvVar) (string, bool) {
	return os.LookupEnv(string(envvar))
}

// AllowedRuntimeField checks if the given field name is valid.
//
// Parameters:
//
//	field (string) - the field name to be checked
//
// Returns:
//
//	bool - field name is valid: true or false
func AllowedRuntimeField(field string) bool {
	return allowedRuntimeFields[field]
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

// ConfigEnvVarsMarkdownTable returns a markdown table with env var state and values.
//
// Parameters:
//
//	none
//
// Returns:
//
//	string - the full markdown table
func ConfigEnvVarsMarkdownTable() string {
	tableData := make([][]string, 0, len(ConfigEnvVars)+1)
	tableData = append(tableData, []string{"Env", "Set", "Value"})

	for _, spec := range ConfigEnvVars {
		val, lookup := GetEnv(spec.Env)
		set := "true"
		if !lookup {
			set = "false"
		}
		if lookup && val == "" {
			val = "<empty>"
		}
		tableData = append(tableData, []string{string(spec.Env), set, val})
	}

	return utils.MarkdownTable(tableData)
}
