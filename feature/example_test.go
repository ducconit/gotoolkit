package feature_test

import (
	"context"
	"fmt"
	"github.com/ducconit/gotoolkit/feature"
)

func Example() {
	// 1. Đăng ký một static flag đơn giản
	feature.Register("maintenance-mode", false)
	feature.Register("new-dashboard-ui", true)

	fmt.Println("Bảo trì:", feature.IsEnabled("maintenance-mode"))
	fmt.Println("Dashboard UI mới:", feature.IsEnabled("new-dashboard-ui"))

	// 2. Đăng ký percentage rollout flag (Chỉ bật tính năng cho 20% user)
	feature.RegisterRollout("beta-search-v3", 20)

	// Kiểm tra xem User A và User B có nằm trong tập 20% rollout không.
	// Nhờ sử dụng thuật toán hash ổn định (stateless) qua FNV-1a:
	// - "user-1" hash % 100 = 95 (>= 20) -> false
	// - "user-2" hash % 100 = 14 (< 20)  -> true
	ctx := context.Background()
	userA := "user-1"
	userB := "user-2"

	fmt.Printf("%s được dùng beta search: %v\n", userA, feature.IsEnabledFor(ctx, "beta-search-v3", userA))
	fmt.Printf("%s được dùng beta search: %v\n", userB, feature.IsEnabledFor(ctx, "beta-search-v3", userB))

	// 3. Đăng ký dynamic flag sử dụng hàm đánh giá tùy biến (ví dụ: chỉ bật cho user VIP)
	feature.RegisterFunc("exclusive-deals", func(ctx context.Context, target string) bool {
		// Ở dự án thực tế, bạn có thể đọc từ context hoặc gọi DB/Cache để kiểm tra user profile
		return target == "vip-user-99"
	})

	fmt.Println("User thường xem deal VIP:", feature.IsEnabledFor(ctx, "exclusive-deals", "normal-user-1"))
	fmt.Println("User VIP xem deal VIP:", feature.IsEnabledFor(ctx, "exclusive-deals", "vip-user-99"))

	// Output:
	// Bảo trì: false
	// Dashboard UI mới: true
	// user-1 được dùng beta search: false
	// user-2 được dùng beta search: true
	// User thường xem deal VIP: false
	// User VIP xem deal VIP: true
}
