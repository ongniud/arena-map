package alloc

import (
	"fmt"
	"testing"
	"unsafe"
)

func TestAllocator_AllocAndFree(t *testing.T) {
	alloc := NewAllocator()

	// Test allocation
	ptr, err := alloc.Alloc(64, 8)
	if err != nil || ptr == nil {
		t.Fatalf("New failed: %v", err)
	}

	// Test freeing the allocated memory
	err = alloc.Free(ptr, 64)
	if err != nil {
		t.Fatalf("Free failed: %v", err)
	}
}

func TestAllocator_Stats(t *testing.T) {
	alloc := NewAllocator()

	// Allocate some memory
	_, err := alloc.Alloc(128, 8)
	if err != nil {
		t.Fatalf("New failed: %v", err)
	}

	// Check stats
	stats := alloc.Stats()
	if stats == "" {
		t.Fatal("Stats should not be empty")
	}
}

func TestAllocator_FreeInvalidPointer(t *testing.T) {
	alloc := NewAllocator()

	// Attempt to free a nil pointer
	err := alloc.Free(nil, 64)
	if err == nil {
		t.Fatal("Expected error when freeing a nil pointer")
	}

	// Attempt to free a random pointer
	var invalidPtr unsafe.Pointer = unsafe.Pointer(new(int))
	err = alloc.Free(invalidPtr, 64)
	if err == nil {
		t.Fatal("Expected error when freeing an invalid pointer")
	}
}

func TestAllocator_Growth(t *testing.T) {
	alloc := NewAllocator()

	// Allocate increasing sizes to trigger slab class growth
	sizes := []int{1, 2, 4, 8, 16, 32, 64, 128, 256}
	for _, size := range sizes {
		_, err := alloc.Alloc(size, 8)
		if err != nil {
			t.Fatalf("New failed for size %d: %v", size, err)
		}
	}

	// Check stats to ensure growth
	stats := alloc.Stats()
	if len(stats) == 0 {
		t.Fatal("Stats should not be empty after allocations")
	}

	fmt.Println(stats)
}
