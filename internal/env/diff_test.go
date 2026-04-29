package env

import (
	"testing"
)

func baseEntries() []Entry {
	return []Entry{
		{Key: "APP_NAME", Value: "myapp"},
		{Key: "APP_ENV", Value: "production"},
		{Key: "DB_HOST", Value: "localhost"},
	}
}

func TestDiffNoChanges(t *testing.T) {
	old := baseEntries()
	new := baseEntries()
	diffs := Diff(old, new)
	if len(diffs) != 0 {
		t.Fatalf("expected no diffs, got %d", len(diffs))
	}
}

func TestDiffAdded(t *testing.T) {
	old := baseEntries()
	new := append(baseEntries(), Entry{Key: "NEW_KEY", Value: "newval"})
	diffs := Diff(old, new)
	if len(diffs) != 1 {
		t.Fatalf("expected 1 diff, got %d", len(diffs))
	}
	if diffs[0].Kind != DiffAdded || diffs[0].Key != "NEW_KEY" {
		t.Errorf("unexpected diff: %+v", diffs[0])
	}
}

func TestDiffRemoved(t *testing.T) {
	old := baseEntries()
	new := []Entry{
		{Key: "APP_NAME", Value: "myapp"},
		{Key: "APP_ENV", Value: "production"},
	}
	diffs := Diff(old, new)
	if len(diffs) != 1 {
		t.Fatalf("expected 1 diff, got %d", len(diffs))
	}
	if diffs[0].Kind != DiffRemoved || diffs[0].Key != "DB_HOST" {
		t.Errorf("unexpected diff: %+v", diffs[0])
	}
}

func TestDiffChanged(t *testing.T) {
	old := baseEntries()
	new := []Entry{
		{Key: "APP_NAME", Value: "myapp"},
		{Key: "APP_ENV", Value: "staging"},
		{Key: "DB_HOST", Value: "localhost"},
	}
	diffs := Diff(old, new)
	if len(diffs) != 1 {
		t.Fatalf("expected 1 diff, got %d", len(diffs))
	}
	d := diffs[0]
	if d.Kind != DiffChanged || d.Key != "APP_ENV" || d.OldValue != "production" || d.NewValue != "staging" {
		t.Errorf("unexpected diff: %+v", d)
	}
}

func TestHasChanges(t *testing.T) {
	old := baseEntries()
	if HasChanges(old, old) {
		t.Error("expected no changes for identical slices")
	}
	new := append(baseEntries(), Entry{Key: "EXTRA", Value: "val"})
	if !HasChanges(old, new) {
		t.Error("expected changes to be detected")
	}
}
