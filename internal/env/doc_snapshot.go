// Package env provides utilities for parsing, serializing, and manipulating
// environment variable entries.
//
// # Snapshot
//
// The Snapshot type captures an immutable copy of a slice of [Entry] values
// at a specific point in time, along with an optional human-readable label
// and a UTC timestamp.
//
// Snapshots are useful for tracking changes to an env file over time, auditing
// drift between environments, or rolling back to a previous configuration.
//
// Use [NewSnapshot] to create a snapshot from a current slice of entries.
// Use [DiffFromSnapshot] to compare a snapshot against a newer set of entries
// and obtain a list of [DiffEntry] values describing what changed.
//
// [SnapshotStore] provides a simple in-memory ordered collection of snapshots
// with helpers to retrieve the latest snapshot or iterate over all stored ones.
package env
