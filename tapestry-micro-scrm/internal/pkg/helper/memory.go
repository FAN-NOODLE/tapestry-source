package helper

import (
	"runtime"
	"time"
)

// FuncExecTime 获取函数执行占用的内存和cpu
func FuncExecTime(f func()) (float64, int, time.Duration) {
	// 记录程序开始时的内存使用情况
	var memStats runtime.MemStats
	runtime.ReadMemStats(&memStats)
	startMem := memStats.Alloc

	// 记录程序开始时的CPU数量
	startCPU := runtime.NumCPU()

	// 记录程序开始执行时间
	startTime := time.Now()

	// 执行函数
	f()

	execTime := time.Since(startTime)

	// 记录程序结束时的内存使用情况
	runtime.ReadMemStats(&memStats)
	endMem := memStats.Alloc

	// 记录程序结束时的CPU数量
	endCPU := runtime.NumCPU()

	useMemory := float64(endMem-startMem) / (1024 * 1024)
	useCpu := endCPU - startCPU

	// 输出结果
	// fmt.Printf("Memory usage: %.2f MB\n", float64(endMem-startMem)/(1024*1024))
	// fmt.Printf("CPU usage: %d cores\n", endCPU-startCPU)
	// fmt.Println("程序执行时长:", elapsed)

	return useMemory, useCpu, execTime
}
