package main

import (
	"fmt"
	// "os"
	"path/filepath"
)

func main() {
	fmt.Println(filepath.IsAbs("/mygo/src"))
}
