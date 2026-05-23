package crypto

import (
	"bytes"
	"testing"
)

func TestEncryptDecrypt(t *testing.T) {
	key, err := GenerateKey()
	if err != nil {
		t.Fatalf("GenerateKey failed: %v", err)
	}

	plaintext := []byte("Hello, Legacy DApp! This is a secret message.")

	encrypted, err := Encrypt(plaintext, key)
	if err != nil {
		t.Fatalf("Encrypt failed: %v", err)
	}

	if encrypted.Ciphertext == "" {
		t.Fatal("Ciphertext should not be empty")
	}
	if encrypted.Nonce == "" {
		t.Fatal("Nonce should not be empty")
	}

	decrypted, err := Decrypt(encrypted, key)
	if err != nil {
		t.Fatalf("Decrypt failed: %v", err)
	}

	if !bytes.Equal(plaintext, decrypted) {
		t.Fatalf("Decrypted text doesn't match original.\nExpected: %s\nGot: %s", plaintext, decrypted)
	}
}

func TestEncryptDecryptWithPassword(t *testing.T) {
	password := "my-secure-password-123"
	plaintext := []byte("Confidential inheritance document content.")

	encrypted, err := EncryptWithPassword(plaintext, password)
	if err != nil {
		t.Fatalf("EncryptWithPassword failed: %v", err)
	}

	if encrypted.Salt == "" {
		t.Fatal("Salt should not be empty for password-based encryption")
	}

	decrypted, err := DecryptWithPassword(encrypted, password)
	if err != nil {
		t.Fatalf("DecryptWithPassword failed: %v", err)
	}

	if !bytes.Equal(plaintext, decrypted) {
		t.Fatal("Decrypted text doesn't match original")
	}
}

func TestDecryptWithWrongPassword(t *testing.T) {
	password := "correct-password"
	plaintext := []byte("Secret data")

	encrypted, err := EncryptWithPassword(plaintext, password)
	if err != nil {
		t.Fatalf("EncryptWithPassword failed: %v", err)
	}

	_, err = DecryptWithPassword(encrypted, "wrong-password")
	if err == nil {
		t.Fatal("Expected error when decrypting with wrong password")
	}
}

func TestDecryptWithWrongKey(t *testing.T) {
	key1, _ := GenerateKey()
	key2, _ := GenerateKey()

	plaintext := []byte("Secret data")
	encrypted, err := Encrypt(plaintext, key1)
	if err != nil {
		t.Fatalf("Encrypt failed: %v", err)
	}

	_, err = Decrypt(encrypted, key2)
	if err == nil {
		t.Fatal("Expected error when decrypting with wrong key")
	}
}

func TestGenerateKey(t *testing.T) {
	key1, err := GenerateKey()
	if err != nil {
		t.Fatalf("GenerateKey failed: %v", err)
	}

	key2, err := GenerateKey()
	if err != nil {
		t.Fatalf("GenerateKey failed: %v", err)
	}

	if len(key1) != KeySize {
		t.Fatalf("Key size should be %d, got %d", KeySize, len(key1))
	}

	if bytes.Equal(key1, key2) {
		t.Fatal("Two generated keys should not be equal")
	}
}

func TestDeriveKey(t *testing.T) {
	salt := make([]byte, SaltSize)
	for i := range salt {
		salt[i] = byte(i)
	}

	key1 := DeriveKey("password1", salt)
	key2 := DeriveKey("password2", salt)
	key3 := DeriveKey("password1", salt)

	if !bytes.Equal(key1, key3) {
		t.Fatal("Same password and salt should produce same key")
	}

	if bytes.Equal(key1, key2) {
		t.Fatal("Different passwords should produce different keys")
	}

	if len(key1) != KeySize {
		t.Fatalf("Derived key size should be %d, got %d", KeySize, len(key1))
	}
}

func TestEncryptEmptyData(t *testing.T) {
	key, _ := GenerateKey()
	plaintext := []byte{}

	encrypted, err := Encrypt(plaintext, key)
	if err != nil {
		t.Fatalf("Encrypting empty data should succeed: %v", err)
	}

	decrypted, err := Decrypt(encrypted, key)
	if err != nil {
		t.Fatalf("Decrypting empty data should succeed: %v", err)
	}

	if len(decrypted) != 0 {
		t.Fatal("Decrypted empty data should be empty")
	}
}

func TestEncryptLargeData(t *testing.T) {
	key, _ := GenerateKey()

	// 1MB of data
	plaintext := make([]byte, 1024*1024)
	for i := range plaintext {
		plaintext[i] = byte(i % 256)
	}

	encrypted, err := Encrypt(plaintext, key)
	if err != nil {
		t.Fatalf("Encrypt large data failed: %v", err)
	}

	decrypted, err := Decrypt(encrypted, key)
	if err != nil {
		t.Fatalf("Decrypt large data failed: %v", err)
	}

	if !bytes.Equal(plaintext, decrypted) {
		t.Fatal("Large data roundtrip failed")
	}
}
