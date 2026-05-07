package env

import (
	"fmt"
	"time"
)

// HistoryEntry records a single change event applied to a set of env entries.
type HistoryEntry struct {
	Timestamp time.Time
	Label     string
	Diff      []DiffEntry
}

// History maintains an ordered log of change events for a set of env entries.
type History struct {
	entries []HistoryEntry
}

// NewHistory returns an empty History.
func NewHistory() *History {
	return &History{}
}

// Record appends a new HistoryEntry if the diff contains changes.
// label is a human-readable description of the change (e.g. "load staging").
func (h *History) Record(label string, before, after []Entry) {
	d := Diff(before, after)
	if !HasChanges(d) {
		return
	}
	h.entries = append(h.entries, HistoryEntry{
		Timestamp: time.Now().UTC(),
		Label:     label,
		Diff:      d,
	})
}

// Entries returns all recorded history entries in chronological order.
func (h *History) Entries() []HistoryEntry {
	out := make([]HistoryEntry, len(h.entries))
	copy(out, h.entries)
	return out
}

// Len returns the number of recorded history entries.
func (h *History) Len() int {
	return len(h.entries)
}

// Last returns the most recent HistoryEntry and true, or a zero value and
// false when the history is empty.
func (h *History) Last() (HistoryEntry, bool) {
	if len(h.entries) == 0 {
		return HistoryEntry{}, false
	}
	return h.entries[len(h.entries)-1], true
}

// Clear removes all recorded history entries.
func (h *History) Clear() {
	h.entries = nil
}

// String returns a compact human-readable summary of the history.
func (h *History) String() string {
	if len(h.entries) == 0 {
		return "history: no entries"
	}
	return fmt.Sprintf("history: %d entries, latest: %q at %s",
		len(h.entries),
		h.entries[len(h.entries)-1].Label,
		h.entries[len(h.entries)-1].Timestamp.Format(time.RFC3339),
	)
}
