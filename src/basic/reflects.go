package basic

import (
	"fmt"
	"reflect"
)

type nested struct {
	cc string
}

type foo struct {
	a string
	b int
	c nested
}

type sample struct {
	name    string
	age     int
	sallery float32
	subject []string
	fo      foo
}

func (s *sample) Say() {}

func ident(depth int) {
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
		len := vf.Len()
		fmt.Printf(" -> ")
		slice := vf.Slice(0, len)
		for i := 0; i < len; i++ {
			fmt.Printf("[%d]: %s", i, slice.Index(i).String())
			if i != len-1 {
				fmt.Printf(", ")
			} else {
				fmt.Println()
			}
		}
	case reflect.Struct:
		fmt.Printf(" <struct>:\n")
		ref(vf, depth+1)
	default:
		fmt.Printf(" -> %v\n", vf)
	}
}

func refInternal(t reflect.Type, v reflect.Value, depth int) {
	switch t.Kind() {
	case reflect.Struct:
		n := t.NumField()
		// ident(depth)
		// fmt.Printf("type: %v, methods: %d, fields: %d\n", t, t.NumMethod(), n)
		for i := 0; i < n; i++ {
			ident(depth)
			tf := t.Field(i)
			fmt.Printf("%s.%s[%v]", t.Name(), tf.Name, tf.Type)
			printValue(v.Field(i), depth)
		}
	case reflect.Func:
		fmt.Printf("func name is '%v'\n", v)
		fmt.Printf("t.NumIn(): %v\n", t.NumIn())
		fmt.Printf("t.NumOut(): %v\n", t.NumOut())
		fmt.Printf("t.In(0): %v\n", t.In(0))
		fmt.Printf("t.Out(0): %v\n", t.Out(0))
	default:
		fmt.Printf("%s[%v]", t.Name(), v.Type())
		printValue(v, depth)
	}
}

func ref(o any, depth int) {
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

func fooo(i int) string {
	return "hello foo"
}

func reflectsRun() {
	println("\nreflects run..")

	s := sample{
		name:    "tom",
		age:     24,
		sallery: 9000.0,
		subject: []string{"java", "go", "c++"},
		fo: foo{
			a: "aaa",
			b: 111,
			c: nested{
				cc: "nested struct depth 2",
			},
		},
	}
	ref(s, 0)

	ref("hello", 0)
	ref(float32(9), 0)
	ref(float64(88), 0)
	ref(nil, 0)
	ref(true, 0)
	ref(2, 0)
	ref(int32(23), 0)

	// ano := func(a int) string {
	// 	return "sdkkd"
	// }
	ref(fooo, 0)
}
