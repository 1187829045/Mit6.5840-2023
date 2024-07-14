package main

import "time"
import "math/rand"

func main() {
	rand.Seed(time.Now().UnixNano()) // 设置随机数种子为当前时间的纳秒级Unix时间戳

	count := 0    // 计数已接收到的投票数量
	finished := 0 // 计数已完成的goroutine数量

	for i := 0; i < 100; i++ {
		go func() {
			vote := requestVote1() // 发起投票请求
			if vote {
				count++ //如果投票成功，则增加计数
			}
			finished++ //标记goroutine已完成
		}()
	}

	for count < 5 && finished != 10 {
		// 等待，直到收到至少5个投票或所有goroutine完成
	}

	if count >= 5 {
		println("received 5+ votes!") // 如果收到至少5个投票，则打印信息
	} else {
		println("lost") // 如果没有收到足够的投票，则打印信息
	}
}

func requestVote1() bool {
	// Go 语言中用于表示时间间隔的常量，表示一毫秒的时间单位
	time.Sleep(time.Duration(rand.Intn(100)) * time.Millisecond) // 模拟随机等待时间
	return rand.Int()%2 == 0                                     // 模拟投票结果，随机返回true或false
}
