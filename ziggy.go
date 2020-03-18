package main

import (
	"flag"
	"fmt"

	"github.com/AtomicMalloc/debugger/cli"
)

func main() {
	cmd := cli.CLI{}

	var domainid int
	var filename string
	flag.StringVar(&filename, "path", "FUCK", "Path to the os's binary")
	flag.IntVar(&domainid, "id", -1, "The domain id to connect")
	flag.Parse()
	fmt.Println(filename)

	err := cmd.Init(uint32(domainid), filename)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println("Welcome to the ziggy debugger!")
	for {
		input := cmd.ReadInput()
		cmd.ProcessInput(input)
	}
}
