package main

import (
	"fmt"
	"os"
)

func (s *Shell) ReportToFile() {
	f, err := os.Create("csopesy-log.txt")
	if err != nil {
		fmt.Println("Error writing report:", err)
		return
	}
	defer f.Close()

	fmt.Fprintf(f, "=== CSOPESY REPORT ===\n")
	fmt.Fprintf(f, "Total Memory: %d KB\n", s.MemMgr.TotalMemoryKB)
	fmt.Fprintf(f, "Used Memory : %d KB\n", s.MemMgr.UsedMemoryKB)
	fmt.Fprintf(f, "Free Memory : %d KB\n", s.MemMgr.TotalMemoryKB-s.MemMgr.UsedMemoryKB)
	fmt.Fprintf(f, "Paged In    : %d\n", s.MemMgr.PagedInCount)
	fmt.Fprintf(f, "Paged Out   : %d\n", s.MemMgr.PagedOutCount)
	fmt.Fprintf(f, "\n--- Process List ---\n")
	for name, p := range s.Processes {
		status := "RUNNING"
		if p.Finished {
			status = "FINISHED"
		}
		mem := "IN MEMORY"
		if !p.InMemory {
			mem = "PAGED OUT"
		}
		fmt.Fprintf(f, "- %s (ID: %d) | %d KB | %s | %s\n", name, p.ID, p.MemoryRequired, mem, status)
	}
	fmt.Println("Report written to csopesy-log.txt")
}
