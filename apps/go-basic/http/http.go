package http

import (
	"fmt"
	"os"
)

func Main() {
	fmt.Printf("\nhttp run, args: %+v\n", os.Args)

	var flag = "-c"
	if len(os.Args) > 1 {
		flag = os.Args[1]
	}

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
