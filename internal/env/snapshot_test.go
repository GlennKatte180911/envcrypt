package env

import (
	"testing"
	"time"
)

var snapBase = []Entry{
	{Key: "HOST", Value: "localhost"},
	{Key: "PORT", Value: "5432"},
}

func TestNewSnapshotCopiesEntries(t *testing.T) {
	snap := NewSnapshot(snapBase, "initial")
	if len(snap.Entries) != len(snapBase) {
		t.Fatalf("expected %d entries, got %d", len(snapBase), len(snap.Entries))
	}
	// Mutate original; snapshot should be unaffected.
	snapBase[0].Value = "changed"
	if snap.Entries[0].Value == "changed" {
		t.Error("snapshot entries should be independent of original slice")
	}
	snapBase[0].Value = "localhost" // restore
}

func TestNewSnapshotLabel(t *testing.T) {
	snap := NewSnapshot(snapBase, "v1")
	if snap.Label != "v1" {
		t.Errorf("expected label 'v1', got %q", snap.Label)
	}
}

func TestNewSnapshotTimestamp(t *testing.T) {
	before := time.Now().UTC()
	snap := NewSnapshot(snapBase, "")
	after := time.Now().UTC()
	if snap.CreatedAt.Before(before) || snap.CreatedAt.After(after) {
		t.Error("snapshot timestamp is outside expected range")
	}
}

func TestDiffFromSnapshot(t *testing.T) {
	snap := NewSnapshot(snapBase, "base")
	current := []Entry{
		{Key: "HOST", Value: "remotehost"},
		{Key: "PORT", Value: "5432"},
		{Key: "USER", Value: "admin"},
	}
	diffs := DiffFromSnapshot(snap, current)
	if !HasChanges(diffs) {
		t.Error("expected changes between snapshot and current entries")
	}
}

func TestSnapshotStoreLatestEmpty(t *testing.T) {
	var store SnapshotStore
	_, err := store.Latest()
	if err == nil {
		t.Error("expected error from empty store")
	}
}

func TestSnapshotStoreAddAndLatest(t *testing.T) {
	var store SnapshotStore
	s1 := NewSnapshot(snapBase, "first")
	s2 := NewSnapshot(snapBase, "second")
	store.Add(s1)
	store.Add(s2)
	latest, err := store.Latest()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if latest.Label != "second" {
		t.Errorf("expected latest label 'second', got %q", latest.Label)
	}
}

func TestSnapshotStoreAll(t *testing.T) {
	var store SnapshotStore
	store.Add(NewSnapshot(snapBase, "a"))
	store.Add(NewSnapshot(snapBase, "b"))
	all := store.All()
	if len(all) != 2 {
		t.Fatalf("expected 2 snapshots, got %d", len(all))
	}
}

func TestSnapshotStoreLen(t *testing.T) {
	var store SnapshotStore
	if store.Len() != 0 {
		t.Error("expected empty store length 0")
	}
	store.Add(NewSnapshot(snapBase, "x"))
	if store.Len() != 1 {
		t.Error("expected store length 1 after add")
	}
}
