package env

import (
	"sort"
	"testing"
)

var groupEntries = []Entry{
	{Key: "DB_HOST", Value: "localhost"},
	{Key: "DB_PORT", Value: "5432"},
	{Key: "APP_NAME", Value: "envcrypt"},
	{Key: "APP_ENV", Value: "production"},
	{Key: "SECRET", Value: "xyz"},
}

func TestGroupByPrefix(t *testing.T) {
	groups := GroupByPrefix(groupEntries, "_")

	if len(groups["DB"]) != 2 {
		t.Errorf("expected 2 DB entries, got %d", len(groups["DB"]))
	}
	if len(groups["APP"]) != 2 {
		t.Errorf("expected 2 APP entries, got %d", len(groups["APP"]))
	}
	if len(groups["SECRET"]) != 1 {
		t.Errorf("expected 1 SECRET entry, got %d", len(groups["SECRET"]))
	}
}

func TestGroupBySuffix(t *testing.T) {
	entries := []Entry{
		{Key: "DB_HOST", Value: "localhost"},
		{Key: "APP_HOST", Value: "0.0.0.0"},
		{Key: "DB_PORT", Value: "5432"},
		{Key: "STANDALONE", Value: "yes"},
	}
	groups := GroupBySuffix(entries, "_")

	if len(groups["HOST"]) != 2 {
		t.Errorf("expected 2 HOST entries, got %d", len(groups["HOST"]))
	}
	if len(groups["PORT"]) != 1 {
		t.Errorf("expected 1 PORT entry, got %d", len(groups["PORT"]))
	}
	if len(groups["STANDALONE"]) != 1 {
		t.Errorf("expected 1 STANDALONE entry, got %d", len(groups["STANDALONE"]))
	}
}

func TestGroupCustomKeyFn(t *testing.T) {
	// bucket everything into "odd" or "even" by index
	groups := Group(groupEntries, func(e Entry) string {
		if len(e.Key)%2 == 0 {
			return "even"
		}
		return "odd"
	})
	total := 0
	for _, v := range groups {
		total += len(v)
	}
	if total != len(groupEntries) {
		t.Errorf("expected %d total entries after grouping, got %d", len(groupEntries), total)
	}
}

func TestFlattenPreservesAllEntries(t *testing.T) {
	groups := GroupByPrefix(groupEntries, "_")
	flat := Flatten(groups)

	if len(flat) != len(groupEntries) {
		t.Errorf("expected %d entries after flatten, got %d", len(groupEntries), len(flat))
	}

	// Verify every original key is present.
	keys := make(map[string]bool)
	for _, e := range flat {
		keys[e.Key] = true
	}
	for _, e := range groupEntries {
		if !keys[e.Key] {
			t.Errorf("key %q missing after flatten", e.Key)
		}
	}
}

func TestGroupByPrefixEmptyEntries(t *testing.T) {
	groups := GroupByPrefix([]Entry{}, "_")
	if len(groups) != 0 {
		t.Errorf("expected empty groups, got %d buckets", len(groups))
	}
}

func TestGroupOrderWithinBucket(t *testing.T) {
	groups := GroupByPrefix(groupEntries, "_")
	db := groups["DB"]
	keys := make([]string, len(db))
	for i, e := range db {
		keys[i] = e.Key
	}
	if !sort.StringsAreSorted([]string{keys[0]}) {
		// just verify the two DB entries are present in insertion order
	}
	if keys[0] != "DB_HOST" || keys[1] != "DB_PORT" {
		t.Errorf("expected DB entries in insertion order, got %v", keys)
	}
}
