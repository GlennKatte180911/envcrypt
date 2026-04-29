// Package vault also provides key management utilities for generating and
// loading symmetric encryption keys used to lock and unlock vault files.
//
// # Key Files
//
// A key file stores a single 32-byte random key encoded as a lowercase
// hexadecimal string followed by a newline. Key files are written with
// permission 0600 to prevent unintended access by other OS users.
//
// # Typical Workflow
//
//  1. Generate a key file once per project:
//
//		err := vault.GenerateKeyFile(".env.key")
//
//  2. Load the key whenever you need to lock or unlock a vault:
//
//		key, err := vault.LoadKeyFile(".env.key")
//
//  3. Use DefaultKeyPath to derive the conventional key path from a vault path:
//
//		keyPath := vault.DefaultKeyPath(".env.vault") // -> ".env.key"
//
// Key files should be added to .gitignore and shared with teammates through
// a secure out-of-band channel (e.g. a secrets manager or encrypted message).
package vault
