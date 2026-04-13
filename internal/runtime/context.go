package runtime

import "time"

// Context tracks runtime metadata that is safe to surface at the core level.
type Context struct {
	AdapterID string
	StartedAt time.Time
}
