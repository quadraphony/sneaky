package logx

import (
	"strings"
	"testing"
)

func TestLoggerKeepsBoundedEntries(t *testing.T) {
	logger := New(2)

	logger.Info("one", "first", map[string]string{"step": "1"})
	logger.Error("two", "second", nil)
	logger.Info("three", "third", nil)

	entries := logger.Entries()
	if len(entries) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(entries))
	}
	if entries[0].Event != "two" || entries[1].Event != "three" {
		t.Fatalf("unexpected events: %#v", entries)
	}
}

func TestEntryStringIncludesSortedFields(t *testing.T) {
	entry := Entry{
		Level:   LevelInfo,
		Event:   "manager.start.succeeded",
		Message: "manager start succeeded",
		Fields: map[string]string{
			"b": "2",
			"a": "1",
		},
	}

	got := entry.String()
	if got == "" {
		t.Fatal("expected non-empty entry string")
	}
	if !strings.Contains(got, "a=1") || !strings.Contains(got, "b=2") {
		t.Fatalf("expected fields in string, got %q", got)
	}
}
