package arena

import (
	"fmt"
	"testing"

	"github.com/ongniud/arena-map/alloc"
)

func TestNew(t *testing.T) {
	a := &slabArena{
		alloc: alloc.NewAllocator(alloc.WithSlabSize(1 * alloc.KB)),
	}

	type TestStruct struct {
		A int
		B float64
	}

	// 测试使用非 nil Arena
	ptr := New[TestStruct](a)
	if ptr == nil {
		t.Error("Expected non-nil pointer when using a valid arena")
	}

	// 检查分配的内存是否有效
	if (*ptr).A != 0 || (*ptr).B != 0 {
		t.Error("Expected zero-initialized struct")
	}

	Free[TestStruct](a, ptr)

	// 测试使用 nil Arena
	ptrNil := New[TestStruct](nil)
	if ptrNil == nil {
		t.Error("Expected non-nil pointer when using nil arena")
	}

	// 检查分配的内存是否有效
	if (*ptrNil).A != 0 || (*ptrNil).B != 0 {
		t.Error("Expected zero-initialized struct")
	}
}

func TestMakeSlice(t *testing.T) {
	a := &slabArena{
		alloc: alloc.NewAllocator(),
	}

	// 测试使用非 nil Arena
	slice := MakeSlice[int](a, 5, 10)
	if len(slice) != 5 || cap(slice) != 10 {
		t.Errorf("Expected slice of length 5 and capacity 10, got length %d and capacity %d", len(slice), cap(slice))
	}

	// 测试值初始化
	for i := range slice {
		if slice[i] != 0 {
			t.Errorf("Expected zero-initialized slice, got %d at index %d", slice[i], i)
		}
	}

	// 测试使用 nil Arena
	sliceNil := MakeSlice[int](nil, 5, 10)
	if len(sliceNil) != 5 || cap(sliceNil) != 10 {
		t.Errorf("Expected slice of length 5 and capacity 10, got length %d and capacity %d", len(sliceNil), cap(sliceNil))
	}

	// 检查值初始化
	for i := range sliceNil {
		if sliceNil[i] != 0 {
			t.Errorf("Expected zero-initialized slice, got %d at index %d", sliceNil[i], i)
		}
	}
}

func TestArenaAllocFailure(t *testing.T) {
	a := &slabArena{
		alloc: alloc.NewAllocator(),
	}
	type SmallStruct struct {
		X int
		Y int
	}
	ptr := New[SmallStruct](a)
	ptr.X = 1
	ptr.Y = 2
	fmt.Println(ptr)
	prt2 := New[SmallStruct](a)
	prt2.X = 3
	prt2.Y = 4
	fmt.Println(ptr)
}

func TestMakeSliceAllocFailure(t *testing.T) {
	a := &slabArena{
		alloc: alloc.NewAllocator(),
	}
	// 尝试分配超过 Arena 可用内存的大小
	slice := MakeSlice[int](a, 5, 10)
	if slice != nil {
		t.Error("Expected nil slice when allocating beyond available memory")
	}
}
