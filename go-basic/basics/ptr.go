package basics

import (
	"fmt"
	"unicode/utf8"
)

func ptrRun() {
	println("\nptr run")
	s := "hello"
	// 指针初始值是 nil，打印出来就是 0
	var p *string
	fmt.Printf("'%s' mem addr is 0x%x\n", s, p)
	// 指针复制
	if p == nil {
		p = &s
	}
	// 这阵解引用
	fmt.Printf("s ptr is 0x%x, value is '%s', len is %d\n", p, *p, utf8.RuneCountInString(s))
}
