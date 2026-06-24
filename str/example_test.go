package str_test

import (
	"fmt"
	"github.com/ducconit/gotoolkit/str"
)

func Example() {
	// 1. Xác thực tên hiển thị (Nickname/Full Name)
	validName := "Nguyễn Văn Anh"
	invalidName := "Nguyễn Văn Anh 123" // Chứa số là không hợp lệ

	fmt.Println("Tên hiển thị hợp lệ:", str.IsAlphaSpace(validName))
	fmt.Println("Tên hiển thị không hợp lệ:", str.IsAlphaSpace(invalidName))

	// 2. Xác thực Username hệ thống
	validUser := "user_name_123"
	invalidUser := "user-name" // Chứa gạch ngang là không hợp lệ

	fmt.Println("Username hợp lệ:", str.IsUsername(validUser))
	fmt.Println("Username không hợp lệ:", str.IsUsername(invalidUser))

	// 3. Phát hiện sớm mã độc (HTML/Script Injection)
	cleanComment := "Bài viết rất hay!"
	unsafeComment := "Chào bạn <script>alert('hack')</script>"

	fmt.Println("Bình luận sạch chứa HTML:", str.ContainsHTML(cleanComment))
	fmt.Println("Bình luận độc hại chứa HTML:", str.ContainsHTML(unsafeComment))

	// 4. Chuẩn hóa khoảng trắng
	dirtySpace := "   Nguyễn   Đức   Cường   "
	fmt.Printf("Chuẩn hóa khoảng trắng: %q\n", str.CleanSpace(dirtySpace))

	// 5. Loại bỏ dấu tiếng Việt
	accentText := "Nguyễn Đức Cường"
	fmt.Printf("Xóa dấu: %q\n", str.RemoveAccents(accentText))

	// 6. Tạo URL Slug thân thiện
	slugText := "Nguyễn Đức Cường - Lập trình viên!"
	fmt.Printf("Slug mặc định: %q\n", str.Slugify(slugText))
	fmt.Printf("Slug custom gạch dưới: %q\n", str.Slugify(slugText, "_"))

	// Output:
	// Tên hiển thị hợp lệ: true
	// Tên hiển thị không hợp lệ: false
	// Username hợp lệ: true
	// Username không hợp lệ: false
	// Bình luận sạch chứa HTML: false
	// Bình luận độc hại chứa HTML: true
	// Chuẩn hóa khoảng trắng: "Nguyễn Đức Cường"
	// Xóa dấu: "Nguyen Duc Cuong"
	// Slug mặc định: "nguyen-duc-cuong-lap-trinh-vien"
	// Slug custom gạch dưới: "nguyen_duc_cuong_lap_trinh_vien"
}
