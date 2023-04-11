package basics

import (
	"fmt"
	"os"
)

// 定义结构体
type myError struct {
	errorCode int
	errorStr  string
}

// 实现 error 接口
func (err *myError) Error() string {
	return fmt.Sprintf("my error code: %d, error: %s\n", err.errorCode, err.errorStr)
}

func someWork(val int) (string, error) {
	if val <= 0 {
		return "", &myError{val, "Invalid value"}
	}
	return fmt.Sprintf("val is %d", val), nil
}

func errorRun() {
	println("\nerror Run")

	f, err := os.Open("/hello.txt")
	if err != nil {
		fmt.Printf("open failed: '%s'\n", err.Error())
	} else {
		fmt.Printf("open '%s' success\n", f.Name())
	}

	if s, err := someWork(-1); err != nil {
		fmt.Println(err.Error())
	} else {
		fmt.Println(s)
	}

	if s, err := someWork(1); err == nil {
		fmt.Println(s)
	} else {
		fmt.Println(err.Error())
	}
}
