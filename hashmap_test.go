package amap

import (
	"fmt"
	"testing"

	"github.com/ongniud/arena"
)

func TestArenaHashMapBasicOperations(t *testing.T) {
	a := arena.NewArena()
	defer a.Close()

	hm := NewArenaHashMap[int, string](a)

	// 插入键值对
	hm.Put(1, "one")
	hm.Put(2, "two")
	hm.Put(3, "three")

	if val, ok := hm.Get(1); !ok || val != "one" {
		t.Errorf("expected 'one', got '%v'", val)
	}
	if val, ok := hm.Get(2); !ok || val != "two" {
		t.Errorf("expected 'two', got '%v'", val)
	}
	if val, ok := hm.Get(3); !ok || val != "three" {
		t.Errorf("expected 'three', got '%v'", val)
	}
	if _, ok := hm.Get(4); ok {
		t.Errorf("expected key 4 to not exist")
	}
}

func TestArenaHashMapUpdate(t *testing.T) {
	a := arena.NewArena()
	defer a.Close()

	hm := NewArenaHashMap[int, string](a)
	hm.Put(1, "one")
	hm.Put(2, "two")
	hm.Put(2, "updated_two")
	val, ok := hm.Get(2)
	fmt.Println(val, ok)
	if !ok || val != "updated_two" {
		t.Errorf("expected 'updated_two', got '%v'", val)
	}
}

func TestArenaHashMapDelete(t *testing.T) {
	a := arena.NewArena()
	defer a.Close()

	hm := NewArenaHashMap[int, string](a)
	hm.Put(1, "one")
	hm.Put(2, "two")
	hm.Delete(2)
	if _, ok := hm.Get(2); ok {
		t.Errorf("expected key 2 to be deleted")
	}
	hm.Delete(3)
}

func TestArenaHashMapResize(t *testing.T) {
	a := arena.NewArena()
	defer a.Close()

	hm := NewArenaHashMap[int, string](a)

	numElements := InitialBucketSize * 1
	for i := 0; i < numElements; i++ {
		hm.Put(i, string(rune('a'+i)))
	}

	for i := 0; i < numElements; i++ {
		if val, ok := hm.Get(i); !ok || val != string(rune('a'+i)) {
			t.Errorf("expected '%s', got '%v'", string(rune('a'+i)), val)
		}
	}

	if len(hm.buckets) <= InitialBucketSize {
		t.Errorf("expected buckets to resize, got %d", len(hm.buckets))
	}
}

func TestArenaHashMapEdgeCases(t *testing.T) {
	a := arena.NewArena()
	defer a.Close()

	hm := NewArenaHashMap[string, string](a)

	hm.Put("", "")
	if val, ok := hm.Get(""); !ok || val != "" {
		t.Errorf("expected '', got '%v'", val)
	}

	hm.Put("key", "")
	if val, ok := hm.Get("key"); !ok || val != "" {
		t.Errorf("expected '', got '%v'", val)
	}

	hm.Delete("nonexistent")
}

func TestArenaHashMapFree(t *testing.T) {
	a := arena.NewArena()
	defer a.Close()

	hm := NewArenaHashMap[int, string](a)

	for i := 0; i < 10; i++ {
		hm.Put(i, string(rune('a'+i)))
	}

	hm.Free()
	if hm.buckets != nil {
		t.Errorf("expected buckets to be nil after Free")
	}
	hm.Free()
}
