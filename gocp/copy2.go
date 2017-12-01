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

var lg *log.Logger
var WkWg sync.WaitGroup

type paths struct {
	name    string
	dirPath bool
	exist   bool
}

func (p *paths) initP(s string) error {
	p.name = s
	fi, fiErr := os.Lstat(s)
	if fiErr != nil {
		lg.Fatalln(fiErr)
	}
	if fi == nil {
		p.exist = false
		if strings.HasSuffix(s, "/") {
			p.dirPath = true
			return nil
		}
		p.dirPath = false
		return nil
	}
	p.exist = true
	if fi.IsDir() {
		p.dirPath = true
		return nil
	}
	p.dirPath = false
	return nil
}

var dP paths

func check(d string) bool {
	// sP := paths{}
	// dP = paths{}
	// sP.initP(s)
	dP.initP(d)
	if dP.dirPath && dP.exist {
		return true
	}
	if (!dP.dirPath) && (!dP.exist) {
		return true
	}
	return false
}

func init() {
	src = os.Args[1]
	dst = os.Args[2]
	lg = log.New(os.Stdout, "log", log.Lshortfile)
}

func main() {

	// check(dst)

	srcList, srcGbErr := filepath.Glob(src)
	if srcGbErr != nil {
		lg.Fatalln(srcGbErr)
	}
	srcListLen := len(srcList)
	if srcList == nil || srcListLen == 0 {
		lg.Fatalln("Please give a valid src.")
	}

	for _, s := range srcList {

		sFi, sErr := os.Lstat(s)
		if sErr != nil {
			lg.Fatalln(sErr)
		}

		srcBase := filepath.Base(s)

		if !sFi.IsDir() {
			if dP.dirPath {
				dstName := filepath.Join(dst, srcBase)
			} else {
				dstName := dst
			}
			copyFile(dstName, srcName)
			continue
		}

		dst := dst + srcBase

		dstMkErr := os.MkdirAll(dst, sFi.Mode().Perm())

		if dstMkErr != nil {
			lg.Fatalln(dstMkErr)
		}

		go func(s string) {
			filepath.Walk(s, wkFn)

		}(s)
	}

}

func wkFn(path string, info os.FileInfo, err error) error {
	var fileName string
	if path == src {
		fileName = filepath.Base(path)
	} else {
		fileName = strings.TrimPrefix(path, src)
	}
	dstFileName := filepath.Join(dst, fileName)
	lg.Println(src)
	lg.Println(path)
	lg.Println(dst)
	lg.Println(fileName)
	lg.Println(dstFileName)
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
	lg.Println("in copyfile func")
	lg.Println(dstName)
	lg.Println(srcName)
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
	lg.Println("copyfile end.")
}
