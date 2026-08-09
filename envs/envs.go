package envs

import (
	"os"
	"picochat/utils"
	"picochat/vartypes"
	"strconv"
	"strings"
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
	{Env: "PICOCHAT_BACKEND", Type: vartypes.VarString, Field: "Backend"},
	{Env: "PICOCHAT_URL", Type: vartypes.VarString, Field: "Url"},
	{Env: "PICOCHAT_API_KEY", Type: vartypes.VarString, Field: "ApiKey", Sensitive: true},
	{Env: "PICOCHAT_MODEL", Type: vartypes.VarString, Field: "Model"},
	{Env: "PICOCHAT_CONTEXT", Type: vartypes.VarInt, Field: "Context", Runtime: true},
	{Env: "PICOCHAT_TEMPERATURE", Type: vartypes.VarFloat, Field: "Temperature", Runtime: true},
	{Env: "PICOCHAT_TOP_P", Type: vartypes.VarFloat, Field: "Top_p", Runtime: true},
	{Env: "PICOCHAT_REASONING", Type: vartypes.VarBool, Field: "Reasoning", Runtime: true},
	{Env: "PICOCHAT_EFFORT", Type: vartypes.VarString, Field: "Effort", Runtime: true},
	{Env: "PICOCHAT_VALIDATE", Type: vartypes.VarBool, Field: "Validate", Runtime: true},
	{Env: "PICOCHAT_QUIET", Type: vartypes.VarBool, Field: "Quiet"},
}

var envSpecByField map[string]EnvSpec

func init() {
	envSpecByField = make(map[string]EnvSpec, len(ConfigEnvVars))
	for _, v := range ConfigEnvVars {
		field := strings.ToLower(v.Field)
		envSpecByField[field] = v
	}
}

// GetEnv encapsulates reading environment variables
// and ensures with the use of constants their proper naming.
//
// Parameters:
//
//	envVar (EnvVar) - the name of the environment variable
//
// Returns:
//
//	string - the value of the environment variable
//	bool   - environment variable is actually set (but can be empty)
func GetEnv(envVar EnvVar) (string, bool) {
	return os.LookupEnv(string(envVar))
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
	lowerField := strings.ToLower(field)
	cfg, ok := envSpecByField[lowerField]
	return ok && cfg.Runtime
}

// EnvSpecByField returns the EnvSpec for a given field.
//
// Parameters:
//
//	field (string) - the config field name
//
// Returns:
//
//	EnvSpec - metadata for the field
//	bool    - true if field exists
func EnvSpecByField(field string) (EnvSpec, bool) {
	lowerField := strings.ToLower(field)
	cfg, ok := envSpecByField[lowerField]
	return cfg, ok
}

// ListEnvVars builds a table from env var state and values.
//
// Parameters:
//
//	none
//
// Returns:
//
//	string - the full markdown table
func ListEnvVars() string {
	tableData := make([][]string, 0, len(ConfigEnvVars)+1)
	tableData = append(tableData, []string{"Env", "Type", "Set", "Value"})

	for _, spec := range ConfigEnvVars {
		val, lookup := GetEnv(spec.Env)
		set := strconv.FormatBool(lookup)
		typ := spec.Type.String()

		if lookup && spec.Sensitive && val != "" {
			val = "[hidden]"
		}
		if lookup && val == "" {
			val = "[empty]"
		}
		tableData = append(tableData, []string{string(spec.Env), typ, set, val})
	}

	return utils.MarkdownTable(tableData)
}
