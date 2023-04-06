package basic

import "fmt"

func functionRun() {
	println("\nfunction of golang")

	// 调用函数
	anotherFunction("Hello Function")

	length, area := multiRets(4, 5)
	// length, area := multiRets(4, 5)
	fmt.Println("length:", length, ", area:", area)

	// 若要丢弃某返回值，使用空白符代替
	anotherLength, _ := multiRets2(3, 9)
	fmt.Println("another length:", anotherLength)
}

func anotherFunction(param string) {
	fmt.Println("another function was called with:", param)
}

// 多返回值函数: 方式一匿名返回
// 返回矩形的周长和面积
func multiRets(width int, height int) (int, int) {
	return (width + height) * 2, width * height
}

// 多返回值函数：方式二命名返回
// 返回矩形的周长和面积
func multiRets2(width int, height int) (length, area int) {
	length = (width + height) * 2
	area = width * height
	return
}
