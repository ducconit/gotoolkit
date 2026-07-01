package random

import (
	"math"
	"strings"
	"testing"
	"unicode/utf8"
)

// TestString kiểm tra hàm String với các bộ ký tự khác nhau
func TestString(t *testing.T) {
	t.Run("Length <= 0", func(t *testing.T) {
		if got := String(0, "abc"); got != "" {
			t.Errorf("String(0) = %q, want empty string", got)
		}
		if got := String(-5, "abc"); got != "" {
			t.Errorf("String(-5) = %q, want empty string", got)
		}
	})

	t.Run("Charset empty", func(t *testing.T) {
		if got := String(10, ""); got != "" {
			t.Errorf("String(10, empty) = %q, want empty string", got)
		}
	})

	t.Run("Valid string generation", func(t *testing.T) {
		length := 20
		charset := "abc"
		got := String(length, charset)

		if len(got) != length {
			t.Errorf("String() length = %d, want %d", len(got), length)
		}

		for _, char := range got {
			if !strings.ContainsRune(charset, char) {
				t.Errorf("String() generated char %q not in charset %q", char, charset)
			}
		}
	})
}

// TestPresets kiểm tra các hàm preset sinh chuỗi ngẫu nhiên
func TestPresets(t *testing.T) {
	length := 15

	t.Run("Letters", func(t *testing.T) {
		got := Letters(length)
		if len(got) != length {
			t.Errorf("Letters() length = %d, want %d", len(got), length)
		}
		for _, char := range got {
			if !strings.ContainsRune(charsetLetters, char) {
				t.Errorf("Letters() generated invalid char %q", char)
			}
		}
	})

	t.Run("Lowercases", func(t *testing.T) {
		got := Lowercases(length)
		if len(got) != length {
			t.Errorf("Lowercases() length = %d, want %d", len(got), length)
		}
		for _, char := range got {
			if !strings.ContainsRune(charsetLowercases, char) {
				t.Errorf("Lowercases() generated invalid char %q", char)
			}
		}
	})

	t.Run("Uppercases", func(t *testing.T) {
		got := Uppercases(length)
		if len(got) != length {
			t.Errorf("Uppercases() length = %d, want %d", len(got), length)
		}
		for _, char := range got {
			if !strings.ContainsRune(charsetUppercases, char) {
				t.Errorf("Uppercases() generated invalid char %q", char)
			}
		}
	})

	t.Run("Numbers", func(t *testing.T) {
		got := Numbers(length)
		if len(got) != length {
			t.Errorf("Numbers() length = %d, want %d", len(got), length)
		}
		for _, char := range got {
			if !strings.ContainsRune(charsetNumbers, char) {
				t.Errorf("Numbers() generated invalid char %q", char)
			}
		}
	})

	t.Run("UppercaseNumbers", func(t *testing.T) {
		got := UppercaseNumbers(length)
		if len(got) != length {
			t.Errorf("UppercaseNumbers() length = %d, want %d", len(got), length)
		}
		for _, char := range got {
			if !strings.ContainsRune(charsetUppercaseNumbers, char) {
				t.Errorf("UppercaseNumbers() generated invalid char %q", char)
			}
		}
	})

	t.Run("LowercaseNumbers", func(t *testing.T) {
		got := LowercaseNumbers(length)
		if len(got) != length {
			t.Errorf("LowercaseNumbers() length = %d, want %d", len(got), length)
		}
		for _, char := range got {
			if !strings.ContainsRune(charsetLowercaseNumbers, char) {
				t.Errorf("LowercaseNumbers() generated invalid char %q", char)
			}
		}
	})

	t.Run("Alphanumerics", func(t *testing.T) {
		got := Alphanumerics(length)
		if len(got) != length {
			t.Errorf("Alphanumerics() length = %d, want %d", len(got), length)
		}
		for _, char := range got {
			if !strings.ContainsRune(charsetAlphanumerics, char) {
				t.Errorf("Alphanumerics() generated invalid char %q", char)
			}
		}
	})

	t.Run("Vietnamese", func(t *testing.T) {
		got := Vietnamese(length)
		if utf8.RuneCountInString(got) != length {
			t.Errorf("Vietnamese() rune count = %d, want %d", utf8.RuneCountInString(got), length)
		}
		for _, r := range got {
			if !strings.ContainsRune(charsetVietnamese, r) {
				t.Errorf("Vietnamese() generated invalid rune %q", r)
			}
		}
	})

	t.Run("VietnameseLowercases", func(t *testing.T) {
		got := VietnameseLowercases(length)
		if utf8.RuneCountInString(got) != length {
			t.Errorf("VietnameseLowercases() rune count = %d, want %d", utf8.RuneCountInString(got), length)
		}
		for _, r := range got {
			if !strings.ContainsRune(charsetVietnameseLowercases, r) {
				t.Errorf("VietnameseLowercases() generated invalid rune %q", r)
			}
		}
	})

	t.Run("VietnameseUppercases", func(t *testing.T) {
		got := VietnameseUppercases(length)
		if utf8.RuneCountInString(got) != length {
			t.Errorf("VietnameseUppercases() rune count = %d, want %d", utf8.RuneCountInString(got), length)
		}
		for _, r := range got {
			if !strings.ContainsRune(charsetVietnameseUppercases, r) {
				t.Errorf("VietnameseUppercases() generated invalid rune %q", r)
			}
		}
	})

	t.Run("VietnameseNumbers", func(t *testing.T) {
		got := VietnameseNumbers(length)
		if utf8.RuneCountInString(got) != length {
			t.Errorf("VietnameseNumbers() rune count = %d, want %d", utf8.RuneCountInString(got), length)
		}
		for _, r := range got {
			if !strings.ContainsRune(charsetVietnameseNumbers, r) {
				t.Errorf("VietnameseNumbers() generated invalid rune %q", r)
			}
		}
	})

	t.Run("Runes", func(t *testing.T) {
		customRunes := []rune("💎🔥⭐")
		got := Runes(length, customRunes)
		if utf8.RuneCountInString(got) != length {
			t.Errorf("Runes() rune count = %d, want %d", utf8.RuneCountInString(got), length)
		}
		for _, r := range got {
			if !strings.ContainsRune("💎🔥⭐", r) {
				t.Errorf("Runes() generated invalid rune %q", r)
			}
		}
	})

	t.Run("RuneString", func(t *testing.T) {
		got := RuneString(length, "aàá")
		if utf8.RuneCountInString(got) != length {
			t.Errorf("RuneString() rune count = %d, want %d", utf8.RuneCountInString(got), length)
		}
		for _, r := range got {
			if !strings.ContainsRune("aàá", r) {
				t.Errorf("RuneString() generated invalid rune %q", r)
			}
		}
	})

	t.Run("Any", func(t *testing.T) {
		got := Any(length)
		if len(got) != length {
			t.Errorf("Any() length = %d, want %d", len(got), length)
		}
		for _, char := range got {
			if char < 33 || char > 126 {
				t.Errorf("Any() generated invalid ASCII char %q", char)
			}
		}
	})
}

