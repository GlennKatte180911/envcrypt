package env

import (
	"testing"
)

var tagEntries = []Entry{
	{Key: "DB_HOST", Value: "localhost"},
	{Key: "DB_PASS", Value: "secret"},
	{Key: "APP_PORT", Value: "8080"},
	{Key: "APP_DEBUG", Value: "true"},
}

func TestTagEntryAddsTag(t *testing.T) {
	out := TagEntry(tagEntries, []Tag{"db"}, "DB_HOST", "DB_PASS")
	for _, e := range out {
		if e.Key == "DB_HOST" || e.Key == "DB_PASS" {
			tags := TagsOf(e)
			if len(tags) != 1 || tags[0] != "db" {
				t.Errorf("%s: expected tag 'db', got %v", e.Key, tags)
			}
		} else {
			if len(TagsOf(e)) != 0 {
				t.Errorf("%s: expected no tags", e.Key)
			}
		}
	}
}

func TestTagEntryDoesNotMutateOriginal(t *testing.T) {
	TagEntry(tagEntries, []Tag{"sensitive"}, "DB_PASS")
	if len(TagsOf(tagEntries[1])) != 0 {
		t.Error("original entry should not be mutated")
	}
}

func TestTagsOfNoTags(t *testing.T) {
	e := Entry{Key: "FOO", Value: "bar", Comment: "just a comment"}
	if tags := TagsOf(e); len(tags) != 0 {
		t.Errorf("expected no tags, got %v", tags)
	}
}

func TestTagsOfMultiple(t *testing.T) {
	e := Entry{Key: "X", Value: "1", Comment: "note #tag:alpha,beta"}
	tags := TagsOf(e)
	if len(tags) != 2 || tags[0] != "alpha" || tags[1] != "beta" {
		t.Errorf("unexpected tags: %v", tags)
	}
}

func TestBuildTagIndex(t *testing.T) {
	entries := TagEntry(tagEntries, []Tag{"db"}, "DB_HOST", "DB_PASS")
	entries = TagEntry(entries, []Tag{"app"}, "APP_PORT", "APP_DEBUG")
	idx := BuildTagIndex(entries)

	if len(idx["db"]) != 2 {
		t.Errorf("expected 2 keys for tag 'db', got %d", len(idx["db"]))
	}
	if len(idx["app"]) != 2 {
		t.Errorf("expected 2 keys for tag 'app', got %d", len(idx["app"]))
	}
}

func TestFilterByTag(t *testing.T) {
	entries := TagEntry(tagEntries, []Tag{"secret"}, "DB_PASS")
	result := FilterByTag(entries, "secret")
	if len(result) != 1 || result[0].Key != "DB_PASS" {
		t.Errorf("expected DB_PASS, got %v", result)
	}
}

func TestFilterByTagNoMatch(t *testing.T) {
	result := FilterByTag(tagEntries, "nonexistent")
	if len(result) != 0 {
		t.Errorf("expected empty result, got %v", result)
	}
}

func TestAppendTagsIdempotent(t *testing.T) {
	e := Entry{Key: "K", Value: "v"}
	out1 := TagEntry([]Entry{e}, []Tag{"x"}, "K")
	out2 := TagEntry(out1, []Tag{"x"}, "K")
	tags1 := TagsOf(out1[0])
	tags2 := TagsOf(out2[0])
	if len(tags1) != 1 || len(tags2) != 1 {
		t.Errorf("duplicate tag added: %v vs %v", tags1, tags2)
	}
}
