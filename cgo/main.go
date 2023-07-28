package main

// #include <stdio.h>
/*
int hello(const char *str) {
	printf("hello() called with: %s\n", str);
	return 26;
}
*/

// #include "foo.h"
// #include <stdlib.h>
import "C"
import (
	"fmt"
	"unsafe"
)

func main() {
	fmt.Println("The first cgo example...")
	// cStr := C.CString("Hello cgo way 1!")
	// r := C.hello(cStr)
	// fmt.Printf("c go ret: %v\n", r)
	// C.free(unsafe.Pointer(cStr))

	cStr := C.CString("Hello cgo way 2!")
	r := C.foo_func(cStr)
	fmt.Printf("c go ret: %.1f\n", float64(r))
	C.free(unsafe.Pointer(cStr))
}
