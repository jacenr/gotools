package main

import (
	"bytes"
	"flag"
	"io"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"time"
)

var lg *log.Logger
var wg sync.WaitGroup
var name string
var size string
var modifytime string
var dirs string
var strCh chan string

func init() {
	lg = log.New(os.Stderr, "gofind ", log.Lshortfile)
	flag.StringVar(&name, "n", "", "Optional. Search files by name with \"full\", \"sub\", \"reg\" sub-option.")
	flag.StringVar(&size, "s", "", "Optional. Search files by size with \"=\", \"<\", \"<=\", \">\", \">=\" sub-option.")
	flag.StringVar(&modifytime, "m", "", "Optional. Search files by modify time with \"=\", \"<\", \"<=\", \">\", \">=\" sub-option.") // eg: >=,20171206114930
	flag.StringVar(&dirs, "p", ".", "Optional. The search paths. Separated by comma.")
}

func main() {
	flag.Parse()
	strCh = make(chan string)
	dirsList := strings.Split(dirs, ",")
	for _, dir := range dirsList {
		wg.Add(1)
		go func(dir string) {
			walkErr := filepath.Walk(dir, walkFn)
			if walkErr != nil {
				lg.Fatalln(walkErr)
			}
			wg.Done()
		}(dir)
	}
	go func() {
		wg.Wait()
		close(strCh)
	}()
	newlineBuff := bytes.NewBuffer([]byte("\n"))
	for file := range strCh {
		//		fmt.Println(file)
		fileByte := []byte(file)
		buf := bytes.NewBuffer(fileByte)
		io.Copy(os.Stdout, buf)
		io.Copy(os.Stdout, newlineBuff)
	}
}

func walkFn(path string, info os.FileInfo, err error) error {
	wg.Add(1)
	go func(path string, info os.FileInfo, err error) {
		if err != nil {
			lg.Fatalln(err)
		}
		if byName(info) && bySize(info) && byTime(info) {
			strCh <- path
		}
		wg.Done()
	}(path, info, err)
	return nil
}

func byName(info os.FileInfo) bool {
	if name == "" {
		return true
	}
	nameOpt := strings.Split(name, ",")
	if len(nameOpt) != 2 {
		return false
	}
	fileName := info.Name()
	switch nameOpt[0] {
	case "full":
		return fileName == nameOpt[1]
	case "sub":
		return strings.Contains(fileName, nameOpt[1])
	case "reg":
		matched, matchErr := regexp.MatchString(nameOpt[1], fileName)
		if matchErr != nil {
			return false
		}
		return matched
	}
	return false
}

func bySize(info os.FileInfo) bool {
	if size == "" {
		return true
	}
	sizeOpt := strings.Split(size, ",")
	if len(sizeOpt) != 2 {
		return false
	}
	sizeN, siezErr := strconv.ParseUint(sizeOpt[1], 0, 64)
	if siezErr != nil {
		return false
	}
	fileSize := uint64(info.Size())
	switch sizeOpt[0] {
	case "=":
		return fileSize == sizeN
	case ">":
		return fileSize > sizeN
	case ">=":
		return fileSize >= sizeN
	case "<":
		return fileSize < sizeN
	case "<=":
		return fileSize <= sizeN
	}
	return false
}

func byTime(info os.FileInfo) bool {
	if modifytime == "" {
		return true
	}
	timeOpt := strings.Split(modifytime, ",")
	if len(timeOpt) != 2 {
		return false
	}
	timeForm := "20060102150405"
	timeS, timeErr := time.Parse(timeForm, timeOpt[1])
	if timeErr != nil {
		return false
	}
	fileTime := info.ModTime()
	switch timeOpt[0] {
	case "=":
		return timeS.Equal(fileTime)
	case ">":
		return timeS.After(fileTime)
	case ">=":
		return timeS.After(fileTime) || timeS.Equal(fileTime)
	case "<":
		return timeS.Before(fileTime)
	case "<=":
		return timeS.Before(fileTime) || timeS.Equal(fileTime)
	}
	return false
}
