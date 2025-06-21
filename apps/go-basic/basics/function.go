package basics

import "fmt"

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

func commonFunc() {
	// 调用函数
	anotherFunction("Hello Function")

	length, area := multiRets(4, 5)
	// length, area := multiRets(4, 5)
	fmt.Println("length:", length, ", area:", area)

	// 若要丢弃某返回值，使用空白符代替
	anotherLength, _ := multiRets2(3, 9)
	fmt.Println("another length:", anotherLength)
}

// 自定义函数类型: 名字+签名
type add func(a, b int) int

func superFunc() {
	println("\nsuper function run")

	// 把函数赋值给变量
	fun := func() {
		fmt.Println("func variable run")
	}
	// 调用函数
	fun()

	// 匿名函数调用
	func(param string) {
		fmt.Printf("anonymous func run with '%s'\n", param)
	}("anonymous param")

	// 实现自定义函数
	var a add = func(a, b int) int {
		return a + b
	}
	// 调用自定义函数
	sum := a(3, 9)
	fmt.Println("sum is", sum)
}

// 定义函数，入参为函数类型
func high(convert func(a int) string) {
	s := convert(50)
	fmt.Printf("convert result is '%s'\n", s)
}

// 定义函数，返回值为一个函数
func fun() func(a, b int) string {
	return func(a, b int) string {
		return fmt.Sprintf("sum is %d", a+b)
	}
}

func higherFunc() {
	println("\nhigher function run")

	// 调用高阶函数，传入匿名函数
	high(func(a int) string {
		return fmt.Sprintf("input is %d", a)
	})

	// 调用高阶函数
	fmt.Println("sum of 30 and 90 is", fun()(30, 90))

	// 闭包：匿名函数调用函数体外的变量
	a := 100
	func() {
		fmt.Println("clouser invoke var in outter", a)
	}()
}

func functionRun() {
	println("\nfunction of golang")

	// 普通函数
	commonFunc()

	// 头等函数
	superFunc()

	// 高阶函数
	higherFunc()
}
