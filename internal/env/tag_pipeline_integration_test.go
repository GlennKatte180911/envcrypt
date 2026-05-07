package env_test

import (
	"strings"
	"testing"

	"github.com/yourorg/envcrypt/internal/env"
)

// TestTagThroughPipeline verifies that tagging composes correctly with the
// existing Pipeline, Redact, and Export primitives.
func TestTagThroughPipeline(t *testing.T) {
	base := []env.Entry{
		{Key: "DB_HOST", Value: "db.internal"},
		{Key: "DB_PASS", Value: "s3cr3t"},
		{Key: "APP_PORT", Value: "9000"},
	}

	// Stage 1 – tag secrets.
	tagSecrets := env.PipeFunc(func(in []env.Entry) ([]env.Entry, error) {
		return env.TagEntry(in, []env.Tag{"secret"}, "DB_PASS"), nil
	})

	// Stage 2 – redact secret-tagged entries.
	redactSecrets := env.PipeFunc(func(in []env.Entry) ([]env.Entry, error) {
		secretKeys := make([]string, 0)
		for _, e := range env.FilterByTag(in, "secret") {
			secretKeys = append(secretKeys, e.Key)
		}
		return env.Redact(in, env.WithRedactKeys(secretKeys...)), nil
	})

	p := env.NewPipeline(base)
	out, err := p.Run(tagSecrets, redactSecrets)
	if err != nil {
		t.Fatalf("pipeline error: %v", err)
	}

	for _, e := range out {
		if e.Key == "DB_PASS" {
			if !strings.Contains(e.Value, "*") && e.Value == "s3cr3t" {
				t.Errorf("DB_PASS should have been redacted, got %q", e.Value)
			}
		}
		if e.Key == "APP_PORT" && e.Value != "9000" {
			t.Errorf("APP_PORT should be unchanged, got %q", e.Value)
		}
	}
}

// TestTagIndexAfterMerge ensures BuildTagIndex works correctly on merged
// entry slices where some entries carry pre-existing tags.
func TestTagIndexAfterMerge(t *testing.T) {
	base := env.TagEntry(
		[]env.Entry{{Key: "TOKEN", Value: "abc"}},
		[]env.Tag{"secret"}, "TOKEN",
	)
	override := []env.Entry{
		{Key: "REGION", Value: "us-east-1"},
	}

	merged := env.Merge(base, override)
	idx := env.BuildTagIndex(merged)

	if len(idx["secret"]) != 1 || idx["secret"][0] != "TOKEN" {
		t.Errorf("expected TOKEN in 'secret' index, got %v", idx["secret"])
	}
	if len(idx) != 1 {
		t.Errorf("expected only 'secret' tag in index, got %v", idx)
	}
}
