package main

import (
	"bufio"
	"os"
	"strconv"
	"strings"
)

type Config struct {
	NumCPU        int
	Scheduler     string
	Quantum       int
	BatchFreq     int
	MinIns        int
	MaxIns        int
	DelayPerExec  int
	TotalMemoryKB int
	FrameSizeKB   int
	MinMemPerProc int
	MaxMemPerProc int
}

func LoadConfig(path string) (*Config, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	cfg := &Config{}
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}

		parts := strings.Fields(line)
		if len(parts) != 2 {
			continue
		}

		key := parts[0]
		value := parts[1]

		switch key {
		case "NUM_CPU":
			cfg.NumCPU, _ = strconv.Atoi(value)
		case "SCHEDULER_TYPE":
			cfg.Scheduler = value
		case "TIME_QUANTUM":
			cfg.Quantum, _ = strconv.Atoi(value)
		case "BATCH_FREQUENCY":
			cfg.BatchFreq, _ = strconv.Atoi(value)
		case "MIN_INSTRUCTIONS":
			cfg.MinIns, _ = strconv.Atoi(value)
		case "MAX_INSTRUCTIONS":
			cfg.MaxIns, _ = strconv.Atoi(value)
		case "DELAY_PER_EXEC":
			cfg.DelayPerExec, _ = strconv.Atoi(value)
		case "TOTAL_MEMORY":
			cfg.TotalMemoryKB, _ = strconv.Atoi(value)
		case "FRAME_SIZE":
			cfg.FrameSizeKB, _ = strconv.Atoi(value)
		case "MIN_MEMORY_PER_PROCESS":
			cfg.MinMemPerProc, _ = strconv.Atoi(value)
		case "MAX_MEMORY_PER_PROCESS":
			cfg.MaxMemPerProc, _ = strconv.Atoi(value)
		}
	}

	return cfg, nil
}
