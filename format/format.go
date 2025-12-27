package format

import (
	"strings"
)

// AllowedKeys checks if the argument string is valid.
// Parameters:
//
//	input (string) - the format argument
//
// Returns:
//
//	string - the input string normalized to lower case
//	bool   - ok (true/false) if the input key was valid
func AllowedKeys(input string) (string, bool) {
	arr := []string{"plain", "json", "yaml"}

	normalizedInput := strings.ToLower(input)

	for _, str := range arr {
		if normalizedInput == strings.ToLower(str) {
			return normalizedInput, true
		}
	}

	return normalizedInput, false
}
