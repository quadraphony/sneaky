package runtime

import "context"

// Handle is the minimal runtime contract the core manager controls.
type Handle interface {
	Stop(ctx context.Context) error
	State() State
}

// Session binds a running adapter handle to its runtime metadata.
type Session struct {
	Context Context
	Handle  Handle
}
