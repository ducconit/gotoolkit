// Package str cung cấp các tiện ích xử lý và kiểm tra chuỗi tối giản,
// hiệu năng cao, thread-safe và zero-allocation (tránh sử dụng regular expressions chậm).
package str

import (
	"strings"
	"unicode"
)

// IsAlphaSpace kiểm tra chuỗi đầu vào chỉ chứa ký tự chữ cái (hỗ trợ Unicode/ngôn ngữ có dấu)
// và dấu cách. Không chứa số hoặc ký tự đặc biệt.
// Phù hợp để xác thực tên người dùng (Full Name), Nickname hiển thị.
// Hàm sử dụng fast-path cho ký tự ASCII để tối ưu hóa hiệu năng CPU.
func IsAlphaSpace(s string) bool {
	if s == "" {
		return false
	}

	for _, r := range s {
		if r < 128 {
			if !((r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') || r == ' ') {
				return false
			}
		} else if !unicode.IsLetter(r) && r != ' ' {
			return false
		}
	}
	return true
}

var isUsernameChar = [256]bool{
	'a': true, 'b': true, 'c': true, 'd': true, 'e': true, 'f': true, 'g': true, 'h': true,
	'i': true, 'j': true, 'k': true, 'l': true, 'm': true, 'n': true, 'o': true, 'p': true,
	'q': true, 'r': true, 's': true, 't': true, 'u': true, 'v': true, 'w': true, 'x': true,
	'y': true, 'z': true,
	'0': true, '1': true, '2': true, '3': true, '4': true, '5': true, '6': true, '7': true,
	'8': true, '9': true,
	'_': true,
}

// IsUsername kiểm tra xem chuỗi có phải là một username hợp lệ hay không.
// Username hợp lệ chỉ được phép chứa chữ thường a-z, ký số 0-9 và ký tự gạch dưới _.
// Hàm sử dụng Lookup Table mảng tĩnh để đạt tốc độ xử lý tối đa (zero allocation).
func IsUsername(s string) bool {
	if s == "" {
		return false
	}

	for i := range len(s) {
		if !isUsernameChar[s[i]] {
			return false
		}
	}
	return true
}

// ContainsHTML kiểm tra xem chuỗi đầu vào có chứa mã nguồn HTML, XML hoặc Javascript nhạy cảm không.
// Phù hợp để phát hiện sớm các cuộc tấn công Script Injection / XSS cơ bản trước khi lưu vào DB.
// Hàm sử dụng quét tĩnh không rẽ nhánh để tối ưu hóa CPU tối đa.
func ContainsHTML(s string) bool {
	if s == "" {
		return false
	}

	// 1. Phát hiện các thẻ tag HTML/XML dạng <...>
	open := strings.IndexByte(s, '<')
	if open != -1 {
		close := strings.IndexByte(s[open:], '>')
		if close != -1 {
			return true
		}
	}

	// 2. Phát hiện "javascript:" case-insensitive
	if len(s) >= 11 {
		for i := 0; i <= len(s)-11; i++ {
			c := s[i]
			if c == 'j' || c == 'J' {
				if (s[i+1] == 'a' || s[i+1] == 'A') &&
					(s[i+2] == 'v' || s[i+2] == 'V') &&
					(s[i+3] == 'a' || s[i+3] == 'A') &&
					(s[i+4] == 's' || s[i+4] == 'S') &&
					(s[i+5] == 'c' || s[i+5] == 'C') &&
					(s[i+6] == 'r' || s[i+6] == 'R') &&
					(s[i+7] == 'i' || s[i+7] == 'I') &&
					(s[i+8] == 'p' || s[i+8] == 'P') &&
					(s[i+9] == 't' || s[i+9] == 'T') &&
					s[i+10] == ':' {
					return true
				}
			}
		}
	}

	return false
}

// CleanSpace loại bỏ khoảng trắng ở đầu, cuối và rút gọn các khoảng trắng liên tiếp
// ở giữa chuỗi thành một khoảng trắng duy nhất.
// Hàm sử dụng fast-path và tối ưu sao chép khối (block copy) để giảm thiểu allocations.
func CleanSpace(s string) string {
	if s == "" {
		return ""
	}

	// Kiểm tra ký tự đầu tiên
	if isSpace(s[0]) {
		return cleanSpaceSlow(s, 0)
	}

	inSpace := false
	for i := 1; i < len(s); i++ {
		c := s[i]
		if isSpace(c) {
			if c != ' ' || inSpace {
				return cleanSpaceSlow(s, i)
			}
			inSpace = true
		} else {
			inSpace = false
		}
	}

	// Kiểm tra ký tự cuối cùng
	if isSpace(s[len(s)-1]) {
		return cleanSpaceSlow(s, len(s)-1)
	}

	return s
}

func cleanSpaceSlow(s string, start int) string {
	var builder strings.Builder
	builder.Grow(len(s))

	var inSpace bool
	var firstNonSpace bool

	if start > 0 {
		lastChar := s[start-1]
		if isSpace(lastChar) {
			builder.WriteString(s[:start-1])
			inSpace = true
			firstNonSpace = (start - 1) > 0
		} else {
			builder.WriteString(s[:start])
			inSpace = isSpace(s[start])
			firstNonSpace = true
		}
	} else {
		inSpace = false
		firstNonSpace = false
	}

	for i := start; i < len(s); i++ {
		c := s[i]
		if isSpace(c) {
			if firstNonSpace {
				inSpace = true
			}
		} else {
			if inSpace {
				builder.WriteByte(' ')
				inSpace = false
			}
			builder.WriteByte(c)
			firstNonSpace = true
		}
	}

	return builder.String()
}

func isSpace(c byte) bool {
	return c == ' ' || (c >= 9 && c <= 13)
}

// RemoveAccents loại bỏ các dấu tiếng Việt khỏi chuỗi.
// Hỗ trợ cả bảng mã Unicode dựng sẵn (precomposed) và tổ hợp (decomposed).
// Đạt hiệu năng tối đa (0 allocation nếu chuỗi không có dấu tiếng Việt).
func RemoveAccents(s string) string {
	if s == "" {
		return ""
	}

	hasAccent := false
	var accentIdx int
	for idx, r := range s {
		if r >= 0x0300 && r <= 0x036f { // Combining Diacritical Marks
			hasAccent = true
			accentIdx = idx
			break
		}
		if r > 127 {
			mapped, remove := removeAccentRune(r)
			if remove || mapped != r {
				hasAccent = true
				accentIdx = idx
				break
			}
		}
	}

	if !hasAccent {
		return s
	}

	var builder strings.Builder
	builder.Grow(len(s))
	builder.WriteString(s[:accentIdx])

	for _, r := range s[accentIdx:] {
		mapped, remove := removeAccentRune(r)
		if remove {
			continue
		}
		builder.WriteRune(mapped)
	}

	return builder.String()
}

func removeAccentRune(r rune) (rune, bool) {
	if r >= 0x0300 && r <= 0x036f {
		return 0, true
	}

	switch r {
	case 'á', 'à', 'ả', 'ã', 'ạ', 'ă', 'ắ', 'ằ', 'ẳ', 'ẵ', 'ặ', 'â', 'ấ', 'ầ', 'ẩ', 'ẫ', 'ậ':
		return 'a', false
	case 'Á', 'À', 'Ả', 'Ã', 'Ạ', 'Ă', 'Ắ', 'Ằ', 'Ẳ', 'Ẵ', 'Ặ', 'Â', 'Ấ', 'Ầ', 'Ẩ', 'Ẫ', 'Ậ':
		return 'A', false
	case 'é', 'è', 'ẻ', 'ẽ', 'ẹ', 'ê', 'ế', 'ề', 'ể', 'ễ', 'ệ':
		return 'e', false
	case 'É', 'È', 'Ẻ', 'Ẽ', 'Ẹ', 'Ê', 'Ế', 'Ề', 'Ể', 'Ễ', 'Ệ':
		return 'E', false
	case 'í', 'ì', 'ỉ', 'ĩ', 'ị':
		return 'i', false
	case 'Í', 'Ì', 'Ỉ', 'Ĩ', 'Ị':
		return 'I', false
	case 'ó', 'ò', 'ỏ', 'õ', 'ọ', 'ô', 'ố', 'ồ', 'ổ', 'ỗ', 'ộ', 'ơ', 'ớ', 'ờ', 'ở', 'ỡ', 'ợ':
		return 'o', false
	case 'Ó', 'Ò', 'Ỏ', 'Õ', 'Ọ', 'Ô', 'Ố', 'Ồ', 'Ổ', 'Ỗ', 'Ộ', 'Ơ', 'Ớ', 'Ờ', 'Ở', 'Ỡ', 'Ợ':
		return 'O', false
	case 'ú', 'ù', 'ủ', 'ũ', 'ụ', 'ư', 'ứ', 'ừ', 'ử', 'ữ', 'ự':
		return 'u', false
	case 'Ú', 'Ù', 'Ủ', 'Ũ', 'Ụ', 'Ư', 'Ứ', 'Ừ', 'Ử', 'Ữ', 'Ự':
		return 'U', false
	case 'ý', 'ỳ', 'ỷ', 'ỹ', 'ỵ':
		return 'y', false
	case 'Ý', 'Ỳ', 'Ỷ', 'Ỹ', 'Ỵ':
		return 'Y', false
	case 'đ':
		return 'd', false
	case 'Đ':
		return 'D', false
	}
	return r, false
}

// Slugify tạo một URL slug thân thiện từ chuỗi đầu vào.
// Chuyển toàn bộ thành chữ thường, loại bỏ dấu tiếng Việt và thay thế các ký tự không phải chữ/số bằng separator.
// Tham số sep là optional, mặc định sử dụng "-".
// Hàm được tối ưu hóa chỉ chạy 1 vòng lặp duy nhất và chỉ cấp phát 1 builder để đạt hiệu năng tối đa.
func Slugify(s string, sep ...string) string {
	if s == "" {
		return ""
	}

	separator := "-"
	if len(sep) > 0 {
		separator = sep[0]
	}

	var builder strings.Builder
	builder.Grow(len(s))

	inSeparator := false
	firstChar := false

	for _, r := range s {
		mapped, remove := removeAccentRune(r)
		if remove {
			continue
		}

		// Chuyển thành chữ thường trực tiếp
		if mapped >= 'A' && mapped <= 'Z' {
			mapped = mapped + ('a' - 'A')
		} else if mapped > 127 {
			mapped = unicode.ToLower(mapped)
		}

		isAlphanumeric := (mapped >= 'a' && mapped <= 'z') || (mapped >= '0' && mapped <= '9')

		if isAlphanumeric {
			if inSeparator && firstChar {
				builder.WriteString(separator)
				inSeparator = false
			}
			builder.WriteByte(byte(mapped))
			firstChar = true
		} else {
			if firstChar {
				inSeparator = true
			}
		}
	}

	return builder.String()
}

// IsFalse trả về true nếu chuỗi biểu thị giá trị giả trị (falsy) hoặc rỗng.
// Các giá trị falsy bao gồm (case-insensitive, tự động loại bỏ khoảng trắng ở đầu/cuối):
// "", "0", "f", "false", "no", "n", "off".
// Hàm sử dụng string slicing và so khớp tĩnh để đạt hiệu năng tối đa (0 allocation).
func IsFalse(s string) bool {
	start := 0
	end := len(s)
	for start < end && isSpace(s[start]) {
		start++
	}
	for end > start && isSpace(s[end-1]) {
		end--
	}

	trimmed := s[start:end]
	if len(trimmed) == 0 {
		return true
	}

	switch len(trimmed) {
	case 1:
		c := trimmed[0]
		return c == '0' || c == 'f' || c == 'F' || c == 'n' || c == 'N'
	case 2:
		// no
		return (trimmed[0] == 'n' || trimmed[0] == 'N') && (trimmed[1] == 'o' || trimmed[1] == 'O')
	case 3:
		// off
		return (trimmed[0] == 'o' || trimmed[0] == 'O') &&
			(trimmed[1] == 'f' || trimmed[1] == 'F') &&
			(trimmed[2] == 'f' || trimmed[2] == 'F')
	case 5:
		// false
		return (trimmed[0] == 'f' || trimmed[0] == 'F') &&
			(trimmed[1] == 'a' || trimmed[1] == 'A') &&
			(trimmed[2] == 'l' || trimmed[2] == 'L') &&
			(trimmed[3] == 's' || trimmed[3] == 'S') &&
			(trimmed[4] == 'e' || trimmed[4] == 'E')
	}

	return false
}

// IsTrue trả về true nếu chuỗi biểu thị giá trị chân trị (truthy).
// Hàm ngược lại với IsFalse. Bất kỳ chuỗi nào có giá trị khác falsy thì đều là truthy.
func IsTrue(s string) bool {
	return !IsFalse(s)
}
