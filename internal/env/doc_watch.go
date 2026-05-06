// Package env provides utilities for parsing, transforming, validating,
// and watching .env files.
//
// # File Watching
//
// The [Watcher] type polls a single .env file at a configurable interval
// and emits [WatchEvent] values whenever the file's SHA-256 content hash
// changes. This allows applications or CLI tools to react to out-of-band
// edits without relying on OS-specific filesystem notification APIs.
//
// Basic usage:
//
//	w := env.NewWatcher(".env", 2*time.Second)
//	ch, err := w.Start()
//	if err != nil {
//		log.Fatal(err)
//	}
//	defer w.Stop()
//
//	for ev := range ch {
//		fmt.Printf("env changed at %s (hash %s -> %s)\n",
//			ev.At.Format(time.RFC3339), ev.OldHash[:8], ev.NewHash[:8])
//	}
//
// Watcher.Stop closes the internal stop channel, which causes the
// background goroutine to exit and the event channel to be closed.
package env
