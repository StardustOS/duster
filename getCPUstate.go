package main

import (
	"fmt"

	"github.com/AtomicMalloc/debugger/xen"
)

func main() {
	domain := xen.Xenctrl{}
	domain.Init()
	regs := domain.GetRegisterContext(12, 0)
	rbx, _ := regs.GetRegister("rbx")
	fmt.Println("rbx ", rbx)
	r, err := regs.GetRegister("rflags")
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("rflags ", r)
	r |= 0x00000100
	regs.AddRegister("rflags", r)
	domain.SetRegisterContext(regs, 12, 0)
	// domain.Pause(7)
	// time.Sleep(time.Second * 30)
	// domain.UnPause(7)
}