// TestOTP kiểm tra sinh OTP
func TestOTP(t *testing.T) {
	t.Run("OTP Int", func(t *testing.T) {
		for i := 0; i < 1000; i++ {
			got := OTP()
			if got < 100000 || got > 999999 {
				t.Errorf("OTP() = %d, want in range [100000, 999999]", got)
			}
		}
	})

	t.Run("OTP String", func(t *testing.T) {
		for i := 0; i < 1000; i++ {
			got := OTPString()
			if len(got) != 6 {
				t.Errorf("OTPString() length = %d, want 6", len(got))
			}
			for _, char := range got {
				if char < '0' || char > '9' {
					t.Errorf("OTPString() generated non-numeric char %q", char)
				}
			}
		}
	})
}

// TestIntRange kiểm tra hàm IntRange
func TestIntRange(t *testing.T) {
	t.Run("Valid range", func(t *testing.T) {
		min, max := 10, 50
		for i := 0; i < 1000; i++ {
			got := IntRange(min, max)
			if got < min || got > max {
				t.Errorf("IntRange(%d, %d) = %d, want in range [%d, %d]", min, max, got, min, max)
			}
		}
	})

	t.Run("Min > Max swapped", func(t *testing.T) {
		min, max := 50, 10
		for i := 0; i < 1000; i++ {
			got := IntRange(min, max)
			if got < 10 || got > 50 {
				t.Errorf("IntRange(%d, %d) swapped = %d, want in range [10, 50]", min, max, got)
			}
		}
	})

	t.Run("Min == Max", func(t *testing.T) {
		if got := IntRange(5, 5); got != 5 {
			t.Errorf("IntRange(5, 5) = %d, want 5", got)
		}
	})

	t.Run("Avoid Overflow (math.MinInt to math.MaxInt)", func(t *testing.T) {
		min, max := math.MinInt, math.MaxInt
		for i := 0; i < 1000; i++ {
			got := IntRange(min, max)
			// Không được panic và phải trả về số nguyên hợp lệ
			if got < min || got > max {
				t.Errorf("IntRange(%d, %d) overflow check failed: %d", min, max, got)
			}
		}
	})
}

