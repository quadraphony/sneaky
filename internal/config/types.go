package config

import "fmt"

type Format string

const (
	FormatJSON Format = "json"
)

const AdapterSingbox = "singbox"

type Input struct {
	Source string
	Raw    []byte
}

type Metadata struct {
	Format       Format
	AdapterID    string
	TopLevelKeys []string
}

type ErrorCode string

const (
	ErrCodeEmptyInput           ErrorCode = "empty_input"
	ErrCodeUnreadableInput      ErrorCode = "unreadable_input"
	ErrCodeInvalidJSON          ErrorCode = "invalid_json"
	ErrCodeInvalidTopLevelType  ErrorCode = "invalid_top_level_type"
	ErrCodeUnknownFormat        ErrorCode = "unknown_format"
	ErrCodeAmbiguousFormat      ErrorCode = "ambiguous_format"
	ErrCodeMissingRequiredField ErrorCode = "missing_required_field"
)

type ValidationError struct {
	Code    ErrorCode
	Source  string
	Message string
	Err     error
}

func (e *ValidationError) Error() string {
	if e == nil {
		return "<nil>"
	}

	switch {
	case e.Source != "" && e.Message != "" && e.Err != nil:
		return fmt.Sprintf("%s: %s: %v", e.Source, e.Message, e.Err)
	case e.Source != "" && e.Message != "":
		return fmt.Sprintf("%s: %s", e.Source, e.Message)
	case e.Message != "" && e.Err != nil:
		return fmt.Sprintf("%s: %v", e.Message, e.Err)
	case e.Message != "":
		return e.Message
	case e.Err != nil:
		return e.Err.Error()
	default:
		return string(e.Code)
	}
}

func (e *ValidationError) Unwrap() error {
	if e == nil {
		return nil
	}

	return e.Err
}
