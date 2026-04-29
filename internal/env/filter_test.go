package env

import (
	"testing"
)

func sampleEntries() []Entry {
	return []Entry{
		{Key: "APP_HOST", Value: "localhost"},
		{Key: "APP_PORT", Value: "8080"},
		{Key: "DB_HOST", Value: "db.local"},
		{Key: "DB_PORT", Value: "5432"},
		{Key: "SECRET_KEY", Value: "abc123"},
	}
}

func TestWithPrefix(t *testing.T) {
	entries := sampleEntries()
	got := WithPrefix(entries, "APP_")
	if len(got) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(got))
	}
	for _, e := range got {
		if e.Key != "APP_HOST" && e.Key != "APP_PORT" {
			t.Errorf("unexpected key %q", e.Key)
		}
	}
}

func TestWithoutPrefix(t *testing.T) {
	entries := sampleEntries()
	got := WithoutPrefix(entries, "DB_")
	if len(got) != 3 {
		t.Fatalf("expected 3 entries, got %d", len(got))
	}
	for _, e := range got {
		if e.Key == "DB_HOST" || e.Key == "DB_PORT" {
			t.Errorf("unexpected DB_ key %q survived filter", e.Key)
		}
	}
}

func TestStripPrefix(t *testing.T) {
	entries := WithPrefix(sampleEntries(), "APP_")
	stripped := StripPrefix(entries, "APP_")
	expected := map[string]string{"HOST": "localhost", "PORT": "8080"}
	for _, e := range stripped {
		if v, ok := expected[e.Key]; !ok {
			t.Errorf("unexpected key %q after strip", e.Key)
		} else if e.Value != v {
			t.Errorf("key %q: expected value %q, got %q", e.Key, v, e.Value)
		}
	}
}

func TestUnique(t *testing.T) {
	entries := []Entry{
		{Key: "FOO", Value: "first"},
		{Key: "BAR", Value: "bar"},
		{Key: "FOO", Value: "second"},
	}
	got := Unique(entries)
	if len(got) != 2 {
		t.Fatalf("expected 2 unique entries, got %d", len(got))
	}
	for _, e := range got {
		if e.Key == "FOO" && e.Value != "second" {
			t.Errorf("expected last value 'second', got %q", e.Value)
		}
	}
}

func TestFilter(t *testing.T) {
	entries := sampleEntries()
	got := Filter(entries, func(key string) bool {
		return key == "SECRET_KEY"
	})
	if len(got) != 1 || got[0].Key != "SECRET_KEY" {
		t.Errorf("expected single SECRET_KEY entry, got %+v", got)
	}
}
