package random_test

import (
	"fmt"
	"github.com/ducconit/gotoolkit/random"
)

func Example() {
	// 1. Sinh chuỗi ngẫu nhiên với charset bất kỳ
	customStr := random.String(10, "XYZ123")
	fmt.Println("Custom String length:", len(customStr))

	// 2. Sinh chuỗi ngẫu nhiên theo các bộ ký tự định sẵn
	letters := random.Letters(16)
	fmt.Println("Letters length:", len(letters))

	numbers := random.Numbers(8)
	fmt.Println("Numbers length:", len(numbers))

	// 3. Sinh OTP 6 chữ số dạng int
	otpInt := random.OTP()
	fmt.Println("OTP Int is valid:", otpInt >= 100000 && otpInt <= 999999)

	// 4. Sinh OTP 6 chữ số dạng string
	otpStr := random.OTPString()
	fmt.Println("OTP String length:", len(otpStr))

	// Output:
	// Custom String length: 10
	// Letters length: 16
	// Numbers length: 8
	// OTP Int is valid: true
	// OTP String length: 6
}
