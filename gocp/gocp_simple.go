package main

import (
	"io"
	"log"
	"os"
)

var lg *log.Logger
var dst string
var src string

func init() {
	lg = log.New(os.Stdout, "copy log", log.Lshortfile)
	src = os.Args[1]
	dst = os.Args[2]
}

func main() {
	dstFile, derr := os.Create(dst)
	defer dstFile.Close()
	if derr != nil {
		lg.Fatalln(derr)
	}
	srcFile, serr := os.Open(src)
	defer srcFile.Close()
	if serr != nil {
		log.Fatalln(serr)
	}
	_, cerr := io.Copy(dstFile, srcFile)
	if cerr != nil {
		log.Fatalln(cerr)
	}
}
