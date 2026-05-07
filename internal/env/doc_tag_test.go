package env_test

import (
	"fmt"
	"testing"

	"github.com/yourorg/envcrypt/internal/env"
)

// TestTagDocExample exercises the canonical usage shown in doc_tag.go.
func TestTagDocExample(t *testing.T) {
	entries := []env.Entry{
		{Key: "DB_HOST", Value: "localhost"},
		{Key: "DB_PASS", Value: "hunter2"},
		{Key: "API_KEY", Value: "abc123"},
		{Key: "APP_PORT", Value: "8080"},
	}

	// Tag sensitive entries.
	tagged := env.TagEntry(entries, []env.Tag{"secret"}, "DB_PASS", "API_KEY")

	// Filter to secrets only.
	secrets := env.FilterByTag(tagged, "secret")
	if len(secrets) != 2 {
		t.Fatalf("expected 2 secret entries, got %d", len(secrets))
	}

	// Build an index.
	idx := env.BuildTagIndex(tagged)
	if len(idx["secret"]) != 2 {
		t.Fatalf("expected index to list 2 keys for 'secret'")
	}

	// Verify non-secret entries are untouched.
	nonSecrets := env.FilterByTag(tagged, "public")
	if len(nonSecrets) != 0 {
		t.Errorf("unexpected public entries: %v", nonSecrets)
	}

	// Demonstrate idiomatic printing (compile-time check).
	for tag, keys := range idx {
		_ = fmt.Sprintf("%s → %v", tag, keys)
	}
}
