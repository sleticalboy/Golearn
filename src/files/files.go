package files

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"
)

func readFiles() {
	// 获取当前目录
	dir, _ := os.Getwd()
	// 构造文件路径
	path := fmt.Sprintf("%s/go.mod", dir)
	fmt.Printf("readFiles() path is '%s'\n", path)

	// 读取文件全部内容
	if content, err := os.ReadFile(path); err == nil {
		fmt.Printf(string(content))
	} else {
		fmt.Println(err)
	}

	// 读取文件部分内容（通过 buffer 读取）
	path = fmt.Sprintf("%s/LICENSE", dir)
	fmt.Printf("readFiles() path is '%s'\n", path)

	f, err := os.Open(path)
	defer func() {
		_ = f.Close()
	}()

	if err != nil {
		fmt.Println(err)
		return
	}
	buf := make([]byte, 256)
	if info, err := f.Stat(); err == nil {
		count := info.Size() / 256
		reset := info.Size() % 256
		fmt.Printf("%d -> %d...%d\n", info.Size(), count, reset)
	}
	for {
		readBytes, err := f.Read(buf)
		if err == io.EOF {
			fmt.Println("readFiles() hit EOF!")
			break
		}
		if err != nil {
			fmt.Println(err)
			break
		}
		// buf 没有被填满，直接打印会出现脏数据
		if readBytes != len(buf) {
			fmt.Printf("%s", string(buf[0:readBytes]))
		} else {
			fmt.Println(string(buf))
		}
	}
	fmt.Println()

	// 逐行读取文件
	f, err = os.Open(path)
	defer func() {
		_ = f.Close()
	}()

	if err != nil {
		fmt.Println(err)
		return
	}
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		text := strings.TrimSpace(scanner.Text())
		if len(text) == 0 {
			continue
		}
		fmt.Println(text)
	}
}

func writeFiles() {
	// 获取当前目录
	dir, _ := os.Getwd()
	// 构造文件路径
	err := os.MkdirAll(fmt.Sprintf("%s/build", dir), 0755)
	if err != nil {
		fmt.Println(err)
		return
	}
	path := fmt.Sprintf("%s/build/test.txt", dir)
	fmt.Printf("\nwriteFiles() path is '%s'\n", path)

	// 创建文件
	f, err := os.Create(path)
	defer func() {
		_ = f.Close()
	}()

	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Printf("writeFiles() '%s'\n", f.Name())
	writeBytes, err := f.WriteString("first line\n")
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Printf("write bytes: %d\n", writeBytes)
	writeBytes, err = fmt.Fprintf(f, "write via fmt.Fprintf: %d\n", 90)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Printf("write bytes: %d\n", writeBytes)

	// 追加文件
	f, err = os.OpenFile(path, os.O_APPEND|os.O_WRONLY, 0644)
	defer func() {
		_ = f.Close()
	}()
	writeBytes, err = f.WriteString("append line to file\n")
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Printf("write bytes: %d\n", writeBytes)
}

func Main() {
	println("\nfiles Run")
	readFiles()
	writeFiles()
}
