package env

import "strings"

// Group partitions a slice of Entry values into named buckets using the
// provided key function. Entries for which keyFn returns an empty string
// are collected under the empty-string bucket.
func Group(entries []Entry, keyFn func(Entry) string) map[string][]Entry {
	result := make(map[string][]Entry)
	for _, e := range entries {
		k := keyFn(e)
		result[k] = append(result[k], e)
	}
	return result
}

// GroupByPrefix partitions entries by the portion of the key that precedes the
// first occurrence of sep. For example, with sep "_" the key "DB_HOST" is
// placed in the "DB" bucket.
//
// Entries whose key does not contain sep are collected under the bucket named
// by their full key.
func GroupByPrefix(entries []Entry, sep string) map[string][]Entry {
	return Group(entries, func(e Entry) string {
		if idx := strings.Index(e.Key, sep); idx >= 0 {
			return e.Key[:idx]
		}
		return e.Key
	})
}

// GroupBySuffix partitions entries by the portion of the key that follows the
// last occurrence of sep.
//
// Entries whose key does not contain sep are collected under the bucket named
// by their full key.
func GroupBySuffix(entries []Entry, sep string) map[string][]Entry {
	return Group(entries, func(e Entry) string {
		if idx := strings.LastIndex(e.Key, sep); idx >= 0 {
			return e.Key[idx+len(sep):]
		}
		return e.Key
	})
}

// Flatten converts a grouped map back into a flat slice of entries. Groups are
// visited in an unspecified order; within each group the original order is
// preserved.
func Flatten(groups map[string][]Entry) []Entry {
	var out []Entry
	for _, entries := range groups {
		out = append(out, entries...)
	}
	return out
}
