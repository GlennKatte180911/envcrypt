// Package crypto provides symmetric encryption and decryption utilities
// for envcrypt. It uses AES-256-GCM for authenticated encryption, ensuring
// both confidentiality and integrity of encrypted .env files.
//
// Key derivation is performed via SHA-256 hashing of a passphrase, producing
// a 32-byte key suitable for AES-256. Each encryption operation generates a
// random 12-byte nonce, which is prepended to the resulting ciphertext.
//
// Basic usage:
//
//	key := crypto.DeriveKey("my-team-passphrase")
//
//	ciphertext, err := crypto.Encrypt(key, []byte("DB_PASSWORD=secret"))
//	if err != nil {
//		log.Fatal(err)
//	}
//
//	plaintext, err := crypto.Decrypt(key, ciphertext)
//	if err != nil {
//		log.Fatal(err)
//	}
package crypto
