package stats

import (
	"fmt"
	"time"

	"sneaky-core/internal/runtime"
)

// Snapshot exposes stable manager-level observability values.
type Snapshot struct {
	State            runtime.State
	AdapterID        string
	StartedAt        time.Time
	LastTransitionAt time.Time
	Uptime           time.Duration
	SessionsStarted  uint64
	StartFailures    uint64
	StopFailures     uint64
}

func (s Snapshot) String() string {
	return fmt.Sprintf(
		"state=%s adapter=%s sessions_started=%d start_failures=%d stop_failures=%d uptime=%s",
		s.State,
		s.AdapterID,
		s.SessionsStarted,
		s.StartFailures,
		s.StopFailures,
		s.Uptime,
	)
}
