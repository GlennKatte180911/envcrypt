package env_test

import (
	"testing"

	"github.com/nicholasgasior/envcrypt/internal/env"
)

// TestSnapshotDocExample verifies the documented usage pattern for snapshots:
// capture a baseline, apply changes, then diff against the snapshot.
func TestSnapshotDocExample(t *testing.T) {
	base := []env.Entry{
		{Key: "APP_ENV", Value: "development"},
		{Key: "LOG_LEVEL", Value: "debug"},
	}

	snap := env.NewSnapshot(base, "before-deploy")

	if snap.Label != "before-deploy" {
		t.Errorf("unexpected label: %q", snap.Label)
	}

	updated := []env.Entry{
		{Key: "APP_ENV", Value: "production"},
		{Key: "LOG_LEVEL", Value: "debug"},
		{Key: "SENTRY_DSN", Value: "https://example.sentry.io"},
	}

	diffs := env.DiffFromSnapshot(snap, updated)
	if !env.HasChanges(diffs) {
		t.Fatal("expected changes to be detected")
	}

	changedCount := 0
	for _, d := range diffs {
		if d.Type != "unchanged" {
			changedCount++
		}
	}
	if changedCount == 0 {
		t.Error("expected at least one non-unchanged diff entry")
	}
}
