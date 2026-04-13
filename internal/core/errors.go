package core

import "fmt"

type ErrorCode string

const (
	ErrCodeInvalidArgument ErrorCode = "invalid_argument"
	ErrCodeInvalidState    ErrorCode = "invalid_state"
	ErrCodeAdapterExists   ErrorCode = "adapter_exists"
	ErrCodeAdapterNotFound ErrorCode = "adapter_not_found"
	ErrCodeStartFailed     ErrorCode = "start_failed"
	ErrCodeStopFailed      ErrorCode = "stop_failed"
	ErrCodeRuntimeExited   ErrorCode = "runtime_exited"
)

// Error is the structured error surface used across the core contracts.
type Error struct {
	Code    ErrorCode
	Op      string
	Message string
	Err     error
}

func (e *Error) Error() string {
	if e == nil {
		return "<nil>"
	}

	switch {
	case e.Op != "" && e.Message != "" && e.Err != nil:
		return fmt.Sprintf("%s: %s: %v", e.Op, e.Message, e.Err)
	case e.Op != "" && e.Message != "":
		return fmt.Sprintf("%s: %s", e.Op, e.Message)
	case e.Op != "" && e.Err != nil:
		return fmt.Sprintf("%s: %v", e.Op, e.Err)
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

func (e *Error) Unwrap() error {
	if e == nil {
		return nil
	}

	return e.Err
}

func newError(code ErrorCode, op, message string, err error) *Error {
	return &Error{
		Code:    code,
		Op:      op,
		Message: message,
		Err:     err,
	}
}
