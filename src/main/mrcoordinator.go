package main

// 启动协调器进程，该进程在 ../mr/coordinator.go 中实现
// go run mrcoordinator.go pg*.txt
// 请不要更改此文件。

import "6.5840/mr"
import "time"
import "os"
import "fmt"

func main() {
	// 检查命令行参数数量是否足够
	if len(os.Args) < 2 {
		// 如果参数不足，输出使用说明到标准错误输出
		fmt.Fprintf(os.Stderr, "Usage: mrcoordinator inputfiles...\n")
		// 退出程序，状态码为1表示异常退出
		os.Exit(1)
	}

	// 创建一个Coordinator实例，传入输入文件列表和10个Reduce任务
	m := mr.MakeCoordinator(os.Args[1:], 10)
	// 循环等待直到Coordinator完成所有任务
	for m.Done() == false {
		// 每秒检查一次任务完成状态
		time.Sleep(time.Second)
	}

	// 等待一秒后程序结束，确保所有输出完成
	time.Sleep(time.Second)
}
