package env

import (
	"strings"
)

// RedactOption controls how values are redacted.
type RedactOption func(*redactConfig)

type redactConfig {
	keys        []string
	suffixes    []string
	placeholder string
}

// WithRedactKeys specifies exact key names whose values should be redacted.
func WithRedactKeys(keys ...string) RedactOption {
	return func(c *redactConfig) {
		c.keys = append(c.keys, keys...)
	}
}

// WithRedactSuffixes specifies key suffixes (case-insensitive) that trigger redaction.
// For example, "_SECRET", "_PASSWORD", "_TOKEN".
func WithRedactSuffixes(suffixes ...string) RedactOption {
	return func(c *redactConfig) {
		for _, s := range suffixes {
			c.suffixes = append(c.suffixes, strings.ToUpper(s))
		}
	}
}

// WithRedactPlaceholder sets the placeholder string used instead of the real value.
// Defaults to "***".
func WithRedactPlaceholder(placeholder string) RedactOption {
	return func(c *redactConfig) {
		c.placeholder = placeholder
	}
}

// Redact returns a copy of entries where sensitive values are replaced with a
// placeholder. Keys are matched either by exact name or by suffix.
func Redact(entries []Entry, opts ...RedactOption) []Entry {
	cfg := &redactConfig{
		placeholder: "***",
	}
	for _, o := range opts {
		o(cfg)
	}

	exactKeys := make(map[string]struct{}, len(cfg.keys))
	for _, k := range cfg.keys {
		exactKeys[k] = struct{}{}
	}

	result := make([]Entry, len(entries))
	for i, e := range entries {
		if shouldRedact(e.Key, exactKeys, cfg.suffixes) {
			result[i] = Entry{Key: e.Key, Value: cfg.placeholder}
		} else {
			result[i] = e
		}
	}
	return result
}

func shouldRedact(key string, exactKeys map[string]struct{}, suffixes []string) bool {
	if _, ok := exactKeys[key]; ok {
		return true
	}
	upper := strings.ToUpper(key)
	for _, s := range suffixes {
		if strings.HasSuffix(upper, s) {
			return true
		}
	}
	return false
}
