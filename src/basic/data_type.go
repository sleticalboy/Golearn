package basic

import (
	"fmt"
	"unsafe"
)

// 常量
const sourceType = 20

func dataTypeRun() {
	fmt.Println("\ndata type of golang.")
	// bool
	// int8, int16, int32, int64, int
	// uint8, uint16, uint32, uint64, uint
	// float32, float64
	// complex64, complex128
	// byte
	// rune
	// string
	age := 89
	fmt.Printf("type of age is %T and size is %d\n", age, unsafe.Sizeof(age))
	// 复数
	var c1 = 3 - 4i
	var c2 = complex(5, 12)
	fmt.Println("c1 + c2 =", c1+c2)
	fmt.Println("c1 * c2 =", c1*c2)

	fmt.Println("source type", sourceType)
}
