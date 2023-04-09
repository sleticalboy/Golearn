package main

import (
	"fmt"
	"golearn/src/files"
	"golearn/src/hello"
)

func init() {
	fmt.Println("go package init exec.")
}

func main() {
	fmt.Println("go main exec.")
	hello.Main()
	//basics.Main()
	files.Main()
}
