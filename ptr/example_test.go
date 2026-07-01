package ptr_test

import (
	"fmt"

	"github.com/ducconit/gotoolkit/ptr"
)

type user struct {
	Name string
}

func ExampleTo() {
	// Lấy con trỏ của một hằng số chuỗi
	pStr := ptr.To("hello world")
	fmt.Println(*pStr)

	// Lấy con trỏ của một struct literal
	pStruct := ptr.To(user{Name: "Tiểu Vi"})
	fmt.Println(pStruct.Name)

	// Output:
	// hello world
	// Tiểu Vi
}

func ExampleFrom() {
	var pNil *int
	pVal := ptr.To(42)

	// Chuyển đổi an toàn từ con trỏ về giá trị
	fmt.Println(ptr.From(pNil))
	fmt.Println(ptr.From(pVal))

	// Output:
	// 0
	// 42
}

func ExampleFromOr() {
	var pNil *string
	pVal := ptr.To("Vi")

	// Chuyển đổi an toàn kèm giá trị dự phòng (fallback)
	fmt.Println(ptr.FromOr(pNil, "Guest"))
	fmt.Println(ptr.FromOr(pVal, "Guest"))

	// Output:
	// Guest
	// Vi
}

func ExampleEqual() {
	p1 := ptr.To(10)
	p2 := ptr.To(10)
	p3 := ptr.To(20)
	var p4 *int

	fmt.Println(ptr.Equal(p1, p2))
	fmt.Println(ptr.Equal(p1, p3))
	fmt.Println(ptr.Equal(p1, p4))
	fmt.Println(ptr.Equal(p4, p4))

	// Output:
	// true
	// false
	// false
	// true
}

func ExampleClone() {
	p1 := ptr.To(100)
	p2 := ptr.Clone(p1)

	fmt.Println(*p2)
	fmt.Println(p1 != p2) // Địa chỉ bộ nhớ khác nhau

	// Output:
	// 100
	// true
}

func ExampleMap() {
	pInt := ptr.To(123)

	// Ánh xạ con trỏ int sang con trỏ string
	pStr := ptr.Map(pInt, func(v int) string {
		return fmt.Sprintf("Number: %d", v)
	})

	fmt.Println(*pStr)

	// Output:
	// Number: 123
}

func ExampleSlice() {
	sValues := []string{"a", "b", "c"}
	sPointers := ptr.Slice(sValues)

	for _, p := range sPointers {
		fmt.Print(*p, " ")
	}
	fmt.Println()

	// Output:
	// a b c 
}

func ExampleSliceFrom() {
	sPointers := []*int{ptr.To(1), nil, ptr.To(3)}
	sValues := ptr.SliceFrom(sPointers)

	fmt.Println(sValues)

	// Output:
	// [1 0 3]
}

func ExampleSliceFromOr() {
	sPointers := []*int{ptr.To(1), nil, ptr.To(3)}
	sValues := ptr.SliceFromOr(sPointers, -1)

	fmt.Println(sValues)

	// Output:
	// [1 -1 3]
}
