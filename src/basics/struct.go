package basics

import "fmt"

// 命名结构体
type Person struct {
	// 命名字段
	// 同类型的字段放在一行是为了紧凑
	name        string
	age, gender int
	// 匿名字段
	float64
}

// 不可导出的结构体
type project struct {
	Path string
}

func NewProject(path string) project {
	return project{Path: path}
}

func (p *project) Build(out string) {
	fmt.Printf("Build() project '%s' to '%s'\n", p.Path, out)
}

func structRun() {
	println("\nstruct run")

	var raw Person
	fmt.Println("raw struct is", raw)
	tom := Person{
		name:   "Tom",
		age:    22,
		gender: 1,
	}
	tom.float64 = 5000.0
	fmt.Println("tom is", tom)

	// 匿名结构体
	jack := struct {
		name        string
		age, gender int
	}{
		name:   "jack",
		age:    23,
		gender: 0,
	}
	fmt.Println("jack is", jack)
	fmt.Println("jack.age is", jack.age)
	jack.age = 34
	fmt.Println("jack is", jack)

	// 结构体指针
	mary := &Person{
		name:   "mary",
		age:    18,
		gender: 0,
	}
	fmt.Println("mary.age is", mary.age)

	p := NewProject("sample-go")
	p.Build(p.Path + "/build")
}
