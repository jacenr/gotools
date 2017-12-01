package main

import (
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"
	"sync"
)

var src string
var dst string
var WkWg sync.WaitGroup
var lg *log.Logger

func init() {
	lg := log.New(os.Stdout, "log", log.Lshortfile)
}

func main() {
	src := os.Args[1]
	dst := os.Args[2]

	//......
	srcList, gbErr := filepath.Glob(src)
	if gbErr != nil {
		lg.Println("Please give a valid src dir.")
		lg.Fatalln(gbErr)
	}

	//......
	dstStatus := struct {
		name  string
		isDir bool
		exist bool
	}{name: dst}
	dstFi, dstFiErr := os.Lstat(dst)
	if dstFi == nil {
		dstStatus.exist = false
	} else {
		dstStatus.exist = true
		dstStatus.isDir = dstFi.IsDir()
	}
	if !dstStatus.exist {
		if strings.HasSuffix(dst, "/") {
			dstStatus.isDir = true
		} else {
			dstStatus.isDir = false
		}
	}

	//......
	srcListlen := len(srcList)
	if dstStatus.isDir {
		if !dstStatus.exist {
			lg.Fatalln("The dst dir must be exist.")
		}
	} else {
		if srcListlen != 1 {
			lg.Fatalln("The number of src != number of dst.")
		} else {
			srcFi, srcFiErr := os.Lstat(srcList[0])
			if srcFiErr != nil {
				lg.Fatalln(srcFiErr)
			}
			if srcFi.IsDir() {
				lg.Fatalln("The src must be a file.")
			}
		}
	}

	for _, srcName := range srcList {
		WkWg.Add(1)
		go func(s string) {
			filepath.Walk(s, wkFn)
			WkWg.Done()
		}(srcName)
	}
	WkWg.Wait()
}

func wkFn(path string, info os.FileInfo, err error) error {
	fileName := strings.TrimPrefix(path, src)
	dstFileName := filepath.Join(dst, fileName)
	pathInfo, pathErr := os.Lstat(path)
	if pathErr != nil {
		lg.Fatalln(pathErr)
	}
	if pathInfo.IsDir() {
		os.MkdirAll(dstFileName, pathInfo.Mode().Perm())
	} else {
		WkWg.Add(1)
		go copyFile(dstFileName, path)
	}
	return nil
}

func copyFile(dstName string, srcName string) {
	defer WkWg.Done()
	dstFile, cErr := os.Create(dstName)
	defer dstFile.Close()
	if cErr != nil {
		lg.Fatalln(cErr)
	}
	srcFile, oErr := os.Open(srcName)
	defer srcFile.Close()
	if oErr != nil {
		lg.Fatalln(oErr)
	}
	_, cpErr := io.Copy(dstFile, srcFile)
	if cpErr != nil {
		lg.Fatalln(cpErr)
	}
}
