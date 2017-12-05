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
	// log.Lshortfile is used to show line number for debug.
	lg = log.New(os.Stdout, "log", log.Lshortfile)
}

func main() {
	// src := strings.TrimSuffix(os.Args[1], "/")
	// dst := strings.TrimSuffix(os.Args[2], "/")

	// // If src is a glob pattern string, to get the matched sources.
	// srcList, _ := filepath.Glob(src)
	// srcListLen := len(srcList)
	// if srcListLen == 0 {
	// 	lg.Fatalln("You must give a valid src.") // If none source is matched, log and exit.
	// }

	argLen := len(os.Args)
	if argLen < 3 {
		lg.Fatalln("Please give the SRC string and DST string.")
	}
	srcList := os.Args[1 : argLen-1]
	srcListLen := len(srcList)
	dst := os.Args[argLen-1]

	// Try to stat the dst file.
	dstFi, dstFiErr := os.Lstat(dst)
	if dstFiErr != nil {
		lg.Println("Warning: the DST may is not exist.")
	}

	// Try to stat the first source for the nexting usage.
	fi, fiErr := os.Lstat(srcList[0])
	if fiErr != nil {
		lg.Fatalln(fiErr)
	}

	// If the dst is not exist.
	if srcListLen == 1 && dstFi == nil {
		if fi.IsDir() {
			os.MkdirAll(dst, fi.Mode().Perm()) // If the src is a directory, make the dst directory. DIR to NEW_DIR.
		} else {
			copyFile(dst, srcList[0]) // If the src is a file, just copy it to the dst. FILE to NEW_FILE.
			lg.Println("Done.")
			os.Exit(0)
		}
	}
	if srcListLen > 1 && dstFi == nil {
		os.MkdirAll(dst, fi.Mode().Perm()) // If there are many sources, make the dst directory. [DIRs, FILEs, ...] to NEW_DIR.
	}

	// If the dst is exist.
	if dstFi != nil {
		if srcListLen == 1 && !dstFi.IsDir() {
			if fi.IsDir() {
				lg.Fatalln("Can't copy a dir to a file.") // ERROR: DIR to FILE.
			} else {
				lg.Println("Warning: The dst file will be overwriten.") // WARNING: FILE to EXIST_FILE.
				copyFile(dst, srcList[0])
				lg.Println("Done.")
				os.Exit(0)
			}
		}
		if srcListLen > 1 && !dstFi.IsDir() {
			lg.Fatalln("Can't copy multi-file src to a dst file.") // ERROR: [DIRs, FILEs, ...] to EXIST_FILE.
		}
	}

	//  **Following: [FILEs, DIRs, ...] to DIR.**

	// Difine the walk function which is used to in filepath.Walk function.
	var src string
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

	// Iteration the sources list.
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

// Do the copy from a source to a destination.
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
