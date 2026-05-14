package env_test

import (
	"fmt"
	"testing"

	"github.com/yourorg/envcrypt/internal/env"
)

func TestScopeDocExample(t *testing.T) {
	entries := []env.Entry{
		{Key: "PROD_API_URL", Value: "https://api.prod.example.com"},
		{Key: "PROD_API_KEY", Value: "secret-prod"},
		{Key: "STAGING_API_URL", Value: "https://api.staging.example.com"},
		{Key: "STAGING_API_KEY", Value: "secret-staging"},
	}

	prod := env.NewScope("prod", entries)

	url, ok := prod.Lookup("API_URL")
	if !ok {
		t.Fatal("expected API_URL in prod scope")
	}
	fmt.Println(url) // https://api.prod.example.com

	promoted := prod.Promote()
	if len(promoted) != 2 {
		t.Fatalf("expected 2 promoted entries, got %d", len(promoted))
	}
	for _, e := range promoted {
		if len(e.Key) < 5 || e.Key[:5] != "PROD_" {
			t.Errorf("promoted key missing prefix: %q", e.Key)
		}
	}

	names := env.ScopeNames(entries)
	if len(names) != 2 {
		t.Errorf("expected 2 scope names, got %d", len(names))
	}
}
