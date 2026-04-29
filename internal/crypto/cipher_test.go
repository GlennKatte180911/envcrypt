package crypto

import (
	"bytes"
	"testing"
)

func TestDeriveKey(t *testing.T) {
	key := DeriveKey("mysecretpassphrase")
	if len(key) != keySize {
		t.Errorf("expected key size %d, got %d", keySize, len(key))
	}

	// Same passphrase should produce same key
	key2 := DeriveKey("mysecretpassphrase")
	if !bytes.Equal(key, key2) {
		t.Error("expected deterministic key derivation")
	}

	// Different passphrase should produce different key
	key3 := DeriveKey("differentpassphrase")
	if bytes.Equal(key, key3) {
		t.Error("expected different keys for different passphrases")
	}
}

func TestEncryptDecrypt(t *testing.T) {
	key := DeriveKey("testpassphrase")
	plaintext := []byte("DB_HOST=localhost\nDB_PORT=5432\nSECRET=supersecret")

	ciphertext, err := Encrypt(key, plaintext)
	if err != nil {
		t.Fatalf("Encrypt failed: %v", err)
	}

	if bytes.Equal(ciphertext, plaintext) {
		t.Error("ciphertext should not equal plaintext")
	}

	decrypted, err := Decrypt(key, ciphertext)
	if err != nil {
		t.Fatalf("Decrypt failed: %v", err)
	}

	if !bytes.Equal(decrypted, plaintext) {
		t.Errorf("expected %q, got %q", plaintext, decrypted)
	}
}

func TestEncryptProducesUniqueOutput(t *testing.T) {
	key := DeriveKey("testpassphrase")
	plaintext := []byte("SAME_INPUT=value")

	c1, _ := Encrypt(key, plaintext)
	c2, _ := Encrypt(key, plaintext)

	if bytes.Equal(c1, c2) {
		t.Error("expected unique ciphertexts due to random nonce")
	}
}

func TestDecryptWithWrongKey(t *testing.T) {
	key := DeriveKey("correctpassphrase")
	wrongKey := DeriveKey("wrongpassphrase")
	plaintext := []byte("SECRET=value")

	ciphertext, err := Encrypt(key, plaintext)
	if err != nil {
		t.Fatalf("Encrypt failed: %v", err)
	}

	_, err = Decrypt(wrongKey, ciphertext)
	if err == nil {
		t.Error("expected error when decrypting with wrong key")
	}
}

func TestDecryptShortCiphertext(t *testing.T) {
	key := DeriveKey("testpassphrase")
	_, err := Decrypt(key, []byte("short"))
	if err == nil {
		t.Error("expected error for short ciphertext")
	}
}
