package env

import (
	"testing"
)

var sortEntries = []Entry{
	{Key: "ZEBRA", Value: "last"},
	{Key: "APPLE", Value: "first"},
	{Key: "MANGO", Value: "middle"},
	{Key: "BANANA", Value: "second"},
}

func TestSortByKeyAscending(t *testing.T) {
	result := SortByKey(sortEntries, SortAscending)
	expected := []string{"APPLE", "BANANA", "MANGO", "ZEBRA"}
	for i, e := range result {
		if e.Key != expected[i] {
			t.Errorf("index %d: got %q, want %q", i, e.Key, expected[i])
		}
	}
}

func TestSortByKeyDescending(t *testing.T) {
	result := SortByKey(sortEntries, SortDescending)
	expected := []string{"ZEBRA", "MANGO", "BANANA", "APPLE"}
	for i, e := range result {
		if e.Key != expected[i] {
			t.Errorf("index %d: got %q, want %q", i, e.Key, expected[i])
		}
	}
}

func TestSortByKeyDoesNotMutateOriginal(t *testing.T) {
	originalFirst := sortEntries[0].Key
	_ = SortByKey(sortEntries, SortAscending)
	if sortEntries[0].Key != originalFirst {
		t.Errorf("original slice was mutated: first key is now %q", sortEntries[0].Key)
	}
}

func TestSortByValueAscending(t *testing.T) {
	result := SortByValue(sortEntries, SortAscending)
	expected := []string{"first", "last", "middle", "second"}
	for i, e := range result {
		if e.Value != expected[i] {
			t.Errorf("index %d: got %q, want %q", i, e.Value, expected[i])
		}
	}
}

func TestSortByValueDescending(t *testing.T) {
	result := SortByValue(sortEntries, SortDescending)
	expected := []string{"second", "middle", "last", "first"}
	for i, e := range result {
		if e.Value != expected[i] {
			t.Errorf("index %d: got %q, want %q", i, e.Value, expected[i])
		}
	}
}

func TestStableSortCustomKeyFn(t *testing.T) {
	// Sort by length of key ascending
	result := StableSort(sortEntries, func(e Entry) string {
		return string(rune(len(e.Key)))
	}, SortAscending)
	if len(result) != len(sortEntries) {
		t.Fatalf("expected %d entries, got %d", len(sortEntries), len(result))
	}
	for i := 1; i < len(result); i++ {
		if len(result[i].Key) < len(result[i-1].Key) {
			t.Errorf("entries not sorted by key length at index %d", i)
		}
	}
}

func TestSortEmptySlice(t *testing.T) {
	result := SortByKey([]Entry{}, SortAscending)
	if len(result) != 0 {
		t.Errorf("expected empty slice, got %d entries", len(result))
	}
}
