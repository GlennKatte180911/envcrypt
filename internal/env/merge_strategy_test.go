package env

import (
	"testing"
)

var stratBase = []Entry{
	{Key: "HOST", Value: "localhost"},
	{Key: "PORT", Value: "5432"},
}

var stratOverride = []Entry{
	{Key: "PORT", Value: "9999"},
	{Key: "DEBUG", Value: "true"},
}

func TestMergeWithStrategyPreferBase(t *testing.T) {
	result, err := MergeWithStrategy(stratBase, stratOverride, PreferBase)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	m := ToMap(result)
	if m["PORT"] != "5432" {
		t.Errorf("expected PORT=5432, got %s", m["PORT"])
	}
	if m["DEBUG"] != "true" {
		t.Errorf("expected DEBUG=true, got %s", m["DEBUG"])
	}
}

func TestMergeWithStrategyPreferOverride(t *testing.T) {
	result, err := MergeWithStrategy(stratBase, stratOverride, PreferOverride)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	m := ToMap(result)
	if m["PORT"] != "9999" {
		t.Errorf("expected PORT=9999, got %s", m["PORT"])
	}
	if m["HOST"] != "localhost" {
		t.Errorf("expected HOST=localhost, got %s", m["HOST"])
	}
}

func TestMergeWithStrategyErrorOnConflict(t *testing.T) {
	_, err := MergeWithStrategy(stratBase, stratOverride, ErrorOnConflict)
	if err == nil {
		t.Fatal("expected conflict error, got nil")
	}
	ce, ok := err.(*ConflictError)
	if !ok {
		t.Fatalf("expected *ConflictError, got %T", err)
	}
	if ce.Key != "PORT" {
		t.Errorf("expected conflict on PORT, got %s", ce.Key)
	}
}

func TestMergeWithStrategyNoConflict(t *testing.T) {
	extra := []Entry{{Key: "LOG_LEVEL", Value: "info"}}
	result, err := MergeWithStrategy(stratBase, extra, ErrorOnConflict)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(result) != 3 {
		t.Errorf("expected 3 entries, got %d", len(result))
	}
}

func TestConflictErrorMessage(t *testing.T) {
	e := &ConflictError{Key: "SECRET"}
	want := `env: merge conflict on key "SECRET"`
	if e.Error() != want {
		t.Errorf("got %q, want %q", e.Error(), want)
	}
}
