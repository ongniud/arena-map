package arena

import (
	"fmt"
	"math/rand"
	"runtime"
	"testing"
	"time"

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

	ptr := New[TestStruct](a)
	if ptr == nil {
		t.Error("Expected non-nil pointer when using a valid arena")
	}

	if (*ptr).A != 0 || (*ptr).B != 0 {
		t.Error("Expected zero-initialized struct")
	}

	Free[TestStruct](a, ptr)
	ptrNil := New[TestStruct](nil)
	if ptrNil == nil {
		t.Error("Expected non-nil pointer when using nil arena")
	}
	if (*ptrNil).A != 0 || (*ptrNil).B != 0 {
		t.Error("Expected zero-initialized struct")
	}
}

func TestMakeSlice(t *testing.T) {
	a := &slabArena{
		alloc: alloc.NewAllocator(),
	}

	slice := MakeSlice[int](a, 5, 10)
	if len(slice) != 5 || cap(slice) != 10 {
		t.Errorf("Expected slice of length 5 and capacity 10, got length %d and capacity %d", len(slice), cap(slice))
	}
	for i := range slice {
		if slice[i] != 0 {
			t.Errorf("Expected zero-initialized slice, got %d at index %d", slice[i], i)
		}
	}

	sliceNil := MakeSlice[int](nil, 5, 10)
	if len(sliceNil) != 5 || cap(sliceNil) != 10 {
		t.Errorf("Expected slice of length 5 and capacity 10, got length %d and capacity %d", len(sliceNil), cap(sliceNil))
	}
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
	slice := MakeSlice[int](a, 5, 10)
	if slice != nil {
		t.Error("Expected nil slice when allocating beyond available memory")
	}
}

func TestSliceAppend1(t *testing.T) {
	s := make([]int, 3, 10)
	copy(s, []int{1, 2, 3})
	fmt.Println(s, len(s), cap(s))

	arena := NewArena()
	r, resized := SliceAppend(arena, s, []int{4, 5, 6}...)
	fmt.Println(r, len(r), cap(r), resized)
	if resized {
		FreeSlice(arena, s)
	}

	s = r
	r, resized = SliceAppend(arena, s, []int{7, 8, 9}...)
	fmt.Println(r, len(r), cap(r), resized)
	if resized {
		FreeSlice(arena, s)
	}

	s = r
	r, resized = SliceAppend(arena, s, []int{10, 11, 12}...)
	fmt.Println(r, len(r), cap(r), resized)
	if resized {
		FreeSlice(arena, s)
	}

	if resized {
		FreeSlice(arena, r)
	}
}

// TestSliceAppend tests the SliceAppend function.
func TestSliceAppend(t *testing.T) {
	arena := NewArena()
	tests := []struct {
		initial  []int
		data     []int
		expected []int
		resized  bool
	}{
		{[]int{1, 2, 3}, []int{4, 5}, []int{1, 2, 3, 4, 5}, true},
		{[]int{}, []int{1}, []int{1}, true},
		{[]int{1, 2}, []int{3, 4, 5, 6}, []int{1, 2, 3, 4, 5, 6}, true},
		{[]int{1}, []int{}, []int{1}, false},
		{[]int{1, 2, 3, 4, 5, 6, 7, 8}, []int{9, 10}, []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}, true},
	}

	for _, test := range tests {
		fmt.Printf("Testing with initial slice: %v and data to append: %v\n", test.initial, test.data)
		result, resized := SliceAppend(arena, test.initial, test.data...)
		fmt.Printf("Result: %v, Resized: %v\n", result, resized)
		// Check the result length
		expectedLen := len(test.expected)
		if len(result) != expectedLen {
			t.Errorf("Expected length %d, got %d", expectedLen, len(result))
		}
		// Check the result content
		for i, val := range test.expected {
			if result[i] != val {
				t.Errorf("Expected %d at index %d, got %d", val, i, result[i])
			}
		}
	}
}

func printMemStats(tag string) {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	fmt.Printf("%s: New = %v KB, Sys = %v KB, NumGC = %v\n",
		tag, m.Alloc/1024, m.Sys/1024, m.NumGC)
}

func TestArenaAlloc_GC(t *testing.T) {
	mem := NewArena()
	defer mem.Close()

	fmt.Println(mem.Stats())
	printMemStats("Before ArenaHashMap Put")
	startPut := time.Now()
	for i := 0; i < 1000000; i++ {
		str := RandomString(100)
		Allocate[string](mem, str)
	}
	fmt.Printf("ArenaHashMap Put Time: %v\n", time.Since(startPut))
	printMemStats("After ArenaHashMap Put")
	fmt.Println(mem.Stats())
	runtime.GC()

	for j := 0; j < 10; j++ {
		startGC := time.Now()
		runtime.GC()
		t.Logf("GC time after map with pointers, iteration %d: %v\n", j+1, time.Since(startGC))
	}
	printMemStats("After ArenaHashMap Manual GC")
}

func TestGoAlloc_GC(t *testing.T) {
	data := map[string]struct{}{}
	for i := 0; i < 1000000; i++ {
		str := RandomString(100)
		data[str] = struct{}{}
	}
	runtime.GC()
	for j := 0; j < 10; j++ {
		startGC := time.Now()
		runtime.GC()
		t.Logf("GC time after map with pointers, iteration %d: %v\n", j+1, time.Since(startGC))
	}
	printMemStats("After ArenaHashMap Manual GC")
	runtime.KeepAlive(data)
}

const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

func RandomString(length int) string {
	result := make([]byte, length)
	for i := 0; i < length; i++ {
		result[i] = charset[rand.Intn(len(charset))]
	}
	return string(result)
}
