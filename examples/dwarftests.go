package main

import "github.com/AtomicMalloc/debugger/file"

func main() {
	k := file.File{Name: "file/testfiles/variable_data"}
	k.Init()
	k.Tag()
}
