package main

import (
	"fmt"
	"log"
	"net"
	"net/rpc"
	"sync"
)

//
// 公共RPC请求/回复定义
//

type PutArgs struct {
	Key   string // 键
	Value string // 值
}

type PutReply struct {
}

type GetArgs struct {
	Key string // 键
}

type GetReply struct {
	Value string // 值
}

// 客户端

func connect() *rpc.Client {
	// 建立与 RPC 服务器的连接
	client, err := rpc.Dial("tcp", ":1234")
	if err != nil {
		log.Fatal("连接错误:", err) // 如果连接失败，输出错误并终止程序
	}
	return client // 返回连接对象
}

func get(key string) string {
	client := connect()                         // 建立 RPC 连接
	args := GetArgs{"subject"}                  // 准备获取键为"subject"的值的请求参数
	reply := GetReply{}                         // 准备存放服务器返回的响应结果
	err := client.Call("KV.Get", &args, &reply) // 发送 KV.Get 请求到服务器，并将响应存放在 reply 中
	if err != nil {
		log.Fatal("错误:", err) // 如果调用过程中出错，输出错误并终止程序
	}
	client.Close()     // 关闭 RPC 连接
	return reply.Value // 返回服务器返回的键值对中对应键为"subject"的值
}

func put(key string, val string) {
	client := connect()
	args := PutArgs{"subject", "6.824"} // 准备设置键为"subject"的值为"6.824"
	reply := PutReply{}
	err := client.Call("KV.Put", &args, &reply)
	if err != nil {
		log.Fatal("错误:", err)
	}
	client.Close()
}

// 服务器

type KV struct {
	mu   sync.Mutex
	data map[string]string
}

func server() {
	kv := &KV{data: map[string]string{}}
	rpcs := rpc.NewServer()
	rpcs.Register(kv)                  // 注册KV服务到RPC服务器
	l, e := net.Listen("tcp", ":1234") // 在1234端口监听
	if e != nil {
		log.Fatal("监听错误:", e)
	}
	go func() {
		for {
			conn, err := l.Accept()
			if err == nil {
				go rpcs.ServeConn(conn) // 处理接受到的连接
			} else {
				break
			}
		}
		l.Close()
	}()
}

func (kv *KV) Get(args *GetArgs, reply *GetReply) error {
	kv.mu.Lock()
	defer kv.mu.Unlock()

	reply.Value = kv.data[args.Key] // 获取键args.Key对应的值

	return nil
}

func (kv *KV) Put(args *PutArgs, reply *PutReply) error {
	kv.mu.Lock()
	defer kv.mu.Unlock()

	kv.data[args.Key] = args.Value // 设置键args.Key的值为args.Value

	return nil
}

//
// 主函数
//

func main() {

	server()                 // 启动RPC服务器
	put("subject", "6.5840") // 设置键为"subject"的值为"6.5840"
	fmt.Printf("Put(subject, 6.5840) 完成\n")
	fmt.Printf("get(subject) -> %s\n", get("subject")) // 获取键为"subject"的值并打印

}
