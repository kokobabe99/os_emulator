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
		cfg, err := LoadConfig("config.txt")
		if err != nil {
			fmt.Println("Config error:", err)
			return
		}
		s.cfg = cfg
		s.MemMgr = NewMemoryManager(cfg.TotalMemoryKB, cfg.FrameSizeKB)
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
		fmt.Println("Running processes:")
		for name, p := range s.Processes {
			fmt.Printf("- %s (ID: %d) | Memory: %d KB | InMemory: %v\n", name, p.ID, p.MemoryRequired, p.InMemory)
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
			fmt.Printf("Process %s | ID: %d | Remaining Instructions: %d\n", p.Name, p.ID, p.Instructions)
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
			fmt.Printf("[Auto] Created process %s\n", name)
			processIDCounter++
			if s.schedulerStopped {
				ticker.Stop()
				return
			}
		}
	}()
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
	fmt.Printf("Used Memory: %d KB / %d KB\n", s.MemMgr.UsedMemoryKB, s.MemMgr.TotalMemoryKB)
	fmt.Println("--- Process List ---")
	for _, p := range s.Processes {
		status := "RUNNING"
		if p.Finished {
			status = "FINISHED"
		}
		mem := "IN MEMORY"
		if !p.InMemory {
			mem = "PAGED OUT"
		}
		fmt.Printf("- %s (ID: %d) | %s | %s | %d KB\n", p.Name, p.ID, status, mem, p.MemoryRequired)
	}
}
