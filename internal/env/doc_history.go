// Package env — History
//
// History maintains an ordered, append-only log of diff events that occurred
// to a set of env entries over time. Each event captures the timestamp, a
// human-readable label, and the full DiffEntry slice produced by Diff.
//
// Typical usage:
//
//	h := env.NewHistory()
//
//	// Load an initial set of entries.
//	before := loadEntries()
//
//	// Apply some change and record it.
//	after := applyChange(before)
//	h.Record("apply staging overrides", before, after)
//
//	// Inspect what changed.
//	for _, he := range h.Entries() {
//		fmt.Printf("%s — %s\n", he.Timestamp.Format(time.RFC3339), he.Label)
//		for _, d := range he.Diff {
//			fmt.Printf("  [%s] %s\n", d.Kind, d.Key)
//		}
//	}
//
// If the diff between before and after contains no changes, Record is a
// no-op and no entry is appended.
package env
