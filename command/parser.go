package command

import (
	"errors"
	"fmt"
	"reflect"
	"strconv"
	"strings"

	"picochat/config"
)

var allowedKeys = map[string]string{
	"temperature": "Temperature",
	"top_p":       "TopP",
	"context":     "Context",
	"model":       "Model",
}

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

func applyToConfig(key string, value any) error {
	fieldName, ok := allowedKeys[key]
	if !ok {
		return fmt.Errorf("unsupported config key '%s'", key)
	}

	cfg := config.Get()
	v := reflect.ValueOf(cfg).Elem()  // dereference pointer to Config struct
	field := v.FieldByName(fieldName) // find struct field

	if !field.IsValid() {
		return fmt.Errorf("unsupported config key '%s'", fieldName)
	}
	if !field.CanSet() {
		return fmt.Errorf("cannot set config key '%s'", fieldName)
	}

	valValue := reflect.ValueOf(value)
	if valValue.Type().ConvertibleTo(field.Type()) {
		field.Set(valValue.Convert(field.Type()))
		return nil
	}

	return fmt.Errorf("cannot assign value of type %T to config key '%s'", value, fieldName)
}
