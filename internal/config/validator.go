package config

import (
	"encoding/json"
	"strconv"
)

func Validate(input Input, metadata Metadata) error {
	if metadata.AdapterID == "" {
		return &ValidationError{
			Code:    ErrCodeUnknownFormat,
			Source:  input.Source,
			Message: "config metadata does not declare an adapter",
		}
	}

	switch metadata.AdapterID {
	case AdapterSingbox:
		return validateSingbox(input)
	default:
		return &ValidationError{
			Code:    ErrCodeUnknownFormat,
			Source:  input.Source,
			Message: "config adapter is outside the current validation scope",
		}
	}
}

func DetectAndValidate(input Input) (Metadata, error) {
	metadata, err := Detect(input)
	if err != nil {
		return Metadata{}, err
	}
	if err := Validate(input, metadata); err != nil {
		return Metadata{}, err
	}
	return metadata, nil
}

func validateSingbox(input Input) error {
	object, err := decodeJSONObject(input)
	if err != nil {
		return err
	}

	rawOutbounds, ok := object["outbounds"]
	if !ok {
		return &ValidationError{
			Code:    ErrCodeMissingRequiredField,
			Source:  input.Source,
			Message: "sing-box config requires an outbounds field",
		}
	}

	var outbounds []map[string]json.RawMessage
	if err := json.Unmarshal(rawOutbounds, &outbounds); err != nil {
		return &ValidationError{
			Code:    ErrCodeMissingRequiredField,
			Source:  input.Source,
			Message: "sing-box outbounds must be an array of objects",
			Err:     err,
		}
	}
	if len(outbounds) == 0 {
		return &ValidationError{
			Code:    ErrCodeMissingRequiredField,
			Source:  input.Source,
			Message: "sing-box config requires at least one outbound",
		}
	}

	for idx, outbound := range outbounds {
		if _, ok := outbound["type"]; !ok {
			return &ValidationError{
				Code:    ErrCodeMissingRequiredField,
				Source:  input.Source,
				Message: "sing-box outbound is missing required type field at index " + strconv.Itoa(idx),
			}
		}
	}

	return nil
}
