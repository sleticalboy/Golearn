package main

import "fmt"

func init() {
    fmt.Println("init function called.")
}

func main() {
    fmt.Println("main function called.")

    // 先执行 init 函数，后执行 main 函数
}
