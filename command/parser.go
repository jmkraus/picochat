package command

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

func ParseArgs(args string) (string, any, error) {
	parts := strings.SplitN(args, "=", 2)
	if len(parts) != 2 {
		return "", nil, errors.New("invalid format, expected key=value")
	}

	key := strings.ToLower(strings.TrimSpace(parts[0]))
	value := strings.TrimSpace(parts[1])

	if key == "" {
		return "", nil, errors.New("invalid format, missing key")
	}
	if value == "" {
		return "", nil, errors.New("invalid format, missing value")
	}

	convertedValue, err := validateAndConvert(key, value)
	if err != nil {
		return "", nil, err
	}

	return key, convertedValue, nil
}

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
	case "model":
		return value, nil
	default:
		return nil, fmt.Errorf("unsupported config key '%s'", key)
	}
}
