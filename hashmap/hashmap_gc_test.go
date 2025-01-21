package hashmap

import (
	"fmt"
	"github.com/ongniud/arena-map/arena"
	"runtime"
	"testing"
	"time"
)

const (
	numEntries = 300000 // 数据量
)

func TestArenaHashMap_PutPerformance(t *testing.T) {
	mem := arena.NewArena()
	defer mem.Close()
	hm := NewArenaHashMap[int, string](mem)

	startPut := time.Now()
	for i := 0; i < numEntries; i++ {
		str := RandomString(1)
		hm.Put(i, str)
	}
	putDuration := time.Since(startPut)
	fmt.Printf("ArenaHashMap Put Time: %v\n", putDuration)
	printMemStats("ArenaHashMap Put Mem Stats")

	startGet := time.Now()
	for i := 0; i < numEntries; i++ {
		hm.Get(i)
	}
	getDuration := time.Since(startGet)
	fmt.Printf("ArenaHashMap Get Time: %v\n", getDuration)

	runtime.KeepAlive(hm)
}

func TestArenaHashMapPtr_PutPerformance(t *testing.T) {
	mem := arena.NewArena()
	defer mem.Close()
	hm := NewArenaHashMapPtr[int, string](mem)

	startPut := time.Now()
	for i := 0; i < numEntries; i++ {
		str := RandomString(1)
		hm.Put(i, str)
	}
	putDuration := time.Since(startPut)
	fmt.Printf("ArenaHashMapPtr Put Time: %v\n", putDuration)
	printMemStats("ArenaHashMapPtr Put Mem Stats")

	startGet := time.Now()
	for i := 0; i < numEntries; i++ {
		hm.Get(i)
	}
	getDuration := time.Since(startGet)
	fmt.Printf("ArenaHashMapPtr Get Time: %v\n", getDuration)

	runtime.KeepAlive(hm)
}

func TestStandardMap_PutPerformance(t *testing.T) {
	hm := make(map[int]string)

	startPut := time.Now()
	for i := 0; i < numEntries; i++ {
		str := RandomString(1)
		hm[i] = str
	}
	putDuration := time.Since(startPut)
	fmt.Printf("Standard Map Put Time: %v\n", putDuration)
	printMemStats("Standard Map Put Mem Stats")

	startGet := time.Now()
	for i := 0; i < numEntries; i++ {
		_ = hm[i]
	}
	getDuration := time.Since(startGet)
	fmt.Printf("Standard Map Get Time: %v\n", getDuration)

	runtime.KeepAlive(hm)
}

// 测试 ArenaHashMap 的 GC 效果
func TestArenaHashMap_GC(t *testing.T) {
	mem := arena.NewArena()
	defer mem.Close()
	hm := NewArenaHashMap[int, string](mem)

	startPut := time.Now()
	for i := 0; i < numEntries; i++ {
		str := RandomString(1)
		hm.Put(i, str)
	}
	fmt.Printf("Put Time: %v\n", time.Since(startPut))
	fmt.Println(mem.Stats())
	printMemStats("Put Mem Stats")
	runtime.GC()

	for j := 0; j < 10; j++ {
		startGC := time.Now()
		runtime.GC()
		t.Logf("GC time, iteration %d: %v\n", j+1, time.Since(startGC))
	}
	printMemStats("GC Mem Stats")

	runtime.KeepAlive(hm)
}

// 测试 ArenaHashMap 的 GC 效果
func TestArenaHashMapPtr_GC(t *testing.T) {
	mem := arena.NewArena()
	defer mem.Close()
	hm := NewArenaHashMapPtr[int, string](mem)

	startPut := time.Now()
	for i := 0; i < numEntries; i++ {
		str := RandomString(1)
		hm.Put(i, str)
	}
	fmt.Printf("Put Time: %v\n", time.Since(startPut))
	fmt.Println(mem.Stats())
	printMemStats("Put Mem Stats")
	runtime.GC()

	for j := 0; j < 10; j++ {
		startGC := time.Now()
		runtime.GC()
		t.Logf("GC time, iteration %d: %v\n", j+1, time.Since(startGC))
	}
	printMemStats("GC Mem Stats")

	runtime.KeepAlive(hm)
}

// 测试 ArenaHashMap 的 GC 效果
func TestStandardMap_GC(t *testing.T) {
	hm := make(map[int]string)
	startPut := time.Now()
	for i := 0; i < numEntries; i++ {
		str := RandomString(1)
		hm[i] = str
	}
	fmt.Printf("Put Time: %v\n", time.Since(startPut))
	printMemStats("Put Mem Stats")
	runtime.GC()

	for j := 0; j < 10; j++ {
		startGC := time.Now()
		runtime.GC()
		t.Logf("GC time, iteration %d: %v\n", j+1, time.Since(startGC))
	}
	printMemStats("GC Mem Stats")
	runtime.KeepAlive(hm)
}

func printMemStats(tag string) {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	fmt.Printf("%s: New = %v KB, Sys = %v KB, NumGC = %v\n",
		tag, m.Alloc/1024, m.Sys/1024, m.NumGC)
}
