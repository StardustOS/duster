package main

import "github.com/AtomicMalloc/debugger/debugger"

func main() {
	d := debugger.Debugger{}
	defer d.Cleanup()
	d.Init(64)
	d.SetBreakpoint(9226, 1, 0)
}
