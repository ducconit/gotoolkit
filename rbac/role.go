// Package rbac cung cấp các giải pháp phân quyền (RBAC) tối giản,
// hiệu năng cao và tiết kiệm tài nguyên.
package rbac

import (
	"maps"
	"slices"
	"sync"
	"sync/atomic"
)

// Role đại diện cho một vai trò của người dùng với các quyền hạn (permissions) đi kèm.
// Struct này được thiết kế thread-safe sử dụng cơ chế Copy-on-Write (COW) kết hợp với
// atomic operations, giúp các thao tác đọc (HasPermission, HasAllPermission) đạt tốc độ tối đa
// và không bị lock-contention.
type Role struct {
	id          string
	permissions atomic.Pointer[map[string]struct{}]
	hasWildcard atomic.Bool
	mu          sync.Mutex // Bảo vệ các thao tác ghi đồng thời (AddPermission)
}

// NewRole khởi tạo một vai trò mới với ID và danh sách quyền ban đầu.
// Tham số id là bắt buộc. Nếu id rỗng, hàm sẽ panic.
func NewRole(id string, permissions ...string) *Role {
	if id == "" {
		panic("rbac: role id is required")
	}

	r := &Role{id: id}

	if len(permissions) == 0 {
		return r
	}

	// Tại thời điểm khởi tạo, chưa có goroutine nào khác giữ reference đến Role này,
	// nên có thể build map trực tiếp mà không cần acquire lock qua AddPermission.
	if slices.Contains(permissions, "*") {
		r.hasWildcard.Store(true)
		return r
	}

	m := make(map[string]struct{}, len(permissions))
	for _, p := range permissions {
		if p != "" {
			m[p] = struct{}{}
		}
	}
	r.permissions.Store(&m)

	return r
}

// ID trả về định danh của Role.
func (r *Role) ID() string {
	return r.id
}

// AddPermission thêm một hoặc nhiều quyền vào Role.
// Thao tác này là thread-safe. Nếu quyền "*" được thêm vào, Role sẽ có mọi quyền
// và các quyền cụ thể khác sẽ được giải phóng để tiết kiệm bộ nhớ.
func (r *Role) AddPermission(permissions ...string) {
	if len(permissions) == 0 {
		return
	}

	r.mu.Lock()
	defer r.mu.Unlock()

	// Nếu đã có quyền đặc biệt "*", không cần lưu thêm quyền nào khác
	if r.hasWildcard.Load() {
		return
	}

	// Kiểm tra xem danh sách quyền mới có chứa "*" không
	if slices.Contains(permissions, "*") {
		r.hasWildcard.Store(true)
		r.permissions.Store(nil) // Giải phóng bộ nhớ của map cũ
		return
	}

	// Kiểm tra xem có quyền mới thực sự cần thêm không để tránh COW copy lãng phí
	oldMapPtr := r.permissions.Load()
	if oldMapPtr != nil {
		oldMap := *oldMapPtr
		hasNew := false
		for _, p := range permissions {
			if p == "" {
				continue
			}
			if _, exists := oldMap[p]; !exists {
				hasNew = true
				break
			}
		}
		if !hasNew {
			return // Tất cả quyền đã tồn tại, không cần copy
		}
	}

	// Copy-on-Write: Tạo map mới kế thừa từ map cũ và bổ sung quyền mới
	var oldSize int
	if oldMapPtr != nil {
		oldSize = len(*oldMapPtr)
	}
	newMap := make(map[string]struct{}, oldSize+len(permissions))

	if oldMapPtr != nil {
		maps.Copy(newMap, *oldMapPtr)
	}

	for _, p := range permissions {
		if p != "" {
			newMap[p] = struct{}{}
		}
	}

	r.permissions.Store(&newMap)
}

// HasPermission kiểm tra xem Role có ít nhất một quyền trong danh sách yêu cầu hay không.
// Trả về true nếu Role có quyền "*" hoặc có ít nhất một quyền khớp.
// Trả về false nếu danh sách quyền cần kiểm tra rỗng hoặc không khớp quyền nào.
func (r *Role) HasPermission(permissions ...string) bool {
	if len(permissions) == 0 {
		return false
	}

	if r.hasWildcard.Load() {
		return true
	}

	pMapPtr := r.permissions.Load()
	if pMapPtr == nil {
		return false
	}

	pMap := *pMapPtr
	for _, p := range permissions {
		if _, exists := pMap[p]; exists {
			return true
		}
	}

	return false
}

// HasAllPermission kiểm tra xem Role có đầy đủ tất cả các quyền trong danh sách yêu cầu hay không.
// Trả về true nếu Role có quyền "*" hoặc sở hữu tất cả các quyền được truyền vào.
// Trả về false nếu danh sách quyền cần kiểm tra rỗng hoặc thiếu ít nhất một quyền.
func (r *Role) HasAllPermission(permissions ...string) bool {
	if len(permissions) == 0 {
		return false
	}

	if r.hasWildcard.Load() {
		return true
	}

	pMapPtr := r.permissions.Load()
	if pMapPtr == nil {
		return false
	}

	pMap := *pMapPtr
	for _, p := range permissions {
		if _, exists := pMap[p]; !exists {
			return false
		}
	}

	return true
}
