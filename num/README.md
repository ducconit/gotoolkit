# Num Package

`num` cung cấp các công cụ định dạng và chuyển đổi số thành chữ tiếng Việt/tiếng Anh hiệu năng cao, được tối ưu hóa bộ nhớ cực hạn.

## Các Tính năng Chính

*   **`ToWords(val, lang)`**: Chuyển đổi số nguyên (`int64`) thành chữ viết. Mặc định hỗ trợ Tiếng Việt (`vi`) và Tiếng Anh (`en`). Cho phép đăng ký thêm các bộ dịch cho các ngôn ngữ khác (`Register`).
*   **`ToShorthand(val, lang)`**: Rút gọn số nguyên thành dạng viết tắt dễ đọc.
    *   *Tiếng Anh*: `1.5K`, `2M`, `3B`,...
    *   *Tiếng Việt*: `1.5 nghìn`, `2 triệu`, `3 tỷ`,...
*   **`Format(val, ...)`**: Định dạng số có dấu phẩy phân tách hàng nghìn (ví dụ `1,234,567`) bằng bộ đệm stack tĩnh, đạt hiệu năng cao vượt trội so với các hàm định dạng chuẩn.

## Hướng dẫn Sử dụng Nhanh

```go
package main

import (
	"fmt"
	"github.com/ducconit/gotoolkit/num"
)

func main() {
	// 1. Định dạng phân tách hàng ngàn
	fmt.Println(num.Format(1234567)) // "1,234,567"

	// 2. Rút gọn số (Shorthand)
	fmt.Println(num.ToShorthand(1500, "en"))      // "1.5K"
	fmt.Println(num.ToShorthand(2500000, "vi"))   // "2.5 triệu"

	// 3. Chuyển số thành chữ (Number to Words)
	wordsVI, _ := num.ToWords(123, "vi")
	fmt.Println(wordsVI) // "một trăm hai mươi ba"

	wordsEN, _ := num.ToWords(123, "en")
	fmt.Println(wordsEN) // "one hundred twenty-three"
}
```

## Báo cáo Benchmark

Đo thực tế trên CPU **AMD Ryzen 7 8745H**:

| Hàm Benchmark | Thời gian thực thi (ns/op) | Dung lượng bộ nhớ (B/op) | Allocations (allocs/op) |
| :--- | :---: | :---: | :---: |
| `BenchmarkToShorthand/en` | 47.39 ns | 5 B | 1 allocs |
| `BenchmarkToShorthand/vi` | 52.88 ns | 16 B | 1 allocs |
| `BenchmarkFormat` | 64.76 ns | 16 B | 1 allocs |
| `BenchmarkToWords/en` | 324.40 ns | 384 B | 9 allocs |
| `BenchmarkToWords/vi` | 403.10 ns | 616 B | 9 allocs |
