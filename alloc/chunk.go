package alloc

type chunk struct {
	loc
	next loc
	// for gc
	refs int32
}

type loc struct {
	slabId  int
	chunkId int
}

var nilLoc = loc{-1, -1} // A sentinel.

// IsNil returns true if the loc came from NilLoc().
func (l loc) IsNil() bool {
	return l.slabId < 0 && l.chunkId < 0
}
