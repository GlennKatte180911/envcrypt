package env

import (
	"errors"
	"testing"
)

func pipelineEntries() []Entry {
	return []Entry{
		{Key: "APP_HOST", Value: "localhost"},
		{Key: "app_port", Value: "8080"},
		{Key: "DB_PASS", Value: "secret"},
		{Key: "EMPTY_KEY", Value: ""},
	}
}

func TestPipelineEmpty(t *testing.T) {
	entries := pipelineEntries()
	out, err := NewPipeline().Run(entries)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(out) != len(entries) {
		t.Fatalf("expected %d entries, got %d", len(entries), len(out))
	}
}

func TestPipelineDoesNotMutateInput(t *testing.T) {
	entries := pipelineEntries()
	orig := make([]Entry, len(entries))
	copy(orig, entries)

	_, _ = NewPipeline().Pipe(PipeFunc(UppercaseKeys)).Run(entries)

	for i, e := range entries {
		if e.Key != orig[i].Key {
			t.Errorf("input mutated at index %d: got %q want %q", i, e.Key, orig[i].Key)
		}
	}
}

func TestPipelineChainedStages(t *testing.T) {
	entries := pipelineEntries()

	out, err := NewPipeline().
		Pipe(PipeFunc(UppercaseKeys)).
		Pipe(PipeFunc(DropEmpty)).
		Run(entries)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	// EMPTY_KEY should be dropped; all keys uppercase
	for _, e := range out {
		if e.Key == "EMPTY_KEY" {
			t.Error("expected EMPTY_KEY to be dropped")
		}
		if e.Key != toUpper(e.Key) {
			t.Errorf("key %q should be uppercase", e.Key)
		}
	}
	if len(out) != 3 {
		t.Fatalf("expected 3 entries after drop, got %d", len(out))
	}
}

func TestPipelineStopsOnError(t *testing.T) {
	sentinel := errors.New("stage error")
	var secondCalled bool

	_, err := NewPipeline().
		Pipe(func(e []Entry) ([]Entry, error) {
			return e, sentinel
		}).
		Pipe(func(e []Entry) ([]Entry, error) {
			secondCalled = true
			return e, nil
		}).
		Run(pipelineEntries())

	if !errors.Is(err, sentinel) {
		t.Fatalf("expected sentinel error, got %v", err)
	}
	if secondCalled {
		t.Error("second stage should not have been called after error")
	}
}

func TestPipeFuncWrapsPlainFunc(t *testing.T) {
	out, err := NewPipeline().
		Pipe(PipeFunc(LowercaseKeys)).
		Run(pipelineEntries())

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	for _, e := range out {
		if e.Key != toLower(e.Key) {
			t.Errorf("expected lowercase key, got %q", e.Key)
		}
	}
}

// helpers to avoid importing strings in test file
func toUpper(s string) string { return UppercaseKeys([]Entry{{Key: s}})[0].Key }
func toLower(s string) string { return LowercaseKeys([]Entry{{Key: s}})[0].Key }
