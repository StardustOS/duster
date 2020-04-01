package main

import (
	"fmt"
	"time"

	"github.com/AtomicMalloc/debugger/xen"
)

func main() {
	call := xen.XenCall{}
	call.Init()
	domain := xen.Xenctrl{}
	domain.Init()

	for i := 0; i < 5; i++ {
		domain.Pause(55)
		call.HyperCall(domain, xen.PauseCPU, 55, 0)
		domain.UnPause(55)
		isPaused := domain.IsPaused(55)
		fmt.Println("The system is current is paused: ", isPaused)
		time.Sleep(time.Second * 5)
		domain.Pause(55)
		call.HyperCall(domain, xen.UnPauseCPU, 55, 0)
		domain.UnPause(55)
		isPaused = domain.IsPaused(55)
		fmt.Println("(WE HAVE UNPAUSED) The system is current is paused: ", isPaused)
		time.Sleep(time.Second * 5)
	}
}
