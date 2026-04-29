package env

// MergeStrategy defines how conflicting keys are handled during a merge.
type MergeStrategy int

const (
	// PreferBase keeps the value from the base entries when a key conflict occurs.
	PreferBase MergeStrategy = iota

	// PreferOverride replaces the base value with the override value on conflict.
	PreferOverride

	// ErrorOnConflict returns an error if any key exists in both sets.
	ErrorOnConflict
)

// MergeWithStrategy merges override into base using the provided strategy.
// It returns the merged slice and any error (only possible with ErrorOnConflict).
func MergeWithStrategy(base, override []Entry, strategy MergeStrategy) ([]Entry, error) {
	baseMap := make(map[string]int, len(base))
	result := make([]Entry, len(base))
	copy(result, base)

	for i, e := range result {
		baseMap[e.Key] = i
	}

	for _, e := range override {
		idx, exists := baseMap[e.Key]
		if !exists {
			result = append(result, e)
			baseMap[e.Key] = len(result) - 1
			continue
		}

		switch strategy {
		case PreferBase:
			// keep existing value — do nothing
		case PreferOverride:
			result[idx] = e
		case ErrorOnConflict:
			return nil, &ConflictError{Key: e.Key}
		}
	}

	return result, nil
}

// ConflictError is returned when ErrorOnConflict strategy encounters a duplicate key.
type ConflictError struct {
	Key string
}

func (e *ConflictError) Error() string {
	return "env: merge conflict on key \"" + e.Key + "\""
}
