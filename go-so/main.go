package main

/*
#include <stdio.h>
#include <stdlib.h>

typedef struct { char *name; int age; char *home; } Person;
typedef unsigned long long GoUint64_;
*/
import "C"
import (
	"fmt"
)
import "time"

//export HelloGo
func HelloGo(s *C.char) {
	fmt.Printf("HelloGo, %v!\n", C.GoString(s))
}

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
func Fib(value int) uint64 {
	// fmt.Printf("Fib() in go, value: %d", value)
	if value < 2 {
		return 1
	}

	var f = uint64(1)
	var s = uint64(1)
	var temp uint64
	for i := 0; i < value; i++ {
		if i < 2 {
			continue
		}
		temp = s
		s = s + f
		if i > 90 {
			fmt.Printf("Fib() %d -> %d\n", i, s)
		}
		f = temp
	}
	return s
}

//export GoHome
func GoHome(ming *C.Person) {
	fmt.Printf("GoHome() -> name: %v, age: %v, home: %v\n", C.GoString(ming.name), ming.age, C.GoString(ming.home))
}

//export GoHome2
func GoHome2(ming C.Person) {
	fmt.Printf("GoHome2() -> name: %v, age: %v, home: %v\n", C.GoString(ming.name), ming.age, C.GoString(ming.home))
}

func main() {
	fmt.Println(Fib(93))
}
