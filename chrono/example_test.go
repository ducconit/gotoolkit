package chrono_test

import (
	"fmt"
	"github.com/ducconit/gotoolkit/chrono"
	"time"
)

func Example() {
	// Khởi tạo thời gian mốc cố định: 2026-07-01 12:40:18 UTC
	loc := time.UTC
	testTime := time.Date(2026, time.July, 1, 12, 40, 18, 0, loc)

	// 1. Định dạng theo kiểu PHP (Y-m-d H:i:s)
	formattedPHP := chrono.FormatPHP(testTime, "Y-m-d H:i:s")
	fmt.Println("PHP formatted:", formattedPHP)

	// 2. Định dạng theo kiểu JS/Moment (YYYY-MM-DD HH:mm:ss)
	formattedJS := chrono.FormatJS(testTime, "YYYY-MM-DD HH:mm:ss")
	fmt.Println("JS formatted:", formattedJS)

	// 3. Parse ngược lại từ chuỗi (PHP style)
	parsedTime, err := chrono.ParseInLocationPHP("Y-m-d H:i:s", "2026-07-01 12:40:18", loc)
	if err == nil {
		fmt.Println("PHP parsed:", parsedTime.Format("2006-01-02 15:04:05"))
	}

	// 4. Chuyển đổi thô từ PHP format sang Go layout để dùng trực tiếp với time.Format
	goLayout := chrono.PHP("Y-m-d H:i:s.v O")
	fmt.Println("Go Layout from PHP:", goLayout)

	// 5. Parse timestamp tường minh (giây và mili giây)
	t1 := chrono.ParseTimestamp(1782812818, time.Second)
	t2 := chrono.ParseTimestamp(1782812818000, time.Millisecond)
	fmt.Println("Parse Seconds:", t1.UTC().Format("2006-01-02 15:04:05"))
	fmt.Println("Parse Milliseconds:", t2.UTC().Format("2006-01-02 15:04:05"))

	// 6. Tự động parse timestamp dựa trên độ lớn của giá trị
	t3 := chrono.AutoParseTimestamp(1782812818000000) // Microseconds
	fmt.Println("Auto Parse Microseconds:", t3.UTC().Format("2006-01-02 15:04:05"))

	// 7. Định dạng đa ngôn ngữ (tiếng Việt và tiếng Nhật)
	fmt.Println("Vietnamese:", chrono.FormatLocalePHP(testTime, "l, d F Y", "vi"))
	fmt.Println("Japanese:", chrono.FormatLocaleJS(testTime, "dddd, DD MMMM YYYY", "ja"))

	// Output:
	// PHP formatted: 2026-07-01 12:40:18
	// JS formatted: 2026-07-01 12:40:18
	// PHP parsed: 2026-07-01 12:40:18
	// Go Layout from PHP: 2006-01-02 15:04:05.000 -0700
	// Parse Seconds: 2026-06-30 09:46:58
	// Parse Milliseconds: 2026-06-30 09:46:58
	// Auto Parse Microseconds: 2026-06-30 09:46:58
	// Vietnamese: Thứ Tư, 01 Tháng Bảy 2026
	// Japanese: 水曜日, 01 7月 2026
}
