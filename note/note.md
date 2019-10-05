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
- main 函数
- init 函数
- 多返回值函数（两种方式）

---
## package
- 声明
- 导出与非导出（类似于 java 语言中的`可见性`）

## 语句
### if-else
### 循环
### switch

## 数组与切片

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
