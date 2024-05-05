package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"
	"unsafe"
)

// #include <stdlib.h>
// #include "samples.h"
import "C"

// 实现 c 中定义的 extern void cgoCallback(int value); 函数
//
//export cgoCallback
func cgoCallback(value int) {
	fmt.Printf("cgoCallback() c++ value: %d\n", value)
}

func main() {
	fmt.Println("The first cgo example...")

	var strs []string
	fmt.Printf("strs: %v\n", strs)

	ints := make([]int, 0)
	fmt.Printf("ints: %v\n", ints)

	crashChan := make(chan struct{})

	// 监听指定的信号
	signalCh := make(chan os.Signal, 1)
	// signal.Notify(signalCh, syscall.SIGINT, syscall.SIGTERM, syscall.SIGKILL)
	signal.Notify(signalCh, syscall.SIGABRT, syscall.SIGALRM, syscall.SIGBUS, syscall.SIGCHLD,
		syscall.SIGCLD, syscall.SIGCONT, syscall.SIGFPE, syscall.SIGHUP, syscall.SIGILL, syscall.SIGINT,
		syscall.SIGIO, syscall.SIGIOT, syscall.SIGKILL, syscall.SIGPIPE, syscall.SIGPOLL, syscall.SIGPROF,
		syscall.SIGPWR, syscall.SIGQUIT, syscall.SIGSEGV, syscall.SIGSTKFLT, syscall.SIGSTOP, syscall.SIGSYS,
		syscall.SIGTERM, syscall.SIGTRAP, syscall.SIGTSTP, syscall.SIGTTIN, syscall.SIGTTOU, syscall.SIGUNUSED,
		syscall.SIGURG, syscall.SIGUSR1, syscall.SIGUSR2, syscall.SIGVTALRM, syscall.SIGWINCH, syscall.SIGXCPU,
		syscall.SIGXFSZ)

	go func() {
		// 阻塞等待信号
		sig := <-signalCh
		fmt.Printf("Received signal: %v, -> %d\n", sig, sig)

		// 处理信号
		switch sig {
		case syscall.SIGINT, syscall.SIGTERM:
			// 处理程序终止信号
			// 可以在此进行资源清理等操作
			fmt.Println("Program terminated.")
			crashChan <- struct{}{}
		case syscall.SIGKILL:
			// 处理程序被系统杀掉的情况
			// 可以在此记录日志或进行其他处理
			fmt.Println("Program killed by system.")
			os.Exit(1)
		default:
			fmt.Printf("Receive sig: %v\n", sig)
		}
	}()

	cStr := C.CString("Hello cgo way 2!")
	r := C.hello(cStr)
	fmt.Printf("c go ret: %.1f\n", float64(r))
	C.free(unsafe.Pointer(cStr))
	for i := 0; i < 50; i++ {
		go func() {
			C.start_loop(0)
		}()
		time.Sleep(time.Millisecond * 2)
	}
	// 开始循环
	go func() {
		C.start_loop(3)
		os.Exit(int(syscall.SIGTERM))
	}()

	// 模拟程序异常退出
	// go func() { C.crash() }()

	// 等待程序异常退出信号
	s := <-crashChan
	fmt.Printf("program finished: %v", s)
}
