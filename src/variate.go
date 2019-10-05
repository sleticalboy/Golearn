package main

import "fmt"

func main() {
    // 显示声明
    // var age int = 25
    // 类型推断
    var age = 25
    // 同时声明多个变量(相同类型)
    var weight, height = 75, 175
    fmt.Println("my age is:", age, ", height is:", height, "cm and weight is:", weight, "kg")
    // 同时声明多个不同类型的变量
    var (
        name = "leebin"
        age2 = 25
        height2 = 175
    )
    fmt.Println("my name is", name, "and I'm", age2, "years old, height is", height2, "cm")
}
