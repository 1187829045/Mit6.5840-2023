package mr

import (
	"fmt"
	"log"
	"net/rpc"
)
import "hash/fnv"

//给 coordinator 发送 rpc 请求分配任务；
//给 coordinator 发送 rpc 通知任务完成；
//自身状态控制，准确来说是 coordinator 不需要 worker 工作时，通知 worker 结束运行；

// Map 函数返回 KeyValue 的一个切片。

type KeyValue struct {
	Key   string
	Value string
}

// 使用 ihash(key) % NReduce 来选择 Reduce
// Map 发出的每个 KeyValue 的任务编号。

func ihash(key string) int {
	h := fnv.New32a()
	h.Write([]byte(key))
	return int(h.Sum32() & 0x7fffffff)
}

// main/mrworker.go 调用此函数。

func Worker(mapf func(string, string) []KeyValue,
	reducef func(string, []string) string) {

	// Your worker implementation here.

	// uncomment to send the Example RPC to the coordinator.
	// CallExample()
}

// 示例函数展示如何向协调器发出 RPC 调用。
// RPC 参数和回复类型在 rpc.go 中定义

func CallExample() {

	// 声明一个参数结构体。
	args := ExampleArgs{}

	// 填充参数。
	args.X = 99

	// 声明一个回复结构体。
	reply := ExampleReply{}

	// 发送RPC请求，等待回复。
	// "Coordinator.Example" 告诉接收服务器我们要调用 Coordinator 结构体的 Example() 方法。
	ok := call("Coordinator.Example", &args, &reply)
	if ok {
		// reply.Y 应该是 100。
		fmt.Printf("reply.Y %v\n", reply.Y)
	} else {
		fmt.Printf("调用失败！\n")
	}
}

// 向协调器发送一个RPC请求，并等待响应。
// 通常返回 true。
// 如果出现问题，则返回 false。

func call(rpcname string, args interface{}, reply interface{}) bool {
	// c, err := rpc.DialHTTP("tcp", "127.0.0.1"+":1234")
	sockname := coordinatorSock()
	c, err := rpc.DialHTTP("unix", sockname)
	if err != nil {
		log.Fatal("dialing:", err)
	}
	defer c.Close()

	err = c.Call(rpcname, args, reply)
	if err == nil {
		return true
	}

	fmt.Println(err)
	return false
}
