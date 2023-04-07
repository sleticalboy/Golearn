package basic

import (
	"fmt"
	"sync"
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

func testBufferChan(notify chan int) {
	fmt.Println("testBufferChan() run....")
	for i := 0; i < 5; i++ {
		notify <- i
		fmt.Println("testBufferChan() notify", i)
	}
	close(notify)
	fmt.Println("testBufferChan() exit....")
}

func bufferedChan() {
	// 缓冲信道
	ch := make(chan int, 2)
	go testBufferChan(ch)
	for v := range ch {
		fmt.Printf("read value %d from ch\n", v)
		time.Sleep(500 * time.Millisecond)
	}
}

func work(group * sync.WaitGroup, i int) {
	fmt.Printf("worker %d starts working...\n", i)
	time.Sleep(500 * time.Millisecond)
	fmt.Printf("worker %d finished working...\n", i)
	// 计数器 -1
	group.Done()
}

func waitGroup() {
	fmt.Println("\nwaitGroup() enter")
	// 比较像 java 中的 CountDownLatch
	var wg sync.WaitGroup
	for i := 0; i < 5; i++ {
		// 计数器 +1
		wg.Add(1)
		go work(&wg, i)
	}
	// 等待所有协程结束（计数器为 0 时结束阻塞）
	wg.Wait()
	fmt.Println("waitGroup() exit")
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

	fmt.Println()

	bufferedChan()

	waitGroup()
}
