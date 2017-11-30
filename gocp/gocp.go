package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"sync"
)

var dst string
var src string
var cpWg sync.WaitGroup

func init() {
	src = os.Args[1]
	dst = os.Args[2]

	cwd, gtWdErr := os.Getwd()
	if gtWdErr != nil {
		log.Fatalln(gtWdErr)
	}

	if !filepath.IsAbs(src) {
		src = cwd + src
	}
	if !filepath.IsAbs(dst) {
		dst = cwd + dst
	}
	if os.IsNotExist(src) {
		log.Fatalln("The source file or dir is not exist.")
	}
	if os.IsNotExist(dst) {
		log.Fatalln("The dst file or dir is not exist.")
	}
	dstFI, dstfiErr := os.Stat(dst)
	if dstfiErr != nil {
		log.Fatalln(dstfiErr)
	}
	if !dstFI.Mode().IsDir() {
		log.Fatalln("The dst must be a dir.")
	}
}

func main() {
	var WkWg sync.WaitGroup
	srcList, _ := filepath.Glob(src)
	if srcList == nil {
		log.Fatalln("The source file or dir is not exist.")
	}

	for _, i := range srcList {
		WkWg.Add(1)
		go func(s string) {
			filepath.Walk(s, wkFn)
			WkWg.Done()
		}(i)
	}
	WkWg.Add(1)
	go func() {
		cpWg.Wait()
		fmt.Println("Done.")
		WkWg.Done()
	}()
	WkWg.Wait()
}

func wkFn(path string, info os.FileInfo, err error) error {
	baseName := filepath.Base(path)
	dirName := filepath.Dir(path)
	srcName := filepath.Join(path, info.Name())
	dstName := filepath.Join(baseName, dirName)
	cpWg.Add(1)
	go copyFile(dstName, srcName)
}

func copyFile(dstName string, srcName string) {
	defer cpWg.Done()
	dstFile, cErr := os.Create(dstName)
	defer dstFile.Close()
	if cErr != nil {
		log.Fatalln(cErr)
	}
	srcFile, oErr := os.Open(srcName)
	defer srcFile.Close()
	if oErr != nil {
		log.Fatalln(oErr)
	}
	_, cpErr := io.Copy(dstFile, srcFile)
	if cpErr != nil {
		log.Fatalln(cpErr)
	}
}
