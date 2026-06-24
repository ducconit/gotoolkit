package str

import (
	"testing"
)

// TestIsAlphaSpace kiểm tra tính chính xác của hàm IsAlphaSpace
func TestIsAlphaSpace(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  bool
	}{
		{"Empty string", "", false},
		{"English letters and spaces", "Nguyen Van A", true},
		{"Vietnamese letters and spaces", "Nguyễn Văn Anh", true},
		{"Contains numbers", "Nguyen Van A 123", false},
		{"Contains special chars", "Nguyen Van A!", false},
		{"Only spaces", "   ", true},
		{"Russian (Cyrillic)", "Иван Иванов", true},
		{"Chinese (Han)", "张伟", true},
		{"Japanese (Kanji & Kana & Space)", "山田 たろう", true},
		{"Korean (Hangul)", "김철수", true},
		{"Thai", "สมชาย", true},
		{"Hindi (Devanagari)", "आरव", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsAlphaSpace(tt.input); got != tt.want {
				t.Errorf("IsAlphaSpace(%q) = %v, want %v", tt.input, got, tt.want)
			}
		})
	}
}

// TestIsUsername kiểm tra tính chính xác của hàm IsUsername
func TestIsUsername(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  bool
	}{
		{"Empty string", "", false},
		{"Valid username lowercase", "user_name_123", true},
		{"Contains uppercase", "User_name", false},
		{"Contains special chars", "user-name", false},
		{"Contains space", "user name", false},
		{"Contains unicode letters", "nguyễn_văn", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsUsername(tt.input); got != tt.want {
				t.Errorf("IsUsername(%q) = %v, want %v", tt.input, got, tt.want)
			}
		})
	}
}

// TestContainsHTML kiểm tra khả năng phát hiện HTML/JS Injection
func TestContainsHTML(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  bool
	}{
		{"Normal text", "Hello World", false},
		{"Empty string", "", false},
		{"Simple HTML tag", "Hello <b>World</b>", true},
		{"Script tag", "<script>alert(1)</script>", true},
		{"Javascript protocol", "javascript:alert(1)", true},
		{"Javascript protocol mixed case", "JaVaScRiPt:alert(1)", true},
		{"Incomplete tag", "Hello <script", false}, // Không có ngoặc đóng, tạm thời coi là an toàn
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ContainsHTML(tt.input); got != tt.want {
				t.Errorf("ContainsHTML(%q) = %v, want %v", tt.input, got, tt.want)
			}
		})
	}
}

// === BENCHMARKS ===

func BenchmarkIsAlphaSpace_English(b *testing.B) {
	s := "Nguyen Van A"
	b.ResetTimer()
	for b.Loop() {
		_ = IsAlphaSpace(s)
	}
}

func BenchmarkIsAlphaSpace_Vietnamese(b *testing.B) {
	s := "Nguyễn Văn Anh"
	b.ResetTimer()
	for b.Loop() {
		_ = IsAlphaSpace(s)
	}
}

func BenchmarkIsUsername(b *testing.B) {
	s := "user_name_123"
	b.ResetTimer()
	for b.Loop() {
		_ = IsUsername(s)
	}
}

func BenchmarkContainsHTML_Safe(b *testing.B) {
	s := "Hello World, this is a clean text username."
	b.ResetTimer()
	for b.Loop() {
		_ = ContainsHTML(s)
	}
}

func BenchmarkContainsHTML_Unsafe(b *testing.B) {
	s := "Hello <b>World</b>, this is <script>alert(1)</script>"
	b.ResetTimer()
	for b.Loop() {
		_ = ContainsHTML(s)
	}
}

// TestCleanSpace kiểm tra tính chính xác của hàm CleanSpace
func TestCleanSpace(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		{"Empty string", "", ""},
		{"Clean string", "Hello World", "Hello World"},
		{"Leading spaces", "   Hello World", "Hello World"},
		{"Trailing spaces", "Hello World   ", "Hello World"},
		{"Multiple inner spaces", "Hello   World   Go", "Hello World Go"},
		{"Tabs and Newlines", "\tHello \n World\r", "Hello World"},
		{"Complex spaces mix", " \t  Hello \n\t   World   \r\n", "Hello World"},
		{"Unicode string clean", "  Nguyễn   Văn   Anh  ", "Nguyễn Văn Anh"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := CleanSpace(tt.input); got != tt.want {
				t.Errorf("CleanSpace(%q) = %q, want %q", tt.input, got, tt.want)
			}
		})
	}
}

