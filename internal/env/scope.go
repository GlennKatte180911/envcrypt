package env

import "strings"

// Scope represents a named environment scope (e.g. "production", "staging")
// that carries a filtered, prefixed view of entries.
type Scope struct {
	Name    string
	entries []Entry
}

// NewScope creates a Scope by extracting entries whose keys begin with
// "<name>_" (case-insensitive) and stripping that prefix.
func NewScope(name string, entries []Entry) Scope {
	prefix := strings.ToUpper(name) + "_"
	var scoped []Entry
	for _, e := range entries {
		if strings.HasPrefix(strings.ToUpper(e.Key), prefix) {
			scoped = append(scoped, Entry{
				Key:     e.Key[len(prefix):],
				Value:   e.Value,
				Comment: e.Comment,
			})
		}
	}
	return Scope{Name: name, entries: scoped}
}

// Entries returns the scoped entries with their prefixes stripped.
func (s Scope) Entries() []Entry {
	return Clone(s.entries)
}

// Lookup returns the value for key within the scope, and whether it was found.
func (s Scope) Lookup(key string) (string, bool) {
	for _, e := range s.entries {
		if e.Key == key {
			return e.Value, true
		}
	}
	return "", false
}

// Promote re-attaches the scope prefix to all entries, returning them as
// top-level entries suitable for merging back into a full env set.
func (s Scope) Promote() []Entry {
	prefix := strings.ToUpper(s.Name) + "_"
	out := make([]Entry, len(s.entries))
	for i, e := range s.entries {
		out[i] = Entry{
			Key:     prefix + e.Key,
			Value:   e.Value,
			Comment: e.Comment,
		}
	}
	return out
}

// ScopeNames returns the distinct scope prefixes found in entries.
// A prefix is detected when a key contains an underscore and the portion
// before the first underscore is shared by at least two keys.
func ScopeNames(entries []Entry) []string {
	counts := make(map[string]int)
	for _, e := range entries {
		if idx := strings.Index(e.Key, "_"); idx > 0 {
			counts[strings.ToUpper(e.Key[:idx])]++
		}
	}
	var names []string
	for k, v := range counts {
		if v >= 2 {
			names = append(names, k)
		}
	}
	return names
}
