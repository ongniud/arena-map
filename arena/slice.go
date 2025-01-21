package arena

const growThreshold = 256

// SliceAppend appends elements to a slice of type T using a provided Arena
// for memory allocation if needed.
func SliceAppend[T any](a Arena, s []T, data ...T) ([]T, bool) {
	if a == nil {
		return append(s, data...), false
	}
	if l := len(s) + len(data); l > cap(s) {
		r := expandSlice(a, s, l)
		return append(r, data...), true
	}
	return append(s, data...), false
}

func expandSlice[T any](a Arena, s []T, l int) []T {
	c := cap(s)
	if c == 0 {
		c = 1
	}
	for l > c {
		if c < growThreshold {
			c *= 2
		} else {
			c += c / 4
		}
	}
	s2 := MakeSlice[T](a, len(s), c)
	copy(s2, s)
	return s2
}
