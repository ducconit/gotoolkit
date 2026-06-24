package rbac_test

import (
	"fmt"
	"github.com/ducconit/gotoolkit/rbac"
)

func ExampleRole() {
	// Khởi tạo một role mới
	userRole := rbac.NewRole("user", "post.read", "post.create")

	// Kiểm tra các quyền cơ bản
	fmt.Println("ID:", userRole.ID())
	fmt.Println("Có quyền post.create:", userRole.HasPermission("post.create"))
	fmt.Println("Có quyền post.delete:", userRole.HasPermission("post.delete"))

	// Thêm quyền động
	userRole.AddPermission("post.delete")
	fmt.Println("Có quyền post.delete sau khi thêm:", userRole.HasPermission("post.delete"))

	// Kiểm tra nhiều quyền cùng lúc
	fmt.Println("Có ít nhất 1 quyền (post.update, post.create):", userRole.HasPermission("post.update", "post.create"))
	fmt.Println("Có tất cả quyền (post.read, post.create):", userRole.HasAllPermission("post.read", "post.create"))

	// Role với quyền đặc biệt "*" đại diện cho tất cả
	adminRole := rbac.NewRole("admin", "*")
	fmt.Println("Admin có quyền post.publish:", adminRole.HasPermission("post.publish"))
	fmt.Println("Admin có tất cả các quyền (any, other):", adminRole.HasAllPermission("any", "other"))

	// Output:
	// ID: user
	// Có quyền post.create: true
	// Có quyền post.delete: false
	// Có quyền post.delete sau khi thêm: true
	// Có ít nhất 1 quyền (post.update, post.create): true
	// Có tất cả quyền (post.read, post.create): true
	// Admin có quyền post.publish: true
	// Admin có tất cả các quyền (any, other): true
}
