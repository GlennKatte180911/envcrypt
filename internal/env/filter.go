package env

import "strings"

// Filter returns only the entries whose keys satisfy the predicate fn.
func Filter(entries []Entry, fn func(key string) bool) []Entry {
	out := make([]Entry, 0, len(entries))
	for _, e := range entries {
		if fn(e.Key) {
			out = append(out, e)
		}
	}
	return out
}

// WithPrefix returns entries whose keys start with the given prefix.
func WithPrefix(entries []Entry, prefix string) []Entry {
	return Filter(entries, func(key string) bool {
		return strings.HasPrefix(key, prefix)
	})
}

// WithoutPrefix returns entries whose keys do NOT start with the given prefix.
func WithoutPrefix(entries []Entry, prefix string) []Entry {
	return Filter(entries, func(key string) bool {
		return !strings.HasPrefix(key, prefix)
	})
}

// StripPrefix removes the given prefix from every entry key that carries it.
// Entries without the prefix are returned unchanged.
func StripPrefix(entries []Entry, prefix string) []Entry {
	out := make([]Entry, len(entries))
	for i, e := range entries {
		if strings.HasPrefix(e.Key, prefix) {
			e.Key = strings.TrimPrefix(e.Key, prefix)
		}
		out[i] = e
	}
	return out
}

// Unique deduplicates entries by key, keeping the last occurrence.
func Unique(entries []Entry) []Entry {
	seen := make(map[string]int, len(entries))
	for i, e := range entries {
		seen[e.Key] = i
	}
	out := make([]Entry, 0, len(seen))
	for _, e := range entries {
		if idx, ok := seen[e.Key]; ok && idx == indexOf(entries, e.Key) {
			out = append(out, e)
			delete(seen, e.Key)
		}
	}
	return out
}

// indexOf returns the index of the last entry with the given key.
func indexOf(entries []Entry, key string) int {
	for i := len(entries) - 1; i >= 0; i-- {
		if entries[i].Key == key {
			return i
		}
	}
	return -1
}
