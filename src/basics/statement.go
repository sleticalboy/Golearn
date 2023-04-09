package basics

import "fmt"

func statementRun() {
	println("\nstatement of golang")
	var num = 20
	if num = 15; num > 10 {
		fmt.Println("num is greater than 10")
	} else {
		fmt.Println("num is less than 10")
	}

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
			case 4, 6:
				fmt.Println("more than 2")
				// 将执行下一个 case 语句
				fallthrough
			default:
				fmt.Println("Nothing")
			}
			continue
		}
		fmt.Printf("i=%d ", i)
	}
}
