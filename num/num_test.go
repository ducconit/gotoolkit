package num_test

import (
	"math"
	"testing"

	"github.com/ducconit/gotoolkit/num"
)

func TestToWords(t *testing.T) {
	tests := []struct {
		name    string
		val     int64
		lang    string
		want    string
		wantErr bool
	}{
		// --- TIẾNG ANH (en) ---
		{"en: zero", 0, "en", "zero", false},
		{"en: single digit", 9, "en", "nine", false},
		{"en: teen", 12, "en", "twelve", false},
		{"en: teen 19", 19, "en", "nineteen", false},
		{"en: tens", 20, "en", "twenty", false},
		{"en: tens with ones", 21, "en", "twenty-one", false},
		{"en: hundreds", 100, "en", "one hundred", false},
		{"en: hundreds with ones", 105, "en", "one hundred five", false},
		{"en: thousands", 1005, "en", "one thousand five", false},
		{"en: millions", 1000000, "en", "one million", false},
		{"en: negative number", -15, "en", "minus fifteen", false},
		{"en: MaxInt64", math.MaxInt64, "en", "nine quintillion two hundred twenty-three quadrillion three hundred seventy-two trillion thirty-six billion eight hundred fifty-four million seven hundred seventy-five thousand eight hundred seven", false},
		{"en: MinInt64", math.MinInt64, "en", "minus nine quintillion two hundred twenty-three quadrillion three hundred seventy-two trillion thirty-six billion eight hundred fifty-four million seven hundred seventy-five thousand eight hundred eight", false},

		// --- TIẾNG VIỆT (vi) ---
		{"vi: không", 0, "vi", "không", false},
		{"vi: năm", 5, "vi", "năm", false},
		{"vi: mười lăm", 15, "vi", "mười lăm", false},
		{"vi: hai mươi lăm", 25, "vi", "hai mươi lăm", false},
		{"vi: hai mươi mốt", 21, "vi", "hai mươi mốt", false},
		{"vi: mười một", 11, "vi", "mười một", false},
		{"vi: hai mươi tư", 24, "vi", "hai mươi tư", false},
		{"vi: mười bốn", 14, "vi", "mười bốn", false},
		{"vi: một trăm lẻ năm", 105, "vi", "một trăm lẻ năm", false},
		{"vi: một nghìn không trăm lẻ năm", 1005, "vi", "một nghìn không trăm lẻ năm", false},
		{"vi: một triệu không trăm lẻ năm", 1000005, "vi", "một triệu không trăm lẻ năm", false},
		{"vi: một triệu", 1000000, "vi", "một triệu", false},
		{"vi: một tỷ", 1000000000, "vi", "một tỷ", false},
		{"vi: một tỷ không trăm lẻ năm", 1000000005, "vi", "một tỷ không trăm lẻ năm", false},
		{"vi: số âm", -123, "vi", "âm một trăm hai mươi ba", false},
		{"vi: MaxInt64", math.MaxInt64, "vi", "chín tỷ tỷ hai trăm hai mươi ba triệu tỷ ba trăm bảy mươi hai nghìn tỷ không trăm ba mươi sáu tỷ tám trăm năm mươi tư triệu bảy trăm bảy mươi lăm nghìn tám trăm lẻ bảy", false},
		{"vi: MinInt64", math.MinInt64, "vi", "âm chín tỷ tỷ hai trăm hai mươi ba triệu tỷ ba trăm bảy mươi hai nghìn tỷ không trăm ba mươi sáu tỷ tám trăm năm mươi tư triệu bảy trăm bảy mươi lăm nghìn tám trăm lẻ tám", false},

		// --- ĐĂNG KÝ NGÔN NGỮ TÙY BIẾN / LỖI ---
		{"unsupported language", 100, "fr", "", true},
		{"case insensitive lang", 10, "EN", "ten", false},
	}

	// Đăng ký test trước khi chạy
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := num.ToWords(tt.val, tt.lang)
			if (err != nil) != tt.wantErr {
				t.Errorf("ToWords() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.wantErr {
				return
			}
			if got != tt.want {
				t.Errorf("ToWords() = %q, want %q", got, tt.want)
			}
		})
	}
}

