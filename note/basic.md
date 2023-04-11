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
    return prefix + " " + man.first + "-" + man.last
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

### 定义并实现接口
```go
// 定义接口
type Runnable interface {
    run()
}
// 定义数据结构
type DownloadTask struct {
    url string
}

// 实现接口
func (task *DownloadTask) run() {
    fmt.Printf("download task run %s\n", task.url)
}

func main() {
    dt := &DownloadTask{
        url: "https://www.example.com",
    }
    // 调用接口
    dt.run()
}
```

### 空接口
```go
// 空接口没有实现任何方法，所以任何类都实现了空接口
interface {}
```

### 类型断言：提取接口的底层值（哪个类实现了接口？）
```go
func itfAssert(itf interface{}) {
    switch itf.(type) {
        case string:
        case int:
        default:
    }
}
```

## 并发

### 并发与并行
- 并行：一边写代码一边听歌（同时进行）
- 并发：写代码、运行代码、解决 bug（时间片轮转执行不同的任务）

### go 语言的并发 `协程`
```go
// 定义函数
func run() {
    fmt.Println("running...")
    time.Sleep(1 * time.Second)
}

// 使用 go 关键字开启协程
func main() {
    go run()
}
```

### 协程通信：`chan` 信道

- 信道默认是双向的
- 双向信道可读可写
- 单向信道只能读或写
- 遍历信道使用 `range`
- 读取信道返回值：`v, ok <- ch`，ok == true 表示管道正常工作，否则表示关闭

```go
func download(finished chan bool) {
    fmt.Println("\nI'm downloading big file!")
    time.Sleep(1500 * time.Millisecond)
    fmt.Println("download over, notify other routine")
    // 下载完成，通过信道通知
    finished <- true
}

func main() {
    // 使用信道通讯，创建管道 make(chan type)
    finished := make(chan bool)
    // 执行下载任务
    go download(finished)
    // 从信道读取值，在其他地方没有写入操作时会一直阻塞
    fmt.Println("\nwaiting for download done!")
    <-finished
    fmt.Println("big file is download!")
    close(finished)
}
```

### 管道死锁
- 当管道只有写入没有读取时，就会发生死锁

### 缓冲信道和工作池

- 缓冲信道：有 buffer 的信道，当写入值超出容量时才阻塞，否则不阻塞
  - 写入时，数据会先写入缓冲区，当缓冲区已满时才阻塞
  - 读取时，先从缓冲区读取数据，缓冲区无数据时才阻塞

```go
// 创建缓冲信道
ch := make(chan int, 5/* capacity */)
ch <- 1
ch <- 2
ch <- 3
fmt.Println(<- ch)
fmt.Println(<- ch)
fmt.Println(<- ch)
```

### 工作池 `WorkerPools`

- `WaitGroup` 用于等待一批协程结束，程序会一致阻塞等待这些协程全部结束
- 工作原理比较像 java 中的 `CountDownLatch`
- go 中没有线程池，我们要利用 `WaitGroup` 自己封装线程池

```go
func work(group *sync.WaitGroup, i int) {
    fmt.Printf("worker %d is working...\n", i)
    time.Sleep(1 * time.Second)
    // 计数器 -1
    group.Done()
    fmt.Printf("worker %d works done...\n", i)
}

func main() {
    fmt.Println("main() enter")
    wg := sync.WaitGroup
    for i := 0; i < 3; i ++ {
        // 计数器 +1
        wg.add(1)
        go work(&wg, i)
    }
    // 等待工作结束（计数器为 0 时结束阻塞）
    wg.Wait()
    fmt.Println("main() exit")
}
```

### Select

- 用于在多个发送/接收信道操作中进行选择
- select 语句会一致阻塞，直到发送/接收操作就绪
- 语法类似 `switch` 语句，没有 `default` 时可能会产生死锁问题
- 多个信道就绪的情况下 `select` 是随机选取的，这与 `switch` 不同
- 类似 java nio 编程中的 selector？

### Mutex

- 用于并发编程时控制边界变量的读写控制

## 结构体组合

- go 不支持继承，但支持组合
- 多个结构体嵌套组合成复杂的结构体

## 多态

- go 中的多态：
  - 多个结构体实现同一个接口
  - 声明接口的切片保存所有结构体，然后遍历切片调用接口方法
- 像 java 中的一个类可以实现多个接口

## Defer

- 含有 `defer` 语句的函数，会在该函数返回之前执行 `defer` 之后的函数
- 像 java 中的 `try...finally` 语句
- 像 c/c++ 函数中定义的 label 一样，在函数返回之前去执行

```go
func doLast(i int) {
    fmt.Printf("doLast() called with: %d\n", i)
}

func main() {
    println("\ndefer run")
    i := 20
    // deferRun 函数退出之前会执行，i 已经被读取到 doLast() 栈中，最后会输出 20
    defer doLast(i)
    // 修改变量值，不会影响 doLast() 输出结果
    i = 30
    fmt.Printf("before deferRun() exit with: %d\n", i)
}
```

## 异常

