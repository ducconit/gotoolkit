package chrono

import (
	"testing"
	"time"
)

// TestPHP kiểm tra việc ánh xạ từ PHP format sang Go layout
func TestPHP(t *testing.T) {
	tests := []struct {
		name   string
		format string
		want   string
	}{
		{"Empty string", "", ""},
		{"Normal DateTime", "Y-m-d H:i:s", "2006-01-02 15:04:05"},
		{"Date Only", "Y/m/d", "2006/01/02"},
		{"Short Date", "y-n-j", "06-1-2"},
		{"Month Day Names", "l, F j, Y", "Monday, January 2, 2006"},
		{"AM/PM Time", "h:i:s A", "03:04:05 PM"},
		{"Time with Milliseconds", "H:i:s.v", "15:04:05.000"},
		{"Time with Timezone", "Y-m-d H:i:s O P T", "2006-01-02 15:04:05 -0700 -07:00 MST"},
		{"Plain text with escape", "Y-m-d \\H\\e\\l\\l\\o H:i:s", "2006-01-02 Hello 15:04:05"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := PHP(tt.format); got != tt.want {
				t.Errorf("PHP(%q) = %q, want %q", tt.format, got, tt.want)
			}
		})
	}
}

// TestJS kiểm tra việc ánh xạ từ JS/Moment format sang Go layout
func TestJS(t *testing.T) {
	tests := []struct {
		name   string
		format string
		want   string
	}{
		{"Empty string", "", ""},
		{"Normal DateTime", "YYYY-MM-DD HH:mm:ss", "2006-01-02 15:04:05"},
		{"Date Only", "YYYY/MM/DD", "2006/01/02"},
		{"Short Date", "YY-M-D", "06-1-2"},
		{"Month Day Names", "dddd, MMMM D, YYYY", "Monday, January 2, 2006"},
		{"AM/PM Time", "hh:mm:ss A", "03:04:05 PM"},
		{"Time with Milliseconds", "HH:mm:ss.SSS", "15:04:05.000"},
		{"Time with Timezone", "YYYY-MM-DD HH:mm:ss ZZ Z", "2006-01-02 15:04:05 -0700 -07:00"},
		{"Plain text with escape", "YYYY-MM-DD [Hello] HH:mm:ss", "2006-01-02 Hello 15:04:05"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := JS(tt.format); got != tt.want {
				t.Errorf("JS(%q) = %q, want %q", tt.format, got, tt.want)
			}
		})
	}
}

// TestFormatAndParse kiểm tra hoạt động thực tế của các hàm Format và Parse
func TestFormatAndParse(t *testing.T) {
	// Khởi tạo thời điểm mốc cố định: 2026-07-01 12:40:18 (UTC)
	loc := time.UTC
	testTime := time.Date(2026, time.July, 1, 12, 40, 18, 0, loc)

	t.Run("PHP Format and Parse", func(t *testing.T) {
		phpFormat := "Y-m-d H:i:s"
		formatted := FormatPHP(testTime, phpFormat)
		expectedStr := "2026-07-01 12:40:18"

		if formatted != expectedStr {
			t.Errorf("FormatPHP() = %q, want %q", formatted, expectedStr)
		}

		parsed, err := ParseInLocationPHP(phpFormat, formatted, loc)
		if err != nil {
			t.Fatalf("ParseInLocationPHP() error = %v", err)
		}

		if !parsed.Equal(testTime) {
			t.Errorf("Parsed time = %v, want %v", parsed, testTime)
		}
	})

	t.Run("JS Format and Parse", func(t *testing.T) {
		jsFormat := "YYYY-MM-DD HH:mm:ss"
		formatted := FormatJS(testTime, jsFormat)
		expectedStr := "2026-07-01 12:40:18"

		if formatted != expectedStr {
			t.Errorf("FormatJS() = %q, want %q", formatted, expectedStr)
		}

		parsed, err := ParseInLocationJS(jsFormat, formatted, loc)
		if err != nil {
			t.Fatalf("ParseInLocationJS() error = %v", err)
		}

		if !parsed.Equal(testTime) {
			t.Errorf("Parsed time = %v, want %v", parsed, testTime)
		}
	})
}

