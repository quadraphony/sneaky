package adapter

import (
	"context"

	"sneaky-core/internal/runtime"
)

// StartRequest holds the validated startup input for an adapter.
type StartRequest struct {
	ConfigPath string
	RawConfig  []byte
}

// Adapter defines the internal contract every protocol implementation must satisfy.
type Adapter interface {
	Identity() string
	Capabilities() Capabilities
	ValidateConfig(StartRequest) error
	Start(ctx context.Context, req StartRequest) (runtime.Handle, error)
}
