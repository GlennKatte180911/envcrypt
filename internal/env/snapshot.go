package env

import (
	"fmt"
	"time"
)

// Snapshot captures the state of a set of env entries at a point in time.
type Snapshot struct {
	Entries   []Entry
	CreatedAt time.Time
	Label     string
}

// NewSnapshot creates a new Snapshot from the given entries and an optional label.
func NewSnapshot(entries []Entry, label string) Snapshot {
	copied := make([]Entry, len(entries))
	copy(copied, entries)
	return Snapshot{
		Entries:   copied,
		CreatedAt: time.Now().UTC(),
		Label:     label,
	}
}

// DiffFromSnapshot returns the Diff between a snapshot's entries and a newer set of entries.
func DiffFromSnapshot(snap Snapshot, current []Entry) []DiffEntry {
	return Diff(snap.Entries, current)
}

// SnapshotStore holds an ordered list of snapshots in memory.
type SnapshotStore struct {
	snapshots []Snapshot
}

// Add appends a snapshot to the store.
func (s *SnapshotStore) Add(snap Snapshot) {
	s.snapshots = append(s.snapshots, snap)
}

// Latest returns the most recent snapshot, or an error if the store is empty.
func (s *SnapshotStore) Latest() (Snapshot, error) {
	if len(s.snapshots) == 0 {
		return Snapshot{}, fmt.Errorf("snapshot store is empty")
	}
	return s.snapshots[len(s.snapshots)-1], nil
}

// All returns a copy of all snapshots in the store.
func (s *SnapshotStore) All() []Snapshot {
	out := make([]Snapshot, len(s.snapshots))
	copy(out, s.snapshots)
	return out
}

// Len returns the number of snapshots in the store.
func (s *SnapshotStore) Len() int {
	return len(s.snapshots)
}
