package amap

import (
	"math/rand"
	"strconv"
	"testing"

	"github.com/ongniud/arena"
)

const (
	numOperations = 100
)

func generateTestData(num int) []string {
	data := make([]string, num)
	for i := 0; i < num; i++ {
		data[i] = strconv.Itoa(i)
	}
	return data
}

// Benchmark for Put operation
func BenchmarkArenaHashMap_Put(b *testing.B) {
	mem := arena.NewArena()
	defer mem.Close()

	hmap := NewArenaHashMap[string, string](mem)
	testData := generateTestData(numOperations)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		key := testData[rand.Intn(len(testData))]
		value := "value" + key
		hmap.Put(key, value)
	}
}

// Benchmark for Get operation
func BenchmarkArenaHashMap_Get(b *testing.B) {
	mem := arena.NewArena()
	defer mem.Close()
	hmap := NewArenaHashMap[string, string](mem)
	testData := generateTestData(numOperations)

	// Prepopulate the hashmap
	for _, key := range testData {
		hmap.Put(key, "value"+key)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		key := testData[rand.Intn(len(testData))]
		hmap.Get(key)
	}
}

// Benchmark for Delete operation
func BenchmarkArenaHashMap_Delete(b *testing.B) {
	mem := arena.NewArena()
	defer mem.Close()
	hmap := NewArenaHashMap[string, string](mem)
	testData := generateTestData(numOperations)

	// Prepopulate the hashmap
	for _, key := range testData {
		hmap.Put(key, "value"+key)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		key := testData[rand.Intn(len(testData))]
		hmap.Delete(key)
	}
}

// Comparison with native Go map
func BenchmarkNativeMap_Put(b *testing.B) {
	hmap := make(map[string]string)
	testData := generateTestData(numOperations)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		key := testData[rand.Intn(len(testData))]
		value := "value" + key
		hmap[key] = value
	}
}

func BenchmarkNativeMap_Get(b *testing.B) {
	hmap := make(map[string]string)
	testData := generateTestData(numOperations)

	// Prepopulate the native map
	for _, key := range testData {
		hmap[key] = "value" + key
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		key := testData[rand.Intn(len(testData))]
		_ = hmap[key]
	}
}

func BenchmarkNativeMap_Delete(b *testing.B) {
	hmap := make(map[string]string)
	testData := generateTestData(numOperations)

	// Prepopulate the native map
	for _, key := range testData {
		hmap[key] = "value" + key
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		key := testData[rand.Intn(len(testData))]
		delete(hmap, key)
	}
}

func TestArenaHashMap_Functionality(t *testing.T) {
	mem := arena.NewArena()
	defer mem.Close()
	hmap := NewArenaHashMap[string, string](mem)
	testData := generateTestData(100)

	for _, key := range testData {
		value := "value" + key
		hmap.Put(key, value)

		// Validate Get
		gotValue, ok := hmap.Get(key)
		if !ok || gotValue != value {
			t.Fatalf("Get failed for key %s, got: %s", key, gotValue)
		}
	}

	// Validate Delete
	for _, key := range testData {
		hmap.Delete(key)
		_, ok := hmap.Get(key)
		if ok {
			t.Fatalf("Delete failed for key %s", key)
		}
	}
}
