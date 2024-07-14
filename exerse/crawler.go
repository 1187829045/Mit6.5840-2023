package main

import (
	"fmt"
	"sync"
)

// Go教程中爬虫练习的几种解决方案
// https://tour.golang.org/concurrency/10

// 串行爬虫

func Serial(url string, fetcher Fetcher, fetched map[string]bool) {
	if fetched[url] {
		return // 如果已经抓取过该URL，则返回
	}
	fetched[url] = true             // 标记该URL已经被抓取
	urls, err := fetcher.Fetch(url) // 获取该URL页面中的所有链接
	if err != nil {
		return // 如果抓取出错，则返回
	}
	for _, u := range urls {
		Serial(u, fetcher, fetched) // 递归抓取页面中的每个链接
	}
	return
}

//
// 并发爬虫，使用共享状态和互斥锁
//

type fetchState struct {
	mu      sync.Mutex
	fetched map[string]bool
}

func (fs *fetchState) testAndSet(url string) bool {
	fs.mu.Lock()
	defer fs.mu.Unlock()
	r := fs.fetched[url]   // 检查该URL是否已经被抓取
	fs.fetched[url] = true // 标记该URL已经被抓取
	return r               // 返回之前的抓取状态
}

func ConcurrentMutex(url string, fetcher Fetcher, fs *fetchState) {
	if fs.testAndSet(url) { // 如果已经抓取过该URL，则返回
		return
	}
	urls, err := fetcher.Fetch(url) // 获取该URL页面中的所有链接
	if err != nil {
		return // 如果抓取出错，则返回
	}
	var done sync.WaitGroup
	for _, u := range urls {
		done.Add(1)
		go func(u string) {
			defer done.Done()
			ConcurrentMutex(u, fetcher, fs) // 并发抓取页面中的每个链接
		}(u)
	}
	done.Wait()
	return
}

func makeState() *fetchState {
	return &fetchState{fetched: make(map[string]bool)} // 创建一个新的抓取状态
}

//
// 并发爬虫，使用通道
//

func worker(url string, ch chan []string, fetcher Fetcher) {
	urls, err := fetcher.Fetch(url) // 获取该URL页面中的所有链接
	if err != nil {
		ch <- []string{} // 如果抓取出错，则发送空列表到通道
	} else {
		ch <- urls // 否则，发送抓取到的链接列表到通道
	}
}

func coordinator(ch chan []string, fetcher Fetcher) {
	n := 1
	fetched := make(map[string]bool)
	for urls := range ch {
		for _, u := range urls {
			if fetched[u] == false {
				fetched[u] = true
				n += 1
				go worker(u, ch, fetcher) // 启动工作协程来抓取新的链接
			}
		}
		n -= 1
		if n == 0 {
			break
		}
	}
}

func ConcurrentChannel(url string, fetcher Fetcher) {
	ch := make(chan []string)
	go func() {
		ch <- []string{url} // 将初始URL放入通道
	}()
	coordinator(ch, fetcher)
}

//
// 主函数
//

func main() {
	fmt.Printf("=== 串行爬虫 ===\n")
	Serial("http://golang.org/", fetcher, make(map[string]bool)) // 使用串行爬虫抓取初始URL

	fmt.Printf("=== 并发爬虫，使用互斥锁 ===\n")
	ConcurrentMutex("http://golang.org/", fetcher, makeState()) // 使用并发爬虫和互斥锁抓取初始URL

	fmt.Printf("=== 并发爬虫，使用通道 ===\n")
	ConcurrentChannel("http://golang.org/", fetcher) // 使用并发爬虫和通道抓取初始URL
}

//
// Fetcher 接口
//

type Fetcher interface {
	// Fetch 返回页面中找到的URL列表。
	Fetch(url string) (urls []string, err error)
}

// fakeFetcher 是返回固定结果的 Fetcher。
type fakeFetcher map[string]*fakeResult

type fakeResult struct {
	body string
	urls []string
}

func (f fakeFetcher) Fetch(url string) ([]string, error) {
	if res, ok := f[url]; ok {
		fmt.Printf("找到：   %s\n", url) // 打印找到的URL
		return res.urls, nil
	}
	fmt.Printf("未找到： %s\n", url) // 打印未找到的URL
	return nil, fmt.Errorf("未找到：%s", url)
}

// fetcher 是一个填充好的 fakeFetcher。
var fetcher = fakeFetcher{
	"http://golang.org/": &fakeResult{
		"The Go Programming Language",
		[]string{
			"http://golang.org/pkg/",
			"http://golang.org/cmd/",
		},
	},
	"http://golang.org/pkg/": &fakeResult{
		"Packages",
		[]string{
			"http://golang.org/",
			"http://golang.org/cmd/",
			"http://golang.org/pkg/fmt/",
			"http://golang.org/pkg/os/",
		},
	},
	"http://golang.org/pkg/fmt/": &fakeResult{
		"Package fmt",
		[]string{
			"http://golang.org/",
			"http://golang.org/pkg/",
		},
	},
	"http://golang.org/pkg/os/": &fakeResult{
		"Package os",
		[]string{
			"http://golang.org/",
			"http://golang.org/pkg/",
		},
	},
}
