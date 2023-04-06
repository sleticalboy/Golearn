# Learning Notes
> Go 语言学习笔记

---
## Hello World
```go
// 包名
package  main
// 导包
import (
 "fmt"
)
// 主函数
func main() {
    fmt.Println("Hello World")
}
```

---
## 数据类型
### go 中的基本数据类型
- bool
    - true/false
- 数字类型
    - int8, int16, int32, int64, int（有符号）
    - uint8, uint16, uint32, uint64, uint（无符号，可以认为是非负数）
    - float32, float64（浮点型）
    - complex64, complex128（复数）
    - byte（是 uint8 的别名）
    - rune（是 int32 的别名）
- string

### 数据类型判断
```go
package main

import "fmt"
import "unsafe"

func main()  {
    value := "string"
    // 判断类型
    fmt.Printf("data type %T", value)
    // 判断长度
    fmt.Printf("data size %d", unsafe.Sizeof(value))
}
```

### 类型转换
```go
package main

import "fmt"

func main()  {
    i := 25
    j := 56.7
    fmt.Println("sum is", (i + int(j)))
}
```

---
## 变量、常量
### 变量声明
- 显示声明： var age int = 25
- 类型推断： var age = 25
- 简短声明： age ：= 25 （ 操作符的左边至少有一个变量是尚未声明的）
- 声明多个变量： age, height := 25, 175

### 常量声明
```go
package main

import "fmt"

func main()  {
    // 常量不允许再次被赋值
    const defaultType  = 26
    fmt.Println("the constant value is", defaultType)
}
```

---
## 函数
### 术语
- 参数列表
- 返回值
- 空白符

### 特殊的函数
- main 函数：
    - main 函数是 go 程序的入口；
    - main 函数应该放在 main 包中；
- init 函数
    - 所有的包都可以包含一个 init 函数
    - init 函数无参数和返回值类型
    - 作用： 可用于执行初始化任务，也可用于在开始执行之前验证程序的正确性
- 多返回值函数（两种方式）
- 大写开头的函数名或变量和小写开头的区别（类似于 java 语言中的`可见性`）
    - 大写开头： 包外可见
    - 小写开头： 包内可见

---
## package
- 作用： 包用于组织 Go 源代码，提供了更好的可重用性与可读性
- 声明： 关键字 package
    ```go
    // 声明一个包
    package pkgName
    ```
- 包的导入： import 关键字
    ```go
    package main
    // 导入包 fmt
    import "fmt"
    ```
- 包的初始化顺序（可对比 java 语言中静态代码块的初始化）
    - 先初始化包级别的变量
    - 紧接着调用 init 函数。包可以有多个 init 函数（在一个文件或分布于多个文件中），它们按照编译器解析它们的顺序进行调用
- 自定义包的导入
  - 放到 GOROOT/GOPATH 目录中之后即可像系统包一样导入
  - 使用 `go mod` 管理之后，导入完整包名

## 语句
### if-else
```go
package main
import "fmt"
func main() {
    var num = 20
    if num = 15; num > 10 {
        fmt.Println("num is greater than 10")
    } else {
        fmt.Println("num is less than 10")
    }
}
```
### 循环（go 中的循环只有 for）
```go
package main
import "fmt"
func main() {
    var num = 20
    for i := 0; i < num; i++ {
        if i > 10 {
            break
        }
        if i%2 == 0 {
            switch i {
                case 0:
                 fmt.Println("zero")
            }
            continue
        }
        fmt.Printf("i=%d ", i)
    }
}
```
### switch
```go
package main
import "fmt"
func main() {
    var num = 20
    for i := 0; i < num; i++ {
        if i > 10 {
            break
        }
        if i%2 == 0 {
            switch i {
            case 0:
                fmt.Println("Zero")
            case 2:
                fmt.Println("Two")
            case 4, 6: // 多判断
                fmt.Println("more than 2")
                // 将执行下一个 case 语句
                fallthrough
            default: // 默认 
                fmt.Println("Nothing")
            }
            continue
        }
        fmt.Printf("i=%d ", i)
    }
}
```

## 数组与切片

### 数组声明（数组一旦声明后其长度不可变）
```go
var a [3]int // 通过索引进行复制
var b := [3]int{12} // 默认值为 0
var c := [...]int{12, 13} // 可变长度，编译器自动计算长度
```

### 数组的长度是数组类型的一部分，因此不同长度、类型的数组不能相互赋值
```go
a := [3]string{"a", "b", "c"}
var b [5]string
b = a // 赋值会失败
```

### 值传递：修改后值不会影响前值
```go
a := [...]string{"a", "b", "c"}
b := a
b[0] = "aaaa" // 此时 a[0] 不会被改变
```

### 数组长度：`len(a)`

### 数组遍历
```go
// for 下标遍历
a := [...]string{"a", "b", "c"}
for i := 0; i < len(a): i++ {
    fmt.Println("element %d is %s\n", i, a[i])
}

// 内置 range 遍历
for i, v := range a {
    fmt.Println("element %d is %s\n", i, v)
}
```

