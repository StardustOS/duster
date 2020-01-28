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
	fmt.Println("The word size is %d\n", control.WordSize(24))
	time.Sleep(time.Second * 3)
	fmt.Println("Unpausing")
	control.UnPause(24)

}