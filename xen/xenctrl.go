package xen

/*
#cgo CFLAGS: -g -Wall
#cgo LDFLAGS: -lxenctrl -lxencall
#include <xenctrl.h>
#include <errno.h>
#include <stdio.h>
#include <xencall.h>
#include <string.h>
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
	//printf("rip (from C): %lu\n", context.x64.user_regs.rip);

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
	"fmt"
)

type Uint64 C.ulong

type Xenctrl struct {
	key *C.xc_interface
}

func (control *Xenctrl) Init() error {
	control.key = C.xc_interface_open(nil, nil, 0)
	return nil
}

func (control *Xenctrl) IsPaused(doaminId uint32) bool {
	paused := C.is_paused(control.key, C.uint(doaminId))
	if paused == 0 {
		return false
	} else if paused == 1 {
		return true
	}
	fmt.Println("Something went wrong")
	return false
}

func (control *Xenctrl) Close() error {
	C.xc_interface_close(control.key)
	control.key = nil
	return nil
}

func (control *Xenctrl) Key() *C.xc_interface {
	return control.key
}

func (control *Xenctrl) Pause(domain uint32) error {
	err := C.xc_domain_pause(control.key, C.uint(domain))
	if err != 0 {
		return errors.New("SOMETHING BAD HAPPENDED")
	}
	return nil
}

func (control *Xenctrl) UnPause(domain uint32) error {
	err := C.xc_domain_unpause(control.key, C.uint(domain))
	if err != 0 {
		return errors.New("SOMETHING BAD HAPPENDED")
	}
	return nil
}

func (control *Xenctrl) SetDebug(domain uint32, enable bool) error {
	if enable {
		err := C.xc_domain_setdebugging(control.key, C.uint32_t(domain), 1)
		if err != 0 {
			fmt.Println("Error at the debugging")
		}
	} else {
		C.xc_domain_setdebugging(control.key, C.uint32_t(domain), 0)
	}
	return nil
}

func (control *Xenctrl) WordSize(domainid uint32) WordSize {
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

func (control *Xenctrl) GetRegisterContext(domainid uint32, vcpu uint32) *Register {
	var context C.struct_Regs
	err := C.getRegister(control.key, C.uint32_t(domainid), C.uint32_t(vcpu), &context)
	if err != 0 {
		fmt.Println("Something went wrong in get register")
	}
	//fmt.Println(err)
	//fmt.Println("GetContext rbx")
	register := &Register{}
	register.AddRegister("rax", uint64(context.Rax))
	register.AddRegister("rbx", uint64(context.Rbx))
	register.AddRegister("rbp", uint64(context.Rbp))
	register.AddRegister("rcx", uint64(context.Rcx))
	register.AddRegister("rsp", uint64(context.Rsp))
	register.AddRegister("rip", uint64(context.Rip))
	register.AddRegister("rsi", uint64(context.Rsi))
	register.AddRegister("rdi", uint64(context.Rdi))
	register.AddRegister("r8", uint64(context.R8))
	register.AddRegister("r9", uint64(context.R9))
	register.AddRegister("r10", uint64(context.R10))
	register.AddRegister("r11", uint64(context.R11))
	register.AddRegister("r12", uint64(context.R12))
	register.AddRegister("r13", uint64(context.R13))
	register.AddRegister("r14", uint64(context.R14))
	register.AddRegister("r15", uint64(context.R15))
	register.AddRegister("rflags", uint64(context.Rflags))
	return register
}

func (control *Xenctrl) SetRegisterContext(regs *Register, domainid, vcpu uint32) error {
	if control.key == nil {
		fmt.Println("KEY IS NIL")
		return nil
	}
	err := C.setRegister(control.key, regs.convertC(), C.uint32_t(domainid), C.uint32_t(vcpu))
	fmt.Println("Error from set registers: ", err)
	return nil
}
