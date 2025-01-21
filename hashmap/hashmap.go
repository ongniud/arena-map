package hashmap

import (
	"github.com/ongniud/arena-map/arena"
)

const InitialBucketSize = 16
const DefaultLoadFactor = 0.75

type node[K comparable, V any] struct {
	key   *K
	value *V
	next  *node[K, V]
}

type ArenaHashMap[K comparable, V any] struct {
	buckets []*node[K, V]
	mem     arena.Arena

	size      int
	factor    float64
	threshold int
}

func NewArenaHashMap[K comparable, V any](mem arena.Arena) *ArenaHashMap[K, V] {
	buckets := arena.MakeSlice[*node[K, V]](mem, InitialBucketSize, InitialBucketSize)
	return &ArenaHashMap[K, V]{
		buckets:   buckets,
		mem:       mem,
		size:      0,
		factor:    DefaultLoadFactor,
		threshold: int(DefaultLoadFactor * InitialBucketSize),
	}
}

func (h *ArenaHashMap[K, V]) Free() {
	if h.buckets == nil {
		return
	}
	for _, head := range h.buckets {
		for n := head; n != nil; {
			next := n.next
			arena.Release[K](h.mem, n.key)
			arena.Release[V](h.mem, n.value)
			arena.Release[node[K, V]](h.mem, n)
			n = next
		}
	}
	arena.FreeSlice[*node[K, V]](h.mem, h.buckets)
	h.buckets = nil
}

func (h *ArenaHashMap[K, V]) Put(key K, value V) {
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

	newNode := arena.New[node[K, V]](h.mem)
	newNode.key = newKey
	newNode.value = newValue
	newNode.next = head
	h.buckets[index] = newNode

	if h.size++; h.size > h.threshold {
		h.resize()
	}
}

func (h *ArenaHashMap[K, V]) Get(key K) (V, bool) {
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

func (h *ArenaHashMap[K, V]) Delete(key K) {
	index := hash(key) % len(h.buckets)
	head := h.buckets[index]
	var prev *node[K, V]
	for n := head; n != nil; n = n.next {
		if *n.key == key {
			if prev == nil {
				h.buckets[index] = n.next
			} else {
				prev.next = n.next
			}
			arena.Release[K](h.mem, n.key)
			arena.Release[V](h.mem, n.value)
			arena.Release[node[K, V]](h.mem, n)
			h.size--
			return
		}
		prev = n
	}
}

func (h *ArenaHashMap[K, V]) resize() {
	newSize := len(h.buckets) * 2
	newBuckets := arena.MakeSlice[*node[K, V]](h.mem, newSize, newSize)

	for _, bucket := range h.buckets {
		for n := bucket; n != nil; {
			next := n.next
			index := hash(*n.key) % newSize
			n.next = newBuckets[index]
			newBuckets[index] = n
			n = next
		}
	}

	arena.FreeSlice[*node[K, V]](h.mem, h.buckets)
	h.buckets = newBuckets
	h.threshold = int(float64(newSize) * h.factor)
}
