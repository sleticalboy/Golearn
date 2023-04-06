package basic

import (
	"fmt"
	"time"
)

func doBusyWork() {
	fmt.Println("I'm doing busy work!")
}

func download(finished chan bool) {
	fmt.Println("\nI'm downloading big file!")
	time.Sleep(1500 * time.Millisecond)
	fmt.Println("download over, notify other routine")
	// 下载完成，通过信道通知
	finished <- true
}

func routineRun() {
	println("\nroutine run")

	// 执行后立即返回，不会阻塞主协程
	go doBusyWork()
	// 主协程睡眠，等待子协程执行完毕
	time.Sleep(1 * time.Second)
	fmt.Println("busy work is done!")

	// 使用信道通讯
	finished := make(chan bool)
	// 执行下载任务
	go download(finished)
	// 从信道读取值，在其他地方没有写入操作时会一直阻塞
	fmt.Println("\nwaiting for download done!")
	<-finished
	fmt.Println("big file is download!")
	close(finished)
}
