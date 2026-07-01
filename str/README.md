# Str Package

`str` cung cấp các tiện ích xử lý chuỗi và kiểm tra tính hợp lệ (validation) tối giản, hiệu năng cao, được tối ưu hóa cực hạn để loại bỏ việc cấp phát bộ nhớ (allocation) không cần thiết trên heap.

## Các Tính năng Chính

*   **`CleanSpace(s)`**: Làm sạch khoảng trắng thừa trong chuỗi. Đạt **0 allocation** nếu chuỗi đầu vào đã sạch.
*   **`RemoveAccents(s)`**: Loại bỏ toàn bộ dấu tiếng Việt (bao gồm cả ký tự dựng sẵn Unicode và tổ hợp). Đạt **0 allocation** đối với chuỗi không dấu.
*   **`Slugify(s, sep)`**: Chuyển đổi chuỗi thành slug thân thiện với URL (tự động loại bỏ dấu tiếng Việt, ký tự đặc biệt, làm sạch khoảng trắng và lowercasing). Thực thi trong đúng **1 allocation duy nhất**.
*   **`IsTrue(s)` & `IsFalse(s)`**: Đánh giá giá trị truthy/falsy của chuỗi tương tự như các ngôn ngữ động (hỗ trợ `true`, `1`, `t`, `false`, `0`, `f`, `no`, `off`). Đạt **0 allocation** nhờ tận dụng slicing và trim khoảng trắng trực tiếp.
*   **Các hàm kiểm tra hợp lệ**: `IsAlphaSpace`, `IsUsername`, `ContainsHTML` tối ưu hóa CPU.

## Hướng dẫn Sử dụng Nhanh

```go
package main

import (
	"fmt"
	"github.com/ducconit/gotoolkit/str"
)

func main() {
	// 1. Loại bỏ khoảng trắng thừa
	fmt.Println(str.CleanSpace("   Hello   World   ")) // "Hello World"

	// 2. Loại bỏ dấu tiếng Việt
	fmt.Println(str.RemoveAccents("Xin chào Việt Nam!")) // "Xin chao Viet Nam!"

	// 3. Tạo URL Slug
	fmt.Println(str.Slugify("Học lập trình Golang 1.26!", "-")) // "hoc-lap-trinh-golang-126"

	// 4. Kiểm tra Truthy/Falsy
	fmt.Println(str.IsTrue("  true  ")) // true
	fmt.Println(str.IsFalse("0"))      // true
}
```

## Báo cáo Benchmark

Đo thực tế trên CPU **AMD Ryzen 7 8745H**:

| Hàm Benchmark | Thời gian thực thi (ns/op) | Dung lượng bộ nhớ (B/op) | Allocations (allocs/op) |
| :--- | :---: | :---: | :---: |
| `BenchmarkIsAlphaSpace_English` | 10.56 ns | 0 B | 0 allocs |
| `BenchmarkIsAlphaSpace_Vietnamese` | 31.72 ns | 0 B | 0 allocs |
| `BenchmarkIsUsername` | 7.09 ns | 0 B | 0 allocs |
| `BenchmarkContainsHTML_Safe` | 13.22 ns | 0 B | 0 allocs |
| `BenchmarkContainsHTML_Unsafe` | 6.07 ns | 0 B | 0 allocs |
| `BenchmarkCleanSpace_Clean` | 10.79 ns | 0 B | 0 allocs |
| `BenchmarkCleanSpace_Dirty` | 48.55 ns | 24 B | 1 allocs |
| `BenchmarkRemoveAccents_NoAccents` | 25.84 ns | 0 B | 0 allocs |
| `BenchmarkRemoveAccents_WithAccents` | 274.90 ns | 80 B | 1 allocs |
| `BenchmarkSlugify_Default` | 287.20 ns | 80 B | 1 allocs |
| `BenchmarkIsFalse_True` | 4.80 ns | 0 B | 0 allocs |
| `BenchmarkIsFalse_False` | 2.53 ns | 0 B | 0 allocs |
