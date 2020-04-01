package main

import (
	"fmt"

	"github.com/AtomicMalloc/debugger/debugger"
)

func main() {
	d := debugger.Debugger{}
	defer d.Cleanup()
	err := d.Init(274, "../stardust-experimental/build/mini-os")
	if err != nil {
		fmt.Println(err)
	}

	d.StartSingle(0, true)
	d.SetBreakpoint("startup.c", 74, 0)
	// fmt.Println("Got out of breakpoint")
	// for i := 0; i < 1000; i++ {
	// 	fmt.Println("Before step")
	for {
		d.Step(0)
		fmt.Println(d.GetLineInformation())
	}
	
	// 	fmt.Println("After step")
	// }
	d.StartSingle(0, false)
	d.UnPause()
	// for i := 0; i < 10000; i++ {
	// 	d.UnPause()
	// }
	// d.StartSingle(0, false)
	// d.UnPause()
}
