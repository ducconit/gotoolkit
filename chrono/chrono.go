// Package chrono cung cấp các tiện ích xử lý thời gian tối giản, hiệu năng cao,
// hỗ trợ chuyển đổi các định dạng quen thuộc của PHP và JS/Moment sang layout của Go.
package chrono

import (
	"strings"
	"time"
)

// PHP chuyển đổi chuỗi định dạng thời gian kiểu PHP (ví dụ: "Y-m-d H:i:s")
// sang layout tương ứng của Go (ví dụ: "2006-01-02 15:04:05").
// Để in các chữ cái thông thường tránh bị dịch thành token, hãy sử dụng dấu gạch chéo ngược để escape (ví dụ: "\h\e\l\l\o").
// Hàm sử dụng switch-case tối ưu và strings.Builder để đạt tốc độ tối đa (0 allocation nếu trống).
func PHP(format string) string {
	if format == "" {
		return ""
	}
	var builder strings.Builder
	builder.Grow(len(format) * 2)

	for i := 0; i < len(format); i++ {
		c := format[i]
		// Hỗ trợ escape ký tự bằng dấu gạch chéo ngược \
		if c == '\\' && i+1 < len(format) {
			i++
			builder.WriteByte(format[i])
			continue
		}
		switch c {
		case 'Y':
			builder.WriteString("2006")
		case 'y':
			builder.WriteString("06")
		case 'm':
			builder.WriteString("01")
		case 'n':
			builder.WriteString("1")
		case 'M':
			builder.WriteString("Jan")
		case 'F':
			builder.WriteString("January")
		case 'd':
			builder.WriteString("02")
		case 'j':
			builder.WriteString("2")
		case 'D':
			builder.WriteString("Mon")
		case 'l':
			builder.WriteString("Monday")
		case 'H':
			builder.WriteString("15")
		case 'h':
			builder.WriteString("03")
		case 'g':
			builder.WriteString("3")
		case 'i':
			builder.WriteString("04")
		case 's':
			builder.WriteString("05")
		case 'v':
			builder.WriteString("000")
		case 'u':
			builder.WriteString("000000")
		case 'a':
			builder.WriteString("pm")
		case 'A':
			builder.WriteString("PM")
		case 'O':
			builder.WriteString("-0700")
		case 'P':
			builder.WriteString("-07:00")
		case 'T':
			builder.WriteString("MST")
		default:
			builder.WriteByte(c)
		}
	}
	return builder.String()
}

// JS chuyển đổi chuỗi định dạng thời gian kiểu JS/Moment (ví dụ: "YYYY-MM-DD HH:mm:ss")
// sang layout tương ứng của Go (ví dụ: "2006-01-02 15:04:05").
// Để in các chuỗi chữ thông thường, hãy bọc chúng trong cặp ngoặc vuông (ví dụ: "[Hello]").
// Hàm sử dụng phương pháp lookahead scanning để match các token dài, đạt hiệu năng O(N).
func JS(format string) string {
	if format == "" {
		return ""
	}
	var builder strings.Builder
	builder.Grow(len(format) * 2)

	i := 0
	n := len(format)
	for i < n {
		c := format[i]

		// Hỗ trợ bọc text cố định trong dấu ngoặc vuông [...] để tránh bị dịch thành token
		if c == '[' {
			i++ // Bỏ qua '['
			end := i
			for end < n && format[end] != ']' {
				end++
			}
			builder.WriteString(format[i:end])
			i = end
			if i < n {
				i++ // Bỏ qua ']'
			}
			continue
		}

		// Match tokens dài 4 ký tự
		if i+4 <= n {
			sub := format[i : i+4]
			switch sub {
			case "YYYY":
				builder.WriteString("2006")
				i += 4
				continue
			case "MMMM":
				builder.WriteString("January")
				i += 4
				continue
			case "dddd":
				builder.WriteString("Monday")
				i += 4
				continue
			}
		}
		// Match tokens dài 3 ký tự
		if i+3 <= n {
			sub := format[i : i+3]
			switch sub {
			case "MMM":
				builder.WriteString("Jan")
				i += 3
				continue
			case "ddd":
				builder.WriteString("Mon")
				i += 3
				continue
			case "SSS":
				builder.WriteString("000")
				i += 3
				continue
			}
		}
		// Match tokens dài 2 ký tự
		if i+2 <= n {
			sub := format[i : i+2]
			switch sub {
			case "YY":
				builder.WriteString("06")
				i += 2
				continue
			case "MM":
				builder.WriteString("01")
				i += 2
				continue
			case "DD":
				builder.WriteString("02")
				i += 2
				continue
			case "HH":
				builder.WriteString("15")
				i += 2
				continue
			case "hh":
				builder.WriteString("03")
				i += 2
				continue
			case "mm":
				builder.WriteString("04")
				i += 2
				continue
			case "ss":
				builder.WriteString("05")
				i += 2
				continue
			case "ZZ":
				builder.WriteString("-0700")
				i += 2
				continue
			}
		}
		// Match token 1 ký tự
		switch c {
		case 'M':
			builder.WriteString("1")
		case 'D':
			builder.WriteString("2")
		case 'H':
			builder.WriteString("15")
		case 'h':
			builder.WriteString("3")
		case 'm':
			builder.WriteString("4")
		case 's':
			builder.WriteString("5")
		case 'A':
			builder.WriteString("PM")
		case 'a':
			builder.WriteString("pm")
		case 'Z':
			builder.WriteString("-07:00")
		default:
			builder.WriteByte(c)
		}
		i++
	}
	return builder.String()
}

