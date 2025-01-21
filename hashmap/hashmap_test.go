package hashmap

import (
	"fmt"
	"testing"

	"github.com/ongniud/arena-map/arena"
)

func TestArenaHashMapBasicOperations(t *testing.T) {
	a := arena.NewArena()
	defer a.Close()

	hm := NewArenaHashMap[int, string](a)

	// 插入键值对
	hm.Put(1, "one")
	hm.Put(2, "two")
	hm.Put(3, "three")

	// 测试 Get
	if val, ok := hm.Get(1); !ok || val != "one" {
		t.Errorf("expected 'one', got '%v'", val)
	}
	if val, ok := hm.Get(2); !ok || val != "two" {
		t.Errorf("expected 'two', got '%v'", val)
	}
	if val, ok := hm.Get(3); !ok || val != "three" {
		t.Errorf("expected 'three', got '%v'", val)
	}

	// 测试不存在的 Key
	if _, ok := hm.Get(4); ok {
		t.Errorf("expected key 4 to not exist")
	}
}

func TestArenaHashMapUpdate(t *testing.T) {
	a := arena.NewArena()
	defer a.Close()

	hm := NewArenaHashMap[int, string](a)

	// 插入键值对
	hm.Put(1, "one")
	hm.Put(2, "two")

	// 更新键值对
	hm.Put(2, "updated_two")

	// 测试更新是否成功
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

	// 插入键值对
	hm.Put(1, "one")
	hm.Put(2, "two")

	// 删除键值对
	hm.Delete(2)

	// 测试是否删除成功
	if _, ok := hm.Get(2); ok {
		t.Errorf("expected key 2 to be deleted")
	}

	// 测试删除不存在的键
	hm.Delete(3) // 不应该崩溃
}

func TestArenaHashMapResize(t *testing.T) {
	a := arena.NewArena()
	defer a.Close()

	hm := NewArenaHashMap[int, string](a)

	// 插入多个键值对，触发扩容
	numElements := InitialBucketSize * 1 // 超过阈值以触发扩容
	for i := 0; i < numElements; i++ {
		hm.Put(i, string(rune('a'+i)))
	}

	// 检查所有元素是否都能正确获取
	for i := 0; i < numElements; i++ {
		if val, ok := hm.Get(i); !ok || val != string(rune('a'+i)) {
			t.Errorf("expected '%s', got '%v'", string(rune('a'+i)), val)
		}
	}

	// 检查扩容是否正确
	if len(hm.buckets) <= InitialBucketSize {
		t.Errorf("expected buckets to resize, got %d", len(hm.buckets))
	}
}

func TestArenaHashMapEdgeCases(t *testing.T) {
	a := arena.NewArena()
	defer a.Close()

	hm := NewArenaHashMap[string, string](a)

	// 插入空字符串键和值
	hm.Put("", "")
	if val, ok := hm.Get(""); !ok || val != "" {
		t.Errorf("expected '', got '%v'", val)
	}

	// 插入 nil 值 (模拟 nil 的场景）
	hm.Put("key", "")
	if val, ok := hm.Get("key"); !ok || val != "" {
		t.Errorf("expected '', got '%v'", val)
	}

	// 删除不存在的键
	hm.Delete("nonexistent") // 不应该崩溃
}

func TestArenaHashMapFree(t *testing.T) {
	a := arena.NewArena()
	defer a.Close()

	hm := NewArenaHashMap[int, string](a)

	// 插入多个键值对
	for i := 0; i < 10; i++ {
		hm.Put(i, string(rune('a'+i)))
	}

	// 释放 HashMap
	hm.Free()

	// 确保 buckets 被清空
	if hm.buckets != nil {
		t.Errorf("expected buckets to be nil after Free")
	}

	// 确保再次释放不会崩溃
	hm.Free()
}
