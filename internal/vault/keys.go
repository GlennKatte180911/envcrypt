package vault

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

const keyFileExtension = ".key"

// GenerateKeyFile creates a new random 32-byte symmetric key and writes it
// as a hex-encoded string to the specified file path.
func GenerateKeyFile(path string) error {
	raw := make([]byte, 32)
	if _, err := rand.Read(raw); err != nil {
		return fmt.Errorf("generating key: %w", err)
	}

	encoded := hex.EncodeToString(raw) + "\n"
	if err := os.WriteFile(path, []byte(encoded), 0600); err != nil {
		return fmt.Errorf("writing key file: %w", err)
	}
	return nil
}

// LoadKeyFile reads a hex-encoded key from the given file and returns the raw bytes.
func LoadKeyFile(path string) ([]byte, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("reading key file: %w", err)
	}

	hexStr := strings.TrimSpace(string(data))
	key, err := hex.DecodeString(hexStr)
	if err != nil {
		return nil, fmt.Errorf("decoding key: %w", err)
	}

	if len(key) != 32 {
		return nil, errors.New("key must be exactly 32 bytes")
	}
	return key, nil
}

// DefaultKeyPath returns a conventional key file path next to a vault file.
// e.g. ".env.vault" -> ".env.key"
func DefaultKeyPath(vaultPath string) string {
	ext := filepath.Ext(vaultPath)
	base := strings.TrimSuffix(vaultPath, ext)
	return base + keyFileExtension
}
