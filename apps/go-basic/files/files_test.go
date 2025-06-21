package files

import (
	"fmt"
	"testing"
)

func TestReadFile(t *testing.T) {
	readFiles()
}

func TestWriteFile(t *testing.T) {
	writeFiles()
}

func TestParseDir(t *testing.T) {
	fileMap, err := parseDir("/home/binlee/Downloads/xyt")
	if err != nil {
		fmt.Printf("TestParseDir() failed: %v", err)
		return
	}
	fmt.Printf("TestParseDir() file maps: %v", fileMap)
}
