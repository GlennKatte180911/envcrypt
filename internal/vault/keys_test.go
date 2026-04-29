package vault

import (
	"encoding/hex"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestGenerateKeyFile(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "test.key")

	if err := GenerateKeyFile(path); err != nil {
		t.Fatalf("GenerateKeyFile() error = %v", err)
	}

	info, err := os.Stat(path)
	if err != nil {
		t.Fatalf("key file not created: %v", err)
	}
	if info.Mode().Perm() != 0600 {
		t.Errorf("expected permissions 0600, got %v", info.Mode().Perm())
	}
}

func TestLoadKeyFile(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "test.key")

	if err := GenerateKeyFile(path); err != nil {
		t.Fatalf("GenerateKeyFile() error = %v", err)
	}

	key, err := LoadKeyFile(path)
	if err != nil {
		t.Fatalf("LoadKeyFile() error = %v", err)
	}
	if len(key) != 32 {
		t.Errorf("expected 32-byte key, got %d bytes", len(key))
	}
}

func TestLoadKeyFileInvalidHex(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "bad.key")
	os.WriteFile(path, []byte("notvalidhex\n"), 0600)

	_, err := LoadKeyFile(path)
	if err == nil {
		t.Error("expected error for invalid hex, got nil")
	}
}

func TestLoadKeyFileWrongLength(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "short.key")
	short := hex.EncodeToString([]byte("tooshort"))
	os.WriteFile(path, []byte(short+"\n"), 0600)

	_, err := LoadKeyFile(path)
	if err == nil {
		t.Error("expected error for short key, got nil")
	}
}

func TestDefaultKeyPath(t *testing.T) {
	cases := []struct {
		input string
		want  string
	}{
		{".env.vault", ".env.key"},
		{"secrets.vault", "secrets.key"},
		{"/path/to/app.vault", "/path/to/app.key"},
	}
	for _, tc := range cases {
		got := DefaultKeyPath(tc.input)
		if !strings.HasSuffix(got, ".key") {
			t.Errorf("DefaultKeyPath(%q) = %q, want suffix .key", tc.input, got)
		}
		if got != tc.want {
			t.Errorf("DefaultKeyPath(%q) = %q, want %q", tc.input, got, tc.want)
		}
	}
}
