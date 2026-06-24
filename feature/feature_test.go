package feature

import (
	"context"
	"fmt"
	"hash/fnv"
	"sync"
	"testing"
)

// TestStaticFlags kiểm tra tính năng feature flag tĩnh bật/tắt
func TestStaticFlags(t *testing.T) {
	mgr := NewManager()

	// Mặc định flag chưa đăng ký phải trả về false
	if mgr.IsEnabled("non-exist") {
		t.Error("Flag chua dang ky le ra phai la false")
	}

	mgr.Register("ui-v2", true)
	if !mgr.IsEnabled("ui-v2") {
		t.Error("ui-v2 phai la true sau khi Register")
	}

	mgr.Register("ui-v2", false)
	if mgr.IsEnabled("ui-v2") {
		t.Error("ui-v2 phai la false sau khi cap nhat lai")
	}

	mgr.Unregister("ui-v2")
	if mgr.IsEnabled("ui-v2") {
		t.Error("ui-v2 phai la false sau khi Unregister")
	}
}

// TestDynamicFlags kiểm tra dynamic evaluation
func TestDynamicFlags(t *testing.T) {
	mgr := NewManager()

	mgr.RegisterFunc("premium-feature", func(ctx context.Context, target string) bool {
		// Cho phép nếu target (ví dụ user ID) kết thúc bằng "-vip"
		return len(target) > 4 && target[len(target)-4:] == "-vip"
	})

	tests := []struct {
		target string
		want   bool
	}{
		{"user-1-vip", true},
		{"user-2", false},
		{"vip", false},
		{"", false},
	}

	for _, tt := range tests {
		t.Run(fmt.Sprintf("target=%s", tt.target), func(t *testing.T) {
			got := mgr.IsEnabledFor(context.Background(), "premium-feature", tt.target)
			if got != tt.want {
				t.Errorf("premium-feature for %s = %v, want %v", tt.target, got, tt.want)
			}
		})
	}
}

// TestPercentageRollout kiểm tra tính năng percentage rollout (stateless)
func TestPercentageRollout(t *testing.T) {
	mgr := NewManager()

	// Đăng ký 30% rollout
	mgr.RegisterRollout("new-checkout", 30)

	// Đảm bảo target rỗng thì trả về false
	if mgr.IsEnabledFor(context.Background(), "new-checkout", "") {
		t.Error("Target rong le ra phai luon la false khi check rollout")
	}

	// Chạy thử với 1000 targets giả lập để kiểm tra tính phân phối đều và ổn định
	total := 1000
	enabledCount := 0

	for i := 0; i < total; i++ {
		target := fmt.Sprintf("user-%d", i)
		
		// Kích hoạt flag lần 1
		res1 := mgr.IsEnabledFor(context.Background(), "new-checkout", target)
		
		// Kích hoạt flag lần 2 - phải cho kết quả giống hệt (Consistent)
		res2 := mgr.IsEnabledFor(context.Background(), "new-checkout", target)
		if res1 != res2 {
			t.Fatalf("Tinh nhat quan bi vi pham cho target %s: %v vs %v", target, res1, res2)
		}

		if res1 {
			enabledCount++
		}
	}

	// Tỷ lệ thực tế nên xấp xỉ 30% (cho phép sai số nhỏ 5% do phân phối mẫu nhỏ 1000 phần tử)
	ratio := float64(enabledCount) / float64(total) * 100
	t.Logf("Percentage rollout thuc te: %.2f%% (Muon: 30%%)", ratio)

	if ratio < 25 || ratio > 35 {
		t.Errorf("Ty le rollout thuc te %.2f%% nam ngoai sai so cho phep (25%% - 35%%)", ratio)
	}
}

// TestDefaultManager kiểm tra việc gọi các package-level functions
func TestDefaultManager(t *testing.T) {
	Register("global-flag", true)
	defer Unregister("global-flag")

	if !IsEnabled("global-flag") {
		t.Error("Default manager phai hoat dong dung voi Register va IsEnabled")
	}
}

