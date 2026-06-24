package num_test

import (
	"fmt"

	"github.com/ducconit/gotoolkit/num"
)

func ExampleToWords() {
	// 1. Chuyển đổi số sang tiếng Việt (mặc định)
	strVI, err := num.ToWords(1005, "vi")
	if err != nil {
		fmt.Println("Lỗi:", err)
		return
	}
	fmt.Println("Tiếng Việt:", strVI)

	// 2. Chuyển đổi số sang tiếng Anh (mặc định)
	strEN, err := num.ToWords(1005, "en")
	if err != nil {
		fmt.Println("Lỗi:", err)
		return
	}
	fmt.Println("Tiếng Anh:", strEN)

	// Output:
	// Tiếng Việt: một nghìn không trăm lẻ năm
	// Tiếng Anh: one thousand five
}

func ExampleRegister() {
	// Đăng ký ngôn ngữ mới (ví dụ: Tiếng Pháp - "fr")
	num.Register("fr", func(val int64) string {
		switch val {
		case 1:
			return "un"
		case 2:
			return "deux"
		case 3:
			return "trois"
		default:
			return "inconnu"
		}
	})

	// Sử dụng ngôn ngữ mới vừa đăng ký
	strFR, err := num.ToWords(3, "fr")
	if err != nil {
		fmt.Println("Lỗi:", err)
		return
	}
	fmt.Println("Tiếng Pháp (3):", strFR)

	// Output:
	// Tiếng Pháp (3): trois
}

func ExampleToShorthand() {
	// 1. Dạng viết tắt tiếng Anh (mặc định quốc tế)
	fmt.Println("en (1500):", num.ToShorthand(1500, "en"))
	fmt.Println("en (1250000):", num.ToShorthand(1250000, "en"))

	// 2. Dạng viết tắt tiếng Việt
	fmt.Println("vi (1500):", num.ToShorthand(1500, "vi"))
	fmt.Println("vi (1250000):", num.ToShorthand(1250000, "vi"))

	// Output:
	// en (1500): 1.5K
	// en (1250000): 1.25M
	// vi (1500): 1.5 nghìn
	// vi (1250000): 1.25 triệu
}

func ExampleFormat() {
	val := 1234567.89

	// Định dạng kiểu Việt Nam: 1.234.567,89
	fmt.Println("vi:", num.Format(val, 2, ".", ","))

	// Định dạng kiểu Mỹ/Anh: 1,234,567.89
	fmt.Println("en:", num.Format(val, 2, ",", "."))

	// Định dạng kiểu khoảng trắng: 1 234 567,89
	fmt.Println("space:", num.Format(val, 2, " ", ","))

	// Output:
	// vi: 1.234.567,89
	// en: 1,234,567.89
	// space: 1 234 567,89
}
