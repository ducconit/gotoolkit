// Package ptr cung cấp các hàm tiện ích để làm việc với con trỏ (pointer)
// một cách type-safe sử dụng Go Generics.
package ptr

// To trả về một con trỏ trỏ tới bản sao của giá trị v.
// Hàm này hữu ích khi cần truyền hằng số hoặc kết quả trực tiếp của một hàm vào nơi nhận con trỏ.
func To[T any](v T) *T {
	return new(v)
}

// From trả về giá trị mà con trỏ p trỏ tới.
// Nếu p là nil, hàm trả về giá trị mặc định (zero value) của kiểu T.
func From[T any](p *T) T {
	if p == nil {
		var zero T
		return zero
	}
	return *p
}

// FromOr trả về giá trị mà con trỏ p trỏ tới.
// Nếu p là nil, hàm trả về giá trị fallback được cung cấp.
func FromOr[T any](p *T, fallback T) T {
	if p == nil {
		return fallback
	}
	return *p
}

// Clone tạo một con trỏ mới trỏ tới bản sao của giá trị mà p đang trỏ tới.
// Nếu p là nil, hàm trả về nil.
func Clone[T any](p *T) *T {
	if p == nil {
		return nil
	}
	v := *p
	return &v
}

// Equal so sánh giá trị của hai con trỏ.
// Trả về true nếu:
//   - Cả hai con trỏ đều nil.
//   - Cả hai con trỏ đều không nil và trỏ tới hai giá trị bằng nhau.
//
// Trả về false trong các trường hợp còn lại.
func Equal[T comparable](p1, p2 *T) bool {
	if p1 == nil && p2 == nil {
		return true
	}
	if p1 == nil || p2 == nil {
		return false
	}
	return *p1 == *p2
}

// Map ánh xạ con trỏ p (kiểu T) sang một con trỏ mới (kiểu U) thông qua hàm biến đổi fn.
// Nếu p là nil, hàm trả về nil mà không gọi fn.
func Map[T, U any](p *T, fn func(T) U) *U {
	if p == nil {
		return nil
	}
	v := fn(*p)
	return &v
}

// Slice chuyển đổi một slice chứa các giá trị kiểu T thành một slice chứa các con trỏ trỏ tới bản sao của chúng.
func Slice[T any](s []T) []*T {
	if s == nil {
		return nil
	}
	res := make([]*T, len(s))
	for i := range s {
		res[i] = &s[i]
	}
	return res
}

// SliceFrom chuyển đổi một slice chứa các con trỏ kiểu *T thành một slice chứa các giá trị kiểu T.
// Nếu một con trỏ trong slice là nil, giá trị tương ứng trong kết quả sẽ là zero value của kiểu T.
func SliceFrom[T any](s []*T) []T {
	if s == nil {
		return nil
	}
	res := make([]T, len(s))
	for i, p := range s {
		if p != nil {
			res[i] = *p
		}
	}
	return res
}

// SliceFromOr chuyển đổi một slice chứa các con trỏ kiểu *T thành một slice chứa các giá trị kiểu T.
// Nếu một con trỏ trong slice là nil, giá trị tương ứng trong kết quả sẽ là giá trị fallback được cung cấp.
func SliceFromOr[T any](s []*T, fallback T) []T {
	if s == nil {
		return nil
	}
	res := make([]T, len(s))
	for i, p := range s {
		if p != nil {
			res[i] = *p
		} else {
			res[i] = fallback
		}
	}
	return res
}
