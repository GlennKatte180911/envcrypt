package env

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"os"
	"sync"
	"time"
)

// WatchEvent describes a change detected in a watched env file.
type WatchEvent struct {
	Path    string
	OldHash string
	NewHash string
	At      time.Time
}

// Watcher polls an env file for changes and emits events on a channel.
type Watcher struct {
	path     string
	interval time.Duration
	lastHash string
	mu       sync.Mutex
	stopCh   chan struct{}
}

// NewWatcher creates a Watcher for the given file path and poll interval.
func NewWatcher(path string, interval time.Duration) *Watcher {
	return &Watcher{
		path:     path,
		interval: interval,
		stopCh:   make(chan struct{}),
	}
}

// Start begins polling and sends WatchEvents to the returned channel.
// Call Stop to terminate the watcher.
func (w *Watcher) Start() (<-chan WatchEvent, error) {
	hash, err := hashFile(w.path)
	if err != nil {
		return nil, fmt.Errorf("watch: initial hash: %w", err)
	}
	w.mu.Lock()
	w.lastHash = hash
	w.mu.Unlock()

	ch := make(chan WatchEvent, 4)
	go func() {
		defer close(ch)
		ticker := time.NewTicker(w.interval)
		defer ticker.Stop()
		for {
			select {
			case <-w.stopCh:
				return
			case t := <-ticker.C:
				newHash, err := hashFile(w.path)
				if err != nil {
					continue
				}
				w.mu.Lock()
				old := w.lastHash
				if newHash != old {
					w.lastHash = newHash
					w.mu.Unlock()
					ch <- WatchEvent{Path: w.path, OldHash: old, NewHash: newHash, At: t}
				} else {
					w.mu.Unlock()
				}
			}
		}
	}()
	return ch, nil
}

// Stop terminates the polling goroutine.
func (w *Watcher) Stop() {
	close(w.stopCh)
}

func hashFile(path string) (string, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return "", err
	}
	sum := sha256.Sum256(data)
	return hex.EncodeToString(sum[:]), nil
}
