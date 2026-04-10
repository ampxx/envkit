package encrypt

import (
	"strings"
	"testing"
)

func TestEncryptDecrypt_RoundTrip(t *testing.T) {
	plaintext := "DATABASE_URL=postgres://user:secret@localhost/db"
	passphrase := "supersecret"

	encoded, err := Encrypt(plaintext, passphrase)
	if err != nil {
		t.Fatalf("Encrypt() error: %v", err)
	}

	if encoded == "" {
		t.Fatal("Encrypt() returned empty string")
	}

	decoded, err := Decrypt(encoded, passphrase)
	if err != nil {
		t.Fatalf("Decrypt() error: %v", err)
	}

	if decoded != plaintext {
		t.Errorf("expected %q, got %q", plaintext, decoded)
	}
}

func TestEncrypt_ProducesUniqueOutput(t *testing.T) {
	plaintext := "API_KEY=abc123"
	passphrase := "mypassphrase"

	a, _ := Encrypt(plaintext, passphrase)
	b, _ := Encrypt(plaintext, passphrase)

	// AES-GCM uses a random nonce, so two encryptions must differ.
	if a == b {
		t.Error("expected two encryptions of the same plaintext to differ (random nonce)")
	}
}

func TestDecrypt_WrongPassphrase(t *testing.T) {
	encoded, err := Encrypt("SECRET=value", "correctpass")
	if err != nil {
		t.Fatalf("Encrypt() error: %v", err)
	}

	_, err = Decrypt(encoded, "wrongpass")
	if err == nil {
		t.Fatal("expected error when decrypting with wrong passphrase, got nil")
	}
}

func TestDecrypt_InvalidBase64(t *testing.T) {
	_, err := Decrypt("not-valid-base64!!!", "pass")
	if err == nil {
		t.Fatal("expected error on invalid base64 input, got nil")
	}
}

func TestDecrypt_TruncatedData(t *testing.T) {
	// A valid base64 string that is too short to contain a nonce.
	_, err := Decrypt("YWJj", "pass")
	if err == nil {
		t.Fatal("expected error on truncated ciphertext, got nil")
	}
	if !strings.Contains(err.Error(), "too short") {
		t.Errorf("unexpected error message: %v", err)
	}
}

func TestEncryptDecrypt_EmptyPlaintext(t *testing.T) {
	// Ensure that encrypting an empty string round-trips correctly.
	plaintext := ""
	passphrase := "somepassphrase"

	encoded, err := Encrypt(plaintext, passphrase)
	if err != nil {
		t.Fatalf("Encrypt() error: %v", err)
	}

	decoded, err := Decrypt(encoded, passphrase)
	if err != nil {
		t.Fatalf("Decrypt() error: %v", err)
	}

	if decoded != plaintext {
		t.Errorf("expected empty string, got %q", decoded)
	}
}
