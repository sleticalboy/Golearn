package main

import "fmt"

type gethw struct {
	field interface{}
}

var v gethw

// SetField 外部可以访问
func SetField(val interface{}) {
	v.field = val
	fmt.Printf("SetField() '%v' from plugin\n", val)
}

// GetField 获取字段
func GetField() interface{} {
	return v.field
}
