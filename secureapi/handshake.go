package secureapi

import (
	"crypto/ecdh"
	"crypto/rand"
	"crypto/sha256"
	"fmt"
)

// GenerateEphemeralKeyPair sinh một cặp khóa ECDH (Curve P-256) tạm thời cho Server.
func GenerateEphemeralKeyPair() (*ecdh.PrivateKey, []byte, error) {
	privKey, err := ecdh.P256().GenerateKey(rand.Reader)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to generate server ECDH key: %w", err)
	}
	return privKey, privKey.PublicKey().Bytes(), nil
}

// DeriveSessionKey tính toán shared secret từ Private Key của Server và Public Key (raw bytes) của Client,
// sau đó sử dụng SHA-256 để tạo ra Session Key dài 32 bytes (phù hợp cho AES-256).
func DeriveSessionKey(serverPrivKey *ecdh.PrivateKey, clientPubKeyBytes []byte) ([]byte, error) {
	clientPubKey, err := ecdh.P256().NewPublicKey(clientPubKeyBytes)
	if err != nil {
		return nil, fmt.Errorf("invalid client public key: %w", err)
	}

	sharedSecret, err := serverPrivKey.ECDH(clientPubKey)
	if err != nil {
		return nil, fmt.Errorf("failed to compute ECDH shared secret: %w", err)
	}

	// Derive key bằng cách hash shared secret qua sha256.Sum256 (Zero-allocation)
	sum := sha256.Sum256(sharedSecret)
	return sum[:], nil
}