// FormatPHP định dạng time.Time theo cú pháp PHP.
func FormatPHP(t time.Time, format string) string {
	return t.Format(PHP(format))
}

// ParsePHP chuyển đổi một chuỗi thời gian thành time.Time theo cú pháp PHP.
func ParsePHP(format, value string) (time.Time, error) {
	return time.Parse(PHP(format), value)
}

// ParseInLocationPHP chuyển đổi một chuỗi thời gian thành time.Time theo múi giờ và cú pháp PHP.
func ParseInLocationPHP(format, value string, loc *time.Location) (time.Time, error) {
	return time.ParseInLocation(PHP(format), value, loc)
}

// FormatJS định dạng time.Time theo cú pháp JS/Moment.
func FormatJS(t time.Time, format string) string {
	return t.Format(JS(format))
}

// ParseJS chuyển đổi một chuỗi thời gian thành time.Time theo cú pháp JS/Moment.
func ParseJS(format, value string) (time.Time, error) {
	return time.Parse(JS(format), value)
}

// ParseInLocationJS chuyển đổi một chuỗi thời gian thành time.Time theo múi giờ và cú pháp JS/Moment.
func ParseInLocationJS(format, value string, loc *time.Location) (time.Time, error) {
	return time.ParseInLocation(JS(format), value, loc)
}

// ParseTimestamp chuyển đổi timestamp dạng số nguyên với đơn vị thời gian tường minh thành time.Time.
// Hỗ trợ các đơn vị tiêu chuẩn: time.Second, time.Millisecond, time.Microsecond, time.Nanosecond.
// Đạt hiệu năng tối đa (0 allocation, O(1) CPU).
func ParseTimestamp(val int64, unit time.Duration) time.Time {
	switch unit {
	case time.Second:
		return time.Unix(val, 0)
	case time.Millisecond:
		return time.UnixMilli(val)
	case time.Microsecond:
		return time.UnixMicro(val)
	case time.Nanosecond:
		return time.Unix(0, val)
	default:
		// Mặc định fallback về giây nếu đơn vị không hợp lệ
		return time.Unix(val, 0)
	}
}

// AutoParseTimestamp tự động nhận dạng đơn vị của timestamp dựa trên độ lớn của giá trị
// (Giây, Mili, Micro, Nano) và chuyển đổi thành time.Time tương ứng.
// Giải pháp đạt độ chính xác cao và an toàn cho tới năm 3000.
func AutoParseTimestamp(val int64) time.Time {
	const (
		limitSecond      = 32503680000          // Ngưỡng giây (năm 3000)
		limitMillisecond = limitSecond * 1000   // Ngưỡng mili giây
		limitMicrosecond = limitSecond * 1000000 // Ngưỡng micro giây
	)

	if val < limitSecond {
		return time.Unix(val, 0)
	} else if val < limitMillisecond {
		return time.UnixMilli(val)
	} else if val < limitMicrosecond {
		return time.UnixMicro(val)
	} else {
		return time.Unix(0, val)
	}
}
// Locale đại diện cho định danh ngôn ngữ. Chấp nhận string (ví dụ: "vi", "ja-JP")
// hoặc bất kỳ đối tượng nào có phương thức String() (ví dụ: language.Tag của package golang.org/x/text/language).
type Locale any