// TestRemoveAccents kiểm tra tính chính xác của hàm RemoveAccents
func TestRemoveAccents(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		{"Empty string", "", ""},
		{"English (No accents)", "Hello World 123!", "Hello World 123!"},
		{"Vietnamese precomposed lower", "áàảãạăắằẳẵặâấầẩẫậđ", "aaaaaaaaaaaaaaaaad"},
		{"Vietnamese precomposed upper", "ÁÀẢÃẠĂẮẰẲẴẶÂẤẦẨẪẬĐ", "AAAAAAAAAAAAAAAAAD"},
		{"Vietnamese sentence", "Nguyễn Đức Cường", "Nguyen Duc Cuong"},
		{"Vietnamese decomposed (combining marks)", "Nguyễn Đức Cường", "Nguyen Duc Cuong"}, // Chuỗi chứa tone marks tổ hợp (sửa Nguŷễn -> Nguyễn)
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := RemoveAccents(tt.input); got != tt.want {
				t.Errorf("RemoveAccents(%q) = %q, want %q", tt.input, got, tt.want)
			}
		})
	}
}

// TestSlugify kiểm tra tính chính xác của hàm Slugify
func TestSlugify(t *testing.T) {
	tests := []struct {
		name  string
		input string
		sep   []string
		want  string
	}{
		{"Empty string", "", nil, ""},
		{"Normal English", "Hello World", nil, "hello-world"},
		{"Multiple spaces and uppercase", "  Hello   World  ", nil, "hello-world"},
		{"Vietnamese with accents", "Nguyễn Đức Cường @123", nil, "nguyen-duc-cuong-123"},
		{"Custom separator", "Nguyễn Đức Cường", []string{"_"}, "nguyen_duc_cuong"},
		{"Empty custom separator", "Hello World", []string{""}, "helloworld"},
		{"Leading & Trailing special chars", "---Hello World---", nil, "hello-world"},
		{"Multiple custom separator", "Hello World", []string{"/*/"}, "hello/*/world"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Slugify(tt.input, tt.sep...); got != tt.want {
				t.Errorf("Slugify(%q, %v) = %q, want %q", tt.input, tt.sep, got, tt.want)
			}
		})
	}
}

// === BENCHMARKS NEW ===

func BenchmarkCleanSpace_Clean(b *testing.B) {
	s := "Hello World Go"
	b.ResetTimer()
	for b.Loop() {
		_ = CleanSpace(s)
	}
}

func BenchmarkCleanSpace_Dirty(b *testing.B) {
	s := "   Hello   World   Go   "
	b.ResetTimer()
	for b.Loop() {
		_ = CleanSpace(s)
	}
}

func BenchmarkRemoveAccents_NoAccents(b *testing.B) {
	s := "Hello World Go, this is a normal string."
	b.ResetTimer()
	for b.Loop() {
		_ = RemoveAccents(s)
	}
}

func BenchmarkRemoveAccents_WithAccents(b *testing.B) {
	s := "Nguyễn Đức Cường - Lập trình viên Golang hiệu năng cao."
	b.ResetTimer()
	for b.Loop() {
		_ = RemoveAccents(s)
	}
}

func BenchmarkSlugify_Default(b *testing.B) {
	s := "Nguyễn Đức Cường - Lập trình viên Golang hiệu năng cao."
	b.ResetTimer()
	for b.Loop() {
		_ = Slugify(s)
	}
}

// TestIsFalse kiểm tra tính đúng đắn của hàm IsFalse
func TestIsFalse(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  bool
	}{
		{"Empty string", "", true},
		{"Spaces only", "   \t\n  ", true},
		{"False lowercase", "false", true},
		{"False uppercase", "FALSE", true},
		{"False mixed case", "fAlSe", true},
		{"False with spaces", "  false  ", true},
		{"Zero", "0", true},
		{"Char f", "f", true},
		{"Char F", "F", true},
		{"Char n", "n", true},
		{"Char N", "N", true},
		{"No", "no", true},
		{"No uppercase", "NO", true},
		{"Off", "off", true},
		{"Off uppercase", "OFF", true},
		{"True (not false)", "true", false},
		{"One (not false)", "1", false},
		{"Random string (not false)", "hello", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsFalse(tt.input); got != tt.want {
				t.Errorf("IsFalse(%q) = %v, want %v", tt.input, got, tt.want)
			}
		})
	}
}

// TestIsTrue kiểm tra tính đúng đắn của hàm IsTrue
func TestIsTrue(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  bool
	}{
		{"Empty string", "", false},
		{"False", "false", false},
		{"True", "true", true},
		{"One", "1", true},
		{"Yes", "yes", true},
		{"Random string", "hello_world", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsTrue(tt.input); got != tt.want {
				t.Errorf("IsTrue(%q) = %v, want %v", tt.input, got, tt.want)
			}
		})
	}
}

// === BENCHMARKS ISFALSE/ISTRUE ===

func BenchmarkIsFalse_True(b *testing.B) {
	s := "  false  "
	b.ResetTimer()
	for b.Loop() {
		_ = IsFalse(s)
	}
}

func BenchmarkIsFalse_False(b *testing.B) {
	s := "this_is_a_random_truthful_string"
	b.ResetTimer()
	for b.Loop() {
		_ = IsFalse(s)
	}
}

