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
	JsonField string
	Runtime   bool
	Sensitive bool
}

var ConfigEnvVars = []EnvSpec{
	{Env: "PICOCHAT_BACKEND", Type: vartypes.VarString, Field: "Backend", JsonField: "backend"},
	{Env: "PICOCHAT_URL", Type: vartypes.VarString, Field: "Url", JsonField: "url"},
	{Env: "PICOCHAT_API_KEY", Type: vartypes.VarString, Field: "ApiKey", JsonField: "api_key", Sensitive: true},
	{Env: "PICOCHAT_MODEL", Type: vartypes.VarString, Field: "Model", JsonField: "model"},
	{Env: "PICOCHAT_CONTEXT", Type: vartypes.VarInt, Field: "Context", JsonField: "context", Runtime: true},
	{Env: "PICOCHAT_TEMPERATURE", Type: vartypes.VarFloat, Field: "Temperature", JsonField: "temperature", Runtime: true},
	{Env: "PICOCHAT_TOP_P", Type: vartypes.VarFloat, Field: "Top_p", JsonField: "top_p", Runtime: true},
	{Env: "PICOCHAT_REASONING", Type: vartypes.VarBool, Field: "Reasoning", JsonField: "reasoning", Runtime: true},
	{Env: "PICOCHAT_EFFORT", Type: vartypes.VarString, Field: "Effort", JsonField: "effort", Runtime: true},
	{Env: "PICOCHAT_VALIDATE", Type: vartypes.VarBool, Field: "Validate", JsonField: "validate", Runtime: true},
	{Env: "PICOCHAT_QUIET", Type: vartypes.VarBool, Field: "Quiet", JsonField: "quiet"},
}

var envSpecByField map[string]EnvSpec

func init() {
	envSpecByField = make(map[string]EnvSpec, len(ConfigEnvVars))
	for _, v := range ConfigEnvVars {
		envSpecByField[v.JsonField] = v
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

// AllowedRuntimeField checks if the given canonical lowercase JSON field can
// be set at runtime.
//
// Parameters:
//
//	field (string) - the canonical lowercase JSON field name to be checked
//
// Returns:
//
//	bool - field can be set at runtime: true or false
func AllowedRuntimeField(field string) bool {
	cfg, ok := envSpecByField[field]
	return ok && cfg.Runtime
}

// EnvSpecByField returns the EnvSpec for a given canonical lowercase JSON field.
//
// Parameters:
//
//	field (string) - the canonical lowercase JSON field name
//
// Returns:
//
//	EnvSpec - metadata for the field
//	bool    - true if field exists
func EnvSpecByField(field string) (EnvSpec, bool) {
	cfg, ok := envSpecByField[field]
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
