package main

import (
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"
	"sync"
)

var lg *log.Logger
var wg sync.WaitGroup

func init() {
	lg = log.New(os.Stdout, "log", log.Lshortfile)
}

func main() {
	src := strings.TrimSuffix(os.Args[0], "/")
	lg.Println(src)
	dst := strings.TrimSuffix(os.Args[1], "/")
	lg.Println(dst)

	dstFi, dstFiErr := os.Lstat(dst)
	if dstFiErr != nil {
		lg.Println("Warning: the dst may is not exist.")
	}

	srcList, _ := filepath.Glob(src)
	srcListLen := len(srcList)
	if srcListLen == 0 {
		lg.Fatalln("You must give a valid src.")
	}

	fi, fiErr := os.Lstat(srcList[0])
	if srcListLen == 1 && dstFi == nil {
		if fiErr != nil {
			lg.Fatalln(fiErr)
		}
		if fi.IsDir() {
			os.MkdirAll(dst, fi.Mode().Perm())
		} else {
			copyFile(dst, srcList[0])
			os.Exit(0)
		}
	}
	if srcListLen > 1 && dstFi == nil {
		os.MkdirAll(dst, fi.Mode().Perm())
	}
	if srcListLen == 1 && !dstFi.IsDir() {
		if fi.IsDir() {
			lg.Fatalln("Can't copy a dir to a file.")
		} else {
			lg.Println("The dst file will be overwriten.")
			copyFile(dst, srcList[0])
			os.Exit(0)
		}
	}
	if srcListLen > 1 && !dstFi.IsDir() {
		lg.Fatalln("Can't copy multi-file src to a dst file.")
	}

	wkFn := func(path string, info os.FileInfo, err error) error {
		dirName := filepath.Dir(src)
		fileName := strings.TrimPrefix(path, dirName)
		dstFileName := filepath.Join(dst, fileName)
		pathInfo, pathErr := os.Lstat(path)
		if pathErr != nil {
			lg.Fatalln(pathErr)
		}
		if pathInfo.IsDir() {
			os.MkdirAll(dstFileName, pathInfo.Mode().Perm())
		} else {
			wg.Add(1)
			go func() {
				copyFile(dstFileName, path)
				wg.Done()
			}()
		}
		return nil
	}

	for _, i := range srcList {
		wg.Add(1)
		go func(src string) {
			filepath.Walk(src, wkFn)
			wg.Done()
		}(i)
	}
	wg.Wait()
	lg.Println("Done.")
}

// func wkFn(path string, info os.FileInfo, err error) error {
// 	// var fileName string
// 	// if path == src {
// 	// fileName := filepath.Base(path)
// 	// } else {
// 	// fileName = strings.TrimPrefix(path, src)
// 	// }
// 	// dstFileName := filepath.Join(dst, fileName)
// 	// lg.Println(src)
// 	// lg.Println(path)
// 	// lg.Println(dst)
// 	// lg.Println(fileName)
// 	// lg.Println(dstFileName)
// 	dirName := filepath.Dir(src)
// 	fileName := strings.TrimPrefix(path, dirName)
// 	dstFileName := filepath.Join(dst, fileName)
// 	pathInfo, pathErr := os.Lstat(path)
// 	if pathErr != nil {
// 		lg.Fatalln(pathErr)
// 	}
// 	if pathInfo.IsDir() {
// 		os.MkdirAll(dstFileName, pathInfo.Mode().Perm())
// 	} else {
// 		wg.Add(1)
// 		go copyFile(dstFileName, path)
// 	}
// 	return nil
// }

func copyFile(dstName string, srcName string) {
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
