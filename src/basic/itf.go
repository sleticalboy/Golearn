package basic

import "fmt"

// 定义接口
type Runnable interface {
	run()
}

type Comparable interface {
	compare(other interface{}) int
}

type DownloadTask struct {
	url string
}

// 实现接口
func (task *DownloadTask) run() {
	fmt.Printf("download task run %s\n", task.url)
}

// 实现多个接口
func (task *DownloadTask) compare(other interface{}) int {
	switch v := other.(type) {
	case Comparable:
		return task.compare(v)
	default:
		fmt.Println("wtf!!", other)
		return -1
	}
}

type UpdateTask struct {
	version string
}

func (task *UpdateTask) run() {
	fmt.Printf("update task run %s\n", task.version)
}

func itfRun() {
	println("\ninterface run!")

	// 接口默认值为 nil
	var r Runnable
	fmt.Println("raw r is", r)

	dt := &DownloadTask{
		url: "https://www.example.com",
	}
	fmt.Printf("type of task is %T\n", dt)
	dt.run()
	dt.compare("hello")

	ut := &UpdateTask{
		version: "2.7.9",
	}
	fmt.Printf("type of task is %T\n", ut)
	ut.run()

	// 接口的用法
	tasks := []Runnable{dt, ut}
	for _, t := range tasks {
		t.run()
	}
	// 接口断言
	itfAssert(dt)
	itfAssert("hello")
}

func itfAssert(itf interface{}) {
	switch v := itf.(type) {
	case Runnable:
		v.run()
	default:
		fmt.Printf("I do't known: %T\n", itf)
	}
}
