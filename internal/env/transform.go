package env

import "strings"

// TransformFunc is a function that transforms a single Entry.
// It returns the transformed Entry and whether to keep it.
type TransformFunc func(Entry) (Entry, bool)

// Transform applies a TransformFunc to each Entry in the slice,
// collecting results where the keep flag is true.
func Transform(entries []Entry, fn TransformFunc) []Entry {
	result := make([]Entry, 0, len(entries))
	for _, e := range entries {
		if out, keep := fn(e); keep {
			result = append(result, out)
		}
	}
	return result
}

// UppercaseKeys returns a TransformFunc that converts every key to uppercase.
func UppercaseKeys() TransformFunc {
	return func(e Entry) (Entry, bool) {
		e.Key = strings.ToUpper(e.Key)
		return e, true
	}
}

// LowercaseKeys returns a TransformFunc that converts every key to lowercase.
func LowercaseKeys() TransformFunc {
	return func(e Entry) (Entry, bool) {
		e.Key = strings.ToLower(e.Key)
		return e, true
	}
}

// RenameKey returns a TransformFunc that renames a specific key.
// Entries whose key does not match oldKey are passed through unchanged.
func RenameKey(oldKey, newKey string) TransformFunc {
	return func(e Entry) (Entry, bool) {
		if e.Key == oldKey {
			e.Key = newKey
		}
		return e, true
	}
}

// MaskValues returns a TransformFunc that replaces every value with mask.
// Useful for logging or displaying entries without exposing secrets.
func MaskValues(mask string) TransformFunc {
	return func(e Entry) (Entry, bool) {
		e.Value = mask
		return e, true
	}
}

// DropEmpty returns a TransformFunc that removes entries with empty values.
func DropEmpty() TransformFunc {
	return func(e Entry) (Entry, bool) {
		return e, e.Value != ""
	}
}

// Chain combines multiple TransformFuncs into one, applying them in order.
// If any function signals drop (keep=false), the entry is discarded immediately.
func Chain(fns ...TransformFunc) TransformFunc {
	return func(e Entry) (Entry, bool) {
		for _, fn := range fns {
			var keep bool
			e, keep = fn(e)
			if !keep {
				return e, false
			}
		}
		return e, true
	}
}
