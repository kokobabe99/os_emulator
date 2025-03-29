package main

import (
	"fmt"
	"os"
)

func (s *Shell) ReportToFile() {
	file, err := os.Create("report.txt")
	if err != nil {
		fmt.Println("Error creating report:", err)
		return
	}
	defer file.Close()

	// CPU 信息部分
	fmt.Fprintf(file, "CPU Information:\n")
	fmt.Fprintf(file, "Total CPUs: %d\n", len(s.scheduler.CPUs))

	// 计算 CPU 利用率
	cpuUtil := 0.0
	activeCPUs := 0
	memoryUtil := float64(s.MemMgr.UsedMemoryKB) / float64(s.MemMgr.TotalMemoryKB) * 100

	if s.MemMgr.UsedMemoryKB > 0 {
		cpuUtil = memoryUtil
		if cpuUtil > 100 {
			cpuUtil = 100
		}
		// 确保活跃 CPU 数量不超过总数
		activeCPUs = int(float64(len(s.scheduler.CPUs)) * memoryUtil / 100)
		if activeCPUs > len(s.scheduler.CPUs) {
			activeCPUs = len(s.scheduler.CPUs)
		}
	}

	fmt.Fprintf(file, "Active CPUs: %d\n", activeCPUs)
	fmt.Fprintf(file, "CPU Utilization: %.2f%%\n\n", cpuUtil)

	// 内存信息部分
	fmt.Fprintf(file, "Memory Information:\n")
	fmt.Fprintf(file, "Total Memory: %d KB\n", s.MemMgr.TotalMemoryKB)
	fmt.Fprintf(file, "Used Memory: %d KB\n", s.MemMgr.UsedMemoryKB)
	fmt.Fprintf(file, "Free Memory: %d KB\n\n", s.MemMgr.TotalMemoryKB-s.MemMgr.UsedMemoryKB)

	// 运行中的进程
	fmt.Fprintf(file, "=== PROCESSING ===\n")
	for name, p := range s.Processes {
		if p.Instructions > 0 {
			completedInst := p.TotalInstruction - p.Instructions
			fmt.Fprintf(file, "- %s (%s) | ID: %d | Memory: %d KB | InMemory: %v  Instructions: %d/%d\n",
				name,
				p.CreatedAt.Format("01/02/2024 03:04:05PM"),
				p.ID,
				p.MemoryRequired,
				p.InMemory,
				completedInst,
				p.TotalInstruction)
		}
	}

	// 已完成的进程
	fmt.Fprintf(file, "\n=== PROCESS FINISHED ===\n")
	for name, p := range s.Processes {
		if p.Instructions <= 0 || p.Finished {
			fmt.Fprintf(file, "- %s (%s) | ID: %d | Memory: %d KB | InMemory: %v  Instructions: %d/%d\n",
				name,
				p.CreatedAt.Format("01/02/2024 03:04:05PM"),
				p.ID,
				p.MemoryRequired,
				p.InMemory,
				p.TotalInstruction,
				p.TotalInstruction)
		}
	}

	fmt.Println("Report generated: report.txt")
}
