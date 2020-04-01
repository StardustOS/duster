package main

import (
	"time"

	"github.com/AtomicMalloc/debugger/xen"
)

func main() {
	domain := xen.Xenctrl{}
	domain.Init()
	call := xen.XenCall{}
	call.Init()
	domain.Pause(7)
	call.HyperCall(domain, xen.PauseCPU, 7, 0)
	domain.UnPause(7)
	time.Sleep(time.Second * 20)
	domain.Pause(7)
	call.HyperCall(domain, xen.UnPauseCPU, 7, 0)
	domain.UnPause(7)
}
