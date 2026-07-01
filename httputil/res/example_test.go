package res_test

import (
	"fmt"
	"net/http/httptest"
	"github.com/ducconit/gotoolkit/httputil/res"
)

// Giả lập một struct User
type User struct {
	ID    int    `json:"id"`
	Email string `json:"email"`
}

// Giả lập một Gin Context Mock để demo
type GinContextMock struct {
	Code int
	Body any
}

func (g *GinContextMock) JSON(code int, obj any) {
	g.Code = code
	g.Body = obj
}

func Example() {
	// 1. Dùng với net/http (ghi trực tiếp vào http.ResponseWriter)
	rec := httptest.NewRecorder()
	user := User{ID: 1, Email: "user@example.com"}

	// Thành công (Status 200, code "0", msg mặc định "Success")
	_ = res.WriteSuccess(rec, user)
	fmt.Println("Status:", rec.Code)
	fmt.Println("Body:", rec.Body.String())

	// 2. Dùng với Gin Framework (ghi vào Gin Context qua interface)
	c := &GinContextMock{}
	res.GinSuccess(c, user, "Đăng nhập thành công")
	
	// Ép kiểu Envelope để kiểm tra
	env := c.Body.(res.Envelope)
	fmt.Println("Gin Status:", c.Code)
	fmt.Println("Gin Envelope Code:", env.Code)
	fmt.Println("Gin Envelope Msg:", env.Msg)

	// 3. Lỗi Validation (Status 422, code "VALIDATION_FAILED")
	recVal := httptest.NewRecorder()
	errs := map[string]string{
		"email": "Định dạng email không hợp lệ",
	}
	_ = res.WriteValidationError(recVal, errs)
	fmt.Println("Validation Status:", recVal.Code)
	fmt.Println("Validation Body:", recVal.Body.String())

	// 4. Ghi đè biến global message & code (ví dụ cấu hình lúc khởi chạy app)
	res.DefaultUnauthorizedCode = "40101"
	res.DefaultUnauthorizedMsg = "Unauthorized login"
	recAuth := httptest.NewRecorder()
	_ = res.WriteUnauthorized(recAuth, "") // Truyền rỗng để kích hoạt cấu hình mặc định mới
	fmt.Println("Auth Status:", recAuth.Code)
	fmt.Println("Auth Body:", recAuth.Body.String())

	// 5. Trả về dữ liệu phân trang cursor (có total và cursor)
	recPage := httptest.NewRecorder()
	users := []User{
		{ID: 1, Email: "user1@example.com"},
		{ID: 2, Email: "user2@example.com"},
	}
	_ = res.WritePaginateCursor(recPage, users, 100, "next_token_abc")
	fmt.Println("Page Status:", recPage.Code)
	fmt.Println("Page Body:", recPage.Body.String())

	// Output:
	// Status: 200
	// Body: {"code":"0","msg":"Success","data":{"id":1,"email":"user@example.com"},"extra":null}
	//
	// Gin Status: 200
	// Gin Envelope Code: 0
	// Gin Envelope Msg: Đăng nhập thành công
	// Validation Status: 422
	// Validation Body: {"code":"VALIDATION_FAILED","msg":"Dữ liệu đầu vào không hợp lệ","data":null,"extra":{"email":"Định dạng email không hợp lệ"}}
	//
	// Auth Status: 401
	// Auth Body: {"code":"40101","msg":"Unauthorized login","data":null,"extra":null}
	//
	// Page Status: 200
	// Page Body: {"code":"0","msg":"Success","data":[{"id":1,"email":"user1@example.com"},{"id":2,"email":"user2@example.com"}],"extra":{"total":100,"cursor":"next_token_abc"}}
}
