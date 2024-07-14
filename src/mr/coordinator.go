package mr

import (
	"log"
	"net"
	"net/http"
	"net/rpc"
	"os"
	"sync"
	"time"
)

type Task struct {
	TaskId    int       // 任务的唯一标识符
	TaskType  int       // 任务类型：0 表示映射任务，1 表示归并任务
	TaskState int       // 任务状态：0 表示初始化，1 表示任务运行中，2 表示任务已完成
	NReduce   int       // 归并任务的数量
	StartTime time.Time // 任务开始时间

	Input []string // 任务的输入数据，通常是文件路径列表等
}

const (
	StateMap    = 0 // 定义状态常量：映射阶段
	StateReduce = 1 // 定义状态常量：归并阶段
	StateDone   = 2 // 定义状态常量：已完成
)

const (
	TaskStateInit = 0 // 定义任务状态常量：初始化
	TaskStateRun  = 1 // 定义任务状态常量：运行中
	TaskStateDone = 2 // 定义任务状态常量：已完成
)

type Coordinator struct {
	MapTaskChan    chan *Task    // 用于接收 Map 任务的通道
	ReduceTaskChan chan *Task    // 用于接收 Reduce 任务的通道
	NumReduce      int           // Reduce 任务的数量
	NumMap         int           // Map 任务的数量
	NumDoneReduce  int           // 已完成的 Reduce 任务数量
	NumDoneMap     int           // 已完成的 Map 任务数量
	State          int           // 当前 Coordinator 的状态：0 表示映射阶段，1 表示归并阶段，2 表示全部完成
	mu             sync.Mutex    // 用于保护并发访问的互斥锁
	Timeout        time.Duration // 任务超时时间
	MapTasks       map[int]*Task // 存储所有 Map 任务的映射
	ReduceTasks    map[int]*Task // 存储所有 Reduce 任务的映射
}

func (c *Coordinator) Example(args *ExampleArgs, reply *ExampleReply) error {
	reply.Y = args.X + 1
	return nil
}

func (c *Coordinator) server() {
	rpc.Register(c)
	rpc.HandleHTTP()
	//l, e := net.Listen("tcp", ":1234")
	sockname := coordinatorSock()
	os.Remove(sockname)
	l, e := net.Listen("unix", sockname)
	if e != nil {
		log.Fatal("listen error:", e)
	}
	go http.Serve(l, nil)
}

func (c *Coordinator) Done() bool {
	ret := false

	// Your code here.

	return ret
}

//map 任务初始化；
//rpc handler：回应worker分配任务请求、回应worker任务完成通知；
//自身状态控制，处理 map/reduce 阶段，还是已经全部完成；
//任务超时重新分配；
// 创建一个 Coordinator。
// main/mrcoordinator.go 调用此函数。
// nReduce 是要使用的 Reduce 任务的数量。
//每个Map任务应创建nReduce个中间文件供Reduce任务使用。

func MakeCoordinator(files []string, nReduce int) *Coordinator {
	c := Coordinator{}
	// Your code here.
	c.NumMap = len(files)
	c.NumReduce = nReduce
	c.NumDoneMap = 0
	c.NumDoneReduce = 0
	c.MapTaskChan = make(chan *Task, len(files))
	c.ReduceTaskChan = make(chan *Task, nReduce)
	c.MapTasks = make(map[int]*Task)
	c.ReduceTasks = make(map[int]*Task)
	c.State = StateMap
	c.Timeout = time.Duration(time.Second * 10)
	for i, file := range files {
		input := []string{file}
		task := Task{
			TaskId:    i,
			TaskType:  TaskTypeMap,
			TaskState: TaskStateInit,
			Input:     input,
			NReduce:   nReduce,
			StartTime: time.Now(),
		}
		c.MapTaskChan <- &task
		c.MapTasks[i] = &task
	}

	c.server()
	return &c
}