var weekdaysMap = map[string][7]string{
	"en": {"Sunday", "Monday", "Tuesday", "Wednesday", "Thursday", "Friday", "Saturday"},
	"vi": {"Chủ Nhật", "Thứ Hai", "Thứ Ba", "Thứ Tư", "Thứ Năm", "Thứ Sáu", "Thứ Bảy"},
	"ja": {"日曜日", "月曜日", "火曜日", "水曜日", "木曜日", "金曜日", "土曜日"},
	"ko": {"일요일", "월요일", "화요일", "수요일", "목요일", "금요일", "토요일"},
	"zh": {"星期日", "星期一", "星期二", "星期三", "星期四", "星期五", "星期六"},
}

var weekdaysShortMap = map[string][7]string{
	"en": {"Sun", "Mon", "Tue", "Wed", "Thu", "Fri", "Sat"},
	"vi": {"CN", "T2", "T3", "T4", "T5", "T6", "T7"},
	"ja": {"日", "月", "火", "水", "木", "金", "土"},
	"ko": {"일", "월", "화", "수", "목", "금", "토"},
	"zh": {"周日", "周一", "周二", "周三", "周四", "周五", "周六"},
}

var monthsMap = map[string][12]string{
	"en": {"January", "February", "March", "April", "May", "June", "July", "August", "September", "October", "November", "December"},
	"vi": {"Tháng Một", "Tháng Hai", "Tháng Ba", "Tháng Tư", "Tháng Năm", "Tháng Sáu", "Tháng Bảy", "Tháng Tám", "Tháng Chín", "Tháng Mười", "Tháng Mười Một", "Tháng Mười Hai"},
	"ja": {"1月", "2月", "3月", "4月", "5月", "6月", "7月", "8月", "9月", "10月", "11月", "12月"},
	"ko": {"1월", "2월", "3월", "4월", "5월", "6월", "7월", "8월", "9월", "10월", "11월", "12월"},
	"zh": {"一月", "二月", "三月", "四月", "五月", "六月", "七月", "八月", "九月", "十月", "十一月", "十二月"},
}

var monthsShortMap = map[string][12]string{
	"en": {"Jan", "Feb", "Mar", "Apr", "May", "Jun", "Jul", "Aug", "Sep", "Oct", "Nov", "Dec"},
	"vi": {"Thg 1", "Thg 2", "Thg 3", "Thg 4", "Thg 5", "Thg 6", "Thg 7", "Thg 8", "Thg 9", "Thg 10", "Thg 11", "Thg 12"},
	"ja": {"1月", "2月", "3月", "4月", "5月", "6月", "7月", "8月", "9月", "10月", "11月", "12月"},
	"ko": {"1월", "2월", "3월", "4월", "5월", "6월", "7월", "8월", "9월", "10월", "11월", "12월"},
	"zh": {"1月", "2月", "3月", "4月", "5月", "6月", "7月", "8月", "9月", "10月", "11月", "12月"},
}

func parseLocale(loc Locale) string {
	switch v := loc.(type) {
	case string:
		return normalizeLocale(v)
	case interface{ String() string }:
		return normalizeLocale(v.String())
	default:
		return "en"
	}
}

func normalizeLocale(loc string) string {
	if len(loc) < 2 {
		return "en"
	}
	sub := loc[:2]
	var b [2]byte
	b[0] = sub[0]
	b[1] = sub[1]
	for idx := range 2 {
		if b[idx] >= 'A' && b[idx] <= 'Z' {
			b[idx] = b[idx] + ('a' - 'A')
		}
	}
	normalized := string(b[:])
	switch normalized {
	case "vi", "ja", "ko", "zh":
		return normalized
	default:
		return "en"
	}
}

// WeekdayName trả về tên thứ trong tuần theo ngôn ngữ chỉ định.
func WeekdayName(w time.Weekday, loc Locale, short ...bool) string {
	localeStr := parseLocale(loc)
	isShort := len(short) > 0 && short[0]

	if isShort {
		return weekdaysShortMap[localeStr][w]
	}
	return weekdaysMap[localeStr][w]
}

