package adapter

// Capabilities declares optional surfaces an adapter can expose.
type Capabilities struct {
	ProvidesLogs  bool
	ProvidesStats bool
}
