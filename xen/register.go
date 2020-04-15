package xen

import "C"
import (
	"encoding/binary"
	"fmt"
	"strings"

	"github.com/go-delve/delve/pkg/dwarf/op"
)



//The mapping between hardware registers and DWARF registers is specified
//in the System V ABI AMD64 Architecture Processor Supplement page 57,
//figure 3.36
//https://www.uclibc.org/docs/psABI-x86_64.pdf
//(Taken from: https://github.com/go-delve/delve/blob/master/pkg/proc/amd64_arch.go)
var amd64DwarfToName = map[uint64]string{
	0:  "Rax",
	1:  "Rdx",
	2:  "Rcx",
	3:  "Rbx",
	4:  "Rsi",
	5:  "Rdi",
	6:  "Rbp",
	7:  "Rsp",
	8:  "R8",
	9:  "R9",
	10: "R10",
	11: "R11",
	12: "R12",
	13: "R13",
	14: "R14",
	15: "R15",
	16: "Rip",
	17: "XMM0",
	18: "XMM1",
	19: "XMM2",
	20: "XMM3",
	21: "XMM4",
	22: "XMM5",
	23: "XMM6",
	24: "XMM7",
	25: "XMM8",
	26: "XMM9",
	27: "XMM10",
	28: "XMM11",
	29: "XMM12",
	30: "XMM13",
	31: "XMM14",
	32: "XMM15",
	33: "ST(0)",
	34: "ST(1)",
	35: "ST(2)",
	36: "ST(3)",
	37: "ST(4)",
	38: "ST(5)",
	39: "ST(6)",
	40: "ST(7)",
	49: "Rflags",
	50: "Es",
	51: "Cs",
	52: "Ss",
	53: "Ds",
	54: "Fs",
	55: "Gs",
	58: "Fs_base",
	59: "Gs_base",
	64: "MXCSR",
	65: "CW",
	66: "SW",
}

//Represents the state of a register at current point of time
type Register struct {
	registers map[string]uint64
}

//SetRegister takes the name of a register and sets that register to the value passed
func (register *Register) SetRegister(name string, content uint64) error {
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

//GetRegister returns the content of a register 
func (register *Register) GetRegister(name string) (uint64, error) {
	if content, ok := register.registers[name]; ok {
		return content, nil
	}
	return 0, fmt.Errorf("The register %s could not be found", name)
}

//DwarfRegisters return the register in the format recognised by the 
//DWARF stack machine for expression
func (register *Register) DwarfRegisters() *op.DwarfRegisters {
	r := &op.DwarfRegisters{}


	for key, registerName := range amd64DwarfToName {
		registerName = strings.ToLower(registerName)
		if val, ok := register.registers[registerName]; ok {
			b := make([]byte, 8)
			binary.LittleEndian.PutUint64(b, uint64(val))
			r.AddReg(key, &op.DwarfRegister{val, b})
		} else {
			b := make([]byte, 8)
			r.AddReg(key, &op.DwarfRegister{val, b})
		}
	}
	var regs []*op.DwarfRegister

	for _, reg := range r.Regs {
		if reg != nil {
			regs = append(regs, reg)
		}
	}
	r.Regs = regs
	r.ByteOrder = binary.LittleEndian
	r.CFA = int64(register.registers["rbp"]) + 16
	r.BPRegNum = 6
	r.SPRegNum = 7
	r.PCRegNum = 16
	r.FrameBase = int64(register.registers["rbp"]) + 16
	return r
}
