package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"
	"sync"
)

// var dst string
// var src string
var l *log.Logger
var WkWg sync.WaitGroup

type paths struct {
	name    string
	dirPath bool
	exist   bool
}

// var src paths
var dst paths

func (p *paths) initPath(s string) error {
	p.name = s
	fi, fiErr := os.Lstat(s)
	if fiErr != nil {
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

func getAbsPath(s string) string {
	cwd, gtWdErr := os.Getwd()
	if gtWdErr != nil {
		l.Fatalln(gtWdErr)
	}
	if !filepath.IsAbs(s) {
		s = filepath.Join(cwd, s)
	}
	return s
}

func init() {
	l = log.New(os.Stdout, "log", log.Lshortfile)
	dstpath := getAbsPath(os.Args[2])
	dst.initPath(dstpath)
}

func check(srcP *paths, dstP *paths) {
	// if !srcP.exist {
	// 	l.Fatalln("The src is not exist.")
	// }
	if srcP.dirPath {
		if !dstP.exist {
			l.Fatalln("The dst is not exist.")
		}
		if !dstP.dirPath {
			l.Fatalln("The dst must be a dir.")
		}
	}
	if dstP.dirPath {
		if !dstP.exist {
			l.Fatalln("The dst is not exist.")
		}
	}

}

func main() {
	srcpath := getAbsPath(os.Args[1])
	srcpath = strings.TrimSuffix(srcpath, "/")
	srcList, _ := filepath.Glob(srcpath)
	if srcList == nil {
		l.Fatalln("The source file or dir is not exist.")
	}
	for _, i := range srcList {
		srcI := new(paths)
		srcI.initPath(i)
		check(srcI, dst)
		WkWg.Add(1)
		go func(s string) {
			wkerr := filepath.Walk(s, wkFn)
			if wkerr != nil {
				l.Fatalln(wkerr)
			}
			WkWg.Done()
		}(i)
	}
	WkWg.Wait()
	fmt.Println("Done.")
}

func wkFn(path string, info os.FileInfo, err error) error {
	baseName := filepath.Base(path)
	//srcName := filepath.Join(path, info.Name())
	srcName := path
	dstName := filepath.Join(dst, baseName)
	WkWg.Add(1)
	go copyFile(dstName, srcName)
	return nil
}

func copyFile(dstName string, srcName string) {
	defer WkWg.Done()
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
