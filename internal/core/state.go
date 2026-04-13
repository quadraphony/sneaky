package core

import (
	"time"

	"sneaky-core/internal/runtime"
)

// Snapshot exposes the current manager state without leaking adapter internals.
type Snapshot struct {
	State     runtime.State
	AdapterID string
	StartedAt time.Time
	LastError *Error
	Active    bool
}
