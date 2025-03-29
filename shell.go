package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"time"
)

// === Shell Structure ===
type Shell struct {
	cfg              *Config
	MemMgr           *MemoryManager
	Processes        map[string]*Process
	Initialized      bool
	schedulerStopped bool
	scheduler        *Scheduler // 添加这行
}

func NewShell() *Shell {
	return &Shell{
		Processes: make(map[string]*Process),
	}
}

// 添加 printHeader 函数
func (s *Shell) PrintHeader() {
	fmt.Print(`
   ___________ ____  ____  _____________  __
  / ____/ ___// __ \/ __ \/ ____/ ___/\ \/ /
 / /    \__ \/ / / / /_/ / __/  \__ \  \  /
/ /___ ___/ / /_/ / ____/ /___ ___/ /  / /
\____//____/\____/_/   /_____//____/  /_/
`)
	fmt.Println("\nWelcome!type help")
}

// === Main Command Loop ===
func (s *Shell) Run() {
	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Print("csopesy> ")
		line, _ := reader.ReadString('\n')
		input := strings.TrimSpace(line)
		if input == "exit" {
			fmt.Println("Exiting emulator...")
			return
		}
		s.handleCommand(input)
	}
}

// === Command Handler ===
func (s *Shell) handleCommand(input string) {

	if !s.Initialized && input != "initialize" && input != "exit" && input != "help" {
		fmt.Println("Please initialize the system first.")
		return
	}

	args := strings.Fields(input)
	if len(args) == 0 {
		return
	}
	switch args[0] {
	case "help":
		fmt.Println("Available Commands:")
		fmt.Println("  initialize          - Initialize the system")
		fmt.Println("  screen -s <name>    - Create a new process")
		fmt.Println("  screen -ls          - List all processes")
		fmt.Println("  screen -r <name>    - Enter process shell")
		fmt.Println("  scheduler-test      - Start scheduler test")
		fmt.Println("  scheduler-stop      - Stop scheduler test")
		fmt.Println("  report-util         - Generate system report")
		fmt.Println("  process-smi         - Show process status")
		fmt.Println("  vmstat              - Display memory statistics")
		fmt.Println("  exit                - Exit emulator")
	case "initialize":

		if s.Initialized {
			fmt.Println("System already initialized.")
			return
		}
		cfg, err := LoadConfig("config.txt")
		if err != nil {
			fmt.Println("Config error:", err)
			return
		}
		s.cfg = cfg
		s.MemMgr = NewMemoryManager(cfg.TotalMemoryKB, cfg.FrameSizeKB)
		s.scheduler = NewScheduler(cfg, s.MemMgr) // 添加这行，确保调度器在初始化时就创建
		s.Initialized = true
		s.schedulerStopped = false
		fmt.Println("System initialized.")
	case "screen":
		if len(args) < 2 {
			fmt.Println("Invalid screen command")
			return
		}
		s.handleScreen(args[1:])
	case "scheduler-test":
		s.StartSchedulerTest()
	case "scheduler-stop":
		s.StopSchedulerTest()
	case "report-util":
		s.ReportToFile()
	case "process-smi":
		s.PrintProcessSMI()
	case "vmstat":
		s.PrintVMStat()
	default:
		fmt.Println("Unknown command:", args[0])
	}
}