// MonthName trả về tên tháng theo ngôn ngữ chỉ định.
func MonthName(m time.Month, loc Locale, short ...bool) string {
	localeStr := parseLocale(loc)
	isShort := len(short) > 0 && short[0]
	idx := int(m) - 1

	if idx < 0 || idx > 11 {
		return ""
	}

	if isShort {
		return monthsShortMap[localeStr][idx]
	}
	return monthsMap[localeStr][idx]
}

// itoa chuyển đổi số nguyên thành chuỗi với số lượng chữ số cố định (điền thêm số 0 dẫn đầu).
// Hàm cực kỳ tối ưu, không cấp phát heap nếu ghi vào builder.
func itoa(val int, width int) string {
	var buf [20]byte
	i := len(buf) - 1
	for val >= 10 || width > 1 {
		buf[i] = byte(val%10) + '0'
		i--
		val /= 10
		width--
	}
	buf[i] = byte(val) + '0'
	return string(buf[i:])
}

// formatOffset chuyển đổi múi giờ dạng giây sang chuỗi (+HHMM hoặc +HH:MM)
func formatOffset(offset int, colon bool) string {
	var buf [6]byte
	i := 0
	if offset < 0 {
		buf[i] = '-'
		offset = -offset
	} else {
		buf[i] = '+'
	}
	i++

	hours := offset / 3600
	minutes := (offset % 3600) / 60

	buf[i] = byte(hours/10) + '0'
	buf[i+1] = byte(hours%10) + '0'
	i += 2

	if colon {
		buf[i] = ':'
		i++
	}

	buf[i] = byte(minutes/10) + '0'
	buf[i+1] = byte(minutes%10) + '0'
	i += 2

	return string(buf[:i])
}

// FormatLocalePHP định dạng time.Time theo cú pháp PHP và bản địa hóa tên thứ/tháng theo ngôn ngữ chỉ định.
func FormatLocalePHP(t time.Time, format string, loc Locale) string {
	if format == "" {
		return ""
	}
	localeStr := parseLocale(loc)
	var builder strings.Builder
	builder.Grow(len(format) * 3)

	for i := 0; i < len(format); i++ {
		c := format[i]
		if c == '\\' && i+1 < len(format) {
			i++
			builder.WriteByte(format[i])
			continue
		}
		switch c {
		case 'Y':
			builder.WriteString(itoa(t.Year(), 4))
		case 'y':
			builder.WriteString(itoa(t.Year()%100, 2))
		case 'm':
			builder.WriteString(itoa(int(t.Month()), 2))
		case 'n':
			builder.WriteString(itoa(int(t.Month()), 1))
		case 'M':
			builder.WriteString(MonthName(t.Month(), localeStr, true))
		case 'F':
			builder.WriteString(MonthName(t.Month(), localeStr, false))
		case 'd':
			builder.WriteString(itoa(t.Day(), 2))
		case 'j':
			builder.WriteString(itoa(t.Day(), 1))
		case 'D':
			builder.WriteString(WeekdayName(t.Weekday(), localeStr, true))
		case 'l':
			builder.WriteString(WeekdayName(t.Weekday(), localeStr, false))
		case 'H':
			builder.WriteString(itoa(t.Hour(), 2))
		case 'h':
			h := t.Hour()
			if h == 0 {
				h = 12
			} else if h > 12 {
				h -= 12
			}
			builder.WriteString(itoa(h, 2))
		case 'g':
			h := t.Hour()
			if h == 0 {
				h = 12
			} else if h > 12 {
				h -= 12
			}
			builder.WriteString(itoa(h, 1))
		case 'i':
			builder.WriteString(itoa(t.Minute(), 2))
		case 's':
			builder.WriteString(itoa(t.Second(), 2))
		case 'v':
			builder.WriteString(itoa(t.Nanosecond()/1e6, 3))
		case 'u':
			builder.WriteString(itoa(t.Nanosecond()/1e3, 6))
		case 'a':
			if t.Hour() < 12 {
				builder.WriteString("am")
			} else {
				builder.WriteString("pm")
			}
		case 'A':
			if t.Hour() < 12 {
				builder.WriteString("AM")
			} else {
				builder.WriteString("PM")
			}
		case 'O':
			_, offset := t.Zone()
			builder.WriteString(formatOffset(offset, false))
		case 'P':
			_, offset := t.Zone()
			builder.WriteString(formatOffset(offset, true))
		case 'T':
			zone, _ := t.Zone()
			builder.WriteString(zone)
		default:
			builder.WriteByte(c)
		}
	}
	return builder.String()
}

