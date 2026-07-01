# Feature Package

`feature` cung cấp giải pháp quản lý Feature Flags (Feature Toggles) tối giản, hiệu năng cao, thread-safe và hỗ trợ percentage rollout (stateless) cho các dự án Golang.

## Các Tính năng Chính

*   **Đa dạng Cơ Chế Đánh Giá**:
    *   *Tĩnh (Static Flag)*: Bật/tắt đơn giản.
    *   *Rollout theo Tỷ lệ (Percentage Rollout)*: Rollout dần tính năng cho một tỷ lệ phần trăm user nhất định dựa trên thuật toán băm (stateless FNV-1a hash) **Zero-Allocation**.
    *   *Động (Evaluate Function)*: Đăng ký hàm đánh giá phức tạp dựa trên Context hoặc Target truyền vào.
*   **Thiết kế Lock-Free cho Đọc Flag**: Sử dụng kỹ thuật **Copy-on-Write (COW)** qua `atomic.Pointer` giúp các thao tác đọc flag đạt tốc độ tối đa, an toàn tuyệt đối khi chạy đồng thời với hàng ngàn goroutines.

## Hướng dẫn Sử dụng Nhanh

```go
package main

import (
	"context"
	"fmt"
	"github.com/ducconit/gotoolkit/feature"
)

func main() {
	// Sử dụng manager mặc định
	mgr := feature.Default()

	// 1. Flag tĩnh
	mgr.Register("new-ui", true)
	fmt.Println("New UI:", mgr.IsEnabled("new-ui")) // true

	// 2. Percentage Rollout (10% người dùng)
	mgr.RegisterRollout("beta-feature", 10)
	fmt.Println("User A (rollout):", mgr.IsEnabledFor(context.Background(), "beta-feature", "user-A"))
	fmt.Println("User B (rollout):", mgr.IsEnabledFor(context.Background(), "beta-feature", "user-B"))

	// 3. Hàm đánh giá động dựa trên context
	mgr.RegisterFunc("premium-features", func(ctx context.Context, target string) bool {
		isPremium, _ := ctx.Value("is_premium").(bool)
		return isPremium
	})

	ctx := context.WithValue(context.Background(), "is_premium", true)
	fmt.Println("Premium User:", mgr.IsEnabledFor(ctx, "premium-features", "user-123")) // true
}
```

## Báo cáo Benchmark

Đo thực tế trên CPU **AMD Ryzen 7 8745H**:

| Hàm Benchmark | Thời gian thực thi (ns/op) | Dung lượng bộ nhớ (B/op) | Allocations (allocs/op) |
| :--- | :---: | :---: | :---: |
| `BenchmarkIsEnabled_Static` | 11.73 ns | 0 B | 0 allocs |
| `BenchmarkIsEnabled_Rollout` | 19.74 ns | 0 B | 0 allocs |
| `BenchmarkIsEnabled_Func` | 9.80 ns | 0 B | 0 allocs |
| `BenchmarkConcurrentReadWrite` | 4.31 ns | 0 B | 0 allocs |
