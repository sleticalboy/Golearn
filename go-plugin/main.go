package main

import (
	"fmt"
	"plugin"
)

func main() {
	path := "/home/binlee/code/Golearn/sample-plugin/sample-plugin.so"
	plug, err := plugin.Open(path)
	fmt.Printf("main() plug: %v, err: %v\n", plug, err)
	if err != nil {
		panic(err)
	}

	m, err := plug.Lookup("SetField")
	if err != nil {
		panic(err)
	}

	m.(func(interface{}))(map[string]interface{}{"name": "ben"})

	m, err = plug.Lookup("GetField")
	if err != nil {
		panic(err)
	}

	field := m.(func() interface{})()
	fmt.Printf("field is: %v\n", field)
}
