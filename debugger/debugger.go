package debugger

import "C"

import (
	"fmt"

	"github.com/AtomicMalloc/debugger/xen"
)

const (
	singleStep uint64 = 0x00000100
)

type Debugger struct {
	domainid   uint32
	controller xen.Xenctrl
	memory     xen.Memory
	call       xen.XenCall
}

//Init must be run before any of the other methods
//are used
func (debugger *Debugger) Init(domainid uint32) error {
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
	return err
}

func (debugger *Debugger) StartSingle(vcpu uint32, start bool) error {
	registers := debugger.controller.GetRegisterContext(debugger.domainid, vcpu)
	val, err := registers.GetRegister("rflags")
	if err != nil {
		return err
	}
	//fmt.Println("val read from registers", val)
	if start {
		val |= singleStep
	} else {
		val &= ^singleStep
		//	fmt.Println("val written to registers", val)
	}
	//fmt.Println("About to write", val)
	registers.AddRegister("rflags", val)
	debugger.controller.SetRegisterContext(registers, debugger.domainid, vcpu)
	return nil
}

func (debugger *Debugger) IsPaused() bool {
	return debugger.controller.IsPaused(debugger.domainid)
}

func (debugger *Debugger) UnPause() {
	registers := debugger.controller.GetRegisterContext(debugger.domainid, 0)
	val, _ := registers.GetRegister("rip")
	fmt.Println("RIP", val)
	debugger.controller.UnPause(debugger.domainid)
}

func (debugger *Debugger) Step(vcpu uint32) uint64 {
	registers := debugger.controller.GetRegisterContext(debugger.domainid, 0)
	val, _ := registers.GetRegister("rip")
	debugger.controller.UnPause(debugger.domainid)
	return val
	//debugger.controller.UnPause(debugger.domainid)
	//debugger.controller.Pause(debugger.domainid)
	//debugger.call.HyperCall(debugger.controller, xen.PauseCPU, debugger.domainid, vcpu)
	//debugger.StartSingle(vcpu, true)
	//debugger.call.HyperCall(debugger.controller, xen.UnPauseCPU, debugger.domainid, vcpu)
	//debugger.controller.UnPause(debugger.domainid)

	// for !debugger.controller.IsPaused(debugger.domainid) {
	// 	fmt.Println("Wating for pause")
	// }
	// err := debugger.StartSingle(vcpu, false)

	// if err != nil {
	// 	fmt.Println("err", err)
	// }
	// reg := debugger.controller.GetRegisterContext(debugger.domainid, vcpu)

	// rip, err := reg.GetRegister("rip")
	// if err != nil {
	// 	fmt.Println(err)
	// }
	// fmt.Println("RIP ", rip)
	// return rip
}

//Cleanup must be run before the end of the program
func (debugger *Debugger) Cleanup() {
	debugger.controller.Close()
	debugger.memory.Close()
}