### 异常处理
```go
func main() {
    f, err := os.Open("/hello.txt")
    if err != nil {
        fmt.Printf("open failed: '%s'\n", err.Error())
    } else {
        fmt.Printf("open '%s' success\n", f.Name())
    }
}
```

### 自定义异常

- 需要实现 `error` 接口

```go
// 定义结构体
type myError struct {
    errorCode int
    errorStr  string
}
// 实现 error 接口
func (err *myError) Error() string {
    return fmt.Sprintf("my error code: %d, error: %s\n", err.errorCode, err.errorStr)
}
```

## panic 和 recover

- panic 程序终止运行
  - 发生了一个不可恢复的错误，程序不能继续运行，比如：c10k 问题
  - 编程错误，比如 `1 / 0`
- recover 重获对程序的控制
  - 只有在相同的协程中 `recover` 才会起作用，否则不起作用
  - 打印堆栈 `debug.PrintStack()`

## 函数是一等公民

### 函数变量
```go
func main() {
    // 把函数复制给变量
    fun := func() {
        fmt.Println("func variable run")
    }
    // 调用函数
    fun()
}
```

### 匿名函数
```go
func main() {
    // 匿名函数调用
    func(param string) {
        fmt.Printf("anonymous func run with '%s'\n", param)
    }("anonymous param")
}
```

### 自定义函数类型
```go
// 自定义函数类型: 名字+签名
type add func(a, b int) int

func main() {
    // 实现自定义函数
    var a add = func(a, b int) int {
        return a + b
    }
    // 调用自定义函数
    sum := a(3, 9)
    fmt.Println("sum is", sum)
}
```

### 高阶函数

- 满足以下任意条件即为高阶函数：
  - 接收一个或多个函数作为参数
  - 返回值是一个函数
- 高阶函数常用来设计流式 API，比如：java 中的 stream API

```go
// 定义函数，入参为函数类型
func high(convert func (a int) string) {
    s := convert(50)
    fmt.Printf("convert result is '%s'\n", s)
}

func main() {
    // 调用高阶函数，传入匿名函数
    high(func (a int) string {
        return fmt.Sprintf("input is %d", a)
    })
}
```

### 闭包
```go
func main() {
    // 闭包：匿名函数调用函数体外的变量
    a := 100
    func() {
        fmt.Println("clouser invoke var in outter", a)
    }()
}
```

## 反射

- go 的反射有着比较严格的要求，对一切非导出的字段、方法进行反射都是不允许的

```go
type ISay interface {
    Say()
}

type Nested struct {
    Cc string
}

// 实现接口 ISay
func (n Nested) Say() {
    fmt.Println("sample.Say() run...")
}

// 实现接口 Stringer
func (n Nested) String() string {
    return fmt.Sprintf("Nestedt{Cc: %s}", n.Cc)
}

type foo struct {
    a string
    b int
    // 二级嵌套
    N Nested
}

type sample struct {
    name    string
    age     int
    salary  float32
    subject []string
    // 一级嵌套
    fo      foo
}

func main(o any) {
    // 自定义结构体（稍微复杂点）
    s := sample{
        name:    "tom",
        age:     24,
        salary:  9000.0,
        subject: []string{"java", "go", "c++"},
        fo: foo{
            a: "aaa",
            b: 111,
            N: Nested{
                Cc: "nested struct depth 2",
            },
        },
    }

    // 反射获取类型（java 的 Object.getClass()）
    t := reflect.TypeOf(s)
    // 反射获取值
    v := reflect.ValueOf(s)
    // 返回实际类型（java 的 Class）枚举类型 Kind
    kind := t.Kind()
    // Invalid Kind = iota/Bool
    // /Int/Int8/Int16/Int32/Int64/Uint/Uint8/Uint16/Uint32/Uint64/Uintptr
    // Float32/Float64 float、double 类型
    // Complex64/Complex128 复数类型
    // Array/Slice 数组、切片类型
    // Chan 信道
    // Func 函数
    // Interface 接口
    // Map 字典
    // Pointer 指针
    // UnsafePointer
    // String 字符串
    // Struct 结构体
    if kind == reflect.Struct {
        // 获取方法和字段数量，然后通过索引去反射
        nf := t.NumField()
        for i := 0; i < nf; i++ {
            var tf reflect.StructField = t.Field(i)
            fmt.Printf("%s.%s[%v]", t.Name(), tf.Name, tf.Type)
        }
    }
    // 反射调用无返回值方法
    if m, ok := reflect.TypeOf(s.fo.N).MethodByName("Say"); ok {
        m.Func.Call([]reflect.Value{0: reflect.ValueOf(s.fo.N)})
    }
    // 反射调用有返回值方法
    if m, ok := reflect.TypeOf(s.fo.N).MethodByName("String"); ok {
        ret := m.Func.Call([]reflect.Value{0: reflect.ValueOf(s.fo.N)})
        fmt.Printf("reflect call String() '%s'\n", ret[0].String())
    }
}
```

## 文件操作

