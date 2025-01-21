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

	// inner
	free loc

	objs int
}

func newClass(alloc *allocator, chunkSize int) class {
	chunkCount := alloc.opt.SlabSize / chunkSize
	if chunkCount <= 0 {
		chunkCount = 1
	}
	slabSize := chunkSize * chunkCount
	//fmt.Printf("Creating new slab class: chunkSize=%d, chunkCount=%d, slabSize=%d\n", chunkSize, chunkCount, slabSize)
	return class{
		alloc:     alloc,
		slabSize:  slabSize,
		chunkSize: chunkSize,
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
		//fmt.Printf("Initialized Chunk: slabId=%d, chunkId=%d\n", slabId, i)
	}
	//fmt.Printf("Added new slab: idx=%d, len(memory)=%d, free Chunk=%d\n", slabIdx, sc.freeChunkCount())
	return nil
}

func (sc *class) Alloc(align int) (unsafe.Pointer, error) {
	ck, err := sc.getChunk()
	if err != nil {
		return nil, err
	}
	//fmt.Printf("New: slabId=%d, chunkId=%d\n", ck.loc.slabId, ck.loc.chunkId)
	slb := sc.Slab(ck.slabId)
	if slb == nil {
		return nil, errors.New("invalid slab")
	}
	addr := slb.Addr() + uintptr(ck.chunkId*sc.chunkSize)
	//fmt.Printf("New: ptr=%v, len(slabs)=%d\n", uintptr(ptr), len(sc.slabs))
	sc.objs++
	return unsafe.Pointer(addr), nil
}

func (sc *class) Free(ptr unsafe.Pointer) bool {
	//fmt.Printf("Free: ptr=%v\n", uintptr(ptr))
	si, ci := sc.locateSlabChunk(ptr)
	//fmt.Printf("Free: located slabIdx=%d, chunkIdx=%d\n", si, ci)
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
	//fmt.Printf("locateSlabChunk: ptr=%v\n", uintptr(ptr))
	addr := uintptr(ptr)
	for i := 0; i < len(sc.slabs); i++ {
		//fmt.Printf("Checking slab[%d]: ptr=%v\n", i, sc.slabs[i].Addr())
		offset := addr - sc.slabs[i].Addr()
		if offset >= 0 && offset < uintptr(sc.slabSize) {
			//fmt.Printf("Found slab[%d]: offset=%v\n", i, uintptr(offset))
			return i, int(math.Floor(float64(offset) / float64(sc.chunkSize)))
		}
	}
	return -1, -1
}

func (sc *class) getChunk() (*chunk, error) {
	if sc.free.IsNil() {
		//fmt.Println("No free chunks available, adding a new slab...")
		if err := sc.addSlab(); err != nil {
			return nil, err
		}
	}
	ck := sc.popFreeChunk()
	if ck == nil {
		return nil, errors.New("no free Chunk")
	}
	//free := sc.freeChunkCount()
	//fmt.Printf("Allocated Chunk: slabId=%d, chunkId=%d, free=%d\n", ck.loc.slabId, ck.loc.chunkId, free)
	return ck, nil
}

func (sc *class) refChunk(c *chunk) {
	c.refs++
}

func (sc *class) freeChunk(c *chunk) bool {
	c.refs--
	//fmt.Printf("Decrementing refs: current refs=%d\n", c.refs)
	if c.refs < 0 {
		panic(fmt.Sprintf("refs < 0 during decRef: %#v", sc))
	}
	if c.refs == 0 {
		//fmt.Printf("Chunk is now free: slabId=%d, chunkId=%d\n", c.loc.slabId, c.loc.chunkId)
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
	//fmt.Printf("Pushed Chunk to free list: slabId=%d, chunkId=%d\n", c.self.slabId, c.self.chunkId)
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
	//fmt.Printf("Popped Chunk from free list: slabId=%d, chunkId=%d\n", c.loc.slabId, c.loc.chunkId)
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
