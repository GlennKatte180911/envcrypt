package env_test

import (
	"testing"

	"github.com/your-org/envcrypt/internal/env"
)

// TestHistoryDocExample exercises the documented usage pattern for History.
func TestHistoryDocExample(t *testing.T) {
	before := []env.Entry{
		{Key: "APP_ENV", Value: "development"},
		{Key: "PORT", Value: "8080"},
	}

	after := []env.Entry{
		{Key: "APP_ENV", Value: "production"},
		{Key: "PORT", Value: "8080"},
		{Key: "DEBUG", Value: "false"},
	}

	h := env.NewHistory()
	h.Record("promote to production", before, after)

	if h.Len() != 1 {
		t.Fatalf("expected 1 history entry, got %d", h.Len())
	}

	he, ok := h.Last()
	if !ok {
		t.Fatal("expected Last() to return an entry")
	}
	if he.Label != "promote to production" {
		t.Errorf("unexpected label: %q", he.Label)
	}
	if he.Timestamp.IsZero() {
		t.Error("expected non-zero timestamp")
	}
	if !env.HasChanges(he.Diff) {
		t.Error("expected diff to have changes")
	}

	// A second record with identical entries should be skipped.
	h.Record("no-op", after, after)
	if h.Len() != 1 {
		t.Errorf("expected history length to remain 1, got %d", h.Len())
	}
}
