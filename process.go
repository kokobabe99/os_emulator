package main

type Process struct {
	ID             int
	Name           string
	Instructions   int
	MemoryRequired int
	Pages          int
	InMemory       bool
	Finished       bool
}

var processIDCounter = 1

func NewProcess(name string, cfg *Config) *Process {
	ins := randInt(cfg.MinIns, cfg.MaxIns)
	mem := randInt(cfg.MinMemPerProc, cfg.MaxMemPerProc)
	p := &Process{
		ID:             processIDCounter,
		Name:           name,
		Instructions:   ins,
		MemoryRequired: mem,
		InMemory:       false,
		Finished:       false,
	}
	processIDCounter++

	if cfg.FrameSizeKB > 0 {
		p.Pages = mem / cfg.FrameSizeKB
		if mem%cfg.FrameSizeKB != 0 {
			p.Pages += 1
		}
	}
	return p
}
