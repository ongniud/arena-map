package arena

func Allocate[T any](mem Arena, data T) *T {
	switch v := any(data).(type) {
	case string:
		bytes := MakeSlice[byte](mem, len(v), len(v))
		copy(bytes, v)
		ptr := New[string](mem)
		*ptr = string(bytes)
		return any(ptr).(*T)
	default:
		allocated := New[T](mem)
		*allocated = data
		return allocated
	}
}

func Release[T any](mem Arena, ptr *T) {
	if ptr == nil {
		return
	}
	switch v := any(*ptr).(type) {
	case string:
		bytes := []byte(v)
		FreeSlice[byte](mem, bytes)
		Free[T](mem, ptr)
	default:
		Free[T](mem, ptr)
	}
}
