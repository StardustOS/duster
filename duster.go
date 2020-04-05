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
	flag.StringVar(&filename, "path", "que", "Path to the os's binary")
	flag.IntVar(&domainid, "id", -1, "The domain id to connect")
	flag.Parse()

	err := cmd.Init(uint32(domainid), filename)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println("Welcome to Duster!")
	for {
		input := cmd.ReadInput()
		cmd.ProcessInput(input)
	}
}