// TestFloatRange kiểm tra hàm FloatRange
func TestFloatRange(t *testing.T) {
	t.Run("Valid range", func(t *testing.T) {
		min, max := 1.5, 5.5
		for i := 0; i < 1000; i++ {
			got := FloatRange(min, max)
			if got < min || got >= max {
				t.Errorf("FloatRange(%f, %f) = %f, want in range [%f, %f)", min, max, got, min, max)
			}
		}
	})

	t.Run("Min > Max swapped", func(t *testing.T) {
		min, max := 5.5, 1.5
		for i := 0; i < 1000; i++ {
			got := FloatRange(min, max)
			if got < 1.5 || got >= 5.5 {
				t.Errorf("FloatRange(%f, %f) swapped = %f, want in range [1.5, 5.5)", min, max, got)
			}
		}
	})

	t.Run("Min == Max", func(t *testing.T) {
		if got := FloatRange(5.5, 5.5); got != 5.5 {
			t.Errorf("FloatRange(5.5, 5.5) = %f, want 5.5", got)
		}
	})
}

// TestHash32 kiểm tra hàm băm FNV-1a Hash32
func TestHash32(t *testing.T) {
	t.Run("Deterministic results", func(t *testing.T) {
		key := "user-12345"
		hash1 := Hash32(key)
		hash2 := Hash32(key)
		if hash1 != hash2 {
			t.Errorf("Hash32() should be deterministic, got %d and %d", hash1, hash2)
		}
	})

	t.Run("Different keys different hashes", func(t *testing.T) {
		hash1 := Hash32("user-1")
		hash2 := Hash32("user-2")
		if hash1 == hash2 {
			t.Errorf("Hash32() should yield different results for different inputs, got same hash %d", hash1)
		}
	})
}

// === BENCHMARKS ===

func BenchmarkString_Len10(b *testing.B) {
	for b.Loop() {
		_ = String(10, charsetLetters)
	}
}

func BenchmarkString_Len100(b *testing.B) {
	for b.Loop() {
		_ = String(100, charsetLetters)
	}
}

func BenchmarkLetters_Len32(b *testing.B) {
	for b.Loop() {
		_ = Letters(32)
	}
}

func BenchmarkNumbers_Len10(b *testing.B) {
	for b.Loop() {
		_ = Numbers(10)
	}
}

func BenchmarkOTP(b *testing.B) {
	for b.Loop() {
		_ = OTP()
	}
}

func BenchmarkOTPString(b *testing.B) {
	for b.Loop() {
		_ = OTPString()
	}
}

func BenchmarkIntRange(b *testing.B) {
	for b.Loop() {
		_ = IntRange(10, 100)
	}
}

func BenchmarkFloatRange(b *testing.B) {
	for b.Loop() {
		_ = FloatRange(1.5, 100.5)
	}
}

func BenchmarkAlphanumerics_Len32(b *testing.B) {
	for b.Loop() {
		_ = Alphanumerics(32)
	}
}

func BenchmarkVietnamese_Len32(b *testing.B) {
	for b.Loop() {
		_ = Vietnamese(32)
	}
}

func BenchmarkAny_Len32(b *testing.B) {
	for b.Loop() {
		_ = Any(32)
	}
}

func BenchmarkHash32(b *testing.B) {
	key := "some-random-key-to-hash-for-bucket-evaluation"
	b.ResetTimer()
	for b.Loop() {
		_ = Hash32(key)
	}
}