// === BENCHMARKS ===

func BenchmarkPHP_Conversion(b *testing.B) {
	format := "Y-m-d H:i:s.v O P T"
	b.ResetTimer()
	for b.Loop() {
		_ = PHP(format)
	}
}

func BenchmarkJS_Conversion(b *testing.B) {
	format := "YYYY-MM-DD HH:mm:ss.SSS ZZ Z"
	b.ResetTimer()
	for b.Loop() {
		_ = JS(format)
	}
}

func BenchmarkFormatPHP(b *testing.B) {
	t := time.Now()
	format := "Y-m-d H:i:s"
	b.ResetTimer()
	for b.Loop() {
		_ = FormatPHP(t, format)
	}
}

func BenchmarkFormatJS(b *testing.B) {
	t := time.Now()
	format := "YYYY-MM-DD HH:mm:ss"
	b.ResetTimer()
	for b.Loop() {
		_ = FormatJS(t, format)
	}
}

// TestParseTimestamp kiểm tra hàm ParseTimestamp với các đơn vị thời gian cụ thể
func TestParseTimestamp(t *testing.T) {
	tests := []struct {
		name string
		val  int64
		unit time.Duration
		want time.Time
	}{
		{"Seconds", 1782812818, time.Second, time.Unix(1782812818, 0)},
		{"Milliseconds", 1782812818000, time.Millisecond, time.UnixMilli(1782812818000)},
		{"Microseconds", 1782812818000000, time.Microsecond, time.UnixMicro(1782812818000000)},
		{"Nanoseconds", 1782812818000000000, time.Nanosecond, time.Unix(0, 1782812818000000000)},
		{"Invalid Unit fallback to seconds", 1782812818, time.Hour, time.Unix(1782812818, 0)},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ParseTimestamp(tt.val, tt.unit)
			if !got.Equal(tt.want) {
				t.Errorf("ParseTimestamp(%d, %v) = %v, want %v", tt.val, tt.unit, got, tt.want)
			}
		})
	}
}

// TestAutoParseTimestamp kiểm tra khả năng tự động nhận dạng đơn vị của timestamp
func TestAutoParseTimestamp(t *testing.T) {
	tests := []struct {
		name string
		val  int64
		want time.Time
	}{
		{"Auto Seconds", 1782812818, time.Unix(1782812818, 0)},
		{"Auto Milliseconds", 1782812818123, time.UnixMilli(1782812818123)},
		{"Auto Microseconds", 1782812818123456, time.UnixMicro(1782812818123456)},
		{"Auto Nanoseconds", 1782812818123456789, time.Unix(0, 1782812818123456789)},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := AutoParseTimestamp(tt.val)
			if !got.Equal(tt.want) {
				t.Errorf("AutoParseTimestamp(%d) = %v, want %v", tt.val, got, tt.want)
			}
		})
	}
}

// BenchmarkAutoParseTimestamp đo hiệu năng hàm tự động nhận dạng timestamp
func BenchmarkAutoParseTimestamp(b *testing.B) {
	val := int64(1782812818123) // Millisecond timestamp
	b.ResetTimer()
	for b.Loop() {
		_ = AutoParseTimestamp(val)
	}
}
// MockLanguageTag giả lập đối tượng Tag của golang.org/x/text/language
type MockLanguageTag struct {
	tag string
}

func (m MockLanguageTag) String() string {
	return m.tag
}

// TestWeekdayName kiểm tra việc lấy tên thứ theo các locale khác nhau
func TestWeekdayName(t *testing.T) {
	tests := []struct {
		name   string
		w      time.Weekday
		loc    Locale
		short  bool
		want   string
	}{
		{"English Monday", time.Monday, "en", false, "Monday"},
		{"English Monday Short", time.Monday, "en", true, "Mon"},
		{"Vietnamese Monday", time.Monday, "vi", false, "Thứ Hai"},
		{"Vietnamese Monday Short", time.Monday, "vi", true, "T2"},
		{"Japanese Monday", time.Monday, "ja", false, "月曜日"},
		{"Japanese Monday Short", time.Monday, "ja", true, "月"},
		{"Korean Monday", time.Monday, "ko", false, "월요일"},
		{"Chinese Monday", time.Monday, "zh", false, "星期一"},
		{"Mock Language Tag", time.Monday, MockLanguageTag{"vi-VN"}, false, "Thứ Hai"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := WeekdayName(tt.w, tt.loc, tt.short)
			if got != tt.want {
				t.Errorf("WeekdayName(%v, %v, %v) = %q, want %q", tt.w, tt.loc, tt.short, got, tt.want)
			}
		})
	}
}

