package basics

import "fmt"

func sliceRun() {
	println("\nslice run...")
	// 创建切片：类型，长度，容量
	// a := make([]int, 2, 2)
	s := []string{"a", "b", "c"}
	fmt.Printf("len of s is %d, cap is %d\n", len(s), cap(s))

	// 追加切片元素：执行 append 之后长度加 1，如果长度超出容量，则新容量变为原容量的 2 倍
	s = append(s, "ddd")
	fmt.Printf("len of s is %d, cap is %d\n", len(s), cap(s))

	// 追加切片
	s2 := []string{"e", "f", "g"}
	s = append(s, s2...)
	fmt.Printf("len of s is %d, cap is %d\n", len(s), cap(s))
}
