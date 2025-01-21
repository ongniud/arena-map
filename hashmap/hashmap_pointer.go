package hashmap

import (
	"github.com/ongniud/arena-map/arena"
)

type nodePtr[K comparable, V any] struct {
	key   *K
	value *V
	next  *nodePtr[K, V]
}

type ArenaHashMapPtr[K comparable, V any] struct {
	buckets []*nodePtr[K, V]
	mem     arena.Arena

	size      int
	factor    float64
	threshold int
}

func NewArenaHashMapPtr[K comparable, V any](mem arena.Arena) *ArenaHashMapPtr[K, V] {
	buckets := arena.MakeSlice[*nodePtr[K, V]](mem, InitialBucketSize, InitialBucketSize)
	return &ArenaHashMapPtr[K, V]{
		buckets:   buckets,
		mem:       mem,
		size:      0,
		factor:    DefaultLoadFactor,
		threshold: int(DefaultLoadFactor * InitialBucketSize),
	}
}

func (h *ArenaHashMapPtr[K, V]) Free() {
	if h.buckets == nil {
		return
	}
	for _, head := range h.buckets {
		for n := head; n != nil; {
			next := n.next
			arena.Release[K](h.mem, n.key)
			arena.Release[V](h.mem, n.value)
			arena.Release[nodePtr[K, V]](h.mem, n)
			n = next
		}
	}
	arena.FreeSlice[*nodePtr[K, V]](h.mem, h.buckets)
	h.buckets = nil
}

func (h *ArenaHashMapPtr[K, V]) Put(key K, value V) {
	index := hash(key) % len(h.buckets)
	head := h.buckets[index]

	newKey := arena.Allocate[K](h.mem, key)
	newValue := arena.Allocate[V](h.mem, value)

	for n := head; n != nil; n = n.next {
		if *n.key == key {
			arena.Release[V](h.mem, n.value)
			n.value = newValue
			return
		}
	}

	newNode := arena.New[nodePtr[K, V]](h.mem)
	newNode.key = newKey
	newNode.value = newValue
	newNode.next = head

	h.buckets[index] = newNode

	// 检查负载因子并扩容
	if h.size++; h.size > h.threshold {
		h.resize()
	}
}

func (h *ArenaHashMapPtr[K, V]) Get(key K) (V, bool) {
	index := hash(key) % len(h.buckets)
	head := h.buckets[index]
	for n := head; n != nil; n = n.next {
		if *n.key == key {
			return *n.value, true
		}
	}
	var zero V
	return zero, false
}

func (h *ArenaHashMapPtr[K, V]) Delete(key K) {
	index := hash(key) % len(h.buckets)
	head := h.buckets[index]
	var prev *nodePtr[K, V]
	for n := head; n != nil; n = n.next {
		if *n.key == key {
			if prev == nil {
				h.buckets[index] = n.next
			} else {
				prev.next = n.next
			}
			arena.Release[K](h.mem, n.key)
			arena.Release[V](h.mem, n.value)
			arena.Release[nodePtr[K, V]](h.mem, n)
			h.size--
			return
		}
		prev = n
	}
}

func (h *ArenaHashMapPtr[K, V]) resize() {
	newSize := len(h.buckets) * 2
	newBuckets := arena.MakeSlice[*nodePtr[K, V]](h.mem, newSize, newSize)

	for _, bucket := range h.buckets {
		for n := bucket; n != nil; {
			next := n.next
			index := hash(*n.key) % newSize
			n.next = newBuckets[index]
			newBuckets[index] = n
			n = next
		}
	}

	arena.FreeSlice[*nodePtr[K, V]](h.mem, h.buckets)
	h.buckets = newBuckets
	h.threshold = int(float64(newSize) * h.factor)
}
