package envs

import (
	"os"
	"picochat/vartypes"
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
		set := "false"
		if lookup {
			set = "true"
		}
		if lookup && val == "" {
			val = "[empty]"
		}
		if lookup && spec.Sensitive && val != "" {
			val = "[hidden]"
		}
		tableData = append(tableData, []string{string(spec.Env), set, val})
	}

	return markdownTable(tableData)
}

// markdownTable creates a markdown table from a 2-dim array
// and adjusts the column widths according to content.
//
// Parameters:
//
//	tableData ([][]string) - the rows and columns of the table
//
// Returns:
//
//	string - markdown string with all line breaks
func markdownTable(tableData [][]string) string {
	if len(tableData) == 0 || len(tableData[0]) == 0 {
		return ""
	}

	numColumns := len(tableData[0])
	separator := make([]string, numColumns)
	for i := range separator {
		separator[i] = "---"
	}

	rows := make([][]string, 0, len(tableData)+1)
	rows = append(rows, tableData[0], separator)
	rows = append(rows, tableData[1:]...)

	maxWidths := make([]int, numColumns)
	for _, row := range rows {
		for colIdx := range numColumns {
			col := ""
			if colIdx < len(row) {
				col = row[colIdx]
			}
			if len(col) > maxWidths[colIdx] {
				maxWidths[colIdx] = len(col)
			}
		}
	}

	pad := func(s string, width int, fill byte) string {
		if len(s) >= width {
			return s
		}
		return s + strings.Repeat(string(fill), width-len(s))
	}

	var builder strings.Builder
	for rowIdx, row := range rows {
		fill := byte(' ')
		if rowIdx == 1 {
			fill = '-'
		}

		builder.WriteByte('|')
		for colIdx := range numColumns {
			col := ""
			if colIdx < len(row) {
				col = row[colIdx]
			}
			builder.WriteByte(' ')
			builder.WriteString(pad(col, maxWidths[colIdx], fill))
			builder.WriteString(" |")
		}
		if rowIdx < len(rows)-1 {
			builder.WriteByte('\n')
		}
	}

	return builder.String()
}
