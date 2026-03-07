package convert

import (
	"fmt"
	"strconv"
	"strings"
)

// ParseKeyVal parses a string of the form "key=value" and returns
// the key, converted value, and error.
//
// Parameters:
//
//	args - the input string to parse
//
// Returns:
//
//	key   - the parsed key in lower case
//	value - the parsed and converted value
//	error - error if any
func ParseKeyVal(args string) (string, any, error) {
	parts := strings.SplitN(args, "=", 2)
	if len(parts) != 2 {
		return "", nil, fmt.Errorf("invalid format, expected key=value")
	}

	key := strings.ToLower(strings.TrimSpace(parts[0]))
	value := strings.TrimSpace(parts[1])

	if key == "" {
		return "", nil, fmt.Errorf("invalid format, missing key")
	}
	if value == "" {
		return "", nil, fmt.Errorf("invalid format, missing value")
	}

	convertedValue, err := validateAndConvert(key, value)
	if err != nil {
		return "", nil, err
	}

	return key, convertedValue, nil
}

// ValidateAndConvert validates the key and converts the value to the appropriate type.
//
// Parameters:
//
//	key   - the configuration key
//	value - the string representation of the value
//
// Returns:
//
//	any   - the converted value
//	error - error if any
func validateAndConvert(key, value string) (any, error) {
	switch key {
	case "temperature", "top_p":
		v, err := strconv.ParseFloat(value, 64)
		if err != nil {
			return nil, fmt.Errorf("invalid float value for key '%s'", key)
		}
		return v, nil
	case "context":
		v, err := strconv.Atoi(value)
		if err != nil {
			return nil, fmt.Errorf("invalid integer value for key '%s'", key)
		}
		return v, nil
	case "url", "model":
		return value, nil
	case "reasoning", "quiet":
		v, err := StringToBool(value)
		if err != nil {
			return nil, err
		}
		return v, nil
	default:
		// Don't forget to update config/config.go --> Set()
		return nil, fmt.Errorf("unsupported config key '%s'", key)
	}
}

// TypeConvert converts the value to the given var type.
//
// Parameters:
//
//	varType - the var type
//	value   - the string representation of the value
//
// Returns:
//
//	any   - the converted value
//	error - error if any
func TypeConvert(varType, value string) (any, error) {
	switch varType {
	case "float":
		v, err := strconv.ParseFloat(value, 64)
		if err != nil {
			return nil, fmt.Errorf("invalid float value '%s'", varType)
		}
		return v, nil
	case "int":
		v, err := strconv.Atoi(value)
		if err != nil {
			return nil, fmt.Errorf("invalid integer value '%s'", varType)
		}
		return v, nil
	case "string":
		return value, nil
	case "bool":
		v, err := StringToBool(value)
		if err != nil {
			return nil, err
		}
		return v, nil
	default:
		return nil, fmt.Errorf("unsupported var type '%s'", varType)
	}
}

// StringToBool converts a string representation of a bool to a boolean value.
//
// Parameters:
//
//	s (string) - the string representation
//
// Returns:
//
//	bool  - the boolean value of the string
//	error - error if any
func StringToBool(s string) (bool, error) {
	switch strings.ToLower(s) {
	case "true", "1", "t", "y", "yes":
		return true, nil
	case "false", "0", "f", "n", "no":
		return false, nil
	default:
		return false, fmt.Errorf("invalid boolean value '%s'", s)
	}
}
