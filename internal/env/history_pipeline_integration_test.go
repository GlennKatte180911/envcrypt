package env_test

import (
	"testing"

	"github.com/your-org/envcrypt/internal/env"
)

// TestHistoryThroughPipeline verifies that History correctly captures diffs
// produced by a Pipeline transformation.
func TestHistoryThroughPipeline(t *testing.T) {
	initial := []env.Entry{
		{Key: "app_env", Value: "dev"},
		{Key: "port", Value: "3000"},
	}

	pipe := env.NewPipeline(initial)
	pipe.Pipe(env.PipeFunc(func(entries []env.Entry) ([]env.Entry, error) {
		return env.UppercaseKeys(entries), nil
	}))

	result, err := pipe.Run()
	if err != nil {
		t.Fatalf("pipeline error: %v", err)
	}

	h := env.NewHistory()
	h.Record("uppercase keys", initial, result)

	if h.Len() != 1 {
		t.Fatalf("expected 1 history entry, got %d", h.Len())
	}

	he, _ := h.Last()
	if !env.HasChanges(he.Diff) {
		t.Error("expected changes after uppercasing keys")
	}
}

// TestHistoryMultiplePipelineStages records each intermediate stage.
func TestHistoryMultiplePipelineStages(t *testing.T) {
	initial := []env.Entry{
		{Key: "app_env", Value: ""},
		{Key: "port", Value: "3000"},
	}

	h := env.NewHistory()

	// Stage 1: uppercase keys.
	stage1 := env.UppercaseKeys(initial)
	h.Record("uppercase", initial, stage1)

	// Stage 2: drop empty values.
	stage2 := env.DropEmpty(stage1)
	h.Record("drop empty", stage1, stage2)

	if h.Len() != 2 {
		t.Fatalf("expected 2 history entries, got %d", h.Len())
	}

	entries := h.Entries()
	if entries[0].Label != "uppercase" {
		t.Errorf("expected first label 'uppercase', got %q", entries[0].Label)
	}
	if entries[1].Label != "drop empty" {
		t.Errorf("expected second label 'drop empty', got %q", entries[1].Label)
	}

	// The drop-empty stage should have removed APP_ENV (empty value).
	var removedFound bool
	for _, d := range entries[1].Diff {
		if d.Key == "APP_ENV" && d.Kind == env.DiffRemoved {
			removedFound = true
		}
	}
	if !removedFound {
		t.Error("expected APP_ENV to appear as removed in stage2 diff")
	}
}
