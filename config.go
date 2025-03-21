package main

import (
	"fmt"
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
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	tokens := strings.Fields(string(data))
	if len(tokens) < 11 {
		return nil, fmt.Errorf("invalid config.txt format: expected 11 space-separated values")
	}

	toInt := func(s string) int {
		n, _ := strconv.Atoi(s)
		return n
	}

	return &Config{
		NumCPU:        toInt(tokens[0]),
		Scheduler:     tokens[1],
		Quantum:       toInt(tokens[2]),
		BatchFreq:     toInt(tokens[3]),
		MinIns:        toInt(tokens[4]),
		MaxIns:        toInt(tokens[5]),
		DelayPerExec:  toInt(tokens[6]),
		TotalMemoryKB: toInt(tokens[7]),
		FrameSizeKB:   toInt(tokens[8]),
		MinMemPerProc: toInt(tokens[9]),
		MaxMemPerProc: toInt(tokens[10]),
	}, nil
}
