package files

import (
	"fmt"
	"testing"
)

func TestParseDir(t *testing.T) {
	fileMap, err := ParseDir("/home/binlee/Downloads/xyt")
	if err != nil {
		fmt.Printf("TestParseDir() failed: %v", err)
		return
	}
	fmt.Printf("TestParseDir() file maps: %v", fileMap)
}
