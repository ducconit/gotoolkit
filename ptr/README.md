# ptr

Package `ptr` cung cấp các hàm tiện ích để làm việc với con trỏ (pointer) một cách type-safe sử dụng Go Generics.

## Tại sao cần package này?

Trong Go, bạn không thể trực tiếp lấy địa chỉ bộ nhớ của một hằng số (literal) hay kết quả trả về từ một hàm (ví dụ: `&"hello"`, `&100`, hoặc `&getAge()` sẽ gây lỗi biên dịch). Để lấy con trỏ, bạn thường phải khai báo một biến tạm:

```go
temp := "hello"
ptr := &temp
```

Package `ptr` giúp đơn giản hóa việc này thành một dòng code duy nhất, đồng thời hỗ trợ chuyển đổi ngược lại từ con trỏ về giá trị an toàn (tránh lỗi panic do dereference con trỏ `nil`).

## Cài đặt

Package được tích hợp sẵn trong `gotoolkit`:

```go
import "github.com/ducconit/gotoolkit/ptr"
```

## Các API Chính

### 1. `ptr.To[T any](v T) *T`
Chuyển đổi một giá trị sang con trỏ. Hỗ trợ bất kỳ kiểu dữ liệu nào (struct, string, int, bool...).

```go
pStr := ptr.To("hello") // *string
pInt := ptr.To(123)     // *int
```

### 2. `ptr.From[T any](p *T) T`
Chuyển đổi con trỏ về giá trị. Nếu con trỏ là `nil`, trả về giá trị mặc định (zero value) của kiểu dữ liệu đó.

```go
var p *int
val := ptr.From(p) // val = 0
```

### 3. `ptr.FromOr[T any](p *T, fallback T) T`
Chuyển đổi con trỏ về giá trị. Nếu con trỏ là `nil`, trả về giá trị dự phòng (`fallback`) được chỉ định.

```go
var p *string
val := ptr.FromOr(p, "default_value") // val = "default_value"
```

### 4. `ptr.Clone[T any](p *T) *T`
Sao chép giá trị của một con trỏ sang một địa chỉ vùng nhớ mới. Nếu con trỏ là `nil`, trả về `nil`.

```go
p1 := ptr.To(10)
p2 := ptr.Clone(p1) // p2 trỏ tới vùng nhớ mới chứa giá trị 10
```

### 5. `ptr.Equal[T comparable](p1, p2 *T) bool`
So sánh giá trị của hai con trỏ. Trả về `true` nếu cả hai đều `nil`, hoặc trỏ tới hai giá trị bằng nhau.

```go
p1 := ptr.To(10)
p2 := ptr.To(10)
ptr.Equal(p1, p2) // true
```

### 6. `ptr.Map[T, U any](p *T, fn func(T) U) *U`
Ánh xạ con trỏ kiểu `T` sang con trỏ kiểu `U` qua hàm biến đổi. Trả về `nil` nếu con trỏ nguồn là `nil`.

```go
pInt := ptr.To(100)
pStr := ptr.Map(pInt, func(v int) string {
    return fmt.Sprintf("val: %d", v)
}) // *string trỏ tới "val: 100"
```

### 7. Slice Utilities

- `ptr.Slice[T any](s []T) []*T`: Chuyển một slice giá trị thành slice chứa các con trỏ trỏ tới bản sao của các giá trị đó.
- `ptr.SliceFrom[T any](s []*T) []T`: Chuyển một slice con trỏ thành slice giá trị (dùng zero value nếu con trỏ là `nil`).
- `ptr.SliceFromOr[T any](s []*T, fallback T) []T`: Chuyển một slice con trỏ thành slice giá trị (dùng `fallback` nếu con trỏ là `nil`).
