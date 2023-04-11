package main

import (
	"com.binlee/goweb/samples"
	"os"
)

func main() {
	if len(os.Args) == 1 {
		samples.Html()
	} else if os.Args[1] == "json" {
		samples.Json()
	} else if os.Args[1] == "form" {
		samples.FormLogin()
	} else if os.Args[1] == "others" {
		samples.Others()
	}
}
