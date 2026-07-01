# SecureAPI Package

`secureapi` cung cấp giải pháp bảo mật API toàn diện cho các ứng dụng web và mobile bằng cách tự động mã hóa đầu cuối (End-to-End Encryption - E2EE) cho payload HTTP.

## Các Tính năng Chính

*   **ECDH Key Exchange (Curve P-256)**: Thực hiện bắt tay (handshake) bảo mật để tạo khóa đối xứng chung giữa Client và Server mà không truyền khóa qua mạng.
*   **Mã Hóa Payload AES-GCM (256-bit)**: Mã hóa toàn bộ dữ liệu request/response bằng thuật toán mã hóa đối xứng mạnh nhất, ngăn chặn hoàn toàn tấn công Man-in-the-Middle (MITM).
*   **Session Management**: Quản lý các session key tạm thời thông qua `SessionStore` an toàn với cơ chế tự động dọn dẹp (garbage collection) ngầm.
*   **Tích hợp Gin Middleware**: Thư mục con `gin` chứa middleware tích hợp sẵn giúp tự động hóa quá trình giải mã Request Body và mã hóa Response Body trên Gin framework.

## Hướng dẫn Sử dụng Nhanh

### 1. Quá trình Bắt tay ECDH (Server Side)
```go
package main

import (
	"fmt"
	"github.com/ducconit/gotoolkit/secureapi"
)

func main() {
	// Bước 1: Server sinh cặp khóa tạm thời (ephemeral key)
	privKey, pubKeyBytes, _ := secureapi.GenerateEphemeralKeyPair()

	// Bước 2: Nhận Public Key từ Client gửi lên
	clientPubKeyBytes := []byte{...} // Public Key từ Client

	// Bước 3: Tính toán khóa đối xứng chung (Session Key - 32 bytes)
	sessionKey, _ := secureapi.DeriveSessionKey(privKey, clientPubKeyBytes)
	fmt.Println("Session Key derived successfully!")
}
```

### 2. Sử dụng SessionStore quản lý khóa
```go
package main

import (
	"time"
	"github.com/ducconit/gotoolkit/secureapi"
)

func main() {
	// Tạo Store quản lý session với TTL là 1 giờ
	store := secureapi.NewSessionStore(1 * time.Hour)
	defer store.Close()

	// Lưu session key
	sessionID, _ := store.CreateSession(derivedSessionKey)

	// Lấy session key bằng Session ID nhận từ Header request
	key, err := store.GetSession(sessionID)
}
```

## Báo cáo Benchmark

Đo thực tế trên CPU **AMD Ryzen 7 8745H**:

| Hàm Benchmark | Kích thước Payload | Thời gian thực thi (ns/op) | Dung lượng bộ nhớ (B/op) | Allocations (allocs/op) |
| :--- | :---: | :---: | :---: | :---: |
| `Benchmark_ECDHHandshake` | - | 54,263 ns | 1,104 B | 16 allocs |
| `Benchmark_ECDH_Derive` | - | 42,434 ns | 528 B | 8 allocs |
| `Benchmark_EncryptAESGCM_1KB` | 1 KB | 793 ns | 2,432 B | 3 allocs |
| `Benchmark_EncryptAESGCM_10KB` | 10 KB | 3,569 ns | 12,160 B | 3 allocs |
| `Benchmark_EncryptAESGCM_100KB` | 100 KB | 31,210 ns | 107,776 B | 3 allocs |
| `Benchmark_DecryptAESGCM_1KB` | 1 KB | 736 ns | 2,304 B | 3 allocs |
| `Benchmark_DecryptAESGCM_10KB` | 10 KB | 3,689 ns | 11,520 B | 3 allocs |
| `Benchmark_DecryptAESGCM_100KB` | 100 KB | 33,827 ns | 107,776 B | 3 allocs |