// === screen command ===
func (s *Shell) handleScreen(args []string) {
	switch args[0] {
	case "-s":
		if len(args) < 2 {
			fmt.Println("screen -s <process name>")
			return
		}
		name := args[1]
		p := NewProcess(name, s.cfg)
		s.MemMgr.Allocate(p)
		s.Processes[name] = p
		fmt.Printf("Process %s created with %d instructions and %d KB memory\n", name, p.Instructions, p.MemoryRequired)
	case "-ls":
		fmt.Println("CPU Information:")
		fmt.Printf("Total CPUs: %d\n", len(s.scheduler.CPUs))

		// 检查 CPU 和内存状态
		cpuUtil := 0.0
		activeCPUs := 0

		// 只有当有进程在内存中时才计算 CPU 利用率
		runningProcesses := 0
		// 计算内存使用率

		memoryUtil := float64(s.MemMgr.UsedMemoryKB) / float64(s.MemMgr.TotalMemoryKB) * 100

		// 根据内存使用率计算 CPU 利用率
		if s.MemMgr.UsedMemoryKB > 0 {
			cpuUtil = memoryUtil
			if cpuUtil > 100 {
				cpuUtil = 100
			}
			activeCPUs = int(float64(len(s.scheduler.CPUs)) * memoryUtil / 100)
		}

		// 如果有进程在内存中，则 CPU 处于活跃状态
		if runningProcesses > 0 && s.cfg.NumCPU == 1 {
			activeCPUs = 1 // 单 CPU 系统
			cpuUtil = 1 * 100
		}

		fmt.Printf("Active CPUs: %d\n", activeCPUs)
		fmt.Printf("CPU Utilization: %.2f%%\n\n", cpuUtil)

		fmt.Println("=== PROCESSING ===")
		for name, p := range s.Processes {
			if p.Instructions > 0 {
				completedInst := p.TotalInstruction - p.Instructions
				fmt.Printf("- %s (%s) | ID: %d | Memory: %d KB | InMemory: %v  Instructions: %d/%d\n",
					name,
					p.CreatedAt.Format("01/02/2024 03:04:05PM"),
					p.ID,
					p.MemoryRequired,
					p.InMemory,
					completedInst,
					p.TotalInstruction)
			}
		}

		fmt.Println("\n=== PROCESS FINISHED ===")
		for name, p := range s.Processes {
			if p.Instructions <= 0 || p.Finished {
				fmt.Printf("- %s (%s) | ID: %d | Memory: %d KB | InMemory: %v  Instructions: %d/%d\n",
					name,
					p.CreatedAt.Format("01/02/2024 03:04:05PM"),
					p.ID,
					p.MemoryRequired,
					p.InMemory,
					p.TotalInstruction,
					p.TotalInstruction)
			}
		}
	case "-r":
		if len(args) < 2 {
			fmt.Println("screen -r <process name>")
			return
		}
		s.ResumeScreen(args[1])
	default:
		fmt.Println("Unknown screen command")
	}
}

// === screen -r logic ===
func (s *Shell) ResumeScreen(name string) {
	p, ok := s.Processes[name]
	if !ok || p.Finished {
		fmt.Printf("Process %s not found.\n", name)
		return
	}
	reader := bufio.NewReader(os.Stdin)
	fmt.Printf("[Entering %s Shell] > type 'process-smi' or 'exit'\n", p.Name)
	for {
		fmt.Printf("[%s]> ", p.Name)
		line, _ := reader.ReadString('\n')
		cmd := strings.TrimSpace(line)
		switch cmd {
		case "process-smi":

			fmt.Printf("Process %s | ID: %d | Remaining Instructions: %d", p.Name, p.ID, p.TotalInstruction-p.Instructions)

			if p.Instructions <= 0 {
				fmt.Println("Finished!")
				p.Finished = true
			}
		case "exit":
			return
		default:
			fmt.Println("Unknown command")
		}
	}
}

// === scheduler-test ===
func (s *Shell) StartSchedulerTest() {
	fmt.Println("Starting scheduler test...")

	go func() {
		ticker := time.NewTicker(time.Duration(s.cfg.BatchFreq) * time.Second)
		for range ticker.C {
			name := fmt.Sprintf("p%03d", processIDCounter)
			p := NewProcess(name, s.cfg)

			if !s.MemMgr.Allocate(p) {
				fmt.Println("Memory full. Skipping:", name)
				continue
			}

			s.Processes[name] = p
			s.scheduler.AddProcess(p) // 添加到调度器
			//fmt.Printf("[Auto] Created process %s\n", name)
			processIDCounter++

			if s.schedulerStopped {
				ticker.Stop()
				return
			}
		}
	}()

	// 启动调度器
	go s.scheduler.Start()
}

