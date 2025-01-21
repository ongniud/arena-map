package alloc

import (
	"errors"
	"fmt"
	"math"
	"sort"
	"strings"
	"unsafe"
)

type allocator struct {
	classes []class
	opt     *Options
}

func NewAllocator(opts ...Option) Allocator {
	options := newDefaultOptions()
	for _, opt := range opts {
		opt(options)
	}
	return &allocator{
		opt: options,
	}
}

func (a *allocator) Alloc(size, align int) (unsafe.Pointer, error) {
	// Find the appropriate slab class for the requested size.
	i := a.findClass(size, true)
	cla := &a.classes[i]
	// Allocate memory from the selected slab class.
	return cla.Alloc(align)
}

func (a *allocator) Free(ptr unsafe.Pointer, size int) error {
	// Find the appropriate slab class for the requested size.
	i := a.findClass(size, false)
	if i < 0 {
		return errors.New("not from this alloc")
	}
	cla := &a.classes[i]
	// Free the memory in the slab class.
	cla.Free(ptr)
	return nil
}

func (a *allocator) malloc(size int) []byte {
	return a.opt.Malloc(size)
}

func (a *allocator) findClass(size int, create bool) int {
	i := sort.Search(len(a.classes), func(i int) bool { return size <= a.classes[i].chunkSize })
	if i < len(a.classes) {
		return i
	}
	if !create {
		return -1
	}
	var chunkSize int
	if len(a.classes) == 0 {
		chunkSize = 1
	} else {
		last := len(a.classes) - 1
		chunkSize = int(math.Ceil(float64(a.classes[last].chunkSize) * a.opt.GrowthFactor))
	}
	a.classes = append(a.classes, newClass(a, chunkSize))
	return a.findClass(size, create)
}

func (a *allocator) Stats() string {
	var stats []string
	totalSlabs := 0
	totalObjs := 0
	totalAllocated := 0
	totalFree := 0
	totalUsed := 0

	for i, cla := range a.classes {
		slabCount := len(cla.slabs)
		totalSlabs += slabCount
		objCount := cla.objs
		totalObjs += objCount
		allocated := slabCount * cla.slabSize
		totalAllocated += allocated
		free := cla.freeChunkCount() * cla.chunkSize
		totalFree += free
		used := allocated - free
		totalUsed += used
		stat := fmt.Sprintf(
			"class:%d, slabs:%d, objs:%d, slabSize:%d, chunkSize:%d, allocated:%d, free:%d, used:%d",
			i, slabCount, objCount, cla.slabSize, cla.chunkSize, allocated, free, used,
		)
		stats = append(stats, stat)
	}
	totalStat := fmt.Sprintf(
		"Total slabs: %d, Total objects: %d, Total allocated: %d bytes, Total free: %d bytes, Total used: %d bytes",
		totalSlabs, totalObjs, totalAllocated, totalFree, totalUsed,
	)
	return strings.Join(append(stats, totalStat), "\n")
}
