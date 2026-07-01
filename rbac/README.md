# RBAC Package

`rbac` cung cấp giải pháp phân quyền người dùng dựa trên vai trò (Role-Based Access Control) tối giản, hiệu năng cao, thread-safe tuyệt đối và tiết kiệm tài nguyên bộ nhớ cho Golang.

## Các Tính năng Chính

*   **Đọc Lock-Free Cực Nhanh**: Tận dụng kỹ thuật **Copy-on-Write (COW)** qua `atomic.Pointer` giúp các thao tác kiểm tra quyền (`HasPermission`, `HasAllPermission`) đạt tốc độ tối đa (dưới **1ns**), an toàn tuyệt đối khi chạy đồng thời với hàng ngàn goroutines.
*   **Hỗ Trợ Wildcard (`*`)**: Hỗ trợ ký tự đại diện `*` (admin). Khi wildcard được kích hoạt, hệ thống sẽ tự động giải phóng map bộ nhớ các quyền cũ để tiết kiệm RAM.
*   **Hỗ trợ Go Iterators hiện đại (Go 1.26+)**: Duyệt qua danh sách quyền của một vai trò cực kỳ sạch sẽ và hiệu quả thông qua `iter.Seq` của Go standard library (`Permissions()`).

## Hướng dẫn Sử dụng Nhanh

```go
package main

import (
	"fmt"
	"github.com/ducconit/gotoolkit/rbac"
)

func main() {
	// 1. Khởi tạo một Role mới kèm danh sách quyền
	userRole := rbac.NewRole("user", "article:read", "article:write")

	// 2. Kiểm tra quyền hạn (Lock-Free)
	fmt.Println("Can read:", userRole.HasPermission("article:read"))   // true
	fmt.Println("Can delete:", userRole.HasPermission("article:delete")) // false

	// 3. Thêm quyền (Thread-Safe COW)
	userRole.AddPermission("article:delete")
	fmt.Println("Can delete now:", userRole.HasPermission("article:delete")) // true

	// 4. Wildcard Role (Admin)
	adminRole := rbac.NewRole("admin", "*")
	fmt.Println("Admin can check anything:", adminRole.HasPermission("any:random:permission")) // true

	// 5. Duyệt qua danh sách quyền bằng Iterator
	fmt.Println("Permissions list:")
	for p := range userRole.Permissions() {
		fmt.Println("-", p)
	}
}
```

## Báo cáo Benchmark

Đo thực tế trên CPU **AMD Ryzen 7 8745H**:

| Hàm Benchmark | Thời gian thực thi (ns/op) | Dung lượng bộ nhớ (B/op) | Allocations (allocs/op) |
| :--- | :---: | :---: | :---: |
| `BenchmarkHasPermission_Wildcard` | 0.98 ns | 0 B | 0 allocs |
| `BenchmarkHasPermission_Small` | 6.89 ns | 0 B | 0 allocs |
| `BenchmarkHasPermission_Large` | 0.87 ns | 0 B | 0 allocs |
| `BenchmarkHasAllPermission_Success` | 14.80 ns | 0 B | 0 allocs |
| `BenchmarkAddPermission_COW` | 11.92 ns | 0 B | 0 allocs |
| `BenchmarkConcurrentReadWrite` | 3.67 ns | 0 B | 0 allocs |
| `BenchmarkNewRole_WithPermissions` | 158.20 ns | 312 B | 4 allocs |
| `BenchmarkAddPermission_Duplicate` | 18.23 ns | 0 B | 0 allocs |
| `BenchmarkAddPermission_NewPerm_Large` | 1235.00 ns | 2032 B | 11 allocs |
