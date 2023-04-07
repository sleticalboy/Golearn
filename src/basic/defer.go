package basic

import "fmt"

func doLast(i int) {
	fmt.Printf("doLast() called with: %d\n", i)
}

func deferRun() {
	println("\ndefer run")

	i := 20
	// deferRun 函数退出之前会执行，i 已经被读取到 doLast() 栈中，最后会输出 20
	defer doLast(i)

	// 修改变量值，不会影响 doLast() 输出结果
	i = 30
	fmt.Printf("before deferRun() exit with: %d\n", i)
}