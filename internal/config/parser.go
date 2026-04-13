package config

import (
	"bytes"
	"encoding/json"
	"errors"
	"os"
)

func LoadFile(path string) (Input, error) {
	raw, err := os.ReadFile(path)
	if err != nil {
		return Input{}, &ValidationError{
			Code:    ErrCodeUnreadableInput,
			Source:  path,
			Message: "failed to read config file",
			Err:     err,
		}
	}

	return Parse(raw, path)
}

func Parse(raw []byte, source string) (Input, error) {
	trimmed := bytes.TrimSpace(raw)
	if len(trimmed) == 0 {
		return Input{}, &ValidationError{
			Code:    ErrCodeEmptyInput,
			Source:  source,
			Message: "config input is empty",
		}
	}

	return Input{
		Source: source,
		Raw:    append([]byte(nil), trimmed...),
	}, nil
}

func decodeJSONObject(input Input) (map[string]json.RawMessage, error) {
	var topLevel any
	if err := json.Unmarshal(input.Raw, &topLevel); err != nil {
		var syntaxErr *json.SyntaxError
		if errors.As(err, &syntaxErr) {
			return nil, &ValidationError{
				Code:    ErrCodeInvalidJSON,
				Source:  input.Source,
				Message: "config is not valid JSON",
				Err:     err,
			}
		}

		return nil, &ValidationError{
			Code:    ErrCodeInvalidJSON,
			Source:  input.Source,
			Message: "config is not valid JSON",
			Err:     err,
		}
	}

	if _, ok := topLevel.(map[string]any); !ok {
		return nil, &ValidationError{
			Code:    ErrCodeInvalidTopLevelType,
			Source:  input.Source,
			Message: "config must be a top-level JSON object",
		}
	}

	var object map[string]json.RawMessage
	if err := json.Unmarshal(input.Raw, &object); err != nil {
		return nil, &ValidationError{
			Code:    ErrCodeInvalidJSON,
			Source:  input.Source,
			Message: "config is not valid JSON",
			Err:     err,
		}
	}

	if object == nil {
		return nil, &ValidationError{
			Code:    ErrCodeInvalidTopLevelType,
			Source:  input.Source,
			Message: "config must be a top-level JSON object",
		}
	}

	return object, nil
}
