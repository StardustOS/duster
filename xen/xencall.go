package xen

/*
#cgo CFLAGS: -g -Wall
#cgo LDFLAGS: -lxenctrl -lxencall
#include <xenctrl.h>
#include <errno.h>
#include <stdio.h>
#include <xencall.h>
#include <xenctrl.h>
#include <string.h>
#include <xen/domctl.h>
int pause_cpu(xc_interface* interface, xencall_handle* key, uint32_t domainid, int command, uint32_t vcpu) {
	DECLARE_HYPERCALL_BUFFER(struct xen_domctl, domctl);
	domctl = (struct xen_domctl*)(xc_hypercall_buffer_alloc(interface, domctl, sizeof(*domctl)));
	domctl->domain = domainid;
  	domctl->interface_version = XEN_DOMCTL_INTERFACE_VERSION;
	domctl->cmd = command;
	memset(&domctl->u, 0, sizeof(domctl->u));
	domctl->u.gdbsx_pauseunp_vcpu.vcpu = vcpu;
	int err = xencall1(key, __HYPERVISOR_domctl, HYPERCALL_BUFFER_AS_ARG(domctl));
	if (err) {
		puts(strerror(errno));
	}
	return err;
}
*/
import "C"
import "fmt"

type Hypercall uint

const (
	PauseCPU   Hypercall = 0
	UnPauseCPU Hypercall = 1
)

func (call Hypercall) ConvertToC() C.int {
	switch call {
	case PauseCPU:
		return C.XEN_DOMCTL_gdbsx_pausevcpu
	case UnPauseCPU:
		return C.XEN_DOMCTL_gdbsx_unpausevcpu
	}
	return C.int(0)
}

type XenCall struct {
	key *C.xencall_handle
}

func (call *XenCall) Init() error {
	call.key = C.xencall_open(nil, 0)
	return nil
}

func (call *XenCall) Close() error {
	C.xencall_close(call.key)
	return nil
}

func (call *XenCall) HyperCall(domain Xenctrl, hypercall Hypercall, domainID, vcpu uint32) error {
	switch hypercall {
	case PauseCPU, UnPauseCPU:
		err := C.pause_cpu(domain.Key(), call.key, C.uint32_t(domainID), hypercall.ConvertToC(), C.uint32_t(vcpu))
		fmt.Println("Error from hyper call", err)
	}
	return nil
}
