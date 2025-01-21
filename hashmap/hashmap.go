package hashmap

import (
	"github.com/ongniud/arena-map/arena"
)

const InitialBucketSize = 16 // 固定大小的哈希桶
const DefaultLoadFactor = 0.75

type node[K comparable, V any] struct {
	key   uintptr
	value uintptr
	next  uintptr
}

type ArenaHashMap[K comparable, V any] struct {
	buckets []uintptr
	mem     arena.Arena

	size      int
	factor    float64
	threshold int
}

func NewArenaHashMap[K comparable, V any](mem arena.Arena) *ArenaHashMap[K, V] {
	rawBuckets := arena.MakeSlice[*node[K, V]](mem, InitialBucketSize, InitialBucketSize)
	buckets := arena.MakeSlice[uintptr](mem, InitialBucketSize, InitialBucketSize)
	for i := range rawBuckets {
		buckets[i] = tToPtr(rawBuckets[i])
	}
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
	for _, bucket := range h.buckets {
		for n := ptrToT[node[K, V]](bucket); n != nil; {
			next := n.next
			arena.Release[K](h.mem, ptrToT[K](n.key))
			arena.Release[V](h.mem, ptrToT[V](n.value))
			arena.Release[node[K, V]](h.mem, n)
			n = ptrToT[node[K, V]](next)
		}
	}
	arena.FreeSlice[uintptr](h.mem, h.buckets)
	h.buckets = nil
}

func (h *ArenaHashMap[K, V]) Put(key K, value V) {
	index := hash(key) % len(h.buckets)
	head := h.buckets[index]

	//fmt.Println("put", key, value)
	//fmt.Println("idx", index)

	for n := ptrToT[node[K, V]](head); n != nil; n = ptrToT[node[K, V]](n.next) {
		if *ptrToT[K](n.key) == key {
			arena.Release[V](h.mem, ptrToT[V](n.value))
			n.value = tToPtr(arena.Allocate[V](h.mem, value))
			return
		}
	}

	newKey := arena.Allocate[K](h.mem, key)
	newValue := arena.Allocate[V](h.mem, value)
	newNode := arena.New[node[K, V]](h.mem)
	newNode.key = tToPtr(newKey)
	newNode.value = tToPtr(newValue)
	newNode.next = head
	h.buckets[index] = tToPtr(newNode)

	//fmt.Println("kvs")
	//for n := ptrToT[node[K, V]](h.buckets[index]); n != nil; n = ptrToT[node[K, V]](n.next) {
	//	k := *ptrToT[K](n.key)
	//	v := *ptrToT[V](n.value)
	//	fmt.Println(k, v)
	//}

	if h.size++; h.size > h.threshold {
		h.resize()
	}

	//fmt.Println()
}

func (h *ArenaHashMap[K, V]) Get(key K) (V, bool) {
	index := hash(key) % len(h.buckets)
	head := h.buckets[index]

	//fmt.Println("get", key)
	//fmt.Println("idx", index)
	//fmt.Println("kvs")

	for n := ptrToT[node[K, V]](head); n != nil; n = ptrToT[node[K, V]](n.next) {
		//k := *ptrToT[K](n.key)
		//v := *ptrToT[V](n.value)
		//fmt.Println(k, v)
		if *ptrToT[K](n.key) == key {
			return *ptrToT[V](n.value), true
		}
	}

	//fmt.Println()
	var zero V
	return zero, false
}

func (h *ArenaHashMap[K, V]) Delete(key K) {
	index := hash(key) % len(h.buckets)
	head := h.buckets[index]
	var prev *node[K, V]
	for n := ptrToT[node[K, V]](head); n != nil; n = ptrToT[node[K, V]](n.next) {
		if *ptrToT[K](n.key) == key {
			if prev == nil {
				h.buckets[index] = n.next
			} else {
				prev.next = n.next
			}
			arena.Release[K](h.mem, ptrToT[K](n.key))
			arena.Release[V](h.mem, ptrToT[V](n.value))
			arena.Release[node[K, V]](h.mem, n)
			h.size--
			return
		}
		prev = n
	}
}

func (h *ArenaHashMap[K, V]) resize() {
	newSize := len(h.buckets) * 2
	newBuckets := arena.MakeSlice[uintptr](h.mem, newSize, newSize)

	for _, bucket := range h.buckets {
		for n := ptrToT[node[K, V]](bucket); n != nil; {
			next := n.next
			index := hash(ptrToT[K](n.key)) % newSize
			n.next = newBuckets[index]
			newBuckets[index] = tToPtr(n)
			n = ptrToT[node[K, V]](next)
		}
	}

	arena.FreeSlice[uintptr](h.mem, h.buckets)
	h.buckets = newBuckets
	h.threshold = int(float64(newSize) * h.factor)
}
