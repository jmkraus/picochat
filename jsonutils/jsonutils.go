package jsonutils

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/google/jsonschema-go/jsonschema"
)

var (
	resolvedSchema *jsonschema.Resolved
)

// PrettyPrint reformats a JSON string representation into pretty style.
//
// Parameters:
//
//	jsonStr (string) - the JSON as a single or multiline string
//
// Returns:
//
//	string - the formatted JSON
//	error  - error if any
func PrettyPrint(jsonStr string) (string, error) {
	jsonStr = strings.TrimSpace(jsonStr)
	if jsonStr == "" {
		return "", fmt.Errorf("json string is empty")
	}

	var buf bytes.Buffer
	if err := json.Indent(&buf, []byte(jsonStr), "", "  "); err != nil {
		return "", fmt.Errorf("invalid json - format failed: %w", err)
	}
	return buf.String(), nil
}

// ValidateJSON validates a JSON string representation against a given schema.
//
// Parameters:
//
//	schemaMap (map[string]any) - the schema as a map representation
//	jsonStr (string)           - the JSON as a single or multiline string
//
// Returns:
//
//	error - an error if the validation fails, or nil
func ValidateJSON(schemaMap map[string]any, jsonStr string) error {
	if len(schemaMap) == 0 {
		return fmt.Errorf("schema definition is empty")
	}
	jsonStr = strings.TrimSpace(jsonStr)
	if jsonStr == "" {
		return fmt.Errorf("json string is empty")
	}

	resolved := resolvedSchema

	if resolved == nil {
		schemaBytes, err := json.Marshal(schemaMap)
		if err != nil {
			return fmt.Errorf("marshal schema failed: %w", err)
		}
		var schema jsonschema.Schema
		if err := json.Unmarshal(schemaBytes, &schema); err != nil {
			return fmt.Errorf("unmarshal schema failed: %w", err)
		}

		resolved, err = schema.Resolve(nil)
		if err != nil {
			return fmt.Errorf("resolve schema failed: %w", err)
		}

		resolvedSchema = resolved
	}

	var instance any
	dec := json.NewDecoder(strings.NewReader(jsonStr))
	if err := dec.Decode(&instance); err != nil {
		return fmt.Errorf("invalid json: %w", err)
	}
	if dec.More() {
		return fmt.Errorf("invalid json: trailing data")
	}

	return resolved.Validate(instance)
}
