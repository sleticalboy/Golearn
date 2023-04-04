package main

import (
	"fmt"
	"golearn/src/basic"
	"golearn/src/hello"
)

func init() {
	fmt.Println("go package init exec.")
}

func main() {
	fmt.Println("go main exec.")
	hello.Say()
	basic.Main()
}
