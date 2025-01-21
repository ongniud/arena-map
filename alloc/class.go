package alloc

import (
	"errors"
	"fmt"
	"math"
	"unsafe"
)

type class struct {
	alloc     *allocator
	slabs     []Slab
	slabSize  int
	chunkSize int

	free loc
	objs int
}

func newClass(alloc *allocator, ckSize int) class {
	ckCnt := alloc.opt.SlabSize / ckSize
	if ckCnt <= 0 {
		ckCnt = 1
	}
	slabSize := ckSize * ckCnt
	return class{
		alloc:     alloc,
		slabSize:  slabSize,
		chunkSize: ckSize,
		free:      nilLoc,
	}
}

func (sc *class) Slab(i int) *Slab {
	if sc == nil || i >= len(sc.slabs) {
		return nil
	}
	return &sc.slabs[i]
}

func (sc *class) Chunk(loc loc) *chunk {
	return sc.Slab(loc.slabId).Chunk(loc.chunkId)
}

func (sc *class) addSlab() error {
	slabIdx := len(sc.slabs)
	slb, err := NewSlab(sc.alloc, slabIdx, sc.slabSize, sc.chunkSize)
	if err != nil {
		return err
	}
	sc.slabs = append(sc.slabs, slb)
	for i := range slb.chunks {
		sc.pushFreeChunk(&slb.chunks[i])
	}
	return nil
}

func (sc *class) Alloc(align int) (unsafe.Pointer, error) {
	ck, err := sc.getChunk()
	if err != nil {
		return nil, err
	}
	slb := sc.Slab(ck.slabId)
	if slb == nil {
		return nil, errors.New("invalid slab")
	}
	addr := slb.Addr() + uintptr(ck.chunkId*sc.chunkSize)
	sc.objs++
	return unsafe.Pointer(addr), nil
}

func (sc *class) Free(ptr unsafe.Pointer) bool {
	si, ci := sc.locateSlabChunk(ptr)
	if si < 0 || ci < 0 {
		return false
	}
	slb := &sc.slabs[si]
	freed := sc.freeChunk(&slb.chunks[ci])
	if freed {
		sc.objs--
	}
	return freed
}

func (sc *class) locateSlabChunk(ptr unsafe.Pointer) (int, int) {
	addr := uintptr(ptr)
	for i := 0; i < len(sc.slabs); i++ {
		offset := addr - sc.slabs[i].Addr()
		if offset >= 0 && offset < uintptr(sc.slabSize) {
			return i, int(math.Floor(float64(offset) / float64(sc.chunkSize)))
		}
	}
	return -1, -1
}

func (sc *class) getChunk() (*chunk, error) {
	if sc.free.IsNil() {
		if err := sc.addSlab(); err != nil {
			return nil, err
		}
	}
	ck := sc.popFreeChunk()
	if ck == nil {
		return nil, errors.New("no free Chunk")
	}
	return ck, nil
}

func (sc *class) refChunk(c *chunk) {
	c.refs++
}

func (sc *class) freeChunk(c *chunk) bool {
	c.refs--
	if c.refs < 0 {
		panic(fmt.Sprintf("refs < 0 during decRef: %#v", sc))
	}
	if c.refs == 0 {
		sc.pushFreeChunk(c)
		return true
	}
	return false
}

func (sc *class) pushFreeChunk(c *chunk) {
	if c.refs != 0 {
		panic(fmt.Sprintf("pushFreeChunk() non-zero refs: %v", c.refs))
	}
	c.next = sc.free
	sc.free = c.loc
}

func (sc *class) popFreeChunk() *chunk {
	if sc.free.IsNil() {
		panic("popFreeChunk() when free is nil")
	}
	c := sc.Chunk(sc.free)
	if c == nil || c.refs != 0 {
		panic(fmt.Sprintf("popFreeChunk() non-zero refs: %v", c.refs))
	}
	sc.free = c.next
	c.next = nilLoc
	c.refs = 1
	return c
}

func (sc *class) freeChunkCount() int {
	count := 0
	for current := sc.free; !current.IsNil(); {
		current = sc.Chunk(current).next
		count++
	}
	return count
}
