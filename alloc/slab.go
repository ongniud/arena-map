package alloc

import (
	"errors"
	"unsafe"
)

type Slab struct {
	memory []byte
	chunks []chunk
}

func NewSlab(alloc *allocator, slabIdx int, slabSize, chunkSize int) (Slab, error) {
	mem := alloc.malloc(slabSize)
	if mem == nil {
		return Slab{}, errors.New("malloc fail")
	}
	chunks := make([]chunk, slabSize/chunkSize)
	for i := range chunks {
		chunks[i].loc.slabId = slabIdx
		chunks[i].loc.chunkId = i
	}
	return Slab{
		memory: mem,
		chunks: chunks,
	}, nil
}

func (s *Slab) Addr() uintptr {
	return uintptr(unsafe.Pointer(&s.memory[0]))
}

func (s *Slab) Chunk(i int) *chunk {
	if s == nil || i >= len(s.chunks) {
		return nil
	}
	return &s.chunks[i]
}
