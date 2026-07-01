package ptr

import (
	"reflect"
	"testing"
)

type user struct {
	Name string
	Age  int
}

func TestTo(t *testing.T) {
	// Test primitive type (int)
	vInt := 42
	pInt := To(vInt)
	if pInt == nil {
		t.Fatal("To(int) returned nil")
	}
	if *pInt != vInt {
		t.Errorf("To(int) = %v, want %v", *pInt, vInt)
	}

	// Test primitive type (string)
	vStr := "hello"
	pStr := To(vStr)
	if pStr == nil {
		t.Fatal("To(string) returned nil")
	}
	if *pStr != vStr {
		t.Errorf("To(string) = %q, want %q", *pStr, vStr)
	}

	// Test struct type
	vStruct := user{Name: "Vi", Age: 20}
	pStruct := To(vStruct)
	if pStruct == nil {
		t.Fatal("To(struct) returned nil")
	}
	if *pStruct != vStruct {
		t.Errorf("To(struct) = %+v, want %+v", *pStruct, vStruct)
	}
}

func TestFrom(t *testing.T) {
	tests := []struct {
		name string
		input *int
		want  int
	}{
		{
			name:  "nil pointer returns zero value",
			input: nil,
			want:  0,
		},
		{
			name:  "non-nil pointer returns value",
			input: To(100),
			want:  100,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := From(tt.input)
			if got != tt.want {
				t.Errorf("From() = %v, want %v", got, tt.want)
			}
		})
	}

	// Test with struct
	t.Run("nil struct pointer returns zero struct", func(t *testing.T) {
		var p *user
		got := From(p)
		want := user{}
		if got != want {
			t.Errorf("From(*user) = %+v, want %+v", got, want)
		}
	})

	t.Run("non-nil struct pointer returns struct value", func(t *testing.T) {
		p := &user{Name: "Vi", Age: 20}
		got := From(p)
		want := user{Name: "Vi", Age: 20}
		if got != want {
			t.Errorf("From(*user) = %+v, want %+v", got, want)
		}
	})
}

func TestFromOr(t *testing.T) {
	tests := []struct {
		name     string
		input    *string
		fallback string
		want     string
	}{
		{
			name:     "nil pointer returns fallback value",
			input:    nil,
			fallback: "default",
			want:     "default",
		},
		{
			name:     "non-nil pointer returns value",
			input:    To("hello"),
			fallback: "default",
			want:     "hello",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := FromOr(tt.input, tt.fallback)
			if got != tt.want {
				t.Errorf("FromOr() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestClone(t *testing.T) {
	t.Run("nil pointer returns nil", func(t *testing.T) {
		var p *int
		got := Clone(p)
		if got != nil {
			t.Errorf("Clone(nil) = %v, want nil", got)
		}
	})

	t.Run("non-nil pointer returns a new pointer with the same value", func(t *testing.T) {
		originalVal := 42
		p1 := &originalVal
		p2 := Clone(p1)

		if p2 == nil {
			t.Fatal("Clone() returned nil for non-nil input")
		}
		if p1 == p2 {
			t.Error("Clone() returned the same pointer, expected a new pointer address")
		}
		if *p2 != *p1 {
			t.Errorf("Clone() value = %v, want %v", *p2, *p1)
		}

		// Verify modifying the original doesn't affect the clone
		*p1 = 99
		if *p2 != 42 {
			t.Errorf("Clone value changed to %v after original value modification, expected it to remain 42", *p2)
		}
	})
}

func TestEqual(t *testing.T) {
	tests := []struct {
		name string
		p1   *int
		p2   *int
		want bool
	}{
		{
			name: "both nil",
			p1:   nil,
			p2:   nil,
			want: true,
		},
		{
			name: "p1 nil, p2 non-nil",
			p1:   nil,
			p2:   To(10),
			want: false,
		},
		{
			name: "p1 non-nil, p2 nil",
			p1:   To(10),
			p2:   nil,
			want: false,
		},
		{
			name: "both non-nil, different values",
			p1:   To(10),
			p2:   To(20),
			want: false,
		},
		{
			name: "both non-nil, same values",
			p1:   To(10),
			p2:   To(10),
			want: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Equal(tt.p1, tt.p2)
			if got != tt.want {
				t.Errorf("Equal() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMap(t *testing.T) {
	t.Run("nil pointer returns nil", func(t *testing.T) {
		var p *int
		got := Map(p, func(v int) string {
			return "value"
		})
		if got != nil {
			t.Errorf("Map(nil) = %v, want nil", got)
		}
	})

	t.Run("non-nil pointer applies function and returns new pointer", func(t *testing.T) {
		p := To(10)
		got := Map(p, func(v int) string {
			if v == 10 {
				return "ten"
			}
			return "other"
		})

		if got == nil {
			t.Fatal("Map() returned nil for non-nil input")
		}
		if *got != "ten" {
			t.Errorf("Map() value = %q, want %q", *got, "ten")
		}
	})
}

func TestSlice(t *testing.T) {
	t.Run("nil slice returns nil", func(t *testing.T) {
		var s []int
		got := Slice(s)
		if got != nil {
			t.Errorf("Slice(nil) = %v, want nil", got)
		}
	})

	t.Run("non-nil slice returns slice of pointers", func(t *testing.T) {
		s := []int{10, 20, 30}
		got := Slice(s)

		if len(got) != len(s) {
			t.Fatalf("Slice() length = %d, want %d", len(got), len(s))
		}

		for i, p := range got {
			if p == nil {
				t.Errorf("Slice()[%d] is nil", i)
				continue
			}
			if *p != s[i] {
				t.Errorf("Slice()[%d] value = %v, want %v", i, *p, s[i])
			}
		}
	})
}

func TestSliceFrom(t *testing.T) {
	t.Run("nil slice returns nil", func(t *testing.T) {
		var s []*int
		got := SliceFrom(s)
		if got != nil {
			t.Errorf("SliceFrom(nil) = %v, want nil", got)
		}
	})

	t.Run("slice of pointers returns slice of values", func(t *testing.T) {
		s := []*int{To(10), nil, To(30)}
		got := SliceFrom(s)
		want := []int{10, 0, 30}

		if !reflect.DeepEqual(got, want) {
			t.Errorf("SliceFrom() = %v, want %v", got, want)
		}
	})
}

func TestSliceFromOr(t *testing.T) {
	t.Run("nil slice returns nil", func(t *testing.T) {
		var s []*int
		got := SliceFromOr(s, 99)
		if got != nil {
			t.Errorf("SliceFromOr(nil) = %v, want nil", got)
		}
	})

	t.Run("slice of pointers returns slice of values with fallback", func(t *testing.T) {
		s := []*int{To(10), nil, To(30)}
		got := SliceFromOr(s, 99)
		want := []int{10, 99, 30}

		if !reflect.DeepEqual(got, want) {
			t.Errorf("SliceFromOr() = %v, want %v", got, want)
		}
	})
}
