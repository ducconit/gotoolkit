package secureapi_test

import (
	"crypto/rand"
	"testing"

	"github.com/ducconit/gotoolkit/encrypt"
	"github.com/ducconit/gotoolkit/secureapi"
)

// Benchmark đo toàn bộ luồng xử lý Handshake trên Server (Sinh khóa Server + Tính Shared Secret)
func Benchmark_ECDHHandshake(b *testing.B) {
	_, clientPubKeyBytes, _ := secureapi.GenerateEphemeralKeyPair()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		// 1. Sinh cặp khóa của Server
		serverPrivKey, _, err := secureapi.GenerateEphemeralKeyPair()
		if err != nil {
			b.Fatal(err)
		}

		// 2. Server tính Session Key dựa trên Client Public Key
		_, err = secureapi.DeriveSessionKey(serverPrivKey, clientPubKeyBytes)
		if err != nil {
			b.Fatal(err)
		}
	}
}

// Benchmark đo tốc độ tính toán Shared Secret của thuật toán ECDH P-256
func Benchmark_ECDH_Derive(b *testing.B) {
	serverPrivKey, _, _ := secureapi.GenerateEphemeralKeyPair()
	_, clientPubKeyBytes, _ := secureapi.GenerateEphemeralKeyPair()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := secureapi.DeriveSessionKey(serverPrivKey, clientPubKeyBytes)
		if err != nil {
			b.Fatal(err)
		}
	}
}

// Benchmark mã hóa AES-GCM với các kích thước payload khác nhau
func Benchmark_EncryptAESGCM_1KB(b *testing.B) {
	key := make([]byte, 32)
	_, _ = rand.Read(key)
	payload := make([]byte, 1024) // 1 KB
	_, _ = rand.Read(payload)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := encrypt.AESGCM(payload, key)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func Benchmark_EncryptAESGCM_10KB(b *testing.B) {
	key := make([]byte, 32)
	_, _ = rand.Read(key)
	payload := make([]byte, 10*1024) // 10 KB
	_, _ = rand.Read(payload)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := encrypt.AESGCM(payload, key)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func Benchmark_EncryptAESGCM_100KB(b *testing.B) {
	key := make([]byte, 32)
	_, _ = rand.Read(key)
	payload := make([]byte, 100*1024) // 100 KB
	_, _ = rand.Read(payload)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := encrypt.AESGCM(payload, key)
		if err != nil {
			b.Fatal(err)
		}
	}
}

// Benchmark giải mã AES-GCM với các kích thước payload khác nhau.
// Do hàm DecryptAESGCM sử dụng giải mã in-place (mutates/phá hủy bản mã), 
// chúng ta bắt buộc phải copy ciphertext sang buffer tạm trước mỗi lần giải mã trong vòng lặp test.
func Benchmark_DecryptAESGCM_1KB(b *testing.B) {
	key := make([]byte, 32)
	_, _ = rand.Read(key)
	payload := make([]byte, 1024) // 1 KB
	_, _ = rand.Read(payload)
	ciphertext, _ := encrypt.AESGCM(payload, key)
	
	tempCiphertext := make([]byte, len(ciphertext))

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		copy(tempCiphertext, ciphertext)
		_, err := encrypt.DecryptAESGCM(tempCiphertext, key)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func Benchmark_DecryptAESGCM_10KB(b *testing.B) {
	key := make([]byte, 32)
	_, _ = rand.Read(key)
	payload := make([]byte, 10*1024) // 10 KB
	_, _ = rand.Read(payload)
	ciphertext, _ := encrypt.AESGCM(payload, key)
	
	tempCiphertext := make([]byte, len(ciphertext))

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		copy(tempCiphertext, ciphertext)
		_, err := encrypt.DecryptAESGCM(tempCiphertext, key)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func Benchmark_DecryptAESGCM_100KB(b *testing.B) {
	key := make([]byte, 32)
	_, _ = rand.Read(key)
	payload := make([]byte, 100*1024) // 100 KB
	_, _ = rand.Read(payload)
	ciphertext, _ := encrypt.AESGCM(payload, key)
	
	tempCiphertext := make([]byte, len(ciphertext))

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		copy(tempCiphertext, ciphertext)
		_, err := encrypt.DecryptAESGCM(tempCiphertext, key)
		if err != nil {
			b.Fatal(err)
		}
	}
}
