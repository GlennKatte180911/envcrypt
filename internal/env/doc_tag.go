// Package env provides utilities for parsing, manipulating, and serializing
// environment variable entries.
//
// # Tag
//
// The tag sub-feature lets callers attach lightweight string labels to
// individual entries. Tags are stored inside the entry's Comment field using
// the reserved token "#tag:<comma-separated-list>" so they survive
// serialisation round-trips through [Parse] and [Serialize].
//
// Typical usage:
//
//	entries = env.TagEntry(entries, []env.Tag{"secret"}, "DB_PASS", "API_KEY")
//
//	// Later, isolate all secrets:
//	secrets := env.FilterByTag(entries, "secret")
//
//	// Or build an index of every tag used in the file:
//	idx := env.BuildTagIndex(entries)
//	for tag, keys := range idx {
//		fmt.Printf("%s → %v\n", tag, keys)
//	}
//
// Tags are merged non-destructively: applying the same tag twice does not
// create duplicates, and unrelated comment text is preserved.
package env
