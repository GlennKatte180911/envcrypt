package env

// DiffKind describes the type of change between two sets of entries.
type DiffKind string

const (
	// DiffAdded indicates a key present in the new set but not the old.
	DiffAdded DiffKind = "added"
	// DiffRemoved indicates a key present in the old set but not the new.
	DiffRemoved DiffKind = "removed"
	// DiffChanged indicates a key present in both sets but with a different value.
	DiffChanged DiffKind = "changed"
)

// DiffEntry represents a single difference between two sets of env entries.
type DiffEntry struct {
	Key      string
	Kind     DiffKind
	OldValue string
	NewValue string
}

// Diff compares two slices of Entry and returns the differences.
// oldEntries is the baseline; newEntries is the updated set.
func Diff(oldEntries, newEntries []Entry) []DiffEntry {
	oldMap := ToMap(oldEntries)
	newMap := ToMap(newEntries)

	var result []DiffEntry

	// Detect removed and changed keys.
	for _, e := range oldEntries {
		newVal, exists := newMap[e.Key]
		if !exists {
			result = append(result, DiffEntry{
				Key:      e.Key,
				Kind:     DiffRemoved,
				OldValue: e.Value,
			})
		} else if newVal != e.Value {
			result = append(result, DiffEntry{
				Key:      e.Key,
				Kind:     DiffChanged,
				OldValue: e.Value,
				NewValue: newVal,
			})
		}
	}

	// Detect added keys.
	for _, e := range newEntries {
		if _, exists := oldMap[e.Key]; !exists {
			result = append(result, DiffEntry{
				Key:      e.Key,
				Kind:     DiffAdded,
				NewValue: e.Value,
			})
		}
	}

	return result
}

// HasChanges returns true if Diff produces any differences.
func HasChanges(oldEntries, newEntries []Entry) bool {
	return len(Diff(oldEntries, newEntries)) > 0
}
