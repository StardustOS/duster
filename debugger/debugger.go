package debugger

import "C"

import (
	"bufio"
	"fmt"
	"os"

	"github.com/AtomicMalloc/debugger/file"
	"github.com/AtomicMalloc/debugger/xen"
)

const (
	singleStep uint64 = 0x00000100
	breakInt   byte   = 0xCC
)

type Debugger struct {
	domainid     uint32
	breakpoints  map[uint64]byte
	controller   xen.Xenctrl
	memory       xen.Memory
	call         xen.XenCall
	fileHandler  file.File
	insinglestep bool
}

//Init must be run before any of the other methods
//are used
func (debugger *Debugger) Init(domainid uint32, name string) error {
	debugger.breakpoints = make(map[uint64]byte)
	debugger.domainid = domainid
	err := debugger.controller.Init()
	if err != nil {
		return err
	}
	err = debugger.memory.Init(&debugger.controller)
	if err != nil {
		return err
	}
	debugger.controller.Pause(domainid)
	err = debugger.controller.SetDebug(domainid, true)
	if err != nil {
		return err
	}
	debugger.controller.UnPause(domainid)
	err = debugger.call.Init()
	if err != nil {
		return err
	}
	debugger.fileHandler.Name = name
	err = debugger.fileHandler.Init()
	return err
}

func (debugger *Debugger) StartSingle(vcpu uint32, start bool) error {
	registers := debugger.controller.GetRegisterContext(debugger.domainid, vcpu)
	val, err := registers.GetRegister("rflags")
	if err != nil {
		return err
	}
	if start {
		val |= singleStep
		debugger.insinglestep = true
	} else {
		val &= ^singleStep
		debugger.insinglestep = false

	}
	registers.AddRegister("rflags", val)
	debugger.controller.SetRegisterContext(registers, debugger.domainid, vcpu)
	return nil
}

func (debugger *Debugger) IsPaused() bool {
	return debugger.controller.IsPaused(debugger.domainid)
}

func (debugger *Debugger) UnPause() {
	debugger.controller.UnPause(debugger.domainid)
}

func (debugger *Debugger) GetLineInformation() string {
	filename, lineNo := debugger.fileHandler.CurrentLine()
	file, _ := os.Open(filename)
	reader := bufio.NewReader(file)
	var line string
	for i := 0; i < lineNo; i += 1 {
		line, _ = reader.ReadString('\n')
	}
	return fmt.Sprintf("%s:%d - %s", filename, lineNo, line)
}

func (debugger *Debugger) Check(address uint64) {
	debugger.fileHandler.UpdateLine(address)
	fmt.Println(debugger.GetLineInformation())
}

func (debugger *Debugger) Step(vcpu uint32) uint64 {
	registers := debugger.controller.GetRegisterContext(debugger.domainid, vcpu)
	rip, _ := registers.GetRegister("rip")
	previousRIP := uint64(0)
	for !debugger.fileHandler.UpdateLine(rip) {
		for !debugger.controller.IsPaused(debugger.domainid) {
		}
		previousRIP = rip
		registers := debugger.controller.GetRegisterContext(debugger.domainid, vcpu)
		rip, _ = registers.GetRegister("rip")
		if oldByte, ok := debugger.breakpoints[rip]; ok {
			debugger.memory.Write(rip, 1, []byte{oldByte})
			rip -= 1
			registers.AddRegister("rip", rip)
		} else if _, ok = debugger.breakpoints[previousRIP]; ok {
			debugger.memory.Write(previousRIP, 1, []byte{breakInt})
		}
		debugger.controller.UnPause(debugger.domainid)
	}
	if _, ok := debugger.breakpoints[previousRIP]; ok {
		debugger.memory.Write(previousRIP, 1, []byte{breakInt})
	}
	return rip
}

func (debugger *Debugger) Continue(vcpu uint32) error {
	if debugger.controller.IsPaused(debugger.domainid) {
		if debugger.insinglestep {
			debugger.StartSingle(vcpu, false)
		}
		return debugger.controller.UnPause(debugger.domainid)
	}
	return nil
}

func (debugger *Debugger) SetBreakpoint(filename string, line int, vcpu uint32) error {
	fmt.Println("Filename", filename)
	fmt.Println("line", line)
	debugger.controller.Pause(debugger.domainid)
	address := debugger.fileHandler.GetAddress(filename, line)
	fmt.Println("address", address)
	err := debugger.memory.Map(address, debugger.domainid, 1, vcpu)
	if err != nil {
		return err
	}
	bytes, err := debugger.memory.Read(address, 1)
	debugger.breakpoints[address] = bytes[0]
	if err != nil {
		return err
	}

	debugger.memory.Write(address, 1, []byte{breakInt})
	debugger.controller.UnPause(debugger.domainid)
	return nil
}

// func (debugger *Debugger) RemoveBreakpoint(filename string, line int, vcpu uint32) error {
// 	address := debugger.fileHandler.GetAddress(filename, line)
// 	if val, ok := debugger.breakpoints[address]; ok {
// 		delete(debugger.breakpoints, address)
// 		err := memory.Write(address, 1, []byte{val})
// 		return err
// 	}
// 	return nil
// }

//Cleanup must be run before the end of the program
func (debugger *Debugger) Cleanup() {
	debugger.controller.Close()
	debugger.memory.Close()
}
