package arena

import (
	"unsafe"

	"github.com/ongniud/arena-map/alloc"
)

type Arena interface {
	New(size, align int) unsafe.Pointer
	Free(ptr unsafe.Pointer, size int)
	Close()
	Stats() string
}

func NewArena() Arena {
	return &slabArena{
		alloc: alloc.NewAllocator(alloc.WithSlabSize(16 * alloc.MB)),
	}
}

type slabArena struct {
	alloc alloc.Allocator
}

func (s *slabArena) New(size, align int) unsafe.Pointer {
	ptr, err := s.alloc.Alloc(size, align)
	if err != nil {
		return nil
	}
	return ptr
}

func (s *slabArena) Free(ptr unsafe.Pointer, size int) {
	if err := s.alloc.Free(ptr, size); err != nil {
		_ = err
	}
}

func (s *slabArena) Close() {
}

func (s *slabArena) Stats() string {
	return s.alloc.Stats()
}
