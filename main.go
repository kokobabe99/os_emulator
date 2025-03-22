package main

import (
	"math/rand"
	"time"
)

// === (existing code preserved, omitted here for brevity) ===

/*
*
4 核 CPU
使用 FCFS 调度器
RR 时间片为 3（FCFS 下忽略）
每 5 秒生成一个批处理进程
每个进程 5~10 条指令
每条指令执行耗时 2 tick
总内存为 1024KB，每帧大小为 1024KB（flat memory）
每进程占用 128~256KB
*/
func main() {
	rand.Seed(time.Now().UnixNano())

	shell := NewShell()
	shell.PrintHeader()
	shell.Run()
}
