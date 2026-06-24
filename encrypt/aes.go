package encrypt

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"errors"
	"io"
)

// ErrCiphertextTooShort trả về khi bản mã quá ngắn để chứa Nonce (12 bytes).
var ErrCiphertextTooShort = errors.New("ciphertext too short to contain nonce")

// AESGCM mã hóa plaintext sử dụng thuật toán AES-GCM với Key 32 bytes (AES-256).
// Bản mã trả về sẽ tự động chèn 12 bytes Nonce ngẫu nhiên ở đầu.
// Hàm này được tối ưu hóa để tránh cấp phát bộ nhớ lại (re-allocation) trong lúc mã hóa.
func AESGCM(plaintext []byte, key []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	aesgcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	nonceSize := aesgcm.NonceSize()
	capacity := nonceSize + len(plaintext) + aesgcm.Overhead()

	// Tạo slice với length là nonceSize nhưng capacity lớn bằng dung lượng cuối cùng của bản mã.
	// Việc này giúp aesgcm.Seal ghi thẳng vào phần capacity mở rộng mà không bị malloc heap lần nữa.
	buf := make([]byte, nonceSize, capacity)
	if _, err := io.ReadFull(rand.Reader, buf); err != nil {
		return nil, err
	}

	// Seal(dst, nonce, plaintext, additionalData)
	ciphertext := aesgcm.Seal(buf, buf, plaintext, nil)
	return ciphertext, nil
}

// DecryptAESGCM giải mã bản mã được tạo ra bởi hàm AESGCM.
// Bản mã đầu vào bắt buộc phải chứa Nonce (12 bytes) ở đầu.
// Hàm này an toàn và không làm thay đổi (mutate) bản mã đầu vào.
func DecryptAESGCM(ciphertext []byte, key []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	aesgcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	nonceSize := aesgcm.NonceSize()
	if len(ciphertext) < nonceSize {
		return nil, ErrCiphertextTooShort
	}

	nonce := ciphertext[:nonceSize]
	actualCiphertext := ciphertext[nonceSize:]

	// Cấp phát buffer chứa plaintext để tránh mutate ciphertext ban đầu.
	// Dung lượng ước tính là độ dài ciphertext trừ đi overhead của GCM (16 bytes).
	plaintextLen := len(actualCiphertext) - aesgcm.Overhead()
	if plaintextLen < 0 {
		return nil, errors.New("ciphertext is too short or corrupted")
	}
	dst := make([]byte, 0, plaintextLen)

	plaintext, err := aesgcm.Open(dst, nonce, actualCiphertext, nil)
	if err != nil {
		return nil, err
	}

	return plaintext, nil
}

// DecryptAESGCMInPlace giải mã bản mã in-place để tránh cấp phát bộ nhớ cho plaintext kết quả.
// CHÚ Ý: Hàm này sẽ thay đổi (mutate) trực tiếp dữ liệu trên slice ciphertext truyền vào.
// Chỉ sử dụng khi bạn không cần tái sử dụng bản mã gốc sau khi giải mã.
func DecryptAESGCMInPlace(ciphertext []byte, key []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	aesgcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	nonceSize := aesgcm.NonceSize()
	if len(ciphertext) < nonceSize {
		return nil, ErrCiphertextTooShort
	}

	nonce := ciphertext[:nonceSize]
	actualCiphertext := ciphertext[nonceSize:]

	// Dùng actualCiphertext[:0] làm buffer đích (dst) của Open.
	// Do actualCiphertext có thừa dung lượng để chứa plaintext giải mã, 
	// hàm Open sẽ thực thi in-place và đạt Zero-Allocation cho plaintext.
	plaintext, err := aesgcm.Open(actualCiphertext[:0], nonce, actualCiphertext, nil)
	if err != nil {
		return nil, err
	}

	return plaintext, nil
}

