// Package env provides the Lineage type for tracking the history of changes
// applied to a collection of env entries over time.
//
// # Overview
//
// A Lineage records each mutation as a LineageRecord containing a timestamp,
// a ChangeKind label, a human-readable description, and the computed Diff
// between the previous and updated entry slices. Only mutations that produce
// a non-empty diff are stored, so no-op operations are silently ignored.
//
// # Usage
//
//	l := env.NewLineage()
//
//	base := []env.Entry{{Key: "APP_ENV", Value: "development"}}
//	updated := []env.Entry{{Key: "APP_ENV", Value: "production"}}
//
//	l.Record(env.ChangeKindSet, "promote to production", base, updated)
//
//	for _, r := range l.Records() {
//		fmt.Printf("[%s] %s — %d change(s)\n", r.Kind, r.Label, len(r.Diff))
//	}
//
//	fmt.Println("Keys ever changed:", l.ChangedKeys())
package env
