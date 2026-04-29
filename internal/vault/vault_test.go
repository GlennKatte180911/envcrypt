package vault_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/yourorg/envcrypt/internal/env"
	"github.com/yourorg/envcrypt/internal/vault"
)

func tempVaultPath(t *testing.T) string {
	t.Helper()
	return filepath.Join(t.TempDir(), "test.env.enc")
}

func sampleEntries() []env.Entry {
	return []env.Entry{
		{Key: "DB_HOST", Value: "localhost"},
		{Key: "DB_PORT", Value: "5432"},
		{Key: "SECRET_KEY", Value: "s3cr3t!"},
	}
}

func TestLockAndUnlock(t *testing.T) {
	v := vault.New(tempVaultPath(t))
	entries := sampleEntries()
	passphrase := "correct-horse-battery-staple"

	if err := v.Lock(entries, passphrase); err != nil {
		t.Fatalf("Lock() error = %v", err)
	}

	got, err := v.Unlock(passphrase)
	if err != nil {
		t.Fatalf("Unlock() error = %v", err)
	}

	if len(got) != len(entries) {
		t.Fatalf("entry count: want %d, got %d", len(entries), len(got))
	}
	for i, e := range entries {
		if got[i].Key != e.Key || got[i].Value != e.Value {
			t.Errorf("entry[%d]: want %+v, got %+v", i, e, got[i])
		}
	}
}

func TestUnlockWrongPassphrase(t *testing.T) {
	v := vault.New(tempVaultPath(t))

	if err := v.Lock(sampleEntries(), "right-passphrase"); err != nil {
		t.Fatalf("Lock() error = %v", err)
	}

	_, err := v.Unlock("wrong-passphrase")
	if err == nil {
		t.Fatal("Unlock() expected error with wrong passphrase, got nil")
	}
}

func TestLockCreatesFile(t *testing.T) {
	path := tempVaultPath(t)
	v := vault.New(path)

	if err := v.Lock(sampleEntries(), "passphrase"); err != nil {
		t.Fatalf("Lock() error = %v", err)
	}

	if _, err := os.Stat(path); os.IsNotExist(err) {
		t.Fatalf("Lock() did not create file at %s", path)
	}
}

func TestUnlockMissingFile(t *testing.T) {
	v := vault.New("/nonexistent/path/file.enc")
	_, err := v.Unlock("passphrase")
	if err == nil {
		t.Fatal("Unlock() expected error for missing file, got nil")
	}
}
