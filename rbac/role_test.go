package rbac

import (
	"slices"
	"sync"
	"testing"
)

// TestNewRole kiểm tra việc khởi tạo Role
func TestNewRole(t *testing.T) {
	t.Run("Khoi tao thanh cong", func(t *testing.T) {
		r := NewRole("admin", "read", "write")
		if r.ID() != "admin" {
			t.Errorf("NewRole() ID = %v, want admin", r.ID())
		}
		if !r.HasPermission("read") {
			t.Error("NewRole() thieu quyen 'read'")
		}
		if !r.HasPermission("write") {
			t.Error("NewRole() thieu quyen 'write'")
		}
	})

	t.Run("Panic khi ID rong", func(t *testing.T) {
		defer func() {
			if r := recover(); r == nil {
				t.Errorf("NewRole() voi ID rong le ra phai panic")
			}
		}()
		NewRole("")
	})
}

// TestAddPermission kiểm tra việc thêm quyền động
func TestAddPermission(t *testing.T) {
	tests := []struct {
		name        string
		initial     []string
		toAdd       []string
		checkHas    []string
		wantHas     bool
		checkAll    []string
		wantAll     bool
		hasWildcard bool
	}{
		{
			name:     "Them quyen moi binh thuong",
			initial:  []string{"read"},
			toAdd:    []string{"write", "delete"},
			checkHas: []string{"write"},
			wantHas:  true,
			checkAll: []string{"read", "write", "delete"},
			wantAll:  true,
		},
		{
			name:        "Them quyen wildcard '*'",
			initial:     []string{"read"},
			toAdd:       []string{"*"},
			checkHas:    []string{"any_random_permission"},
			wantHas:     true,
			checkAll:    []string{"a", "b", "c"},
			wantAll:     true,
			hasWildcard: true,
		},
		{
			name:     "Them quyen trung lap va rong",
			initial:  []string{"read"},
			toAdd:    []string{"read", "", "write"},
			checkHas: []string{"write"},
			wantHas:  true,
			checkAll: []string{"read", "write"},
			wantAll:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := NewRole("test-role", tt.initial...)
			r.AddPermission(tt.toAdd...)

			if tt.hasWildcard && !r.hasWildcard.Load() {
				t.Error("Role phai co wildcard flag bang true")
			}

			if got := r.HasPermission(tt.checkHas...); got != tt.wantHas {
				t.Errorf("HasPermission(%v) = %v, want %v", tt.checkHas, got, tt.wantHas)
			}

			if got := r.HasAllPermission(tt.checkAll...); got != tt.wantAll {
				t.Errorf("HasAllPermission(%v) = %v, want %v", tt.checkAll, got, tt.wantAll)
			}
		})
	}
}

// TestHasPermission kiểm tra kiểm tra ít nhất 1 quyền (Table-Driven)
func TestHasPermission(t *testing.T) {
	tests := []struct {
		name        string
		rolePerms   []string
		checkPerms  []string
		want        bool
		description string
	}{
		{
			name:        "Khop dung 1 quyen",
			rolePerms:   []string{"read", "write"},
			checkPerms:  []string{"read"},
			want:        true,
			description: "Co quyen 'read'",
		},
		{
			name:        "Khop 1 trong nhieu quyen",
			rolePerms:   []string{"read", "write"},
			checkPerms:  []string{"delete", "write"},
			want:        true,
			description: "Khop quyen 'write'",
		},
		{
			name:        "Khong khop quyen nao",
			rolePerms:   []string{"read", "write"},
			checkPerms:  []string{"delete", "admin"},
			want:        false,
			description: "Khong co quyen phu hop",
		},
		{
			name:        "Check danh sach rong",
			rolePerms:   []string{"read", "write"},
			checkPerms:  []string{},
			want:        false,
			description: "Danh sach kiem tra rong phai tra ve false",
		},
		{
			name:        "Co quyen wildcard '*'",
			rolePerms:   []string{"*"},
			checkPerms:  []string{"any", "other"},
			want:        true,
			description: "Quyen '*' cho phep tat ca",
		},
		{
			name:        "Co quyen wildcard '*' check rong",
			rolePerms:   []string{"*"},
			checkPerms:  []string{},
			want:        false,
			description: "Co '*' nhung check rong van phai tra ve false",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := NewRole("test", tt.rolePerms...)
			if got := r.HasPermission(tt.checkPerms...); got != tt.want {
				t.Errorf("HasPermission() = %v, want %v (%s)", got, tt.want, tt.description)
			}
		})
	}
}

