package config

import (
	"maps"
	"picochat/utils"
	"sort"
)

type Template struct {
	Description string `toml:"Description"`
	Prompt      string `toml:"Prompt"`
}

var templates map[string]Template

// setTemplates stores loaded templates as an internal copy.
//
// Parameters:
//
//	in (map[string]Template) - loaded templates from config
//
// Returns:
//
//	none
func setTemplates(in map[string]Template) {
	if len(in) == 0 {
		templates = nil
		return
	}

	templates = make(map[string]Template, len(in))
	maps.Copy(templates, in)
}

// GetTemplate returns a template text by key.
//
// Parameters:
//
//	key (string) - template key
//
// Returns:
//
//	string - template prompt (empty string if not found)
func GetTemplate(key string) string {
	tpl, ok := templates[key]
	if !ok {
		return ""
	}
	return tpl.Prompt
}

// GetTemplateDescription returns a template description by key.
//
// Parameters:
//
//	key (string) - template key
//
// Returns:
//
//	string - template description (empty string if not found)
func GetTemplateDescription(key string) string {
	tpl, ok := templates[key]
	if !ok {
		return ""
	}
	return tpl.Description
}

// listTemplateKeys returns all loaded template keys.
//
// Parameters:
//
//	none
//
// Returns:
//
//	[]string - sorted list of template keys
func listTemplateKeys() []string {
	keys := make([]string, 0, len(templates))
	for k := range templates {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}

// ListTemplates returns a markdown table of all loaded templates.
//
// Parameters:
//
//	none
//
// Returns:
//
//	string - markdown table with template key and description
func ListTemplates() string {
	tableData := make([][]string, 0, len(templates)+1)
	tableData = append(tableData, []string{"Key", "Description"})

	for _, key := range listTemplateKeys() {
		tableData = append(tableData, []string{key, templates[key].Description})
	}

	return utils.MarkdownTable(tableData)
}
