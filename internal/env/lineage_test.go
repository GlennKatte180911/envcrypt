package env

import (
	"testing"
	"time"
)

func lineageBase() []Entry {
	return []Entry{
		{Key: "APP_ENV", Value: "development"},
		{Key: "DB_HOST", Value: "localhost"},
	}
}

func TestNewLineageEmpty(t *testing.T) {
	l := NewLineage()
	if l.Len() != 0 {
		t.Fatalf("expected 0 records, got %d", l.Len())
	}
}

func TestLineageRecordNoChangeSkipped(t *testing.T) {
	l := NewLineage()
	base := lineageBase()
	l.Record(ChangeKindSet, "noop", base, base)
	if l.Len() != 0 {
		t.Fatalf("expected no record when diff is empty, got %d", l.Len())
	}
}

func TestLineageRecordAdded(t *testing.T) {
	l := NewLineage()
	base := lineageBase()
	updated := append(base, Entry{Key: "NEW_KEY", Value: "hello"})
	l.Record(ChangeKindSet, "add new key", base, updated)
	if l.Len() != 1 {
		t.Fatalf("expected 1 record, got %d", l.Len())
	}
	rec := l.Records()[0]
	if rec.Kind != ChangeKindSet {
		t.Errorf("expected kind %q, got %q", ChangeKindSet, rec.Kind)
	}
	if rec.Label != "add new key" {
		t.Errorf("unexpected label: %s", rec.Label)
	}
	if rec.Timestamp.IsZero() {
		t.Error("timestamp should not be zero")
	}
	if rec.Timestamp.Location() != time.UTC {
		t.Error("timestamp should be UTC")
	}
}

func TestLineageRecordsIsolated(t *testing.T) {
	l := NewLineage()
	base := lineageBase()
	updated := []Entry{{Key: "APP_ENV", Value: "production"}, {Key: "DB_HOST", Value: "localhost"}}
	l.Record(ChangeKindSet, "promote", base, updated)
	copy1 := l.Records()
	copy1[0].Label = "mutated"
	if l.Records()[0].Label == "mutated" {
		t.Error("Records() should return an isolated copy")
	}
}

func TestLineageChangedKeys(t *testing.T) {
	l := NewLineage()
	base := lineageBase()
	step1 := []Entry{{Key: "APP_ENV", Value: "staging"}, {Key: "DB_HOST", Value: "localhost"}}
	step2 := []Entry{{Key: "APP_ENV", Value: "staging"}, {Key: "DB_HOST", Value: "db.prod"}}
	l.Record(ChangeKindSet, "step1", base, step1)
	l.Record(ChangeKindMerge, "step2", step1, step2)
	keys := l.ChangedKeys()
	if len(keys) != 2 {
		t.Fatalf("expected 2 changed keys, got %d: %v", len(keys), keys)
	}
}

func TestLineageMultipleRecords(t *testing.T) {
	l := NewLineage()
	base := lineageBase()
	updated1 := []Entry{{Key: "APP_ENV", Value: "staging"}, {Key: "DB_HOST", Value: "localhost"}}
	updated2 := []Entry{{Key: "APP_ENV", Value: "production"}, {Key: "DB_HOST", Value: "localhost"}}
	l.Record(ChangeKindSet, "first", base, updated1)
	l.Record(ChangeKindSet, "second", updated1, updated2)
	if l.Len() != 2 {
		t.Fatalf("expected 2 records, got %d", l.Len())
	}
}
