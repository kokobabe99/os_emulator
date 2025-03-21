package main

type CPU struct {
	ID          int
	Running     *Process
	QuantumLeft int // for RR
}

type Scheduler struct {
	CPUs       []*CPU
	Algo       string // "fcfs" or "rr"
	Quantum    int
	ReadyQueue []*Process

	TotalTicks  int
	ActiveTicks int
	IdleTicks   int
}

func NewScheduler(cfg *Config) *Scheduler {
	cpus := make([]*CPU, cfg.NumCPU)
	for i := 0; i < cfg.NumCPU; i++ {
		cpus[i] = &CPU{ID: i}
	}
	return &Scheduler{
		CPUs:       cpus,
		Algo:       cfg.Scheduler,
		Quantum:    cfg.Quantum,
		ReadyQueue: []*Process{},
	}
}

func (s *Scheduler) Tick() {
	s.TotalTicks++

	for _, cpu := range s.CPUs {
		if cpu.Running != nil {
			cpu.Running.Instructions--
			s.ActiveTicks++
			if cpu.Running.Instructions <= 0 {
				cpu.Running.Finished = true
				cpu.Running = nil
			} else if s.Algo == "rr" {
				cpu.QuantumLeft--
				if cpu.QuantumLeft <= 0 {
					// requeue
					s.ReadyQueue = append(s.ReadyQueue, cpu.Running)
					cpu.Running = nil
				}
			}
		} else {
			if len(s.ReadyQueue) > 0 {
				proc := s.ReadyQueue[0]
				s.ReadyQueue = s.ReadyQueue[1:]
				cpu.Running = proc
				if s.Algo == "rr" {
					cpu.QuantumLeft = s.Quantum
				}
			} else {
				s.IdleTicks++
			}
		}
	}
}

func (s *Scheduler) AddProcess(p *Process) {
	s.ReadyQueue = append(s.ReadyQueue, p)
}
