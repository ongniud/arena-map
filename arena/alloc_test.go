package arena

import (
	"fmt"
	"math/rand"
	"runtime"
	"testing"
	"time"
)

func printMemStats(tag string) {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	fmt.Printf("%s: New = %v KB, Sys = %v KB, NumGC = %v\n",
		tag, m.Alloc/1024, m.Sys/1024, m.NumGC)
}

func TestArenaAlloc_GC(t *testing.T) {
	mem := NewArena()
	defer mem.Close()

	fmt.Println(mem.Stats())
	printMemStats("Before ArenaHashMap Put")
	startPut := time.Now()
	for i := 0; i < 1000000; i++ {
		str := RandomString(100)
		Allocate[string](mem, str)
	}
	fmt.Printf("ArenaHashMap Put Time: %v\n", time.Since(startPut))
	printMemStats("After ArenaHashMap Put")
	fmt.Println(mem.Stats())
	runtime.GC()

	for j := 0; j < 10; j++ {
		startGC := time.Now()
		runtime.GC()
		t.Logf("GC time after map with pointers, iteration %d: %v\n", j+1, time.Since(startGC))
	}
	printMemStats("After ArenaHashMap Manual GC")
}

func TestGoAlloc_GC(t *testing.T) {
	data := map[string]struct{}{}
	for i := 0; i < 1000000; i++ {
		str := RandomString(100)
		data[str] = struct{}{}
	}
	runtime.GC()
	for j := 0; j < 10; j++ {
		startGC := time.Now()
		runtime.GC()
		t.Logf("GC time after map with pointers, iteration %d: %v\n", j+1, time.Since(startGC))
	}
	printMemStats("After ArenaHashMap Manual GC")
	runtime.KeepAlive(data)
}

const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

func RandomString(length int) string {
	result := make([]byte, length)
	for i := 0; i < length; i++ {
		result[i] = charset[rand.Intn(len(charset))]
	}
	return string(result)
}
