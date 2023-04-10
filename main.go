package main

import (
	"fmt"
	// "golearn/src/hello"
	// "golearn/src/basics"
	// "golearn/src/files"
	"golearn/src/http"
)

func init() {
	fmt.Println("go package init exec.")
}

func main() {
	fmt.Println("go main exec.")
	// hello.Main()
	// basics.Main()
	// files.Main()
	http.Main()
}
