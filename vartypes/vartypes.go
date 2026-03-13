package vartypes

import (
	"fmt"
	"strconv"
	"strings"
)

// VarType represents supported primitive config value types.
type VarType uint8

const (
	VarUnknown VarType = iota
	VarFloat
	VarInt
	VarString
	VarBool
)

// String returns the string representation of the VarType.
//
// Parameters:
//
//	t (VarType) - the VarType value
//
// Returns:
//
//	string - the string representation of the VarType
func (t VarType) String() string {
	switch t {
	case VarFloat:
		return "float"
	case VarInt:
		return "int"
	case VarString:
		return "string"
	case VarBool:
		return "bool"
	default:
		return "unknown"
	}
}

// Convert converts a string to the target VarType.
//
// Parameters:
//
//	varType (VarType) - the target type for conversion
//	value   (string)  - the raw string value to convert
//
// Returns:
//
//	any   - the converted value
//	error - error if conversion fails or type is unsupported
func Convert(varType VarType, value string) (any, error) {
	switch varType {
	case VarFloat:
		v, err := strconv.ParseFloat(value, 64)
		if err != nil {
			return nil, fmt.Errorf("invalid float value '%s'", value)
		}
		return v, nil
	case VarInt:
		v, err := strconv.Atoi(value)
		if err != nil {
			return nil, fmt.Errorf("invalid integer value '%s'", value)
		}
		return v, nil
	case VarString:
		return value, nil
	case VarBool:
		v, err := stringToBool(value)
		if err != nil {
			return nil, err
		}
		return v, nil
	default:
		return nil, fmt.Errorf("unsupported var type '%s'", varType)
	}
}

// stringToBool converts a string representation of a bool to a boolean value.
//
// Parameters:
//
//	s (string) - the string representation
//
// Returns:
//
//	bool  - the boolean value of the string
//	error - error if any
func stringToBool(s string) (bool, error) {
	switch strings.ToLower(s) {
	case "true", "1", "t", "y", "yes":
		return true, nil
	case "false", "0", "f", "n", "no":
		return false, nil
	default:
		return false, fmt.Errorf("invalid boolean value '%s'", s)
	}
}
