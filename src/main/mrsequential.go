package main

// 简单顺序 MapReduce。
// go run mrsequential.go wc.so pg*.txt

import "fmt"
import "6.5840/mr"
import "plugin"
import "os"
import "log"
import "io/ioutil"
import "sort"

// 按键排序。

type ByKey []mr.KeyValue

// 按键排序。

func (a ByKey) Len() int           { return len(a) }
func (a ByKey) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a ByKey) Less(i, j int) bool { return a[i].Key < a[j].Key }

func main() {
	if len(os.Args) < 3 {
		fmt.Fprintf(os.Stderr, "Usage: mrsequential xxx.so inputfiles...\n")
		os.Exit(1)
	}

	mapf, reducef := loadPlugin(os.Args[1])

	// 读取每个输入文件，
	// 将其传递给 Map，
	// 累积中间 Map 输出。
	intermediate := []mr.KeyValue{}
	for _, filename := range os.Args[2:] {
		file, err := os.Open(filename)
		if err != nil {
			log.Fatalf("cannot open %v", filename)
		}
		content, err := ioutil.ReadAll(file)
		if err != nil {
			log.Fatalf("cannot read %v", filename)
		}
		file.Close()
		kva := mapf(filename, string(content))      //这里形成一个map key=filename value 文件里面的内容
		intermediate = append(intermediate, kva...) //将每组键值对放在切片里面
	}

	// 与真正的 MapReduce 的一个很大的区别是，所有的
	// 中间数据都集中在一个地方，即 middle[],
	// 而不是被分成 NxM 个 bucket。

	sort.Sort(ByKey(intermediate))

	oname := "mr-out-0"
	ofile, _ := os.Create(oname)

	// 对 middle[] 中的每个不同键调用 Reduce，
	// 并将结果打印到 mr-out-0。
	i := 0
	for i < len(intermediate) {
		j := i + 1
		for j < len(intermediate) && intermediate[j].Key == intermediate[i].Key {
			j++
		}
		values := []string{}
		//统计重复的长度
		for k := i; k < j; k++ {
			values = append(values, intermediate[k].Value) //比如key为apple value为“111111....”
		}
		output := reducef(intermediate[i].Key, values)

		// 这是 Reduce 输出每行的正确格式。
		fmt.Fprintf(ofile, "%v %v\n", intermediate[i].Key, output)

		i = j
	}

	ofile.Close()
}

// 从插件文件加载应用程序 Map 和 Reduce 函数
//，例如 ../mrapps/wc.so
// loadPlugin 从插件文件加载应用程序的 Map 和 Reduce 函数
// filename 是插件文件的路径
func loadPlugin(filename string) (func(string, string) []mr.KeyValue, func(string, []string) string) {
	// 打开插件文件
	p, err := plugin.Open(filename)
	if err != nil {
		// 如果插件文件无法打开，记录错误并终止程序
		log.Fatalf("cannot load plugin %v", filename)
	}

	// 查找插件中的 "Map" 函数
	xmapf, err := p.Lookup("Map")
	if err != nil {
		// 如果在插件中找不到 "Map" 函数，记录错误并终止程序
		log.Fatalf("cannot find Map in %v", filename)
	}
	// 将找到的 "Map" 函数转换为指定的函数类型
	mapf := xmapf.(func(string, string) []mr.KeyValue)

	// 查找插件中的 "Reduce" 函数
	xreducef, err := p.Lookup("Reduce")
	if err != nil {
		// 如果在插件中找不到 "Reduce" 函数，记录错误并终止程序
		log.Fatalf("cannot find Reduce in %v", filename)
	}
	// 将找到的 "Reduce" 函数转换为指定的函数类型
	reducef := xreducef.(func(string, []string) string)

	// 返回转换后的 Map 和 Reduce 函数
	return mapf, reducef
}
