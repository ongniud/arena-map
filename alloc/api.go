package alloc

import "unsafe"

type Allocator interface {
	// Alloc allocates memory of the given size and returns a pointer to it.
	// The alignment parameter specifies the alignment of the allocated memory.
	Alloc(size, alignment int) (unsafe.Pointer, error)
	// Free release the memory.
	Free(ptr unsafe.Pointer, size int) error
	Stats() string
}
