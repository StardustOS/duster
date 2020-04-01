package main

import (
	"encoding/binary"
	"fmt"
	"time"

	"github.com/AtomicMalloc/debugger/xen"
)

func main() {
	domainid := uint32(29)
	control := xen.Xenctrl{}
	mem := xen.Memory{}
	control.Init()
	defer control.Close()
	mem.Init(&control)
	defer mem.Close()
	fmt.Println("Pausing")
	control.Pause(domainid)
	control.SetDebug(domainid, true)
	r := control.GetRegisterContext(domainid, 0)
	// fmt.Println("The word size is ", control.WordSize(domainid))
	// fmt.Println(r.Rflags)

	mem.Map(884816, domainid, 1, 0)
	buffer := mem.Read(884816, 4096)
	if buffer != nil {
		fmt.Println(buffer[80:89])
		data := binary.LittleEndian.Uint64(buffer[80:89])
		fmt.Println("The value stored at the address is: ", data)
		d := int(data)
		fmt.Println("The value stored at the address is: ", d)
	}
	//mem.Write(52696, []byte{0xCC})
	mem.UnMap(884816, 1)
	time.Sleep(time.Second * 5)
	fmt.Println("Unpausing")
	//control.SetDebug(10, false)
	control.UnPause(domainid)

}
