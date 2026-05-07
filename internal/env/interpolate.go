package env

import (
	"fmt"
	"strings"
)

// InterpolateError is returned when a variable reference cannot be resolved.
type InterpolateError struct {
	Key string
}

func (e *InterpolateError) Error() string {
	return fmt.Sprintf("interpolate: unresolved variable reference: %q", e.Key)
}

// Interpolate expands ${VAR} and $VAR references within entry values using
// other entries in the slice as the source of substitutions. References to
// unknown keys are returned as an InterpolateError. Entries are processed in
// order, so earlier entries cannot reference later ones.
func Interpolate(entries []Entry) ([]Entry, error) {
	resolved := make(map[string]string, len(entries))
	out := make([]Entry, 0, len(entries))

	for _, e := range entries {
		val, err := expand(e.Value, resolved)
		if err != nil {
			return nil, err
		}
		resolved[e.Key] = val
		out = append(out, Entry{Key: e.Key, Value: val, Comment: e.Comment})
	}
	return out, nil
}

// expand replaces all $VAR and ${VAR} tokens in s using the lookup map.
func expand(s string, lookup map[string]string) (string, error) {
	var sb strings.Builder
	i := 0
	for i < len(s) {
		if s[i] != '$' {
			sb.WriteByte(s[i])
			i++
			continue
		}
		// consume '$'
		i++
		if i >= len(s) {
			sb.WriteByte('$')
			break
		}
		var key string
		if s[i] == '{' {
			// ${VAR} form
			i++ // skip '{'
			start := i
			for i < len(s) && s[i] != '}' {
				i++
			}
			key = s[start:i]
			if i < len(s) {
				i++ // skip '}'
			}
		} else {
			// $VAR form
			start := i
			for i < len(s) && isVarChar(s[i]) {
				i++
			}
			key = s[start:i]
		}
		if key == "" {
			sb.WriteByte('$')
			continue
		}
		val, ok := lookup[key]
		if !ok {
			return "", &InterpolateError{Key: key}
		}
		sb.WriteString(val)
	}
	return sb.String(), nil
}

func isVarChar(c byte) bool {
	return (c >= 'A' && c <= 'Z') || (c >= 'a' && c <= 'z') ||
		(c >= '0' && c <= '9') || c == '_'
}
