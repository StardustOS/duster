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

func (control *Xenctrl) SetDebug(domain uint, enable bool) error {
	if enable {
		C.xc_domain_setdebugging(control.key, C.uint(domain), 1)
	} else {
		C.xc_domain_setdebugging(control.key, C.uint(domain), 0)
	}
	return nil
}

func (control *Xenctrl) WordSize(domainid uint) WordSize {
	var size C.uint
	C.xc_domain_get_guest_width(control.key, C.uint(domainid), &size)
	switch WordSize(size) {
	case SixtyFourBit:
		return SixtyFourBit
	case ThirtyTwoBit:
		return ThirtyTwoBit
	}
	return 0
}