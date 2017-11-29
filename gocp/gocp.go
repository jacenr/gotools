package main

import (
	"fmt"
	// "os"
	"path/filepath"
)

func main() {
	fmt.Println(filepath.IsAbs("/test/mygo/src"))
	absPath, _ := filepath.Abs("/mygo/src")
	fmt.Println(absPath)
}
