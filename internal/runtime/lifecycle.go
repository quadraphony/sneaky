package runtime

// State captures the manager-visible lifecycle state for a tunnel session.
type State string

const (
	StateStopped  State = "stopped"
	StateStarting State = "starting"
	StateRunning  State = "running"
	StateStopping State = "stopping"
)

func (s State) String() string {
	return string(s)
}

func (s State) IsActive() bool {
	return s == StateStarting || s == StateRunning || s == StateStopping
}

func (s State) CanStart() bool {
	return s == StateStopped
}

func (s State) CanStop() bool {
	return s == StateRunning
}
