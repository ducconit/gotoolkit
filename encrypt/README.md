# Encrypt Package

`encrypt` cung cấp các thuật toán mã hóa dữ liệu đối xứng an toàn, hiệu năng cao, tập trung tối đa vào việc giảm allocations bộ nhớ (hỗ trợ các API giải mã in-place).

## Các Tính năng Chính

*   **AES-GCM (256-bit)**: Sử dụng cấu hình chuẩn AES-GCM với Key 32 bytes cực kỳ an toàn.
*   **Nonce tự động tích hợp**: Bản mã hóa sinh ra sẽ tự động chèn 12 bytes Nonce ngẫu nhiên (`crypto/rand`) ở đầu và tự động tách ra lúc giải mã.
*   **Tối ưu Allocations**:
    *   `AESGCM`: Sử dụng bộ nhớ đệm được cấp phát trước dung lượng (capacity) tương ứng với bản mã cuối cùng để loại bỏ malloc heap của thư viện mã hóa chuẩn Go.
    *   `DecryptAESGCMInPlace`: Giải mã trực tiếp trên slice dữ liệu bản mã đầu vào (**Zero-Allocation**), tiết kiệm bộ nhớ khi xử lý các file/dữ liệu payload siêu lớn.

## Hướng dẫn Sử dụng Nhanh

```go
package main

import (
	"fmt"
	"github.com/ducconit/gotoolkit/encrypt"
)

func main() {
	key := []byte("a-very-secret-key-32-bytes-long!") // Key 32 bytes (AES-256)
	plaintext := []byte("Dữ liệu nhạy cảm cần được bảo vệ!")

	// 1. Mã hóa
	ciphertext, err := encrypt.AESGCM(plaintext, key)
	if err != nil {
		panic(err)
	}

	// 2. Giải mã thông thường (an toàn, không đổi ciphertext gốc)
	decrypted, err := encrypt.DecryptAESGCM(ciphertext, key)
	if err == nil {
		fmt.Println("Giải mã:", string(decrypted))
	}

	// 3. Giải mã in-place (Zero-Allocation, ghi đè lên slice ciphertext)
	decryptedInPlace, err := encrypt.DecryptAESGCMInPlace(ciphertext, key)
	if err == nil {
		fmt.Println("Giải mã in-place:", string(decryptedInPlace))
	}
}
```

## Báo cáo Benchmark
*(Chưa thiết lập bộ kiểm thử benchmark tự động cho package này)*
