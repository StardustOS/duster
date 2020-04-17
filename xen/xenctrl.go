package xen

import (
	"github.com/AtomicMalloc/debugger/debugger"
)

/*
#cgo CFLAGS: -g -Wall
#cgo LDFLAGS: -lxenctrl -lxencall
#include <xenctrl.h>
#include <errno.h>
#include <stdio.h>
#include <xencall.h>
#include <string.h>
// Also the capitialisation here is not careless even though it 
// is not good style in C. Since in Go we need upper case to export 
// a struct or attribute, and if it is not capitiliased then Go considers 
// the C a different "package". Thus if there is no uppercase letter then 
// we cannot access any of the data
struct Regs {
	uint64_t Rax;
	uint64_t Rbx;
	uint64_t Rcx;
	uint64_t Rdx;
	uint64_t Rsp;
	uint64_t Rbp;
	uint64_t Rsi;
	uint64_t Rdi;
	uint64_t R8;
	uint64_t R9;
	uint64_t R10;
	uint64_t R11;
	uint64_t R12;
	uint64_t R13;
	uint64_t R14;
	uint64_t R15;
	uint64_t Rip;
	uint64_t Rflags;
	uint64_t fs;
	uint64_t gs;
	uint64_t ds;
	uint64_t ss;
	uint64_t es;
	uint64_t cs;
};

// We need these helper functions (i.e. we can't xc_vcpu_get/setcontext directly in go). This is because
// handling unions in Go is impossible (or such a pain it's to handle and I couldn't find a solution).
// So dealing with the union in C is far more convenient.
int getRegister(xc_interface* key, uint32_t domainid, uint32_t vcpu, struct Regs* buffer) {
	vcpu_guest_context_any_t context;
	int err = xc_vcpu_getcontext(key, domainid, vcpu, &context);
	if (err) {
		puts(strerror(errno));
		return err;
	}
	// puts("Reading out");
	// printf("RAX: %lu\n", context.x64.user_regs.rax);
	// printf("RBX: %lu\n", context.x64.user_regs.rbx);
	// printf("RBP: %lu\n", context.x64.user_regs.rbp);
	// printf("RSI: %lu\n", context.x64.user_regs.rsi);
	buffer->Rax = context.x64.user_regs.rax;
	buffer->Rbx = context.x64.user_regs.rbx;
	buffer->Rcx = context.x64.user_regs.rcx;
	buffer->Rdx = context.x64.user_regs.rdx;
	buffer->Rsp = context.x64.user_regs.rsp;
	buffer->Rbp = context.x64.user_regs.rbp;
	buffer->Rsi = context.x64.user_regs.rsi;
	buffer->Rdi = context.x64.user_regs.rdi;
	buffer->R8 = context.x64.user_regs.r8;
	buffer->R9 = context.x64.user_regs.r9;
	buffer->R10 = context.x64.user_regs.r10;
	buffer->R11 = context.x64.user_regs.r11;
	buffer->R12 = context.x64.user_regs.r12;
	buffer->R13 = context.x64.user_regs.r13;
	buffer->R14 = context.x64.user_regs.r14;
	buffer->R15 = context.x64.user_regs.r15;
	buffer->Rflags = context.x64.user_regs.rflags;
	buffer->fs = context.x64.user_regs.fs;
	buffer->gs = context.x64.user_regs.gs;
	buffer->ds = context.x64.user_regs.ds;
	buffer->ss = context.x64.user_regs.ss;
	buffer->es = context.x64.user_regs.es;
	buffer->cs = context.x64.user_regs.cs;
	buffer->Rip = context.x64.user_regs.rip;

	return 0;
}
int setRegister(xc_interface* key, struct Regs regs, uint32_t domainid, uint32_t vcpu) {
	//printf("DomainID: %d and vcpu: %d\n", domainid, vcpu);
	//printf("fs: %lu\n", regs.fs);
	vcpu_guest_context_any_t context; //= malloc(sizeof(vcpu_guest_context_any_t));
	int err = xc_vcpu_getcontext(key, domainid, vcpu, &context);
	if (err) {
		puts(strerror(errno));
		return err;
	}
	context.x64.user_regs.rax = regs.Rax;
	context.x64.user_regs.rbx = regs.Rbx;
	context.x64.user_regs.rcx = regs.Rcx;
	context.x64.user_regs.rdx = regs.Rdx;
	context.x64.user_regs.rsp = regs.Rsp;
	context.x64.user_regs.rbp = regs.Rbp;
	context.x64.user_regs.rsi = regs.Rsi;
	context.x64.user_regs.rdi = regs.Rdi;
	context.x64.user_regs.r8 = regs.R8;
	context.x64.user_regs.r9 = regs.R9;
	context.x64.user_regs.r10 = regs.R10;
	context.x64.user_regs.r11 = regs.R11;
	context.x64.user_regs.r12 = regs.R12;
	context.x64.user_regs.r13 = regs.R13;
	context.x64.user_regs.r14 = regs.R14;
	context.x64.user_regs.r15 = regs.R15;
	context.x64.user_regs.rflags = regs.Rflags;
	// context.x64.user_regs.fs = regs.fs;
	// context.x64.user_regs.gs = regs.gs;
	// context.x64.user_regs.ds = regs.ds;
	// context.x64.user_regs.ss = regs.ss;
	// context.x64.user_regs.es = regs.es;
	// context.x64.user_regs.cs = regs.cs;
	context.x64.user_regs.rip = regs.Rip;
	//printf("FROM C: %lu\n", context.x64.user_regs.rflags);
	//xc_vcpu_getcontext(key, domainid, vcpu, &context);
	//printf("FROM C: %lu\n", context.x64.user_regs.rflags);
	// puts("Writing back");
	// printf("RAX: %lu\n", context.x64.user_regs.rax);
	// printf("RBX: %lu\n", context.x64.user_regs.rbx);
	// printf("RBP: %lu\n", context.x64.user_regs.rbp);
	// printf("RSI: %lu\n", context.x64.user_regs.rsi);

	err = xc_vcpu_setcontext(key, domainid, vcpu, &context);
	if (err) {
		puts(xc_strerror(key, errno));
		return err;
	}
	return 0;
}
int is_paused(xc_interface *xch, uint32_t domaind) {
	xc_dominfo_t info;
	int no = xc_domain_getinfo(xch, domaind, 1, &info);
	if (no == -1) {
		puts(xc_strerror(xch, errno));
		return -1;
	}
	return info.paused;
}
*/
import "C"