### 读取文件全部内容
```go
func main() {
    // 读取文件全部内容
    if content, err := os.ReadFile(path); err == nil {
        fmt.Printf(string(content))
    } else {
        fmt.Println(err)
    }
}
```

### 读取文件部分内容
```go
func main() {
    // 读取文件部分内容（通过 buffer 读取）
    f, err := os.Open(path)
    // 关闭文件
    defer func() {
        _ = f.Close()
    }()

    if err != nil {
        fmt.Println(err)
        return
    }
    buf := make([]byte, 256)
    if info, err := f.Stat(); err == nil {
        count := info.Size() / 256
        reset := info.Size() % 256
        fmt.Printf("%d -> %d...%d\n", info.Size(), count, reset)
    }
    for {
        readBytes, err := f.Read(buf)
        if err == io.EOF {
            fmt.Println("readFiles() hit EOF!")
            break
        }
        if err != nil {
            fmt.Println(err)
            break
        }
        // buf 没有被填满，直接打印会出现脏数据
        if readBytes != len(buf) {
            fmt.Printf(string(buf[0:readBytes]))
        } else {
            fmt.Printf(string(buf))
        }
    }
}
```

### 逐行读取文件
```go
func main() {
    // 打开文件
    f, err = os.Open(path)
    defer func() {
        _ = f.Close()
    }()
    if err != nil {
        fmt.Println(err)
        return
    }
    // 通过 scanner 逐行读取文件
    scanner := bufio.NewScanner(f)
    for scanner.Scan() {
        text := strings.TrimSpace(scanner.Text())
        if len(text) == 0 {
            continue
        }
        fmt.Println(text)
    }
}
```

### 写入文件
```go
func main() {
    // 创建文件
    f, err := os.Create(path)
    defer func() {
        _ = f.Close()
    }()
    if err != nil {
        fmt.Println(err)
        return
    }
    fmt.Printf("writeFiles() '%s'\n", f.Name())
    // 写入字符串（字节数组）
    writeBytes, err := f.WriteString("first line\n")
    if err != nil {
        fmt.Println(err)
        return
    }
    fmt.Printf("write bytes: %d\n", writeBytes)
    // 通过 fmt.Fprintf() 将数据写入到文件中
    writeBytes, err = fmt.Fprintf(f, "write via fmt.Fprintf: %d\n", 90)
    if err != nil {
        fmt.Println(err)
        return
    }
    fmt.Printf("write bytes: %d\n", writeBytes)
}
```

### 追加文件
```go
func main() {
    // 追加文件（打开时以 append 方式打开）
    f, err = os.OpenFile(path, os.O_APPEND|os.O_WRONLY, 0644)
    defer func() {
        _ = f.Close()
    }()
    writeBytes, err = f.WriteString("append line to file\n")
    if err != nil {
        fmt.Println(err)
        return
    }
    fmt.Printf("write bytes: %d\n", writeBytes)
}
```

## 网络编程

### 服务端
```go
var (
    handlers = map[string]func(w http.ResponseWriter, r *http.Request){
        "/api/v1/hello_world": func(response http.ResponseWriter, request *http.Request) {
            if request.Method == "GET" {
                response.WriteHeader(200)
                str := "Hello Go(GET) Http World!\n"
                _, _ = response.Write(bytes.NewBufferString(str).Bytes())
            } else if request.Method == "POST" {
                body := make([]byte, request.ContentLength)
                _, _ = request.Body.Read(body)
                fmt.Printf("request body is %s\n", string(body))
                str := "Hello Go(POST) Http World!\n"
                _, _ = response.Write(bytes.NewBufferString(str).Bytes())
                fmt.Printf("response body is %s\n", str)
            }
        },
    }
)

func httpServerRun() {
    // 注册 api handler
    for api, handler := range handlers {
        http.HandleFunc(api, handler)
    }
    // 创建服务器
    server := &http.Server{
        Addr: "127.0.0.1:8099",
        Handler: nil,
    }
    // 启动服务器
    err := server.ListenAndServe()
    // err = http.ListenAndServe(server.Addr, nil)
    if err != nil {
        fmt.Printf("Listen http://127.0.0.1.8099 failed: %e\n", err)
    }
}
```

### 客户端
```go
func httpClientRun(method string) {
    requestUrl := "http://127.0.0.1:8099/api/v1/hello_world"
    jsonStr := []byte(`{"hello": "world"}`)
    var (
        response *http.Response
        err      error
    )
    if method == "POST" {
        response, err = http.Post(requestUrl, "application/json", bytes.NewBuffer(jsonStr))
    } else {
        response, err = http.Get(requestUrl)
    }
    if err != nil {
        fmt.Println(err)
        return
    }
    defer func() {
        _ = response.Body.Close()
    }()
    // 记得处理 EOF
    // body := make([]byte, response.ContentLength)
    // _, err = response.Body.Read(body)
    // if err != nil && err != io.EOF {
    //      fmt.Println(err)
    //      return
    // }
    // 使用系统 API 时，系统内部会自动消费掉 EOF
    body, err := io.ReadAll(response.Body)
    if err != nil {
        return
    }
    fmt.Println(string(body))
}
```