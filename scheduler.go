package main

import "time"

type CPU struct {
	ID          int
	Running     *Process
	QuantumLeft int // for RR
	IsActive    bool
	ActiveTime  int64
}

type Scheduler struct {
	CPUs       []*CPU
	Algo       string // "fcfs" or "rr"
	Quantum    int
	ReadyQueue []*Process

	TotalTicks   int
	ActiveTicks  int
	IdleTicks    int
	DelayPerExec int            // 添加执行延迟参数
	MemMgr       *MemoryManager // 添加内存管理器
}

func NewScheduler(cfg *Config, memMgr *MemoryManager) *Scheduler { // 修改构造函数
	cpus := make([]*CPU, cfg.NumCPU)
	for i := 0; i < cfg.NumCPU; i++ {
		cpus[i] = &CPU{ID: i}
	}
	return &Scheduler{
		CPUs:         cpus,
		Algo:         cfg.Scheduler,
		Quantum:      cfg.Quantum,
		ReadyQueue:   []*Process{},
		DelayPerExec: cfg.DelayPerExec,
		MemMgr:       memMgr, // 初始化内存管理器
	}
}

func (s *Scheduler) Tick() {
	s.TotalTicks++

	activeThisTick := 0

	for _, cpu := range s.CPUs {
		if cpu.Running != nil {
			// 检查是否需要等待延迟
			if cpu.Running.DelayCount >= s.DelayPerExec {
				cpu.Running.Instructions--
				cpu.Running.DelayCount = 0 // 重置延迟计数
			} else {
				cpu.Running.DelayCount++ // 增加延迟计数
			}
			activeThisTick++
			cpu.IsActive = true
			cpu.ActiveTime++
			s.ActiveTicks++

			// 检查进程是否完成
			if cpu.Running.Instructions <= 0 {
				cpu.Running.Instructions = 0
				cpu.Running.Finished = true
				cpu.Running.InMemory = false     // 进程完成时设置为不在内存中
				s.MemMgr.Deallocate(cpu.Running) // 从内存中移除
				cpu.Running = nil
				cpu.IsActive = false
			} else if s.Algo == "rr" {
				cpu.QuantumLeft--
				if cpu.QuantumLeft <= 0 {
					// 时间片用完，重新入队
					s.ReadyQueue = append(s.ReadyQueue, cpu.Running)
					cpu.Running = nil
					cpu.IsActive = false
				}
			}
		} else {
			if len(s.ReadyQueue) > 0 {
				proc := s.ReadyQueue[0]
				s.ReadyQueue = s.ReadyQueue[1:]
				cpu.Running = proc
				cpu.IsActive = true
				if s.Algo == "rr" {
					cpu.QuantumLeft = s.Quantum
				}
			} else {
				s.IdleTicks++
			}
		}
	}
	s.ActiveTicks += activeThisTick // 只增加当前时钟周期实际活跃的 CPU 数
}
func (s *Scheduler) AddProcess(p *Process) {
	s.ReadyQueue = append(s.ReadyQueue, p)
}

// 添加以下方法
func (s *Scheduler) Start() {
	ticker := time.NewTicker(time.Millisecond * 10)
	for range ticker.C {
		s.Tick()
	}
}

func (s *Scheduler) Stop() {
	// 清理所有CPU上运行的进程
	for _, cpu := range s.CPUs {
		if cpu.Running != nil {
			cpu.Running = nil
		}
	}
	// 清空就绪队列
	s.ReadyQueue = nil
}

func (s *Scheduler) GetActiveCPUs() int {
	active := 0
	for _, cpu := range s.CPUs {
		if cpu.IsActive {
			active++
		}
	}
	return active
}
