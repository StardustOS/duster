package main

import (
	"compress/gzip"
	"encoding/binary"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"os/user"
	"strings"

	"github.com/AtomicMalloc/debugger/cli"
	"github.com/AtomicMalloc/debugger/debugger"
	"github.com/AtomicMalloc/debugger/file"
	"github.com/AtomicMalloc/debugger/xen"
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
	currentUser, err := user.Current()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// 0 means we're running as root (annoyingly go treats this has a string)
	if strings.Compare(currentUser.Uid, "0") != 0 {
		fmt.Println("Error: debugger must run as root (so it can interface with the Xen API)")
		os.Exit(1)
	}

	cmd := cli.CLI{}

	var id int
	var filename string
	flag.StringVar(&filename, "path", "", "Path to the Operating System's binary")
	flag.IntVar(&id, "id", 0, "The domain id to connect (can be found by running sudo xl list)")
	flag.Parse()

	if id == -1 {
		fmt.Println("Error: no domain id passed")
		os.Exit(1)
	} else if len(filename) == 0 {
		fmt.Println("Error: no image was passed!")
		os.Exit(1)
	}

	if _, err := os.Stat(filename); os.IsNotExist(err) {
		fmt.Printf("Error: %s does not exist\n", filename)
		os.Exit(1)
	}

	domainid := uint32(id)

	cntrl := &xen.Xenctrl{DomainID: domainid}
	err = cntrl.Init()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	err = cntrl.SetDebug(domainid, true)
	if err != nil {
		fmt.Println(err)
		fmt.Println("The above error is most likely caused by the domain id not existing")
		os.Exit(1)
	}

	if !cntrl.IsPaused() {
		fmt.Println("Error: the VM should be paused before the debugger is run!")
		fmt.Println("(To pause the VM on startup use the command: sudo xl create -p <vm name>")
		os.Exit(1)
	}

	mem := &xen.Memory{Domainid: domainid, Vcpu: 0}
	err = mem.Init(cntrl)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	filename = unzipFile(filename)

	p, err := file.NewSymbolicInformation(filename, binary.LittleEndian)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

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
