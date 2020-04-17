package main

import (
	"flag"
	"fmt"
	"compress/gzip"
	"io/ioutil"
	"strings"
	"os"
	"encoding/binary"

	"github.com/AtomicMalloc/debugger/cli"
	"github.com/AtomicMalloc/debugger/debugger"
	"github.com/AtomicMalloc/debugger/xen"
	"github.com/AtomicMalloc/debugger/file"
)

func unzipFile(filename string) string {
	if strings.Contains(filename, ".gz") {
		file, err := os.Open(filename)
		if err != nil {
			fmt.Println(err)
		}
		r, err := gzip.NewReader(file)
		if err != nil {
			fmt.Println(err)
		}
		bytes := make([]byte, 100)
		tmpfile, err := ioutil.TempFile("", "os")
		if err != nil {
			fmt.Println(err)
		}
		for n, _ := r.Read(bytes); n != 0; n, _ = r.Read(bytes) {
			tmpfile.Write(bytes[:n])
		}
		filename = tmpfile.Name()
	}
	return filename
}

func main() {
	cmd := cli.CLI{}

	var id int
	var filename string
	flag.StringVar(&filename, "path", "", "Path to the os's binary")
	flag.IntVar(&id, "id", -1, "The domain id to connect")
	flag.Parse()

	if id == -1 {
		fmt.Println("Error: no domain id passed")
		os.Exit(1)
	} else if len(filename) == 0 {
		fmt.Println("Error: no image was passed!")
		os.Exit(1)
	}

	domainid := uint32(id)

	cntrl := &xen.Xenctrl{DomainID: domainid}
	cntrl.Init()
	cntrl.SetDebug(domainid, true)
	
	mem := &xen.Memory{Domainid: domainid, Vcpu: 0}
	mem.Init(cntrl)
	filename = unzipFile(filename)
	
	p, err := file.NewSymbolicInformation(filename, binary.LittleEndian)
	f := &file.LineInformation{Name: filename}
	err = f.Init()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	dbg := debugger.NewDebugger(mem, cntrl, f, cntrl, p)
	cmd.Init(dbg)
	

	fmt.Println("Welcome to Duster!")
	for {
		input := cmd.ReadInput()
		cmd.ProcessInput(input)
	}
}
