# CSOPESY OS Emulator

A simple operating system emulator that simulates memory management and process scheduling.

## Overview

This emulator implements core OS concepts including:

- Memory management with both flat and paging allocation
- Process scheduling with FCFS and Round Robin algorithms
- Multi-CPU support
- Memory page replacement
- Process state management

## Features

### Memory Management

- Flat Memory Allocation
  - Continuous memory blocks
  - Automatic oldest process swapping
  - Memory fragmentation handling
- Paging Memory Management
  - Fixed-size frames
  - Page table implementation
  - Page-in/Page-out tracking
  - Backing store support

### Process Scheduling

- Multi-CPU Support (Default: 4 cores)
- FCFS (First-Come-First-Serve) Scheduling
  - Processes executed in arrival order
  - No preemption
- RR (Round-Robin) Scheduling
  - Time quantum based execution
  - Process preemption
  - Ready queue management

### Process Management

- Process States:
  - RUNNING: Currently executing
  - FINISHED: Execution completed
  - IN MEMORY: Loaded in main memory
  - PAGED OUT: Swapped to backing store
- Automatic batch process generation
- Process resource tracking

## Commands

### Basic Operations

```bash
initialize          # Initialize system
help               # Show command help
exit               # Exit emulator
```

### Process Management

```bash
screen -s <name>   # Create new process
screen -ls         # List all processes
screen -r <name>   # Enter process shell

```

### Process Management

```bash
vmstat             # View memory statistics
process-smi        # Show process status
report-util        # Generate system report
```

### Scheduler Control

```bash
scheduler-test     # Start scheduler test
scheduler-stop     # Stop scheduler test
```

### Configuration

- The system can be configured using the `config.txt` file. The following parameters can be set:

```plaintext
NUM_CPU 4                  # Number of CPUs
SCHEDULER_TYPE FCFS        # FCFS or RR
TIME_QUANTUM 3             # For RR scheduling
BATCH_FREQUENCY 5          # Auto process creation interval
MIN_INSTRUCTIONS 5         # Min instructions per process
MAX_INSTRUCTIONS 10        # Max instructions per process
DELAY_PER_EXEC 2          # Instruction execution time
TOTAL_MEMORY 1024         # Total memory in KB
FRAME_SIZE 1024           # Frame size in KB
MIN_MEMORY_PER_PROCESS 128    # Min process memory
MAX_MEMORY_PER_PROCESS 256    # Max process memory
```

### Requirements

    - Go 1.16 or higher
    - MacOS/Linux/Windows supported

### ## Building and Running

```bash
# Run directly
go run .

# Build binary
go build -o csopesy

# Run binary
./csopesy
```
