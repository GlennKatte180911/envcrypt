package env

import "fmt"

// FlattenOption configures how nested keys are flattened.
type FlattenOption func(*flattenConfig)

type flattenConfig struct {
	separator string
	prefix    string
}

func defaultFlattenConfig() *flattenConfig {
	return &flattenConfig{separator: "_"}
}

// WithSeparator sets the separator used when joining key segments.
func WithSeparator(sep string) FlattenOption {
	return func(c *flattenConfig) { c.separator = sep }
}

// WithFlattenPrefix prepends a prefix to every flattened key.
func WithFlattenPrefix(prefix string) FlattenOption {
	return func(c *flattenConfig) { c.prefix = prefix }
}

// FlattenNested takes a map of maps (e.g. from JSON config) and converts it
// into a flat []Entry slice by joining nested keys with the configured separator.
// Only string leaf values are included; non-string values are formatted with %v.
func FlattenNested(nested map[string]interface{}, opts ...FlattenOption) []Entry {
	cfg := defaultFlattenConfig()
	for _, o := range opts {
		o(cfg)
	}

	var entries []Entry
	var walk func(prefix string, val interface{})
	walk = func(prefix string, val interface{}) {
		switch v := val.(type) {
		case map[string]interface{}:
			for k, child := range v {
				newKey := k
				if prefix != "" {
					newKey = prefix + cfg.separator + k
				}
				walk(newKey, child)
			}
		case string:
			entries = append(entries, Entry{Key: prefix, Value: v})
		default:
			entries = append(entries, Entry{Key: prefix, Value: fmt.Sprintf("%v", v)})
		}
	}

	topKey := cfg.prefix
	walk(topKey, nested)
	return SortByKey(entries)
}

// UnflattenToMap reconstructs a nested map[string]interface{} from a flat []Entry
// by splitting keys on the given separator.
func UnflattenToMap(entries []Entry, sep string) map[string]interface{} {
	if sep == "" {
		sep = "_"
	}
	root := make(map[string]interface{})
	for _, e := range entries {
		setNested(root, e.Key, e.Value, sep)
	}
	return root
}

func setNested(m map[string]interface{}, key, value, sep string) {
	for i := 0; i < len(key); i++ {
		if key[i] == sep[0] && key[:i+1][i:] == sep {
			head := key[:i]
			tail := key[i+len(sep):]
			if _, ok := m[head]; !ok {
				m[head] = make(map[string]interface{})
			}
			if child, ok := m[head].(map[string]interface{}); ok {
				setNested(child, tail, value, sep)
				return
			}
			return
		}
	}
	m[key] = value
}
