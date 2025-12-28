package format

import (
	"encoding/json"
	"fmt"
	"gopkg.in/yaml.v3"
	"io"
	"picochat/chat"
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
	arr := []string{"plain", "json", "json-pretty", "yaml"}

	normalizedInput := strings.ToLower(input)

	for _, str := range arr {
		if normalizedInput == strings.ToLower(str) {
			return normalizedInput, true
		}
	}

	return normalizedInput, false
}

// RenderResult formats and writes the inference output to screen.
// Parameters:
//
//	w (io.Writer)       - the io output
//	result (ChatResult) - a struct with all chat information
//	outputFmt (string)  - the output format of the current session
//
// Returns:
//
//	error - any error that might have occurred
func RenderResult(w io.Writer, result *chat.ChatResult, outputFmt string) error {

	switch outputFmt {

	case "json":
		enc := json.NewEncoder(w)
		enc.SetEscapeHTML(false)
		return enc.Encode(result)

	case "json-pretty":
		data, err := json.MarshalIndent(result, "", "  ")
		if err != nil {
			return err
		}
		_, err = fmt.Fprintln(w, string(data))
		return err

	case "yaml":
		enc := yaml.NewEncoder(w)
		defer enc.Close()
		return enc.Encode(result)

	case "plain":
		_, err := fmt.Printf(
			"\nElapsed (mm:ss): %s | Tokens/sec: %.1f",
			result.Elapsed,
			result.TokensPS,
		)
		return err

	default:
		return fmt.Errorf("unknown output format: %s", outputFmt)
	}
}