// TestConcurrency kiểm tra an toàn đa luồng (Data Race detector)
func TestConcurrency(t *testing.T) {
	mgr := NewManager()
	var wg sync.WaitGroup

	// Đăng ký trước một số flag
	mgr.Register("flag-1", true)
	mgr.RegisterRollout("flag-2", 50)

	// 10 goroutines liên tục cập nhật và xóa flags
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			mgr.Register(fmt.Sprintf("dynamic-flag-%d", id), true)
			mgr.RegisterRollout("flag-2", id*10)
			mgr.Unregister(fmt.Sprintf("dynamic-flag-%d", id))
		}(i)
	}

	// 10 goroutines liên tục đọc flags
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			_ = mgr.IsEnabled("flag-1")
			_ = mgr.IsEnabledFor(context.Background(), "flag-2", fmt.Sprintf("user-%d", id))
		}(i)
	}

	wg.Wait()
}

// TestHashFNV1aConsistency kiểm chứng tính chính xác của hàm FNV-1a tối ưu so với thư viện chuẩn
func TestHashFNV1aConsistency(t *testing.T) {
	tests := []struct {
		key    string
		target string
	}{
		{"beta-search-v3", "user-1"},
		{"beta-search-v3", "user-2"},
		{"new-checkout", "user-999"},
		{"", ""},
		{"a", "b"},
	}

	for _, tt := range tests {
		t.Run(fmt.Sprintf("key=%s,target=%s", tt.key, tt.target), func(t *testing.T) {
			// Tính bằng hash/fnv stdlib
			h := fnv.New32a()
			_, _ = h.Write([]byte(tt.key))
			_, _ = h.Write([]byte(tt.target))
			want := h.Sum32()

			// Tính bằng hàm tối ưu hóa không allocation
			got := hashFNV1a(tt.key, tt.target)

			if got != want {
				t.Errorf("hashFNV1a() = %v, want %v", got, want)
			}
		})
	}
}

// TestIteratorAll kiểm tra iterator duyệt qua toàn bộ flags
func TestIteratorAll(t *testing.T) {
	mgr := NewManager()
	mgr.Register("f1", true)
	mgr.Register("f2", false)

	expected := map[string]bool{
		"f1": true,
		"f2": false,
	}

	count := 0
	for k, v := range mgr.All() {
		count++
		wantEnabled, exists := expected[k]
		if !exists {
			t.Errorf("Flag %s khong co trong danh sach mong doi", k)
		}
		if v.Enabled != wantEnabled {
			t.Errorf("Flag %s co Enabled = %v, mong doi %v", k, v.Enabled, wantEnabled)
		}
	}

	if count != 2 {
		t.Errorf("Iterator All tra ve %d flags, mong doi 2", count)
	}
}

// TestHasFlag kiểm tra hàm check sự tồn tại của flag
func TestHasFlag(t *testing.T) {
	mgr := NewManager()
	if mgr.HasFlag("test-key") {
		t.Error("Khong duoc ton tai truoc khi dang ky")
	}
	mgr.Register("test-key", true)
	if !mgr.HasFlag("test-key") {
		t.Error("Phai ton tai sau khi dang ky")
	}
}

// === BENCHMARK TESTS ===

// BenchmarkIsEnabled_Static đo tốc độ đọc static flag đơn giản
func BenchmarkIsEnabled_Static(b *testing.B) {
	mgr := NewManager()
	mgr.Register("test-flag", true)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = mgr.IsEnabled("test-flag")
	}
}

// BenchmarkIsEnabled_Rollout đo tốc độ đọc percentage rollout flag (có hashing tối ưu 0 allocation)
func BenchmarkIsEnabled_Rollout(b *testing.B) {
	mgr := NewManager()
	mgr.RegisterRollout("test-rollout", 25)
	ctx := context.Background()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = mgr.IsEnabledFor(ctx, "test-rollout", "user-123456")
	}
}

// BenchmarkIsEnabled_Func đo tốc độ đọc dynamic function flag
func BenchmarkIsEnabled_Func(b *testing.B) {
	mgr := NewManager()
	mgr.RegisterFunc("test-func", func(ctx context.Context, target string) bool {
		return target == "admin"
	})
	ctx := context.Background()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = mgr.IsEnabledFor(ctx, "test-func", "user-123456")
	}
}

// BenchmarkConcurrentReadWrite đo tốc độ khi đọc/ghi đồng thời
func BenchmarkConcurrentReadWrite(b *testing.B) {
	mgr := NewManager()
	mgr.Register("test-flag", true)
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		i := 0
		for pb.Next() {
			if i%100 == 0 {
				mgr.Register("test-flag", i%2 == 0) // Ghi hiếm (1%) - sẽ kích hoạt tối ưu skip duplicate
			} else {
				_ = mgr.IsEnabled("test-flag") // Đọc thường xuyên (99%)
			}
			i++
		}
	})
}
