package env

// Clone returns a deep copy of the given entries slice.
// Modifications to the returned slice do not affect the original.
func Clone(entries []Entry) []Entry {
	if entries == nil {
		return nil
	}
	out := make([]Entry, len(entries))
	copy(out, entries)
	return out
}

// CloneMap returns a shallow copy of the given string map.
func CloneMap(m map[string]string) map[string]string {
	if m == nil {
		return nil
	}
	out := make(map[string]string, len(m))
	for k, v := range m {
		out[k] = v
	}
	return out
}

// Subset returns a new slice containing only the entries whose keys
// appear in the provided keys list. Order follows the original slice.
func Subset(entries []Entry, keys []string) []Entry {
	set := make(map[string]struct{}, len(keys))
	for _, k := range keys {
		set[k] = struct{}{}
	}
	var out []Entry
	for _, e := range entries {
		if _, ok := set[e.Key]; ok {
			out = append(out, e)
		}
	}
	return out
}

// Exclude returns a new slice omitting any entries whose keys appear
// in the provided keys list. Order follows the original slice.
func Exclude(entries []Entry, keys []string) []Entry {
	set := make(map[string]struct{}, len(keys))
	for _, k := range keys {
		set[k] = struct{}{}
	}
	var out []Entry
	for _, e := range entries {
		if _, ok := set[e.Key]; !ok {
			out = append(out, e)
		}
	}
	return out
}
