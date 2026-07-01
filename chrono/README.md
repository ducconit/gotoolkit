# Chrono Package

`chrono` cung cấp các công cụ xử lý thời gian tối giản, hiệu năng cao. Package này giải quyết điểm yếu lớn nhất của Go trong việc định dạng thời gian bằng cách tự động dịch các định dạng quen thuộc của PHP (`Y-m-d H:i:s`) và JS/Moment (`YYYY-MM-DD HH:mm:ss`) sang Go layout tương ứng, đồng thời hỗ trợ bản địa hóa đa ngôn ngữ (i18n) tiếng Việt, Nhật, Hàn, Trung, Anh.

## Các Tính năng Chính

*   **Bộ dịch Layout PHP & JS sang Go**: Không cần nhớ layout `2006-01-02 15:04:05` phiền phức của Go.
*   **Định dạng Đa Ngôn Ngữ (i18n)**: Định dạng và tự động dịch tên thứ, tháng sang Tiếng Việt (`vi`), Tiếng Nhật (`ja`), Tiếng Hàn (`ko`), Tiếng Trung (`zh`), Tiếng Anh (`en`).
*   **Đồng bộ Hóa Thư Viện Tiêu Chuẩn**: Tương thích hoàn hảo với `language.Tag` của package `golang.org/x/text/language` bằng giải pháp **zero-dependency**.
*   **Parse Timestamp Thông Minh**: Hỗ trợ parse timestamp tường minh theo đơn vị, hoặc tự động đoán nhận đơn vị (giây, mili, micro, nano) cực kỳ an toàn.
*   **Tối Ưu Hiệu Năng Vượt Trội**: Zero-regex, sử dụng scan byte thủ công và `strings.Builder` để đạt tốc độ tối đa và tối thiểu allocations.

## Hướng dẫn Sử dụng Nhanh

### 1. Định dạng kiểu PHP và JS thường dùng
```go
package main

import (
	"fmt"
	"time"
	"github.com/ducconit/gotoolkit/chrono"
)

func main() {
	now := time.Now()

	// Định dạng PHP
	fmt.Println(chrono.FormatPHP(now, "Y-m-d H:i:s")) // 2026-07-01 12:40:18

	// Định dạng JS/Moment
	fmt.Println(chrono.FormatJS(now, "YYYY-MM-DD HH:mm:ss.SSS")) // 2026-07-01 12:40:18.000
}
```

### 2. Định dạng Đa Ngôn Ngữ (i18n)
```go
package main

import (
	"fmt"
	"time"
	"github.com/ducconit/gotoolkit/chrono"
)

func main() {
	now := time.Now()

	// Định dạng tiếng Việt
	fmt.Println(chrono.FormatLocalePHP(now, "l, d F Y", "vi")) 
	// Output: Thứ Tư, 01 Tháng Bảy 2026

	// Định dạng tiếng Nhật
	fmt.Println(chrono.FormatLocaleJS(now, "dddd, DD MMMM YYYY", "ja"))
	// Output: 水曜日, 01 7月 2026
}
```

### 3. Parse Timestamp Tự động
```go
package main

import (
	"fmt"
	"github.com/ducconit/gotoolkit/chrono"
)

func main() {
	// Hệ thống tự động đoán nhận timestamp là milliseconds và parse
	t := chrono.AutoParseTimestamp(1782812818000)
	fmt.Println(t.UTC().Format("2006-01-02 15:04:05")) // 2026-06-30 09:46:58
}
```

## Báo cáo Benchmark

Đo thực tế trên CPU **AMD Ryzen 7 8745H**:

| Hàm Benchmark | Thời gian thực thi (ns/op) | Dung lượng bộ nhớ (B/op) | Allocations (allocs/op) |
| :--- | :---: | :---: | :---: |
| `BenchmarkPHP_Conversion` | 62.43 ns | 48 B | 1 allocs |
| `BenchmarkJS_Conversion` | 92.85 ns | 64 B | 1 allocs |
| `BenchmarkFormatPHP` | 129.0 ns | 48 B | 2 allocs |
| `BenchmarkFormatJS` | 158.2 ns | 72 B | 2 allocs |
| `BenchmarkAutoParseTimestamp` | 1.47 ns | 0 B | 0 allocs |
| `BenchmarkFormatLocalePHP` | 206.5 ns | 86 B | 6 allocs |
