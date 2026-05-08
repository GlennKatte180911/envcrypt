package env

import (
	"testing"
)

func patchBase() []Entry {
	return []Entry{
		{Key: "HOST", Value: "localhost"},
		{Key: "PORT", Value: "8080"},
		{Key: "DEBUG", Value: "false"},
	}
}

func TestPatchSet_UpdateExisting(t *testing.T) {
	out, err := Patch(patchBase(), []PatchOp{{Op: "set", Key: "PORT", Value: "9090"}})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out[1].Value != "9090" {
		t.Errorf("expected 9090, got %s", out[1].Value)
	}
}

func TestPatchSet_AddsNewKey(t *testing.T) {
	out, err := Patch(patchBase(), []PatchOp{{Op: "set", Key: "NEW_KEY", Value: "hello"}})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(out) != 4 {
		t.Errorf("expected 4 entries, got %d", len(out))
	}
	if out[3].Key != "NEW_KEY" || out[3].Value != "hello" {
		t.Errorf("unexpected last entry: %+v", out[3])
	}
}

func TestPatchDelete(t *testing.T) {
	out, err := Patch(patchBase(), []PatchOp{{Op: "delete", Key: "DEBUG"}})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(out) != 2 {
		t.Errorf("expected 2 entries, got %d", len(out))
	}
	for _, e := range out {
		if e.Key == "DEBUG" {
			t.Error("DEBUG should have been deleted")
		}
	}
}

func TestPatchDelete_KeyNotFound(t *testing.T) {
	_, err := Patch(patchBase(), []PatchOp{{Op: "delete", Key: "MISSING"}})
	if err == nil {
		t.Fatal("expected error for missing key")
	}
}

func TestPatchRename(t *testing.T) {
	out, err := Patch(patchBase(), []PatchOp{{Op: "rename", Key: "HOST", NewKey: "HOSTNAME"}})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out[0].Key != "HOSTNAME" {
		t.Errorf("expected HOSTNAME, got %s", out[0].Key)
	}
}

func TestPatchRename_KeyNotFound(t *testing.T) {
	_, err := Patch(patchBase(), []PatchOp{{Op: "rename", Key: "NOPE", NewKey: "X"}})
	if err == nil {
		t.Fatal("expected error for missing key")
	}
}

func TestPatchUnknownOp(t *testing.T) {
	_, err := Patch(patchBase(), []PatchOp{{Op: "flip", Key: "HOST"}})
	if err == nil {
		t.Fatal("expected error for unknown op")
	}
}

func TestPatchDoesNotMutateInput(t *testing.T) {
	base := patchBase()
	origVal := base[1].Value
	_, _ = Patch(base, []PatchOp{{Op: "set", Key: "PORT", Value: "1111"}})
	if base[1].Value != origVal {
		t.Error("Patch mutated the original entries slice")
	}
}

func TestPatchChained(t *testing.T) {
	ops := []PatchOp{
		{Op: "set", Key: "HOST", Value: "example.com"},
		{Op: "delete", Key: "DEBUG"},
		{Op: "rename", Key: "PORT", NewKey: "APP_PORT"},
	}
	out, err := Patch(patchBase(), ops)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(out) != 2 {
		t.Errorf("expected 2 entries, got %d", len(out))
	}
	if out[0].Value != "example.com" {
		t.Errorf("expected example.com, got %s", out[0].Value)
	}
	if out[1].Key != "APP_PORT" {
		t.Errorf("expected APP_PORT, got %s", out[1].Key)
	}
}
