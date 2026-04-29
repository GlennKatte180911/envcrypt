package env

import "fmt"

// ToMap converts a slice of Entry values into a map for convenient key-based
// lookup. Duplicate keys are not allowed; an error is returned if any key
// appears more than once.
func ToMap(entries []Entry) (map[string]string, error) {
	m := make(map[string]string, len(entries))
	for _, e := range entries {
		if _, exists := m[e.Key]; exists {
			return nil, fmt.Errorf("duplicate key %q", e.Key)
		}
		m[e.Key] = e.Value
	}
	return m, nil
}

// FromMap converts a map into an ordered slice of Entry values.
// The provided keys slice controls the output order; any key present in keys
// but absent from m is silently skipped.
func FromMap(m map[string]string, keys []string) []Entry {
	entries := make([]Entry, 0, len(keys))
	for _, k := range keys {
		if v, ok := m[k]; ok {
			entries = append(entries, Entry{Key: k, Value: v})
		}
	}
	return entries
}

// Keys returns the keys from a slice of entries in order.
func Keys(entries []Entry) []string {
	keys := make([]string, len(entries))
	for i, e := range entries {
		keys[i] = e.Key
	}
	return keys
}

// Merge combines base entries with override entries. Keys present in overrides
// replace those in base; new keys from overrides are appended at the end.
func Merge(base, overrides []Entry) []Entry {
	result := make([]Entry, len(base))
	copy(result, base)

	index := make(map[string]int, len(base))
	for i, e := range result {
		index[e.Key] = i
	}

	for _, o := range overrides {
		if i, exists := index[o.Key]; exists {
			result[i].Value = o.Value
		} else {
			index[o.Key] = len(result)
			result = append(result, o)
		}
	}
	return result
}
