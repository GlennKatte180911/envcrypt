// Package vault manages encrypted .env vaults on disk.
//
// A vault is a single file that stores AES-GCM encrypted environment variable
// data. The plaintext format is the standard KEY=VALUE .env syntax handled by
// the env package; the ciphertext is produced by the crypto package.
//
// Basic usage:
//
//	v := vault.New(".env.enc")
//
//	// Encrypt and write entries to disk.
//	err := v.Lock(entries, passphrase)
//
//	// Read and decrypt entries from disk.
//	entries, err := v.Unlock(passphrase)
//
// File permissions
//
// Vault files are written with mode 0600 so that only the owning user can read
// them. Teams should add *.enc files to version control and share the
// passphrase through a secure out-of-band channel.
package vault
