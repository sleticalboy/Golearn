package http

import (
	"fmt"
	"os"
)

func Main() {
	println("\nhttp run")

	flag := os.Args[1]

	if flag == "-s" {
		fmt.Println("run as server")
		httpServerRun()
	} else if flag == "-c" {
		fmt.Println("run as client")
		method := "GET"
		if len(os.Args) > 2 && os.Args[2] == "-p" {
			method = "POST"
		}
		httpClientRun(method)
	} else {
		fmt.Println("error args")
	}
}
