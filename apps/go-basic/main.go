package main

import (
	"fmt"
	"github.com/sleticalboy/golearn/apps/go-basic/basics"
	"github.com/sleticalboy/golearn/apps/go-basic/files"
	"github.com/sleticalboy/golearn/apps/go-basic/hello"
	"github.com/sleticalboy/golearn/apps/go-basic/http"
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
