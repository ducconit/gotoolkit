// Package random cung cấp các tiện ích sinh chuỗi và số ngẫu nhiên tối giản,
// hiệu năng cao, thread-safe và tiết kiệm tài nguyên bằng cách tận dụng math/rand/v2 (Go 1.22+).
package random

import (
	"math/rand/v2"
	"unsafe"
)

const (
	charsetLetters          = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	charsetLowercases       = "abcdefghijklmnopqrstuvwxyz"
	charsetUppercases       = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	charsetNumbers          = "0123456789"
	charsetUppercaseNumbers = "ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	charsetLowercaseNumbers = "abcdefghijklmnopqrstuvwxyz0123456789"
)

// String sinh một chuỗi ngẫu nhiên có độ dài length từ bộ ký tự charset cho trước.
// Trả về chuỗi rỗng nếu length <= 0 hoặc charset rỗng.
// Giải thuật tối ưu hóa sử dụng Bitmasking để giảm thiểu số lần gọi generator ngẫu nhiên,
// và unsafe.String để đạt hiệu năng 1 allocation duy nhất.
func String(length int, charset string) string {
	if length <= 0 || charset == "" {
		return ""
	}

	b := make([]byte, length)
	charsetLen := len(charset)

	// Tính số bit cần thiết để biểu diễn chỉ số của charset.
	// Ví dụ: charsetLen = 62, ta cần 6 bit (2^6 = 64).
	var letterIdxBits uint = 1
	for 1<<letterIdxBits < charsetLen {
		letterIdxBits++
	}
	var letterIdxMask uint64 = 1<<letterIdxBits - 1
	letterIdxMax := 63 / letterIdxBits // Số lượng chỉ số có thể trích xuất từ 1 số uint64

	for i, cache, remain := 0, rand.Uint64(), letterIdxMax; i < length; {
		if remain == 0 {
			cache, remain = rand.Uint64(), letterIdxMax
		}
		if idx := int(cache & letterIdxMask); idx < charsetLen {
			b[i] = charset[idx]
			i++
		}
		cache >>= letterIdxBits
		remain--
	}

	return unsafe.String(unsafe.SliceData(b), len(b))
}

// Letters sinh chuỗi ngẫu nhiên chỉ chứa chữ cái (cả hoa và thường) với độ dài cho trước.
func Letters(length int) string {
	return String(length, charsetLetters)
}

// Lowercases sinh chuỗi ngẫu nhiên chỉ chứa chữ cái viết thường với độ dài cho trước.
func Lowercases(length int) string {
	return String(length, charsetLowercases)
}

// Uppercases sinh chuỗi ngẫu nhiên chỉ chứa chữ cái viết hoa với độ dài cho trước.
func Uppercases(length int) string {
	return String(length, charsetUppercases)
}

// Numbers sinh chuỗi ngẫu nhiên chỉ chứa các ký số (0-9) với độ dài cho trước.
func Numbers(length int) string {
	return String(length, charsetNumbers)
}

// UppercaseNumbers sinh chuỗi ngẫu nhiên chỉ chứa chữ cái viết hoa và ký số với độ dài cho trước.
func UppercaseNumbers(length int) string {
	return String(length, charsetUppercaseNumbers)
}

// LowercaseNumbers sinh chuỗi ngẫu nhiên chỉ chứa chữ cái viết thường và ký số với độ dài cho trước.
func LowercaseNumbers(length int) string {
	return String(length, charsetLowercaseNumbers)
}

// Any sinh chuỗi ngẫu nhiên có độ dài cho trước chứa bất kỳ ký tự ASCII in được
// không bao gồm khoảng trắng (từ '!' đến '~', mã ASCII 33-126), bao gồm cả chữ, số và ký tự đặc biệt.
// Sử dụng Bitmasking tối ưu hóa tối đa và unsafe.String để đạt đúng 1 allocation duy nhất.
func Any(length int) string {
	if length <= 0 {
		return ""
	}

	b := make([]byte, length)
	const charsetLen = 94 // 126 - 33 + 1
	const letterIdxBits = 7 // 1<<7 = 128 >= 94
	const letterIdxMask = 1<<letterIdxBits - 1
	const letterIdxMax = 63 / letterIdxBits // 9 chỉ số trên mỗi uint64

	for i, cache, remain := 0, rand.Uint64(), letterIdxMax; i < length; {
		if remain == 0 {
			cache, remain = rand.Uint64(), letterIdxMax
		}
		if idx := int(cache & letterIdxMask); idx < charsetLen {
			b[i] = byte(idx + 33)
			i++
		}
		cache >>= letterIdxBits
		remain--
	}

	return unsafe.String(unsafe.SliceData(b), len(b))
}

// OTP sinh mã OTP gồm 6 chữ số dưới dạng số nguyên (từ 100000 đến 999999).
// Thao tác này là lock-free, zero allocation và có hiệu năng tối đa.
func OTP() int {
	return int(rand.Int32N(900000) + 100000)
}

// OTPString sinh mã OTP gồm 6 chữ số dưới dạng string (ví dụ: "012345", "999999").
// Tối ưu hóa vượt trội bằng cách chỉ gọi sinh số ngẫu nhiên 1 lần duy nhất và tính toán toán học,
// loại bỏ hoàn toàn vòng lặp, chỉ tốn 1 allocation duy nhất.
func OTPString() string {
	val := rand.Uint32N(1000000)
	b := make([]byte, 6)
	b[0] = '0' + byte(val/100000)
	b[1] = '0' + byte((val/10000)%10)
	b[2] = '0' + byte((val/1000)%10)
	b[3] = '0' + byte((val/100)%10)
	b[4] = '0' + byte((val/10)%10)
	b[5] = '0' + byte(val%10)
	return unsafe.String(unsafe.SliceData(b), 6)
}

// IntRange sinh một số nguyên ngẫu nhiên trong khoảng đóng [min, max].
// Nếu min > max, hai giá trị sẽ được hoán đổi cho nhau.
// Đã được cải tiến để tránh tràn số (overflow) khi khoảng cách giữa min và max quá lớn.
func IntRange(min, max int) int {
	if min == max {
		return min
	}
	if min > max {
		min, max = max, min
	}

	// Tránh tràn số bằng cách tính toán khoảng cách dưới dạng uint64
	diff := uint64(max) - uint64(min)
	if diff == ^uint64(0) { // Trường hợp bao phủ toàn bộ dải uint64
		return int(uint64(min) + rand.Uint64())
	}
	return int(uint64(min) + rand.Uint64N(diff+1))
}

// FloatRange sinh một số thập phân ngẫu nhiên trong khoảng bán mở [min, max).
// Nếu min > max, hai giá trị sẽ được hoán đổi cho nhau.
func FloatRange(min, max float64) float64 {
	if min == max {
		return min
	}
	if min > max {
		min, max = max, min
	}
	return min + rand.Float64()*(max-min)
}

// Hash32 tính toán mã băm FNV-1a 32-bit của chuỗi key.
// Đây là hàm băm phi mã hóa (non-cryptographic), zero-allocation và lock-free,
// lý tưởng cho việc phân nhóm ngẫu nhiên ổn định (deterministic bucket rollout / A/B testing).
func Hash32(key string) uint32 {
	const (
		offset32 = 2166136261
		prime32  = 16777619
	)
	hash := uint32(offset32)
	for i := range len(key) {
		hash ^= uint32(key[i])
		hash *= prime32
	}
	return hash
}
