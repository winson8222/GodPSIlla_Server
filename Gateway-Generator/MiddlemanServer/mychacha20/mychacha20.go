package mychacha20

import (
	"crypto/rand"
	"golang.org/x/crypto/chacha20"
	"golang.org/x/crypto/chacha20poly1305"
)

func GenerateChaCha20Key() ([]byte, error) {
	key := make([]byte, chacha20.KeySize) // 32 bytes for a ChaCha20 key

	_, err := rand.Read(key)
	if err != nil {
		return nil, err
	}

	return key, nil
}

func GenerateChaCha20Nonce() ([]byte, error) {
	nonce := make([]byte, chacha20poly1305.NonceSizeX)

	_, err := rand.Read(nonce)
	if err != nil {
		return nil, err
	}

	return nonce, nil
}

func Encrypt(key []byte, nonce []byte, plaintexts [][]byte) error {
	// Encrypt each plaintext in the list
	for _, plaintext := range plaintexts {
		// Create a new ChaCha20 cipher
		cipher, err := chacha20.NewUnauthenticatedCipher(key, nonce)
		if err != nil {
			return err
		}
		cipher.XORKeyStream(plaintext, plaintext)
	}

	return nil
}

func Decrypt(key []byte, nonce []byte, ciphertexts [][]byte) error {
	// Decrypt each plaintext in the list
	for _, ciphertext := range ciphertexts {
		// Create a new ChaCha20 cipher
		cipher, err := chacha20.NewUnauthenticatedCipher(key, nonce)
		if err != nil {
			return err
		}
		// Decrypt the ciphertext
		cipher.XORKeyStream(ciphertext, ciphertext)
	}

	return nil
}
