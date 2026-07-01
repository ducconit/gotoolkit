# Random Package

`random` cung cấp các công cụ sinh số và chuỗi ngẫu nhiên tối giản, hiệu năng cao, thread-safe bằng cách tận dụng tối đa thư viện `math/rand/v2` (từ Go 1.22+).

## Các Tính năng Chính

*   **Sinh Chuỗi Ngẫu Nhiên Tối Ưu**: Sử dụng kỹ thuật **Bitmasking** để giảm số lần gọi hàm sinh số ngẫu nhiên, kết hợp `unsafe.String` để đạt đúng **1 allocation duy nhất** cho chuỗi kết quả.
    *   Hỗ trợ đa dạng tập ký tự: `Letters`, `Lowercases`, `Uppercases`, `Numbers`, `UppercaseNumbers`, `LowercaseNumbers`, `Any` (chữ + số + ký tự đặc biệt).
*   **Sinh OTP siêu tốc**:
    *   `OTP(digits)`: Trả về OTP dạng số nguyên (`int`) (**0 allocations**).
    *   `OTPString(digits)`: Trả về OTP dạng chuỗi (**1 allocation**).
*   **Sinh số theo khoảng**: `IntRange(min, max)`, `FloatRange(min, max)`.
*   **Hàm băm tốc độ cao**: `Hash32(s)` tính toán FNV-1a hash cực nhanh (**0 allocations**).

## Hướng dẫn Sử dụng Nhanh

```go
package main

import (
	"fmt"
	"github.com/ducconit/gotoolkit/random"
)

func main() {
	// 1. Sinh chuỗi ngẫu nhiên (chữ + số)
	fmt.Println("Alphanumeric:", random.UppercaseNumbers(16)) // "A1B2C3D4E5F6G7H8"

	// 2. Sinh chuỗi ngẫu nhiên chứa ký tự đặc biệt
	fmt.Println("Special Chars:", random.Any(32))

	// 3. Sinh OTP
	fmt.Println("OTP Number:", random.OTP(6))        // 938472 (int)
	fmt.Println("OTP String:", random.OTPString(6))  // "938472" (string)

	// 4. Sinh số trong khoảng
	fmt.Println("Range Int:", random.IntRange(1, 100))    // Số nguyên 1-99
	fmt.Println("Range Float:", random.FloatRange(0, 1))  // Số thực 0-1
}
```

## Báo cáo Benchmark

Đo thực tế trên CPU **AMD Ryzen 7 8745H**:

| Hàm Benchmark | Thời gian thực thi (ns/op) | Dung lượng bộ nhớ (B/op) | Allocations (allocs/op) |
| :--- | :---: | :---: | :---: |
| `BenchmarkString_Len10` | 49.30 ns | 16 B | 1 allocs |
| `BenchmarkString_Len100` | 326.60 ns | 112 B | 1 allocs |
| `BenchmarkLetters_Len32` | 116.50 ns | 32 B | 1 allocs |
| `BenchmarkNumbers_Len10` | 69.62 ns | 16 B | 1 allocs |
| `BenchmarkOTP` | 6.33 ns | 0 B | 0 allocs |
| `BenchmarkOTPString` | 17.78 ns | 8 B | 1 allocs |
| `BenchmarkIntRange` | 6.25 ns | 0 B | 0 allocs |
| `BenchmarkFloatRange` | 5.49 ns | 0 B | 0 allocs |
| `BenchmarkAny_Len32` | 138.60 ns | 32 B | 1 allocs |
| `BenchmarkHash32` | 24.28 ns | 0 B | 0 allocs |
