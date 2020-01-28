package main

import (
	"github.com/AtomicMalloc/debugger/xen"
	"time"
	"fmt"
)

func main() {
	control := xen.Xenctrl{}
	control.Init()
	fmt.Println("Pausing")
	control.Pause(24)
	time.Sleep(time.Second * 60)
	fmt.Println("Unpausing")
	control.UnPause(24)

}