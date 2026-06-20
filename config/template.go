package config

import (
	"maps"
	"sort"
)

var templates map[string]string

// setTemplates stores loaded templates as an internal copy.
//
// Parameters:
//
//	in (map[string]string) - loaded templates from config
//
// Returns:
//
//	none
func setTemplates(in map[string]string) {
	if len(in) == 0 {
		templates = nil
		return
	}

	templates = make(map[string]string, len(in))
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
//	string - template value (empty string if not found)
func GetTemplate(key string) string {
	return templates[key]
}

// ListTemplates returns all loaded template keys.
//
// Parameters:
//
//	none
//
// Returns:
//
//	[]string - sorted list of template keys
func ListTemplates() []string {
	keys := make([]string, 0, len(templates))
	for k := range templates {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}
