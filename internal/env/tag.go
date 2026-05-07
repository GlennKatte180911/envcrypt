package env

import (
	"fmt"
	"sort"
	"strings"
)

// Tag represents a string label attached to an entry.
type Tag = string

// TagIndex maps each tag to the set of keys that carry it.
type TagIndex map[Tag][]string

// TagEntry attaches one or more tags to entries whose keys match the
// provided key set. Entries not in keys are left unchanged.
func TagEntry(entries []Entry, tags []Tag, keys ...string) []Entry {
	keySet := make(map[string]struct{}, len(keys))
	for _, k := range keys {
		keySet[k] = struct{}{}
	}

	out := make([]Entry, len(entries))
	for i, e := range entries {
		if _, ok := keySet[e.Key]; ok {
			tagged := appendTags(e.Comment, tags)
			out[i] = Entry{Key: e.Key, Value: e.Value, Comment: tagged}
		} else {
			out[i] = e
		}
	}
	return out
}

// TagsOf returns the tags embedded in an entry's comment field.
// Tags are stored as "#tag:foo,bar" within the comment string.
func TagsOf(e Entry) []Tag {
	const prefix = "#tag:"
	for _, part := range strings.Split(e.Comment, " ") {
		if strings.HasPrefix(part, prefix) {
			raw := strings.TrimPrefix(part, prefix)
			if raw == "" {
				return nil
			}
			return strings.Split(raw, ",")
		}
	}
	return nil
}

// BuildTagIndex constructs a TagIndex from a slice of entries.
func BuildTagIndex(entries []Entry) TagIndex {
	idx := make(TagIndex)
	for _, e := range entries {
		for _, t := range TagsOf(e) {
			idx[t] = append(idx[t], e.Key)
		}
	}
	return idx
}

// FilterByTag returns only entries that carry the given tag.
func FilterByTag(entries []Entry, tag Tag) []Entry {
	var out []Entry
	for _, e := range entries {
		for _, t := range TagsOf(e) {
			if t == tag {
				out = append(out, e)
				break
			}
		}
	}
	return out
}

// appendTags merges new tags into an existing comment string.
func appendTags(comment string, tags []Tag) string {
	const prefix = "#tag:"
	existing := []string{}
	other := []string{}

	for _, part := range strings.Fields(comment) {
		if strings.HasPrefix(part, prefix) {
			existing = append(existing, strings.Split(strings.TrimPrefix(part, prefix), ",")...)
		} else {
			other = append(other, part)
		}
	}

	seen := make(map[string]struct{})
	for _, t := range existing {
		seen[t] = struct{}{}
	}
	for _, t := range tags {
		if _, ok := seen[t]; !ok {
			existing = append(existing, t)
			seen[t] = struct{}{}
		}
	}
	sort.Strings(existing)

	tagPart := fmt.Sprintf("%s%s", prefix, strings.Join(existing, ","))
	parts := append(other, tagPart)
	return strings.Join(parts, " ")
}
