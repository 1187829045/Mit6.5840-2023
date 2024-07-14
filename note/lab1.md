实现一个工作进程来调用应用程序的 Map 和 Reduce 函数，并处理文件的读取和写入；同时还要实现一个协调器进程，负责将任务分发给工作进程，并处理失败的工作进程。
mrsequential.go 会将输出存储在 mr-out-0 文件中。输入来自名为 pg-xxx.txt 的文本文件。
你的任务是实现一个分布式 MapReduce，包括两个程序：协调器和工作进程。只有一个协调器进程，和一个或多个并行执行的工作进程.
工作进程通过 RPC 与协调器通信。每个工作进程将向协调器请求任务，从一个或多个文件中读取任务的输入，执行任务，并将任务的输出写入一个或多个文件。协调器
应该能够在合理的时间内（本实验中使用十秒）发现工作进程没有完成任务，并将同一个任务分配给另一个工作进程。
实现放在 mr/coordinator.go、mr/worker.go 和 mr/rpc.go 中。
$ cd ~/6.5840/src/main
$ bash test-mr.sh
*** Starting wc test.
You can change ret := false to true in the Done function in mr/coordinator.go so that the coordinator exits immediately. Then:

$ bash test-mr.sh
Map 阶段应将中间键分成 nReduce 个桶，其中 nReduce 是传递给 MakeCoordinator() 的 reduce 任务数。每个 Mapper 应创建 nReduce 个中间文件，供 reduce 任务使用。
Worker 实现应将第 X 个 reduce 任务的输出放入名为 mr-out-X 的文件中。每行格式应为 Go 的 "%v %v" 格式，调用时传递键和值。可以参考 main/mrsequential.go 中的注释以确保格式正确，否则测试脚本会失败。
Worker 应将 Map 阶段的中间输出放入当前目录中的文件，在后续作为 Reduce 任务的输入时可以读取它们。
main/mrcoordinator.go 期望 mr/coordinator.go 实现一个 Done() 方法，当 MapReduce 作业完全完成时返回 true；此时，mrcoordinator.go 将退出。
作业完全完成时，工作进程应退出。一个简单的实现方法是使用 call() 的返回值：如果工作进程无法联系协调器，它可以假设协调器已因为作业完成而退出，因此工作进程也可以终止。根据你的设计，可能还会发现给工作进程一个“请退出”的伪任务也很有帮助。
提示：可以从 mr/worker.go 的 Worker() 函数开始修改，向协调器发送 RPC 请求以获取任务。然后修改协调器以响应一个尚未开始的 Map 任务的文件名。然后修改 Worker 以读取该文件并调用应用的 Map 函数，类似于 mrsequential.go 的操作方式。
应用程序的 Map 和 Reduce 函数在运行时使用 Go 的插件包加载，文件名以 .so 结尾。