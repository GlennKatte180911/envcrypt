package env

// Compare provides utilities for comparing two sets of env entries
// beyond simple diffing — including equality checks, symmetric difference,
// intersection, and scoring similarity between two environments.

import (
	"fmt"
	"strings"
)

// CompareResult holds the full comparison between a base and target set of entries.
type CompareResult struct {
	OnlyInBase   []Entry // keys present in base but not in target
	OnlyInTarget []Entry // keys present in target but not in base
	Changed      []EntryPair // keys present in both but with different values
	Unchanged    []Entry // keys present in both with identical values
}

// EntryPair holds a base and target entry for the same key.
type EntryPair struct {
	Base   Entry
	Target Entry
}

// Equal reports whether base and target contain exactly the same keys and values.
func Equal(base, target []Entry) bool {
	r := Compare(base, target)
	return len(r.OnlyInBase) == 0 && len(r.OnlyInTarget) == 0 && len(r.Changed) == 0
}

// Compare performs a full structural comparison between base and target entries.
// Order of entries does not affect the result.
func Compare(base, target []Entry) CompareResult {
	baseMap := ToMap(base)
	targetMap := ToMap(target)

	var result CompareResult

	for _, e := range base {
		tv, ok := targetMap[e.Key]
		if !ok {
			result.OnlyInBase = append(result.OnlyInBase, e)
		} else if tv != e.Value {
			result.Changed = append(result.Changed, EntryPair{
				Base:   e,
				Target: Entry{Key: e.Key, Value: tv},
			})
		} else {
			result.Unchanged = append(result.Unchanged, e)
		}
	}

	for _, e := range target {
		if _, ok := baseMap[e.Key]; !ok {
			result.OnlyInTarget = append(result.OnlyInTarget, e)
		}
	}

	return result
}

// Intersection returns entries whose keys and values are identical in both sets.
func Intersection(base, target []Entry) []Entry {
	return Compare(base, target).Unchanged
}

// SymmetricDifference returns entries that appear in one set but not the other
// (by key), regardless of value.
func SymmetricDifference(base, target []Entry) []Entry {
	r := Compare(base, target)
	out := make([]Entry, 0, len(r.OnlyInBase)+len(r.OnlyInTarget))
	out = append(out, r.OnlyInBase...)
	out = append(out, r.OnlyInTarget...)
	return out
}

// SimilarityScore returns a value in [0.0, 1.0] representing how similar two
// entry sets are. A score of 1.0 means they are identical; 0.0 means they
// share no keys at all.
func SimilarityScore(base, target []Entry) float64 {
	if len(base) == 0 && len(target) == 0 {
		return 1.0
	}
	r := Compare(base, target)
	total := len(r.OnlyInBase) + len(r.OnlyInTarget) + len(r.Changed) + len(r.Unchanged)
	if total == 0 {
		return 1.0
	}
	return float64(len(r.Unchanged)) / float64(total)
}

// Summary returns a human-readable summary string of a CompareResult.
func (r CompareResult) Summary() string {
	var sb strings.Builder
	fmt.Fprintf(&sb, "unchanged=%d added=%d removed=%d changed=%d",
		len(r.Unchanged), len(r.OnlyInTarget), len(r.OnlyInBase), len(r.Changed))
	return sb.String()
}

// HasDifferences reports whether any differences exist between the two sets.
func (r CompareResult) HasDifferences() bool {
	return len(r.OnlyInBase) > 0 || len(r.OnlyInTarget) > 0 || len(r.Changed) > 0
}