// FormatLocaleJS định dạng time.Time theo cú pháp JS/Moment và bản địa hóa tên thứ/tháng theo ngôn ngữ chỉ định.
func FormatLocaleJS(t time.Time, format string, loc Locale) string {
	if format == "" {
		return ""
	}
	localeStr := parseLocale(loc)
	var builder strings.Builder
	builder.Grow(len(format) * 3)

	i := 0
	n := len(format)
	for i < n {
		c := format[i]

		// Bỏ qua escape bọc trong dấu ngoặc vuông [...]
		if c == '[' {
			i++ // Bỏ qua '['
			end := i
			for end < n && format[end] != ']' {
				end++
			}
			builder.WriteString(format[i:end])
			i = end
			if i < n {
				i++ // Bỏ qua ']'
			}
			continue
		}

		// Match tokens dài 4 ký tự
		if i+4 <= n {
			sub := format[i : i+4]
			switch sub {
			case "YYYY":
				builder.WriteString(itoa(t.Year(), 4))
				i += 4
				continue
			case "MMMM":
				builder.WriteString(MonthName(t.Month(), localeStr, false))
				i += 4
				continue
			case "dddd":
				builder.WriteString(WeekdayName(t.Weekday(), localeStr, false))
				i += 4
				continue
			}
		}
		// Match tokens dài 3 ký tự
		if i+3 <= n {
			sub := format[i : i+3]
			switch sub {
			case "MMM":
				builder.WriteString(MonthName(t.Month(), localeStr, true))
				i += 3
				continue
			case "ddd":
				builder.WriteString(WeekdayName(t.Weekday(), localeStr, true))
				i += 3
				continue
			case "SSS":
				builder.WriteString(itoa(t.Nanosecond()/1e6, 3))
				i += 3
				continue
			}
		}
		// Match tokens dài 2 ký tự
		if i+2 <= n {
			sub := format[i : i+2]
			switch sub {
			case "YY":
				builder.WriteString(itoa(t.Year()%100, 2))
				i += 2
				continue
			case "MM":
				builder.WriteString(itoa(int(t.Month()), 2))
				i += 2
				continue
			case "DD":
				builder.WriteString(itoa(t.Day(), 2))
				i += 2
				continue
			case "HH":
				builder.WriteString(itoa(t.Hour(), 2))
				i += 2
				continue
			case "hh":
				h := t.Hour()
				if h == 0 {
					h = 12
				} else if h > 12 {
					h -= 12
				}
				builder.WriteString(itoa(h, 2))
				i += 2
				continue
			case "mm":
				builder.WriteString(itoa(t.Minute(), 2))
				i += 2
				continue
			case "ss":
				builder.WriteString(itoa(t.Second(), 2))
				i += 2
				continue
			case "ZZ":
				_, offset := t.Zone()
				builder.WriteString(formatOffset(offset, false))
				i += 2
				continue
			}
		}
		// Match token 1 ký tự
		switch c {
		case 'M':
			builder.WriteString(itoa(int(t.Month()), 1))
		case 'D':
			builder.WriteString(itoa(t.Day(), 1))
		case 'H':
			builder.WriteString(itoa(t.Hour(), 1))
		case 'h':
			h := t.Hour()
			if h == 0 {
				h = 12
			} else if h > 12 {
				h -= 12
			}
			builder.WriteString(itoa(h, 1))
		case 'm':
			builder.WriteString(itoa(t.Minute(), 1))
		case 's':
			builder.WriteString(itoa(t.Second(), 1))
		case 'A':
			if t.Hour() < 12 {
				builder.WriteString("AM")
			} else {
				builder.WriteString("PM")
			}
		case 'a':
			if t.Hour() < 12 {
				builder.WriteString("am")
			} else {
				builder.WriteString("pm")
			}
		case 'Z':
			_, offset := t.Zone()
			builder.WriteString(formatOffset(offset, true))
		default:
			builder.WriteByte(c)
		}
		i++
	}
	return builder.String()
}
