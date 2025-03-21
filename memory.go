package main

type AllocationMode int

const (
	Flat AllocationMode = iota
	Paging
)

type MemoryManager struct {
	TotalMemoryKB int
	FrameSizeKB   int
	Allocation    AllocationMode
	UsedMemoryKB  int
	Frames        []bool
	PagedInCount  int
	PagedOutCount int
	BackingStore  map[int]*Process
}

func NewMemoryManager(total, frame int) *MemoryManager {
	mode := Flat
	if total > frame {
		mode = Paging
	}
	return &MemoryManager{
		TotalMemoryKB: total,
		FrameSizeKB:   frame,
		Allocation:    mode,
		Frames:        make([]bool, total/frame),
		BackingStore:  make(map[int]*Process),
	}
}

func (m *MemoryManager) Allocate(p *Process) bool {
	if m.Allocation == Flat {
		if m.UsedMemoryKB+p.MemoryRequired > m.TotalMemoryKB {
			m.SwapOutOldest()
		}
		m.UsedMemoryKB += p.MemoryRequired
		p.InMemory = true
		return true
	} else {
		requiredFrames := p.Pages
		free := 0
		for _, f := range m.Frames {
			if !f {
				free++
			}
		}
		if free < requiredFrames {
			m.SwapOutOldest()
		}
		count := 0
		for i := range m.Frames {
			if !m.Frames[i] {
				m.Frames[i] = true
				count++
			}
			if count == requiredFrames {
				break
			}
		}
		m.UsedMemoryKB += p.MemoryRequired
		m.PagedInCount += requiredFrames
		p.InMemory = true
		return true
	}
}

func (m *MemoryManager) SwapOutOldest() {
	for pid, p := range m.BackingStore {
		if p.InMemory {
			m.UsedMemoryKB -= p.MemoryRequired
			p.InMemory = false
			m.PagedOutCount += p.Pages
			delete(m.BackingStore, pid)
			break
		}
	}
}
