package env

import (
	"sort"
	"strings"
)

// SortOrder defines the ordering direction for entries.
type SortOrder int

const (
	// SortAscending orders entries from A to Z.
	SortAscending SortOrder = iota
	// SortDescending orders entries from Z to A.
	SortDescending
)

// SortByKey returns a new slice of entries sorted alphabetically by key.
// The original slice is not modified.
func SortByKey(entries []Entry, order SortOrder) []Entry {
	out := make([]Entry, len(entries))
	copy(out, entries)
	sort.SliceStable(out, func(i, j int) bool {
		ki := strings.ToLower(out[i].Key)
		kj := strings.ToLower(out[j].Key)
		if order == SortDescending {
			return ki > kj
		}
		return ki < kj
	})
	return out
}

// SortByValue returns a new slice of entries sorted alphabetically by value.
// The original slice is not modified.
func SortByValue(entries []Entry, order SortOrder) []Entry {
	out := make([]Entry, len(entries))
	copy(out, entries)
	sort.SliceStable(out, func(i, j int) bool {
		vi := strings.ToLower(out[i].Value)
		vj := strings.ToLower(out[j].Value)
		if order == SortDescending {
			return vi > vj
		}
		return vi < vj
	})
	return out
}

// StableSort returns a new slice of entries sorted by the provided key function.
// Entries that produce equal keys retain their original relative order.
func StableSort(entries []Entry, keyFn func(Entry) string, order SortOrder) []Entry {
	out := make([]Entry, len(entries))
	copy(out, entries)
	sort.SliceStable(out, func(i, j int) bool {
		ki := strings.ToLower(keyFn(out[i]))
		kj := strings.ToLower(keyFn(out[j]))
		if order == SortDescending {
			return ki > kj
		}
		return ki < kj
	})
	return out
}
