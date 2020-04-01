package main

import "github.com/AtomicMalloc/debugger/file"

func main() {
	f := file.File{Name: "./file/a.out"}
	f.Init()
}
