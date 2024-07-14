package main

// 一个用于 MapReduce 的词频统计应用 "plugin"。
// 使用以下命令构建plugin：
// go build -buildmode=plugin wc.go

import "6.5840/mr"
import "unicode"
import "strings"
import "strconv"

// map 函数会对每个输入文件调用一次。第一个参数是输入文件的名称，
// 第二个参数是文件的全部内容。你应该忽略输入文件的名称，只需关
// 注 contents 参数。返回值是一个键/值对的切片。

func Map(filename string, contents string) []mr.KeyValue {
	// 用于检测单词分隔符的函数。
	ff := func(r rune) bool { return !unicode.IsLetter(r) }

	// 将内容拆分成单词数组。
	words := strings.FieldsFunc(contents, ff)

	kva := []mr.KeyValue{}
	for _, w := range words {
		kv := mr.KeyValue{w, "1"}
		kva = append(kva, kv)
	}
	return kva
}

// reduce 函数会对 map 任务生成的每个键调用一次，传入一个由所有 map 任务
// 为该键创建的值的列表。

func Reduce(key string, values []string) string {
	// 返回该单词的出现次数。
	return strconv.Itoa(len(values))
}
