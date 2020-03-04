package xen

import "C"
import (
	"fmt"
)

type WordSize uint

const (
	SixtyFourBit WordSize = 8
	ThirtyTwoBit WordSize = 4
)

type Register struct {
	Type      WordSize
	registers map[string]uint64
}

func (register *Register) AddRegister(name string, content uint64) error {
	if register.registers == nil {
		register.registers = make(map[string]uint64)
	}
	register.registers[name] = content
	return nil
}

func (register *Register) convertC() C.struct_Regs {
	var regs C.struct_Regs
	regs.Rax = C.ulong(register.registers["rax"])
	regs.Rbx = C.ulong(register.registers["rbx"])
	regs.Rip = C.ulong(register.registers["rip"])
	regs.Rcx = C.ulong(register.registers["rcx"])
	regs.Rbp = C.ulong(register.registers["rbp"])
	regs.Rsp = C.ulong(register.registers["rsp"])
	regs.Rsi = C.ulong(register.registers["rsi"])
	regs.Rdi = C.ulong(register.registers["rdi"])
	regs.R8 = C.ulong(register.registers["r8"])
	regs.R9 = C.ulong(register.registers["r9"])
	regs.R10 = C.ulong(register.registers["r10"])
	regs.R11 = C.ulong(register.registers["r11"])
	regs.R12 = C.ulong(register.registers["r12"])
	regs.R13 = C.ulong(register.registers["r13"])
	regs.R14 = C.ulong(register.registers["r14"])
	regs.R15 = C.ulong(register.registers["r15"])
	regs.Rflags = C.ulong(register.registers["rflags"])
	return regs
}

func (register *Register) GetRegister(name string) (uint64, error) {
	if content, ok := register.registers[name]; ok {
		return content, nil
	}
	return 0, fmt.Errorf("The register %s could not be found", name)
}