// TestHasAllPermission kiểm tra yêu cầu đầy đủ quyền (Table-Driven)
func TestHasAllPermission(t *testing.T) {
	tests := []struct {
		name        string
		rolePerms   []string
		checkPerms  []string
		want        bool
		description string
	}{
		{
			name:        "Co du tat ca cac quyen",
			rolePerms:   []string{"read", "write", "delete"},
			checkPerms:  []string{"read", "write"},
			want:        true,
			description: "Co ca 'read' va 'write'",
		},
		{
			name:        "Thieu 1 quyen",
			rolePerms:   []string{"read", "write"},
			checkPerms:  []string{"read", "write", "delete"},
			want:        false,
			description: "Thieu quyen 'delete'",
		},
		{
			name:        "Check danh sach rong",
			rolePerms:   []string{"read", "write"},
			checkPerms:  []string{},
			want:        false,
			description: "Danh sach kiem tra rong phai tra ve false",
		},
		{
			name:        "Co quyen wildcard '*'",
			rolePerms:   []string{"*"},
			checkPerms:  []string{"read", "write", "anything"},
			want:        true,
			description: "Wildcard '*' bao phu tat ca",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := NewRole("test", tt.rolePerms...)
			if got := r.HasAllPermission(tt.checkPerms...); got != tt.want {
				t.Errorf("HasAllPermission() = %v, want %v (%s)", got, tt.want, tt.description)
			}
		})
	}
}

// TestPermissionsIterator kiem tra Iterator Go 1.23+
func TestPermissionsIterator(t *testing.T) {
	t.Run("Iterator voi quyen binh thuong", func(t *testing.T) {
		r := NewRole("test-role", "read", "write", "delete")
		
		var gathered []string
		for p := range r.Permissions() {
			gathered = append(gathered, p)
		}

		slices.Sort(gathered)
		expected := []string{"delete", "read", "write"}

		if !slices.Equal(gathered, expected) {
			t.Errorf("Permissions() iterator returned %v, want %v", gathered, expected)
		}
	})

	t.Run("Iterator voi quyen wildcard toan cuc", func(t *testing.T) {
		r := NewRole("admin", "*")
		var gathered []string
		for p := range r.Permissions() {
			gathered = append(gathered, p)
		}

		if len(gathered) != 1 || gathered[0] != "*" {
			t.Errorf("Permissions() iterator for admin returned %v, want [*]", gathered)
		}
	})
}

// TestConcurrency kiêm tra an toan dong thoi (Data Race)
func TestConcurrency(t *testing.T) {
	r := NewRole("concurrent-role", "init-perm")
	var wg sync.WaitGroup

	// 10 goroutines lien tuc them quyen
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			r.AddPermission("perm-a", "perm-b")
		}(i)
	}

	// 10 goroutines lien tuc doc kiem tra quyen
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			_ = r.HasPermission("init-perm")
			_ = r.HasAllPermission("perm-a", "perm-b")
			// Kiem tra iterator concurrent read
			for range r.Permissions() {
			}
		}()
	}

	wg.Wait()
}

// === BENCHMARK TESTS ===

func BenchmarkHasPermission_Wildcard(b *testing.B) {
	r := NewRole("admin", "*")
	b.ResetTimer()
	for b.Loop() {
		_ = r.HasPermission("user.create", "user.delete")
	}
}

func BenchmarkHasPermission_Small(b *testing.B) {
	r := NewRole("user", "user.read", "user.write")
	b.ResetTimer()
	for b.Loop() {
		_ = r.HasPermission("user.write")
	}
}

func BenchmarkHasPermission_Large(b *testing.B) {
	perms := make([]string, 100)
	for i := 0; i < 100; i++ {
		perms[i] = string(rune(i))
	}
	r := NewRole("power-user", perms...)
	b.ResetTimer()
	for b.Loop() {
		_ = r.HasPermission("nonexistent-permission")
	}
}

func BenchmarkHasAllPermission_Success(b *testing.B) {
	r := NewRole("user", "read", "write", "delete")
	b.ResetTimer()
	for b.Loop() {
		_ = r.HasAllPermission("read", "write")
	}
}

func BenchmarkAddPermission_COW(b *testing.B) {
	r := NewRole("editor", "read")
	b.ResetTimer()
	for b.Loop() {
		r.AddPermission("write")
	}
}

func BenchmarkConcurrentReadWrite(b *testing.B) {
	r := NewRole("bench", "read")
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		i := 0
		for pb.Next() {
			if i%100 == 0 {
				r.AddPermission("write") // Write hiem (1%)
			} else {
				_ = r.HasPermission("read") // Read thuong xuyen (99%)
			}
			i++
		}
	})
}

func BenchmarkNewRole_WithPermissions(b *testing.B) {
	for b.Loop() {
		_ = NewRole("user", "read", "write", "delete")
	}
}

func BenchmarkAddPermission_Duplicate(b *testing.B) {
	r := NewRole("editor", "read", "write")
	b.ResetTimer()
	for b.Loop() {
		r.AddPermission("read", "write")
	}
}

func BenchmarkAddPermission_NewPerm_Large(b *testing.B) {
	perms := make([]string, 20)
	for i := range perms {
		perms[i] = "perm-" + string(rune('a'+i))
	}
	b.ResetTimer()
	for b.Loop() {
		r := NewRole("user", perms...)
		r.AddPermission("new-perm")
	}
}