import (
	"errors"
)

type Uint64 C.ulong

type Xenctrl struct {
	key      *C.xc_interface
	DomainID uint32
}

//Init gets the handler for the xen domain
func (control *Xenctrl) Init() error {
	control.key = C.xc_interface_open(nil, nil, 0)
	return nil
}

//IsPaused returns whether the domain is paused or not
func (control *Xenctrl) IsPaused() bool {
	paused := C.is_paused(control.key, C.uint(control.DomainID))
	if paused == 0 {
		return false
	} else {
		return true
	}
}

//Close destories the handler required to access 
//Xen control API
func (control *Xenctrl) Close() error {
	C.xc_interface_close(control.key)
	control.key = nil
	return nil
}

//Key gets the handler for the Xen control
func (control *Xenctrl) Key() *C.xc_interface {
	return control.key
}

//Pause - pauses the domain
func (control *Xenctrl) Pause() error {
	err := C.xc_domain_pause(control.key, C.uint(control.DomainID))
	if err != 0 {
		return errors.New("Error: could not pause domain")
	}
	return nil
}

//Unpause - unpauses the domain
func (control *Xenctrl) Unpause() error {
	err := C.xc_domain_unpause(control.key, C.uint(control.DomainID))
	if err != 0 {
		return errors.New("Error: could not unpause domain")
	}
	return nil
}

//SetDebug - puts the domain into debug mode or not
func (control *Xenctrl) SetDebug(domain uint32, enable bool) error {
	var err C.int 
	if enable {
		err = C.xc_domain_setdebugging(control.key, C.uint32_t(control.DomainID), 1)
	} else {
		err = C.xc_domain_setdebugging(control.key, C.uint32_t(domain), 0)
	}
	if err != 0 {
		return errors.New("Error: could not put or take the domain out of debug mode")
	}
	return nil
}

//GetRegisters gets the registers from the domain and puts them in Go format
func (control *Xenctrl) GetRegisters(vcpu uint32) (debugger.Registers, error) {
	var context C.struct_Regs
	err := C.getRegister(control.key, C.uint32_t(control.DomainID), C.uint32_t(vcpu), &context)
	if err != 0 {
		return nil, nil
	}
	register := &Register{}
	register.SetRegister("rax", uint64(context.Rax))
	register.SetRegister("rbx", uint64(context.Rbx))
	register.SetRegister("rbp", uint64(context.Rbp))
	register.SetRegister("rcx", uint64(context.Rcx))
	register.SetRegister("rsp", uint64(context.Rsp))
	register.SetRegister("rip", uint64(context.Rip))
	register.SetRegister("rsi", uint64(context.Rsi))
	register.SetRegister("rdi", uint64(context.Rdi))
	register.SetRegister("r8", uint64(context.R8))
	register.SetRegister("r9", uint64(context.R9))
	register.SetRegister("r10", uint64(context.R10))
	register.SetRegister("r11", uint64(context.R11))
	register.SetRegister("r12", uint64(context.R12))
	register.SetRegister("r13", uint64(context.R13))
	register.SetRegister("r14", uint64(context.R14))
	register.SetRegister("r15", uint64(context.R15))
	register.SetRegister("rflags", uint64(context.Rflags))
	return register, nil
}

//SetRegisters - sets a register for a particular vpcu
func (control *Xenctrl) SetRegisters(vcpu uint32, regs debugger.Registers) error {
	r := regs.(*Register)
	if control.key == nil {
		return errors.New("Error: control does not have a handle")
	}
	err := C.setRegister(control.key, r.convertC(), C.uint32_t(control.DomainID), C.uint32_t(vcpu))
	if err != 0 {
		return errors.New("Error: could not set registers")
	}
	return nil
}
