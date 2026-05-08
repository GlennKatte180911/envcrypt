package env

import (
	"testing"
)

func TestCloneIndependence(t *testing.T) {
	orig := []Entry{{Key: "A", Value: "1"}, {Key: "B", Value: "2"}}
	cloned := Clone(orig)
	if len(cloned) != len(orig) {
		t.Fatalf("expected len %d, got %d", len(orig), len(cloned))
	}
	cloned[0] = Entry{Key: "X", Value: "99"}
	if orig[0].Key != "A" {
		t.Errorf("Clone mutated original")
	}
}

func TestCloneNil(t *testing.T) {
	if Clone(nil) != nil {
		t.Error("expected nil for nil input")
	}
}

func TestCloneMapIndependence(t *testing.T) {
	orig := map[string]string{"K": "V"}
	cloned := CloneMap(orig)
	cloned["K"] = "changed"
	if orig["K"] != "V" {
		t.Error("CloneMap mutated original")
	}
}

func TestCloneMapNil(t *testing.T) {
	if CloneMap(nil) != nil {
		t.Error("expected nil for nil input")
	}
}

func TestSubsetKeepsMatchingKeys(t *testing.T) {
	entries := []Entry{{Key: "A"}, {Key: "B"}, {Key: "C"}}
	result := Subset(entries, []string{"A", "C"})
	if len(result) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(result))
	}
	if result[0].Key != "A" || result[1].Key != "C" {
		t.Errorf("unexpected keys: %v", result)
	}
}

func TestSubsetEmptyKeys(t *testing.T) {
	entries := []Entry{{Key: "A"}, {Key: "B"}}
	result := Subset(entries, nil)
	if len(result) != 0 {
		t.Errorf("expected empty result, got %d entries", len(result))
	}
}

func TestExcludeRemovesMatchingKeys(t *testing.T) {
	entries := []Entry{{Key: "A"}, {Key: "B"}, {Key: "C"}}
	result := Exclude(entries, []string{"B"})
	if len(result) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(result))
	}
	for _, e := range result {
		if e.Key == "B" {
			t.Error("excluded key B still present")
		}
	}
}

func TestExcludeNoMatch(t *testing.T) {
	entries := []Entry{{Key: "A"}, {Key: "B"}}
	result := Exclude(entries, []string{"Z"})
	if len(result) != 2 {
		t.Errorf("expected 2 entries unchanged, got %d", len(result))
	}
}
