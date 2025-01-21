package arena

import (
	"fmt"
	"math/rand"
	"runtime"
	"testing"
	"time"
)

// 打印内存使用信息
func printMemStats(tag string) {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	fmt.Printf("%s: New = %v KB, Sys = %v KB, NumGC = %v\n",
		tag, m.Alloc/1024, m.Sys/1024, m.NumGC)
}

// 测试 ArenaHashMap 的 GC 效果
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

	// 强制多次触发 GC，观察时间是否稳定
	for j := 0; j < 10; j++ {
		startGC := time.Now()
		runtime.GC()
		t.Logf("GC time after map with pointers, iteration %d: %v\n", j+1, time.Since(startGC))
	}

	//// 手动触发 GC
	//startGC := time.Now()
	//runtime.GC()
	//fmt.Printf("ArenaHashMap GC Time: %v\n", time.Since(startGC))
	printMemStats("After ArenaHashMap Manual GC")
}

func TestGoAlloc_GC(t *testing.T) {
	data := map[string]struct{}{}
	for i := 0; i < 1000000; i++ {
		str := RandomString(100)
		data[str] = struct{}{}
	}
	runtime.GC()

	// 强制多次触发 GC，观察时间是否稳定
	for j := 0; j < 10; j++ {
		startGC := time.Now()
		runtime.GC()
		t.Logf("GC time after map with pointers, iteration %d: %v\n", j+1, time.Since(startGC))
	}
	printMemStats("After ArenaHashMap Manual GC")
	runtime.KeepAlive(data)
}

// 定义字符集，包含大写字母、小写字母和数字
const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

// RandomString 生成指定长度的随机字符串
func RandomString(length int) string {
	result := make([]byte, length) // 创建指定长度的字节切片
	// 从字符集随机选取字符填充切片
	for i := 0; i < length; i++ {
		result[i] = charset[rand.Intn(len(charset))] // 从字符集随机选择一个字符
	}
	return string(result) // 将字节切片转为字符串返回
}

// MyStruct 定义了一个简单的结构体，用于测试
type MyStruct struct {
	ID   int
	Name string
}

func TestGCPerformance(t *testing.T) {
	// 初始 GC 执行
	runtime.GC()

	// 记录初始内存状态
	var memStatsBefore runtime.MemStats
	runtime.ReadMemStats(&memStatsBefore)
	t.Logf("Before allocation: New = %v, TotalAlloc = %v, Sys = %v\n",
		memStatsBefore.Alloc, memStatsBefore.TotalAlloc, memStatsBefore.Sys)

	// 模拟插入大量数据到 map 中
	mapWithPointers := make(map[int]*MyStruct)
	for i := 0; i < 10000000; i++ {
		mapWithPointers[i] = &MyStruct{
			ID:   i,
			Name: fmt.Sprintf("Name_%d", i),
		}
	}

	// 强制多次触发 GC，观察时间是否稳定
	for j := 0; j < 10; j++ {
		startGC := time.Now()
		runtime.GC()
		t.Logf("GC time after map with pointers, iteration %d: %v\n", j+1, time.Since(startGC))
	}

	// 记录 GC 后的内存状态
	var memStatsAfterGC runtime.MemStats
	runtime.ReadMemStats(&memStatsAfterGC)
	t.Logf("After GC: New = %v, TotalAlloc = %v, Sys = %v\n",
		memStatsAfterGC.Alloc, memStatsAfterGC.TotalAlloc, memStatsAfterGC.Sys)
}
