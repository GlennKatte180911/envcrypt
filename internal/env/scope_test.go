package env

import (
	"sort"
	"testing"
)

var scopeEntries = []Entry{
	{Key: "PROD_DB_HOST", Value: "db.prod.example.com"},
	{Key: "PROD_DB_PORT", Value: "5432"},
	{Key: "STAGING_DB_HOST", Value: "db.staging.example.com"},
	{Key: "STAGING_DB_PORT", Value: "5432"},
	{Key: "APP_NAME", Value: "envcrypt"},
}

func TestNewScopeFiltersCorrectly(t *testing.T) {
	s := NewScope("prod", scopeEntries)
	if s.Name != "prod" {
		t.Fatalf("expected Name prod, got %s", s.Name)
	}
	entries := s.Entries()
	if len(entries) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(entries))
	}
}

func TestNewScopeStripsPrefix(t *testing.T) {
	s := NewScope("staging", scopeEntries)
	for _, e := range s.Entries() {
		if len(e.Key) == 0 || e.Key[:8] == "STAGING_" {
			t.Errorf("prefix not stripped from key %q", e.Key)
		}
	}
}

func TestScopeLookupFound(t *testing.T) {
	s := NewScope("prod", scopeEntries)
	v, ok := s.Lookup("DB_HOST")
	if !ok {
		t.Fatal("expected key DB_HOST to be found")
	}
	if v != "db.prod.example.com" {
		t.Errorf("unexpected value %q", v)
	}
}

func TestScopeLookupMissing(t *testing.T) {
	s := NewScope("prod", scopeEntries)
	_, ok := s.Lookup("NONEXISTENT")
	if ok {
		t.Fatal("expected key to be missing")
	}
}

func TestScopePromoteReattachesPrefix(t *testing.T) {
	s := NewScope("prod", scopeEntries)
	promoted := s.Promote()
	for _, e := range promoted {
		if len(e.Key) < 5 || e.Key[:5] != "PROD_" {
			t.Errorf("expected PROD_ prefix on key %q", e.Key)
		}
	}
}

func TestScopeEntriesIsACopy(t *testing.T) {
	s := NewScope("prod", scopeEntries)
	e := s.Entries()
	e[0].Value = "mutated"
	v, _ := s.Lookup("DB_HOST")
	if v == "mutated" {
		t.Error("Entries() should return a copy, not a reference")
	}
}

func TestScopeNames(t *testing.T) {
	names := ScopeNames(scopeEntries)
	sort.Strings(names)
	if len(names) != 2 {
		t.Fatalf("expected 2 scope names, got %d: %v", len(names), names)
	}
	if names[0] != "PROD" || names[1] != "STAGING" {
		t.Errorf("unexpected scope names: %v", names)
	}
}

func TestScopeNamesMinimumTwo(t *testing.T) {
	single := []Entry{
		{Key: "PROD_ONLY", Value: "x"},
		{Key: "APP_NAME", Value: "y"},
	}
	names := ScopeNames(single)
	if len(names) != 0 {
		t.Errorf("expected no scope names for single-key prefix, got %v", names)
	}
}
