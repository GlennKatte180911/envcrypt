package env

import (
	"testing"
)

var historyBase = []Entry{
	{Key: "APP_ENV", Value: "development"},
	{Key: "PORT", Value: "8080"},
}

func TestHistoryRecordNoChangeSkipped(t *testing.T) {
	h := NewHistory()
	h.Record("no-op", historyBase, historyBase)
	if h.Len() != 0 {
		t.Fatalf("expected 0 entries, got %d", h.Len())
	}
}

func TestHistoryRecordAdded(t *testing.T) {
	h := NewHistory()
	after := append(historyBase, Entry{Key: "DEBUG", Value: "true"})
	h.Record("add DEBUG", historyBase, after)
	if h.Len() != 1 {
		t.Fatalf("expected 1 entry, got %d", h.Len())
	}
	last, ok := h.Last()
	if !ok {
		t.Fatal("expected Last() to return entry")
	}
	if last.Label != "add DEBUG" {
		t.Errorf("unexpected label: %q", last.Label)
	}
}

func TestHistoryRecordRemoved(t *testing.T) {
	h := NewHistory()
	after := historyBase[:1]
	h.Record("remove PORT", historyBase, after)
	if h.Len() != 1 {
		t.Fatalf("expected 1 entry, got %d", h.Len())
	}
}

func TestHistoryMultipleRecords(t *testing.T) {
	h := NewHistory()
	after1 := append(historyBase, Entry{Key: "DEBUG", Value: "true"})
	after2 := append(after1, Entry{Key: "LOG_LEVEL", Value: "info"})
	h.Record("step1", historyBase, after1)
	h.Record("step2", after1, after2)
	if h.Len() != 2 {
		t.Fatalf("expected 2 entries, got %d", h.Len())
	}
	entries := h.Entries()
	if entries[0].Label != "step1" || entries[1].Label != "step2" {
		t.Errorf("unexpected labels: %q %q", entries[0].Label, entries[1].Label)
	}
}

func TestHistoryLastEmpty(t *testing.T) {
	h := NewHistory()
	_, ok := h.Last()
	if ok {
		t.Fatal("expected Last() to return false on empty history")
	}
}

func TestHistoryClear(t *testing.T) {
	h := NewHistory()
	after := append(historyBase, Entry{Key: "X", Value: "1"})
	h.Record("add X", historyBase, after)
	h.Clear()
	if h.Len() != 0 {
		t.Fatalf("expected 0 entries after Clear, got %d", h.Len())
	}
}

func TestHistoryEntriesIsolated(t *testing.T) {
	h := NewHistory()
	after := append(historyBase, Entry{Key: "X", Value: "1"})
	h.Record("add X", historyBase, after)
	got := h.Entries()
	got[0].Label = "mutated"
	if last, _ := h.Last(); last.Label == "mutated" {
		t.Error("Entries() should return a copy, not expose internal slice")
	}
}

func TestHistoryString(t *testing.T) {
	h := NewHistory()
	if s := h.String(); s != "history: no entries" {
		t.Errorf("unexpected empty string: %q", s)
	}
	after := append(historyBase, Entry{Key: "X", Value: "1"})
	h.Record("add X", historyBase, after)
	if s := h.String(); s == "history: no entries" {
		t.Error("expected non-empty string after recording")
	}
}
