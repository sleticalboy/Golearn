package basics

import (
	"fmt"
	"reflect"
)

type ISay interface {
	Say()
}

type Nested struct {
	Cc string
}

func (n Nested) Say() {
	fmt.Println("sample.Say() run...")
}

func (n Nested) String() string {
	return fmt.Sprintf("Nestedt{Cc: %s}", n.Cc)
}

func (n Nested) What() (a, b, c int, reason string) {
	return 1, 2, 3, "nil"
}

type foo struct {
	a string
	b int
	n Nested
}

type sample struct {
	name    string
	age     int
	salary  float32
	subject []string
	fo      foo
}

func indent(depth int) {
	for i := 0; i < depth; i++ {
		fmt.Printf("  ")
	}
}

func printValue(vf reflect.Value, depth int) {
	switch vf.Type().Kind() {
	case reflect.String:
		fmt.Printf(" -> %s\n", vf.String())
	case reflect.Int:
		fmt.Printf(" -> %d\n", vf.Int())
	case reflect.Float32, reflect.Float64:
		fmt.Printf(" -> %f\n", vf.Float())
	case reflect.Slice, reflect.Array:
		length := vf.Len()
		fmt.Printf(" -> ")
		slice := vf.Slice(0, length)
		for i := 0; i < length; i++ {
			fmt.Printf("[%d]: %v", i, slice.Index(i))
			if i != length-1 {
				fmt.Printf(", ")
			} else {
				fmt.Println()
			}
		}
	case reflect.Struct:
		fmt.Printf(" <struct>:\n")
		refWithDepth(vf, depth+1)
	default:
		fmt.Printf(" -> %v\n", vf)
	}
}

func refInternal(t reflect.Type, v reflect.Value, depth int) {
	switch k := t.Kind(); k {
	case reflect.Struct:
		n := 0
		if n := t.NumField(); n > 0 {
			indent(depth)
			fmt.Printf("type: %v, methods: %d, fields: %d\n", t, t.NumMethod(), n)
			indent(depth)
			fmt.Printf("struct fields >>>>>>\n")
			for i := 0; i < n; i++ {
				indent(depth)
				tf := t.Field(i)
				fmt.Printf("%s.%s[%v]", t.Name(), tf.Name, tf.Type)
				printValue(v.Field(i), depth)
			}
		}
		// 值函数才能反射出来？
		if n = t.NumMethod(); n > 0 {
			indent(depth)
			fmt.Printf("struct methods >>>>>>\n")
			for i := 0; i < n; i++ {
				indent(depth)
				tm := t.Method(i)
				if no := tm.Type.NumOut(); no > 0 {
					fmt.Printf("%s.%s()", t.Name(), tm.Name)
					for j := 0; j < no; j++ {
						fmt.Printf(" %s", tm.Type.Out(j).String())
					}
					fmt.Printf("\n")
				} else {
					fmt.Printf("%s.%s()\n", t.Name(), tm.Name)
				}
			}
		}
	case reflect.Func:
		fmt.Printf("sig(%s)\n", t.String())
		fmt.Printf("t.NumIn(): %v\n", t.NumIn())
		fmt.Printf("t.NumOut(): %v\n", t.NumOut())
		fmt.Printf("t.In(0): %v\n", t.In(0))
		fmt.Printf("t.Out(0): %v\n", t.Out(0))
	case reflect.Interface:
		fmt.Printf("%s is a interface\n", t.Name())
	default:
		fmt.Printf("%s[%v]", t.Name(), v.Type())
		printValue(v, depth)
	}
}

func ref(o any) {
	refWithDepth(o, 0)
}

func refWithDepth(o any, depth int) {
	if o == nil {
		fmt.Println("nil")
		return
	}
	switch f := o.(type) {
	case reflect.Value:
		refInternal(f.Type(), f, depth)
	default:
		refInternal(reflect.TypeOf(o), reflect.ValueOf(o), depth)
	}
}

func reflectsRun() {
	println("\nreflects run..")

	s := sample{
		name:    "tom",
		age:     24,
		salary:  9000.0,
		subject: []string{"java", "go", "c++"},
		fo: foo{
			a: "aaa",
			b: 111,
			n: Nested{
				Cc: "nested struct depth 2",
			},
		},
	}
	ref(s)

	ref("hello")
	ref(float32(9))
	ref(float64(88))
	ref(nil)
	ref(true)
	ref(2)
	ref(int32(23))
	ref(func(a int) string {
		return "sdkkd"
	})
	ref([]string{"1", "2"})
}
