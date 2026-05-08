// Package env — clone and subset utilities.
//
// Clone and Subset provide safe, non-mutating helpers for working with
// slices of Entry values.
//
// # Clone
//
// Clone returns a deep copy of an []Entry slice so that downstream
// transformations cannot accidentally modify the caller's original data.
//
//	original := []env.Entry{{Key: "DB_URL", Value: "postgres://..."}}
//	safe := env.Clone(original)
//
// # Subset / Exclude
//
// Subset and Exclude allow callers to project an entry slice down to
// only the keys they care about (or to drop sensitive keys before
// passing entries to a third-party component).
//
//	public := env.Exclude(all, []string{"SECRET_KEY", "DB_PASSWORD"})
package env
