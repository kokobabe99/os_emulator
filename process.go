package main

import "time"

type Process struct {
	ID               int
	Name             string
	TotalInstruction int // 添加这个字段来记录总指令数
	Instructions     int
	MemoryRequired   int
	Pages            int
	InMemory         bool
	Finished         bool
	DelayCount       int
	CreatedAt        time.Time // 添加创建时间字段
}

var processIDCounter = 1

func NewProcess(name string, cfg *Config) *Process {
	ins := randInt(cfg.MinIns, cfg.MaxIns)
	mem := randInt(cfg.MinMemPerProc, cfg.MaxMemPerProc)
	p := &Process{
		ID:               processIDCounter,
		Name:             name,
		Instructions:     ins,
		TotalInstruction: ins,
		MemoryRequired:   mem,
		InMemory:         false,
		Finished:         false,
		CreatedAt:        time.Now(), // 初始化创建时间
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
