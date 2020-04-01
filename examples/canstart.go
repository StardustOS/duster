package main

import "github.com/AtomicMalloc/debugger/xen"

func main() {
	xen.StartDomain("../stardust-experimental/mini-os.conf")
}
