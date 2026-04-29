// Package vault provides functionality for reading and writing encrypted
// .env files ("vaults") to disk using the crypto and env packages.
package vault

import (
	"fmt"
	"os"

	"github.com/yourorg/envcrypt/internal/crypto"
	"github.com/yourorg/envcrypt/internal/env"
)

const defaultFileMode = 0600

// Vault represents an encrypted .env file on disk.
type Vault struct {
	Path string
}

// New returns a Vault backed by the file at path.
func New(path string) *Vault {
	return &Vault{Path: path}
}

// Lock encrypts the given entries with passphrase and writes them to the
// vault file, creating or overwriting it.
func (v *Vault) Lock(entries []env.Entry, passphrase string) error {
	plaintext := []byte(env.Serialize(entries))

	key, err := crypto.DeriveKey(passphrase, nil)
	if err != nil {
		return fmt.Errorf("vault lock: derive key: %w", err)
	}

	ciphertext, err := crypto.Encrypt(key, plaintext)
	if err != nil {
		return fmt.Errorf("vault lock: encrypt: %w", err)
	}

	if err := os.WriteFile(v.Path, ciphertext, defaultFileMode); err != nil {
		return fmt.Errorf("vault lock: write file: %w", err)
	}

	return nil
}

// Unlock reads the vault file, decrypts it with passphrase, and returns the
// parsed entries.
func (v *Vault) Unlock(passphrase string) ([]env.Entry, error) {
	ciphertext, err := os.ReadFile(v.Path)
	if err != nil {
		return nil, fmt.Errorf("vault unlock: read file: %w", err)
	}

	key, err := crypto.DeriveKey(passphrase, nil)
	if err != nil {
		return nil, fmt.Errorf("vault unlock: derive key: %w", err)
	}

	plaintext, err := crypto.Decrypt(key, ciphertext)
	if err != nil {
		return nil, fmt.Errorf("vault unlock: decrypt: %w", err)
	}

	entries, err := env.Parse(string(plaintext))
	if err != nil {
		return nil, fmt.Errorf("vault unlock: parse: %w", err)
	}

	return entries, nil
}
