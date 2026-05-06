package env

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

func writeTempEnvFile(t *testing.T, content string) string {
	t.Helper()
	dir := t.TempDir()
	path := filepath.Join(dir, ".env")
	if err := os.WriteFile(path, []byte(content), 0o600); err != nil {
		t.Fatalf("writeTempEnvFile: %v", err)
	}
	return path
}

func TestWatcherDetectsChange(t *testing.T) {
	path := writeTempEnvFile(t, "KEY=original\n")
	w := NewWatcher(path, 20*time.Millisecond)
	ch, err := w.Start()
	if err != nil {
		t.Fatalf("Start: %v", err)
	}
	defer w.Stop()

	time.Sleep(40 * time.Millisecond)
	if err := os.WriteFile(path, []byte("KEY=changed\n"), 0o600); err != nil {
		t.Fatalf("WriteFile: %v", err)
	}

	select {
	case ev := <-ch:
		if ev.Path != path {
			t.Errorf("expected path %q, got %q", path, ev.Path)
		}
		if ev.OldHash == ev.NewHash {
			t.Error("expected hashes to differ")
		}
		if ev.At.IsZero() {
			t.Error("expected non-zero timestamp")
		}
	case <-time.After(500 * time.Millisecond):
		t.Fatal("timed out waiting for watch event")
	}
}

func TestWatcherNoEventWhenUnchanged(t *testing.T) {
	path := writeTempEnvFile(t, "KEY=stable\n")
	w := NewWatcher(path, 20*time.Millisecond)
	ch, err := w.Start()
	if err != nil {
		t.Fatalf("Start: %v", err)
	}
	defer w.Stop()

	select {
	case ev := <-ch:
		t.Errorf("unexpected event: %+v", ev)
	case <-time.After(120 * time.Millisecond):
		// expected: no change
	}
}

func TestWatcherStartMissingFile(t *testing.T) {
	w := NewWatcher("/nonexistent/.env", 20*time.Millisecond)
	_, err := w.Start()
	if err == nil {
		t.Fatal("expected error for missing file")
	}
}

func TestWatcherStop(t *testing.T) {
	path := writeTempEnvFile(t, "A=1\n")
	w := NewWatcher(path, 20*time.Millisecond)
	ch, err := w.Start()
	if err != nil {
		t.Fatalf("Start: %v", err)
	}
	w.Stop()
	// channel should close shortly after stop
	select {
	case _, ok := <-ch:
		if ok {
			t.Log("received event before close (acceptable)")
		}
	case <-time.After(300 * time.Millisecond):
		t.Fatal("channel did not close after Stop")
	}
}
