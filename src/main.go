package main

import (
	"com.binlee/golearn/basics"
	"com.binlee/golearn/files"
	"com.binlee/golearn/hello"
	"com.binlee/golearn/http"
	"fmt"
)

func init() {
	fmt.Println("go package init exec.")
}

func main() {
	fmt.Println("go main exec.")
	hello.Main()
	basics.Main()
	files.Main()
	http.Main()
}
