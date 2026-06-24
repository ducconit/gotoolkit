// Package feature cung cấp giải pháp quản lý Feature Flags (Feature Toggles) tối giản,
// hiệu năng cao, thread-safe và hỗ trợ percentage rollout (stateless).
package feature

import (
	"context"
	"iter"
	"maps"
	"reflect"
	"sync"
	"sync/atomic"
)

// Flag đại diện cho một cấu hình feature flag.
type Flag struct {
	Key          string
	Enabled      bool
	Percentage   int                                           // Tỷ lệ phần trăm rollout (0-100). Nếu > 0, sẽ kích hoạt chế độ percentage rollout.
	EvaluateFunc func(ctx context.Context, target string) bool // Hàm đánh giá động tùy biến.
}

// Manager quản lý danh sách các feature flags.
// Thiết kế sử dụng Copy-on-Write (COW) qua atomic.Pointer giúp các thao tác đọc flag
// đạt tốc độ tối đa (lock-free) và an toàn tuyệt đối khi chạy đồng thời.
type Manager struct {
	flags atomic.Pointer[map[string]Flag]
	mu    sync.Mutex
}

// NewManager khởi tạo một Manager mới.
func NewManager() *Manager {
	m := make(map[string]Flag)
	mgr := &Manager{}
	mgr.flags.Store(&m)
	return mgr
}

var defaultManager = NewManager()

// Default trả về manager mặc định của package.
func Default() *Manager {
	return defaultManager
}

// Register đăng ký hoặc cập nhật một flag tĩnh bật/tắt đơn giản.
func (m *Manager) Register(key string, enabled bool) {
	m.setFlag(Flag{Key: key, Enabled: enabled})
}

// RegisterRollout đăng ký một flag với chế độ percentage rollout.
// percentage nhận giá trị từ 0 đến 100.
func (m *Manager) RegisterRollout(key string, percentage int) {
	if percentage < 0 {
		percentage = 0
	} else if percentage > 100 {
		percentage = 100
	}
	m.setFlag(Flag{Key: key, Percentage: percentage})
}

// RegisterFunc đăng ký một flag với hàm đánh giá động nhận context và target.
func (m *Manager) RegisterFunc(key string, fn func(ctx context.Context, target string) bool) {
	m.setFlag(Flag{Key: key, EvaluateFunc: fn})
}

// IsEnabled kiểm tra xem flag có được bật hay không (không kèm target/context).
func (m *Manager) IsEnabled(key string) bool {
	return m.IsEnabledFor(context.Background(), key, "")
}

// IsEnabledFor kiểm tra xem flag có được bật cho một target cụ thể (ví dụ: User ID)
// dựa trên các quy tắc tĩnh, percentage rollout, hoặc hàm đánh giá động.
func (m *Manager) IsEnabledFor(ctx context.Context, key string, target string) bool {
	flagsPtr := m.flags.Load()
	if flagsPtr == nil {
		return false
	}

	flag, exists := (*flagsPtr)[key]
	if !exists {
		return false
	}

	// 1. Nếu có hàm đánh giá động
	if flag.EvaluateFunc != nil {
		return flag.EvaluateFunc(ctx, target)
	}

	// 2. Nếu cấu hình percentage rollout (chỉ áp dụng khi target không rỗng)
	if flag.Percentage > 0 {
		if target == "" {
			return false
		}
		if flag.Percentage >= 100 {
			return true
		}
		// TỐI ƯU HÓA: Sử dụng thuật toán FNV-1a hash tự viết không phân bổ bộ nhớ (zero-allocation)
		// giúp tăng hiệu năng CPU và loại bỏ áp lực lên bộ nhớ (GC).
		hashVal := hashFNV1a(key, target)
		return int(hashVal%100) < flag.Percentage
	}

	// 3. Quy tắc tĩnh thông thường
	return flag.Enabled
}

// HasFlag kiểm tra xem flag đã được đăng ký hay chưa.
func (m *Manager) HasFlag(key string) bool {
	flagsPtr := m.flags.Load()
	if flagsPtr == nil {
		return false
	}
	_, exists := (*flagsPtr)[key]
	return exists
}

