package main

import (
	"fmt"

	"github.com/AtomicMalloc/debugger/debugger"
)

func main() {
	d := debugger.Debugger{}
	defer d.Cleanup()
	err := d.Init(64)
	if err != nil {
		fmt.Println(err)
	}

	d.StartSingle(0, true)
	for i := 0; i < 10000; i++ {
		d.Step(0)
	}
	d.StartSingle(0, false)
	d.UnPause()
	// for i := 0; i < 10000; i++ {
	// 	d.UnPause()
	// }
	// d.StartSingle(0, false)
	// d.UnPause()
}