// TestMonthName kiểm tra việc lấy tên tháng theo các locale khác nhau
func TestMonthName(t *testing.T) {
	tests := []struct {
		name   string
		m      time.Month
		loc    Locale
		short  bool
		want   string
	}{
		{"English Jan", time.January, "en", false, "January"},
		{"English Jan Short", time.January, "en", true, "Jan"},
		{"Vietnamese Jan", time.January, "vi", false, "Tháng Một"},
		{"Vietnamese Jan Short", time.January, "vi", true, "Thg 1"},
		{"Japanese Jan", time.January, "ja", false, "1月"},
		{"Korean Jan", time.January, "ko", false, "1월"},
		{"Chinese Jan", time.January, "zh", false, "一月"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := MonthName(tt.m, tt.loc, tt.short)
			if got != tt.want {
				t.Errorf("MonthName(%v, %v, %v) = %q, want %q", tt.m, tt.loc, tt.short, got, tt.want)
			}
		})
	}
}

// TestFormatLocalePHP kiểm tra hàm định dạng đa ngôn ngữ PHP
func TestFormatLocalePHP(t *testing.T) {
	loc := time.UTC
	testTime := time.Date(2026, time.July, 1, 12, 40, 18, 0, loc)

	tests := []struct {
		name   string
		format string
		locale Locale
		want   string
	}{
		{"Vietnamese DateTime", "l, d F Y H:i:s", "vi", "Thứ Tư, 01 Tháng Bảy 2026 12:40:18"},
		{"Japanese DateTime", "l, d F Y H:i:s", "ja", "水曜日, 01 7月 2026 12:40:18"},
		{"English DateTime", "l, d F Y H:i:s", "en", "Wednesday, 01 July 2026 12:40:18"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := FormatLocalePHP(testTime, tt.format, tt.locale)
			if got != tt.want {
				t.Errorf("FormatLocalePHP(%v, %q, %v) = %q, want %q", testTime, tt.format, tt.locale, got, tt.want)
			}
		})
	}
}

// TestFormatLocaleJS kiểm tra hàm định dạng đa ngôn ngữ JS
func TestFormatLocaleJS(t *testing.T) {
	loc := time.UTC
	testTime := time.Date(2026, time.July, 1, 12, 40, 18, 0, loc)

	tests := []struct {
		name   string
		format string
		locale Locale
		want   string
	}{
		{"Vietnamese DateTime JS", "dddd, DD MMMM YYYY HH:mm:ss", "vi", "Thứ Tư, 01 Tháng Bảy 2026 12:40:18"},
		{"Japanese DateTime JS", "dddd, DD MMMM YYYY HH:mm:ss", "ja", "水曜日, 01 7月 2026 12:40:18"},
		{"English DateTime JS", "dddd, DD MMMM YYYY HH:mm:ss", "en", "Wednesday, 01 July 2026 12:40:18"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := FormatLocaleJS(testTime, tt.format, tt.locale)
			if got != tt.want {
				t.Errorf("FormatLocaleJS(%v, %q, %v) = %q, want %q", testTime, tt.format, tt.locale, got, tt.want)
			}
		})
	}
}

// BenchmarkFormatLocalePHP đo hiệu năng định dạng đa ngôn ngữ
func BenchmarkFormatLocalePHP(b *testing.B) {
	t := time.Now()
	format := "l, d F Y H:i:s"
	b.ResetTimer()
	for b.Loop() {
		_ = FormatLocalePHP(t, format, "vi")
	}
}
