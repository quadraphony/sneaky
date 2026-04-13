package logx

import (
	"fmt"
	"sort"
	"strings"
	"time"
)

type Level string

const (
	LevelInfo  Level = "info"
	LevelError Level = "error"
)

// Entry is the stable structured log record surfaced by the core.
type Entry struct {
	Time    time.Time
	Level   Level
	Event   string
	Message string
	Fields  map[string]string
}

func (e Entry) String() string {
	parts := []string{
		e.Time.UTC().Format(time.RFC3339),
		string(e.Level),
		e.Event,
		e.Message,
	}

	if len(e.Fields) == 0 {
		return strings.Join(parts, " ")
	}

	keys := make([]string, 0, len(e.Fields))
	for key := range e.Fields {
		keys = append(keys, key)
	}
	sort.Strings(keys)

	fieldParts := make([]string, 0, len(keys))
	for _, key := range keys {
		fieldParts = append(fieldParts, fmt.Sprintf("%s=%s", key, e.Fields[key]))
	}

	return strings.Join(parts, " ") + " " + strings.Join(fieldParts, " ")
}
