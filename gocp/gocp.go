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
var l *log.Logger

func init() {
	l = log.New(os.Stdout, "copy log", log.Lshortfile)
	src = os.Args[1]
	dst = os.Args[2]

	cwd, gtWdErr := os.Getwd()
	if gtWdErr != nil {
		l.Fatalln(gtWdErr)
	}

	if !filepath.IsAbs(src) {
		src = filepath.Join(cwd, src)
	}
	if !filepath.IsAbs(dst) {
		dst = filepath.Join(cwd, src)
	}
	// if os.IsNotExist(src) {
	// 	l.Fatalln("The source file or dir is not exist.")
	// }
	// if os.IsNotExist(dst) {
	// 	l.Fatalln("The dst file or dir is not exist.")
	// }
	dstFI, dstfiErr := os.Stat(dst)
	if dstfiErr != nil {
		l.Fatalln(dstfiErr)
	}
	if !dstFI.Mode().IsDir() {
		l.Fatalln("The dst must be a dir.")
	}
}

func main() {
	var WkWg sync.WaitGroup
	srcList, _ := filepath.Glob(src)
	if srcList == nil {
		l.Fatalln("The source file or dir is not exist.")
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
	return nil
}

func copyFile(dstName string, srcName string) {
	l.Println(dstName)
	l.Println(srcName)
	defer cpWg.Done()
	dstFile, cErr := os.Create(dstName)
	defer dstFile.Close()
	if cErr != nil {
		l.Fatalln(cErr)
	}
	srcFile, oErr := os.Open(srcName)
	defer srcFile.Close()
	if oErr != nil {
		l.Fatalln(oErr)
	}
	_, cpErr := io.Copy(dstFile, srcFile)
	if cpErr != nil {
		l.Fatalln(cpErr)
	}
}
