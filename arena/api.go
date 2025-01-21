package arena

import (
	"reflect"
	"unsafe"
)

// New allocates memory for a value of type T using the provided Arena.
// If the arena is non-nil, it returns a  *T pointer with memory allocated from the arena.
// If passed arena is nil, it allocates memory using Go's built-in new function.
func New[T any](a Arena) *T {
	if a == nil {
		return new(T)
	}
	elem := reflect.TypeOf((*T)(nil)).Elem()
	size, align := elem.Size(), elem.Align()
	ptr := a.New(int(size), align)
	if ptr == nil {
		return nil
	}
	return (*T)(ptr)
}

// Free releases the memory allocated for a value of type T using the provided Arena.
// It assumes that the pointer was obtained from the New function.
func Free[T any](a Arena, ptr *T) {
	if a != nil && ptr != nil {
		a.Free(unsafe.Pointer(ptr), int(unsafe.Sizeof(*ptr)))
	}
}

// FreeAny releases the memory allocated for a value using the provided Arena.
// It assumes that the pointer was obtained from the New function.
func FreeAny(a Arena, ptr any) {
	if a != nil && ptr != nil {
		val := reflect.ValueOf(ptr)
		if val.Kind() == reflect.Ptr {
			size := int(val.Elem().Type().Size())
			a.Free(unsafe.Pointer(val.Pointer()), size)
		}
	}
}

// Copy allocates memory for a value of type T using the provided Arena,
// then copies the value from the source pointer to the new memory location.
// If the arena is non-nil, it returns a pointer to the newly allocated value.
// If the arena is nil, it returns nil.
func Copy[T any](a Arena, src *T) *T {
	if a == nil || src == nil {
		return nil
	}
	elem := reflect.TypeOf((*T)(nil)).Elem()
	size, align := elem.Size(), elem.Align()
	ptr := a.New(int(size), align)
	if ptr == nil {
		return nil
	}
	dst := (*T)(ptr)
	*dst = *src
	return dst
}

// CopyAny allocates memory for the value of the provided source variable
// using the provided Arena, and copies the value into the newly allocated memory.
// It returns the pointer to the copied value or nil if allocation fails.
func CopyAny(a Arena, src any) any {
	if a == nil || src == nil {
		return nil
	}
	srcVal := reflect.ValueOf(src)
	if srcVal.Kind() != reflect.Ptr || !srcVal.IsValid() {
		return nil
	}
	elemType := srcVal.Elem().Type()
	size, align := elemType.Size(), elemType.Align()
	ptr := a.New(int(size), align)
	if ptr == nil {
		return nil
	}
	dst := reflect.NewAt(elemType, ptr)
	dst.Elem().Set(srcVal.Elem())
	return dst.Interface()
}

// MakeSlice creates a slice of type T with a given length and capacity,
// using the provided Arena for memory allocation.
// If the arena is non-nil, it returns a slice with memory allocated from the arena.
// Otherwise, it returns a slice using Go's built-in make function.
func MakeSlice[T any](a Arena, len, cap int) []T {
	if a == nil {
		return make([]T, len, cap)
	}
	elem := reflect.TypeOf((*T)(nil)).Elem()
	size := int(elem.Size()) * cap
	ptr := a.New(size, elem.Align())
	if ptr == nil {
		return nil
	}
	s := unsafe.Slice((*T)(ptr), cap)
	return s[:len]
}

// FreeSlice releases the memory allocated for a slice of type T using the provided Arena.
// It assumes that the pointer was obtained from the MakeSlice function.
func FreeSlice[T any](a Arena, slice []T) {
	if a != nil && len(slice) > 0 {
		val := reflect.ValueOf(slice)
		ptr := val.UnsafePointer()
		size := int(val.Type().Elem().Size()) * val.Cap()
		a.Free(ptr, size)
	}
}

// FreeAnySlice releases the memory allocated for a slice using the provided Arena.
// It assumes that the pointer was obtained from the MakeSlice function.
func FreeAnySlice(a Arena, slice any) {
	if a != nil && slice != nil {
		val := reflect.ValueOf(slice)
		if val.Kind() == reflect.Slice && val.Len() > 0 {
			ptr := val.UnsafePointer()
			elemSize := val.Type().Elem().Size()
			size := int(elemSize) * val.Cap()
			a.Free(ptr, size)
		}
	}
}
