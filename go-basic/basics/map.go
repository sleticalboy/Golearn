package basics

import "fmt"

func mapRun() {
	// 创建 map
	m := make(map[string]int)
	m["a"] = 1
	fmt.Println("map is", m)

	// 创建 map 并初始化元素
	m2 := map[string]int{
		"a": 111,
		"b": 222,
		"c": 333,
	}
	fmt.Println("map is", m2)
	fmt.Println("map.c is", m2["c"])

	// 遍历
	for k, v := range m2 {
		fmt.Printf("%s: %d\n", k, v)
	}
}
