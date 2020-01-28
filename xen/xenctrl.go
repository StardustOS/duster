package xen
// #cgo CFLAGS: -g -Wall 
//#cgo LDFLAGS: -lxenctrl
// #include <xenctrl.h>
import "C"

type Xenctrl struct {
	key *C.xc_interface
}

func (control *Xenctrl) Init() error {
	control.key = C.xc_interface_open(nil, nil, 0)
	return nil
}

func (control *Xenctrl) Pause(domain uint) error {
	C.xc_domain_pause(control.key, C.uint(domain))
	return nil 
}

func (control *Xenctrl) UnPause(domain uint) error {
	C.xc_domain_unpause(control.key, C.uint(domain))
	return nil
}