### 多维数组
```go
a2d := [2][3]string{
    {"a", "b", "c"},
    {"d", "e", "f"}
}
```

### 数组切片（类似 python 列表、字符串切片）
```go
a := [5]int{1, 2, 3, 4, 5}
//  切片：含头不含尾
b := a[2:5] // -> [3, 4, 5]，当尾等于数组长度时可省略为 b := a[2:]
// 修改切片值时原数组中的值会被修改
b[0] = 333 // 导致 a[2] = 333
```

### 切片（长度可变）
```go
// 创建切片：类型，长度，容量
a := make([]int, 5, 5)
b := []string{"a", "b", "c"}
// 追加切片元素：执行 append 之后长度加 1，如果长度超出容量，则新容量变为原容量的 2 倍
a = append(a, 10)
b = append(b, "dddd")
// 追加切片
c = []string{"e", "f", "g"}
b = append(b, c...)
```

## 函数中的可变参数

### 可变参数被作为`切片`使用

## Maps

### map 创建
```go
// make(map[key-type]value-type)
m := make(map[string]int)
m2 := map[string]int {
    "key": 111,
}
```

### map 增删改查
```go
m := make(map[string]int)
// 添加元素
m["a"] = 111
m["b"] = 222
m["c"] = 333
// 删除元素
delete(m, "a")
// 修改元素
m["c"] = 3
// 查询元素：返回 key 对应的 value 和是否存在
v, b := m["b"]
if b == true {
    fmt.Println("m.b is", v)
} else {
    fmt.Println("m.b not found")
}

```

### map 遍历
```go
m := map[string]int {
    "key": 111,
    "a": 222,
    "b": 333,
    "c": 444,
}
// 遍历结果与存储顺序不一定一致（hash 表）
for k, v := range m {
    fmt.Printf("map[%s] is %d\n", k, v)
}
```

### 获取 map 长度：`len(map)`

## 字符串

### 字符串与字节切片

### 码点 `code point`

### 字符串遍历与 `rune`
```go
s := "hello"
for i, rune := range s {
    fmt.Printf("s[%d] is %c\n", i, rune)
}
```

### 字符串的不可变性（类似 java）

## 指针

### 指针指向一块内存的地址（c）
```go
a := 100
p := &a
fmt.Println("a memory addres is", p)
```

### 函数中数组传参推荐使用切片而不是指针

### 指针运算---不支持！！！

### 指针解引用
```go
a := 100
p := &a
fmt.Println("a is", *p)
```

## 结构体

### 命名结构体
```go
type Person struct {
    // 命名字段
    // 同类型的字段放在一行是为了紧凑
    name, address string
    age, gender int
    // 匿名字段
    float64
}
```

### 匿名结构体
```go
tom := struct {
    name, address string
    age, gender int
}{
    name: "tom",
    age: 18,
    gender: 1
}
```

### 结构体指针
```go
jack := &Person {
    name: "jack",
    age: 18
}
fmt.Println("jack.age is", jack.age)
```

### 嵌套结构体
```go
type Group struct {
    tl, mp Person
}
```

### 提升字段：结构体中有匿名的结构体字段
```go
type Other struct {
    tv, floor string
}
type House struct {
    width, height, length int
    Other
}

var h House
h.width = 30
h.height = 40
h.length = 50
h.Other = Other {
    tv: "ChangHong",
    floor: "fake"
}
```

### 导出结构体：结构体中的字段首字母大写则对外可见

### 结构体的可比较性：结构体中的所有字段均具有可比较性

## 方法与函数

### 方法的定义、与函数的区别
```go
// 像 java 类中的成员方法
func commonFunction() {
	fmt.Println(">>>function is like the Object method in java>>>")
}

// 像 kotlin 的扩展函数，类型用来确定属于谁的方法，类型名用于访问该类型中的字段
// 与 kotlin 不同的是，kotlin 不需要类型名默认就可以访问类型中的字段
func (man Man) getFullName(prefix string) string {
	fmt.Println(">>>method is like the ext function in kotlin>>>")
	return	prefix + " " + man.first + "-" + man.last
}
```

### 值接收器与指针接收器
```go
// 定义结构体
type Man struct {
	first, last string
}

// 值接收器：在方法内修改了类型的字段，在其他地方无效
// 相当于把这个结构体直接 copy 了一份过来，当一个结构体很大时用这种方式就不合适了
func (man Man) changeName(prefix string) {
	man.first = prefix + "-" + man.first
	fmt.Println("changed man in value recv", man)
}

// 指针接收器：在方法内修改了类型的字段，其他地方会生效
func (man *Man) changeLast(prefix string) {
	man.last = prefix + "-" + man.last
	fmt.Println("changed man in pointer recv", man)
}

func main() {
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

	(&tom).changeLast("old")
	fmt.Println("changed tom ", tom)
}
```

## 接口

## 并发

## 协程

## 信道： channel

## 缓冲信道和工作池

## Select

## Mutex

## 结构体取代类

## 多态

## Defer

## 异常
### 异常处理
### 自定义异常

## panic 和 recover

## 函数是一等公民

## 反射

## 文件操作
### 读取文件
### 写入文件
