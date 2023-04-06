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

## Maps

## 字符串

## 指针

## 结构体

## 方法与函数

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