// Test Dynamic Registration
func TestRegister(t *testing.T) {
	// Đăng ký ngôn ngữ tiếng Pháp "fr" giả lập
	num.Register("fr", func(val int64) string {
		if val == 1 {
			return "un"
		}
		if val == 2 {
			return "deux"
		}
		return "inconnu"
	})

	got, err := num.ToWords(1, "fr")
	if err != nil {
		t.Fatalf("Unexpected error for registered 'fr': %v", err)
	}
	if got != "un" {
		t.Errorf("Expected 'un', got %q", got)
	}

	got, err = num.ToWords(2, "FR") // case insensitive check
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
	if got != "deux" {
		t.Errorf("Expected 'deux', got %q", got)
	}
}

func TestToShorthand(t *testing.T) {
	tests := []struct {
		name string
		val  int64
		lang string
		want string
	}{
		// --- TIẾNG ANH (mặc định) ---
		{"en: dưới 1000", 999, "en", "999"},
		{"en: 1K", 1000, "en", "1K"},
		{"en: 1.5K", 1500, "en", "1.5K"},
		{"en: 1.55K", 1550, "en", "1.55K"},
		{"en: làm tròn 1.56K", 1556, "en", "1.56K"},
		{"en: 1M", 1000000, "en", "1M"},
		{"en: 1.25M", 1250000, "en", "1.25M"},
		{"en: 1B", 1000000000, "en", "1B"},
		{"en: 1T", 1000000000000, "en", "1T"},
		{"en: âm 1.5K", -1500, "en", "-1.5K"},
		{"en: MaxInt64", math.MaxInt64, "en", "9.22Qi"},
		{"en: MinInt64", math.MinInt64, "en", "-9.22Qi"},
		{"en: case-insensitive lang", 1500, "EN", "1.5K"},
		{"en: unsupported fallback", 1500, "unknown", "1.5K"},

		// --- TIẾNG VIỆT (vi) ---
		{"vi: dưới 1000", 999, "vi", "999"},
		{"vi: 1 nghìn", 1000, "vi", "1 nghìn"},
		{"vi: 1.5 nghìn", 1500, "vi", "1.5 nghìn"},
		{"vi: 1.55 nghìn", 1550, "vi", "1.55 nghìn"},
		{"vi: 1 triệu", 1000000, "vi", "1 triệu"},
		{"vi: 1.25 triệu", 1250000, "vi", "1.25 triệu"},
		{"vi: 1 tỷ", 1000000000, "vi", "1 tỷ"},
		{"vi: 1 nghìn tỷ", 1000000000000, "vi", "1 nghìn tỷ"},
		{"vi: 1 triệu tỷ", 1000000000000000, "vi", "1 triệu tỷ"},
		{"vi: 1 tỷ tỷ", 1000000000000000000, "vi", "1 tỷ tỷ"},
		{"vi: âm 1.5 nghìn", -1500, "vi", "-1.5 nghìn"},
		{"vi: MaxInt64", math.MaxInt64, "vi", "9.22 tỷ tỷ"},
		{"vi: MinInt64", math.MinInt64, "vi", "-9.22 tỷ tỷ"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := num.ToShorthand(tt.val, tt.lang)
			if got != tt.want {
				t.Errorf("ToShorthand() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestFormat(t *testing.T) {
	tests := []struct {
		name        string
		val         float64
		precision   int
		thousandSep string
		decimalSep  string
		want        string
	}{
		{"VN currency format", 1234567.89, 2, ".", ",", "1.234.567,89"},
		{"US currency format", 1234567.89, 2, ",", ".", "1,234,567.89"},
		{"Round up to integer", 1234567.89, 0, ".", ",", "1.234.568"},
		{"Round up to thousand", 999.9, 0, ",", ".", "1,000"},
		{"Space grouping format", 1234.5678, 3, " ", ".", "1 234.568"},
		{"Zero value", 0.0, 2, ",", ".", "0.00"},
		{"Negative value", -1234.56, 2, ",", ".", "-1,234.56"},
		{"No thousands sep", 123.4, 4, "", ".", "123.4000"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := num.Format(tt.val, tt.precision, tt.thousandSep, tt.decimalSep)
			if got != tt.want {
				t.Errorf("Format() = %q, want %q", got, tt.want)
			}
		})
	}
}

func BenchmarkToShorthand(b *testing.B) {
	b.Run("en", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_ = num.ToShorthand(1250000, "en")
		}
	})
	b.Run("vi", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_ = num.ToShorthand(1250000, "vi")
		}
	})
}

func BenchmarkFormat(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = num.Format(1234567.89, 2, ".", ",")
	}
}

func BenchmarkToWords(b *testing.B) {
	b.Run("en", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_, _ = num.ToWords(1234567, "en")
		}
	})
	b.Run("vi", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_, _ = num.ToWords(1234567, "vi")
		}
	})
}

