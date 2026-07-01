# Res Package

`res` nằm trong thư mục con `httputil/res`, cung cấp giải pháp định dạng JSON response chuẩn REST API tối giản, hiệu năng cao, zero-dependency, tương thích tốt với thư viện tiêu chuẩn `net/http` và đồng thời có khả năng tích hợp linh hoạt với bất kỳ bên thứ ba nào (như Gin framework).

## Các Tính năng Chính

*   **Envelope Chuẩn REST API Tối Giản**: Đưa mọi dữ liệu response về cấu trúc thống nhất gồm: `code`, `msg`, `data`, `extra`.
*   **Cấu Hình Toàn Cục Linh Hoạt**: Cung cấp các biến toàn cục cho cả **mã Code** và **câu Message** mặc định. Ứng dụng dễ dàng ghi đè cấu hình tại hàm `main()` hoặc `init()` trước khi khởi chạy server.
*   **Các Helper Phân Trang Tường Minh**:
    *   **Phân trang thông thường** (`Paginate`): Chỉ chứa `total` trong `extra` (Ví dụ: `WritePaginate`, `GinPaginate`).
    *   **Phân trang Cursor** (`PaginateCursor`): Chứa cả `total` và `cursor` trong `extra` (Ví dụ: `WritePaginateCursor`, `GinPaginateCursor`).
    *   Lược bỏ hoàn toàn `limit` và `offset` để tối ưu hóa hiệu năng cơ sở dữ liệu lớn.
*   **Hỗ Trợ Lỗi Validation Chuyên Biệt**: Tự động trả về HTTP Status `422`, code `"VALIDATION_FAILED"`, `extra` là map lỗi của các trường và `msg` tùy chọn.
*   **Tích Hợp Sẵn Bộ Helper Lỗi Phổ Biến**:
    *   Hỗ trợ đầy đủ các mã lỗi: `BadRequest` (400), `Unauthorized` (401), `Forbidden` (403), `NotFound` (404), `Conflict` (409), `PageExpired` (419), `RateLimit` (429), `InternalServerError` (500), `Maintenance` (503).
    *   Cung cấp các hàm riêng biệt cho net/http (tiền tố `Write`) và Gin (tiền tố `Gin`).
*   **Tương Thích Gin Tuyệt Đối**: Cung cấp interface `GinContext` duck-typing giúp truyền thẳng `*gin.Context` vào các hàm `res.GinSuccess` hay `res.GinError` mà không cần import thư viện Gin bên ngoài.
*   **Hiệu Năng Cao**: Ghi trực tiếp (encode) JSON response vào luồng stream của `http.ResponseWriter` giúp tối giản allocations.

## Hướng dẫn Sử dụng Nhanh

### 1. Dùng với net/http
```go
package main

import (
	"net/http"
	"github.com/ducconit/gotoolkit/httputil/res"
)

func UserHandler(w http.ResponseWriter, r *http.Request) {
	user := map[string]string{"id": "1", "email": "user@example.com"}

	// Thành công (HTTP 200, code "0", msg "Success")
	res.WriteSuccess(w, user)
}
```

### 2. Dùng với Gin và Trả về Phân trang (Cursor & Offset)
```go
package main

import (
	"github.com/gin-gonic/gin"
	"github.com/ducconit/gotoolkit/httputil/res"
)

func ListUsers(c *gin.Context) {
	users := []string{"user1", "user2"}
	
	// 1. Phân trang thông thường (chỉ chứa total)
	res.GinPaginate(c, users, 100)

	// 2. Hoặc phân trang bằng cursor (chứa total & cursor)
	// res.GinPaginateCursor(c, users, 100, "next_cursor_token")
}
```

## Báo cáo Benchmark

Đo thực tế trên CPU **AMD Ryzen 7 8745H**:

| Hàm Benchmark | Thời gian thực thi (ns/op) | Dung lượng bộ nhớ (B/op) | Allocations (allocs/op) |
| :--- | :---: | :---: | :---: |
| `BenchmarkWriteSuccess` | 481.30 ns | 224 B | 7 allocs |
