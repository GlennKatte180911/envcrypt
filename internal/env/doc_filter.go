// Package env provides utilities for parsing, serialising, and manipulating
// .env files.
//
// # Entry
//
// Entry is the fundamental unit of a .env file. Each non-blank, non-comment
// line produces one Entry with a Key and Value. Comment-only lines are
// preserved as Entry values with IsComment set to true, enabling lossless
// round-trips.
//
// # Filter helpers
//
// Filter, WithPrefix, WithoutPrefix, StripPrefix, and Unique operate on
// []Entry slices without mutating the originals, making it straightforward
// to extract environment subsets — for example, pulling all DB_ variables
// from a shared .env before passing them to a service:
//
//	dbEntries := env.WithPrefix(all, "DB_")
//	clean := env.StripPrefix(dbEntries, "DB_")
//
// All filter functions return new slices; the input slice is never modified.
package env