// Unregister xóa một flag khỏi hệ thống.
func (m *Manager) Unregister(key string) {
	m.mu.Lock()
	defer m.mu.Unlock()

	flagsPtr := m.flags.Load()
	if flagsPtr == nil {
		return
	}

	if _, exists := (*flagsPtr)[key]; !exists {
		return // Không tồn tại, không cần làm gì
	}

	newMap := make(map[string]Flag, len(*flagsPtr)-1)
	for k, v := range *flagsPtr {
		if k != key {
			newMap[k] = v
		}
	}
	m.flags.Store(&newMap)
}

// setFlag cập nhật hoặc thêm mới một flag (COW).
func (m *Manager) setFlag(flag Flag) {
	if flag.Key == "" {
		return
	}

	m.mu.Lock()
	defer m.mu.Unlock()

	flagsPtr := m.flags.Load()
	var oldSize int
	if flagsPtr != nil {
		// TỐI ƯU HÓA: Kiểm tra xem flag mới có giống hệt flag cũ đã đăng ký hay không.
		// Nếu giống hệt, ta bỏ qua việc tạo map mới và ghi đè (tiết kiệm allocations và khóa mutex).
		if oldFlag, exists := (*flagsPtr)[flag.Key]; exists {
			if oldFlag.Enabled == flag.Enabled &&
				oldFlag.Percentage == flag.Percentage &&
				compareFuncs(oldFlag.EvaluateFunc, flag.EvaluateFunc) {
				return
			}
		}
		oldSize = len(*flagsPtr)
	}

	newMap := make(map[string]Flag, oldSize+1)
	if flagsPtr != nil {
		maps.Copy(newMap, *flagsPtr)
	}

	newMap[flag.Key] = flag
	m.flags.Store(&newMap)
}

// All trả về một iterator để duyệt qua toàn bộ danh sách các flags đã đăng ký.
// Iterator này an toàn cho concurrency và hoạt động trên snapshot tại thời điểm gọi.
func (m *Manager) All() iter.Seq2[string, Flag] {
	flagsPtr := m.flags.Load()
	if flagsPtr == nil {
		return func(yield func(string, Flag) bool) {}
	}
	return func(yield func(string, Flag) bool) {
		for k, v := range *flagsPtr {
			if !yield(k, v) {
				return
			}
		}
	}
}

// hashFNV1a tính toán FNV-1a 32-bit hash của chuỗi key + target
// mà không phân bổ bộ nhớ (0 allocations).
func hashFNV1a(key, target string) uint32 {
	const (
		offset32 = 2166136261
		prime32  = 16777619
	)
	hash := uint32(offset32)
	for i := 0; i < len(key); i++ {
		hash ^= uint32(key[i])
		hash *= prime32
	}
	for i := 0; i < len(target); i++ {
		hash ^= uint32(target[i])
		hash *= prime32
	}
	return hash
}

// compareFuncs so sánh hai hàm EvaluateFunc một cách an toàn.
func compareFuncs(f1, f2 func(ctx context.Context, target string) bool) bool {
	if f1 == nil && f2 == nil {
		return true
	}
	if f1 == nil || f2 == nil {
		return false
	}
	return reflect.ValueOf(f1).Pointer() == reflect.ValueOf(f2).Pointer()
}

// === Package-level functions sử dụng default manager ===

// Register đăng ký flag tĩnh trên default manager.
func Register(key string, enabled bool) {
	defaultManager.Register(key, enabled)
}

// RegisterRollout đăng ký percentage rollout flag trên default manager.
func RegisterRollout(key string, percentage int) {
	defaultManager.RegisterRollout(key, percentage)
}

// RegisterFunc đăng ký dynamic flag trên default manager.
func RegisterFunc(key string, fn func(ctx context.Context, target string) bool) {
	defaultManager.RegisterFunc(key, fn)
}

// IsEnabled kiểm tra flag trên default manager.
func IsEnabled(key string) bool {
	return defaultManager.IsEnabled(key)
}

// IsEnabledFor kiểm tra flag kèm context và target trên default manager.
func IsEnabledFor(ctx context.Context, key string, target string) bool {
	return defaultManager.IsEnabledFor(ctx, key, target)
}

// HasFlag kiểm tra xem flag đã tồn tại trên default manager chưa.
func HasFlag(key string) bool {
	return defaultManager.HasFlag(key)
}

// Unregister xóa flag trên default manager.
func Unregister(key string) {
	defaultManager.Unregister(key)
}

// All trả về iterator cho default manager.
func All() iter.Seq2[string, Flag] {
	return defaultManager.All()
}
