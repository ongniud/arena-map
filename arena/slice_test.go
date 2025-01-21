package arena

import (
	"fmt"
	"testing"
)

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
