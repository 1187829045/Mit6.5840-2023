package mr

// RPC 定义。
// 请记住将所有名称大写。

import "os"
import "strconv"

// 示例展示如何声明参数
// 并回复 RPC。

type ExampleArgs struct {
	X int
}

type ExampleReply struct {
	Y int
}

// Add your RPC definitions here.

type WorkerArgs struct {
	WorkerId    int
	WorkerState int // init\done\fail
	Task        *Task
}

type WorkerReply struct {
	WorkerId    int
	WorkerState int // init\done\fail
	Task        *Task
}

const (
	TaskTypeMap    = 0
	TaskTypeReduce = 1
)

const (
	WorkerStateInit = 0
	WorkerStateDone = 1
	WorkerStateFail = 2
)

// 在此处添加您的 RPC 定义。
// 在 /var/tmp 中为协调器编写一个独特的 UNIX 域套接字名称
// 无法使用当前目录，因为
// Athena AFS 不支持 UNIX 域套接字。

func coordinatorSock() string {
	s := "/var/tmp/5840-mr-"
	s += strconv.Itoa(os.Getuid())
	return s
}
