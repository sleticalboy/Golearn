package main

import "fmt"

func main() {
    var a [3]int
    a[0] = 12
    a[1] = 13
    a[2] = 13
    fmt.Println("this array is", a)
    // 不赋值的元素默认值为 0
    b := [3]int{3}
    fmt.Println("another array is", b)
    // 不指定数组长度
    c := [...]int{2, 8, 0}
    fmt.Println("dynamic array is", c)

    // 数组的大小是类型的一部分
    // 数组不能调整大小

    // 数组遍历
    iterateArray()

    // 二维数组
    twoDimensionalitiesArray()
}

func iterateArray() {
    a := [...]float64{67.7, 89.8, 21, 78}
    for i := 0; i < len(a); i++ { // looping from 0 to the length of the array
        fmt.Printf("%d th element of a is %.2f\n", i, a[i])
    }
    fmt.Println("calculate the summary of the array")
    sum := float64(0)
    for _, v := range a {//range returns both the index and value
        // fmt.Printf("%d the element of a is %.2f\n", i, v)
        sum += v
    }
    fmt.Println("sum of all elements of a",sum)
}

func twoDimensionalitiesArray() {
    a := [3][3]string{
        {"hello", "world"},
        {"I", "like", "golang"},
    }
    fmt.Println("the tow dimensionalities array is", a)
}
