package main

import "C"
import "fmt"
import "time"

//export GoFib
func GoFib() {
	start := time.Now().UnixMilli()
	j := 0
	for i := 0; i < 100000; i++ {
		// fmt.Printf("%d", i)
		j++
	}
	end := time.Now().UnixMilli()
	fmt.Printf("go cost: %d ms\n", end-start)
}

//export Fib
func Fib(value int) int {
	// fmt.Printf("Fib() in go, value: %d", value)
	if value == 1 || value == 2 {
		return value
	}
	return Fib(value-1) + Fib(value-2)
}

func main() {}
