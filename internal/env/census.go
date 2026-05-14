package env

// Census provides aggregate statistics over a slice of Entry values.
// It is useful for reporting, dashboards, and validation summaries.

// CensusResult holds the computed statistics for a set of entries.
type CensusResult struct {
	// Total is the number of entries examined.
	Total int
	// EmptyValues is the count of entries whose value is the empty string.
	EmptyValues int
	// UniqueKeys is the count of distinct keys.
	UniqueKeys int
	// DuplicateKeys is the count of keys that appear more than once.
	DuplicateKeys int
	// UppercaseKeys is the count of keys that are entirely uppercase (and non-empty).
	UppercaseKeys int
	// LowercaseKeys is the count of keys that are entirely lowercase (and non-empty).
	LowercaseKeys int
	// AverageValueLen is the mean byte-length of all entry values.
	AverageValueLen float64
}

// TakeCensus computes aggregate statistics over entries.
func TakeCensus(entries []Entry) CensusResult {
	if len(entries) == 0 {
		return CensusResult{}
	}

	keyCounts := make(map[string]int, len(entries))
	totalValueLen := 0
	uppercase := 0
	lowercase := 0
	empty := 0

	for _, e := range entries {
		keyCounts[e.Key]++
		totalValueLen += len(e.Value)
		if e.Value == "" {
			empty++
		}
		if isAllUppercase(e.Key) {
			uppercase++
		}
		if isAllLowercase(e.Key) {
			lowercase++
		}
	}

	duplicates := 0
	for _, count := range keyCounts {
		if count > 1 {
			duplicates++
		}
	}

	return CensusResult{
		Total:           len(entries),
		EmptyValues:     empty,
		UniqueKeys:      len(keyCounts),
		DuplicateKeys:   duplicates,
		UppercaseKeys:   uppercase,
		LowercaseKeys:   lowercase,
		AverageValueLen: float64(totalValueLen) / float64(len(entries)),
	}
}

func isAllUppercase(s string) bool {
	if s == "" {
		return false
	}
	for _, r := range s {
		if r >= 'a' && r <= 'z' {
			return false
		}
	}
	return true
}

func isAllLowercase(s string) bool {
	if s == "" {
		return false
	}
	for _, r := range s {
		if r >= 'A' && r <= 'Z' {
			return false
		}
	}
	return true
}
