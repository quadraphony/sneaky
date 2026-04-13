package logx

import (
	"sync"
	"time"
)

// Logger stores a bounded in-memory sequence of structured log entries.
type Logger struct {
	mu         sync.Mutex
	maxEntries int
	entries    []Entry
}

func New(maxEntries int) *Logger {
	if maxEntries <= 0 {
		maxEntries = 128
	}

	return &Logger{
		maxEntries: maxEntries,
		entries:    make([]Entry, 0, maxEntries),
	}
}

func (l *Logger) Info(event, message string, fields map[string]string) {
	l.append(LevelInfo, event, message, fields)
}

func (l *Logger) Error(event, message string, fields map[string]string) {
	l.append(LevelError, event, message, fields)
}

func (l *Logger) Entries() []Entry {
	l.mu.Lock()
	defer l.mu.Unlock()

	out := make([]Entry, 0, len(l.entries))
	for _, entry := range l.entries {
		out = append(out, cloneEntry(entry))
	}
	return out
}

func (l *Logger) append(level Level, event, message string, fields map[string]string) {
	l.mu.Lock()
	defer l.mu.Unlock()

	entry := Entry{
		Time:    time.Now().UTC(),
		Level:   level,
		Event:   event,
		Message: message,
		Fields:  cloneFields(fields),
	}

	if len(l.entries) == l.maxEntries {
		copy(l.entries, l.entries[1:])
		l.entries[len(l.entries)-1] = entry
		return
	}

	l.entries = append(l.entries, entry)
}

func cloneEntry(entry Entry) Entry {
	entry.Fields = cloneFields(entry.Fields)
	return entry
}

func cloneFields(fields map[string]string) map[string]string {
	if len(fields) == 0 {
		return nil
	}

	out := make(map[string]string, len(fields))
	for key, value := range fields {
		out[key] = value
	}
	return out
}
