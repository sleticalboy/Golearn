package basics

import "fmt"

type Man struct {
	first, last string
}

// 像 java 类中的成员方法
func commonFunction() {
	fmt.Println(">>>function is like the Object method in java>>>")
}

// 像 kotlin 的扩展函数，接收器类型用来确定属于谁的方法，类型名用于访问该类型中的字段
// 与 kotlin 不同的是，kotlin 不需要类型名默认就可以访问类型中的字段
func (man Man) getFullName(prefix string) string {
	fmt.Println(">>>method is like the ext function in kotlin>>>")
	return prefix + " " + man.first + "-" + man.last
}

// 方法的指针接收器与值接收器

// 值接收器：在方法内修改了类型的字段，在其他地方无效
// 相当于把这个结构体直接 copy 了一份过来，当一个结构体很大时用这种方式就不合适了
func (man Man) changeName(prefix string) {
	man.first = prefix + "-" + man.first
	fmt.Println(">>>changed man in value recv", man)
}

// 指针接收器：在方法内修改了类型的字段，其他地方会生效
// 把结构体的内存地址传过来了，因此对结构体的任何修改都会反应到其他地方
func (man *Man) changeLast(prefix string) {
	man.last = prefix + "-" + man.last
	fmt.Println(">>>changed man in pointer recv", man)
}

func methodRun() {
	println("\nmehtod & function run")

	commonFunction()

	var man Man
	man.first = "first"
	man.last = "last"
	fmt.Println("man full name is", man.getFullName("prefix"))

	tom := Man{
		first: "tom",
		last:  "jack",
	}
	fmt.Println("raw tom ", tom)
	tom.changeName("young")
	fmt.Println("changed tom ", tom)

	p := &tom
	p.changeLast("old")
	fmt.Println("changed tom ", *p)
}
