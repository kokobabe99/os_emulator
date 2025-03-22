package main

import (
	"math/rand"
	"time"
)

// === Main Function ===
func main() {
	rand.Seed(time.Now().UnixNano())
	shell := NewShell()
	shell.PrintHeader()
	shell.Run()
}