func (s *Shell) StopSchedulerTest() {
	s.schedulerStopped = true
	fmt.Println("Stopped scheduler test.")
}

func (s *Shell) PrintVMStat() {
	fmt.Println("=== vmstat ===")
	fmt.Printf("Total Memory: %d KB\n", s.MemMgr.TotalMemoryKB)
	fmt.Printf("Used Memory : %d KB\n", s.MemMgr.UsedMemoryKB)
	fmt.Printf("Free Memory : %d KB\n", s.MemMgr.TotalMemoryKB-s.MemMgr.UsedMemoryKB)
	fmt.Printf("Paged In    : %d\n", s.MemMgr.PagedInCount)
	fmt.Printf("Paged Out   : %d\n", s.MemMgr.PagedOutCount)
	fmt.Println("--- Process States ---")
	for _, p := range s.Processes {
		status := "RUNNING"
		if p.Finished {
			status = "FINISHED"
		}
		memStatus := "IN MEMORY"
		if !p.InMemory {
			memStatus = "PAGED OUT"
		}
		fmt.Printf("- %s (ID: %d) | %s | %s | %d KB\n", p.Name, p.ID, status, memStatus, p.MemoryRequired)
	}
}

func (s *Shell) PrintProcessSMI() {

	fmt.Println("=== process-smi ===")

	// CPU 和内存信息
	fmt.Printf("CPU Information:\n")
	fmt.Printf("Total CPUs: %d\n", len(s.scheduler.CPUs))

	cpuUtil := 0.0
	activeCPUs := 0
	memoryUtil := float64(s.MemMgr.UsedMemoryKB) / float64(s.MemMgr.TotalMemoryKB) * 100
	if s.MemMgr.UsedMemoryKB > 0 {
		cpuUtil = memoryUtil
		if cpuUtil > 100 {
			cpuUtil = 100
		}
		activeCPUs = int(float64(len(s.scheduler.CPUs)) * memoryUtil / 100)
	}

	fmt.Printf("Active CPUs: %d\n", activeCPUs)
	fmt.Printf("CPU Utilization: %.2f%%\n\n", cpuUtil)

	fmt.Printf("Memory Usage: %d KB / %d KB (%.2f%%)\n\n",
		s.MemMgr.UsedMemoryKB,
		s.MemMgr.TotalMemoryKB,
		memoryUtil)

	// 运行中的进程
	// 运行中的进程
	fmt.Println("=== PROCESSING ===")
	for _, p := range s.Processes {
		if p.Instructions > 0 {
			completedInst := p.TotalInstruction - p.Instructions
			memStatus := "IN MEMORY"
			if !p.InMemory {
				memStatus = "PAGED OUT"
			}
			fmt.Printf("- %s (%s) | ID: %d | %s | %s | Memory: %d KB | Progress: %d/%d\n",
				p.Name,
				p.CreatedAt.Format("01/02/2024 03:04:05PM"),
				p.ID,
				"RUNNING",
				memStatus,
				p.MemoryRequired,
				completedInst,
				p.TotalInstruction)
		}
	}

	// 已完成的进程
	fmt.Println("\n=== PROCESS FINISHED ===")
	for _, p := range s.Processes {
		if p.Instructions <= 0 || p.Finished {
			memStatus := "IN MEMORY"
			if !p.InMemory {
				memStatus = "PAGED OUT"
			}
			fmt.Printf("- %s (%s) | ID: %d | %s | %s | Memory: %d KB | Progress: %d/%d\n",
				p.Name,
				p.CreatedAt.Format("01/02/2024 03:04:05PM"),
				p.ID,
				"FINISHED",
				memStatus,
				p.MemoryRequired,
				p.TotalInstruction,
				p.TotalInstruction)
		}
	}
}
