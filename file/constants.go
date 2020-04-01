package file

import (
	"fmt"
)

type Opcode byte
type OperandType uint

const (
	// DW_OP_litx just means push the literal value x on to the stack
	// (e.g. DW_OP_lit0 just means push zero onto the stack)
	DW_OP_lit0    Opcode = 0x30
	DW_OP_lit31   Opcode = 0x4f
	DW_OP_addr    Opcode = 0x03
	DW_OP_const1u Opcode = 0x08
	DW_OP_const2u Opcode = 0x0a
	DW_OP_const4u Opcode = 0x0c
	DW_OP_const8u Opcode = 0x0e
	DW_OP_const1s Opcode = 0x09
	DW_OP_const2s Opcode = 0x0b
	DW_OP_const4s Opcode = 0x0d
	DW_OP_const8s Opcode = 0x0f
	DW_OP_constu  Opcode = 0x10
	DW_OP_consts  Opcode = 0x11
	DW_OP_fbreg   Opcode = 0x91
	DW_OP_breg0   Opcode = 0x50
	DW_OP_breg31  Opcode = 0x8f
	DW_OP_bregx   Opcode = 0x92

	unsignedLEBI           OperandType = 1
	signedLEBI             OperandType = 2
	signedInt              OperandType = 3
	unsignedInt            OperandType = 4
	unsignedThenSignedLEBI OperandType = 4
)

func (op Opcode) OperandSize() (size uint) {
	switch op {
	case DW_OP_addr, DW_OP_const8s, DW_OP_const8u:
		size = 8
	case DW_OP_const1u, DW_OP_const1s:
		size = 1
	case DW_OP_const2u, DW_OP_const2s:
		size = 2
	case DW_OP_const4u, DW_OP_const4s:
		size = 4
	}
	return
}

func (op Opcode) OperandType() (t []OperandType) {
	if op >= DW_OP_breg0 && op <= DW_OP_breg31 {
		t = append(t, signedLEBI)
		return
	}
	switch op {
	case DW_OP_consts, DW_OP_fbreg:
		t = append(t, signedLEBI)
	case DW_OP_constu:
		t = append(t, unsignedLEBI)
	case DW_OP_addr:
		fallthrough
	case DW_OP_const1u:
		fallthrough
	case DW_OP_const2u:
		fallthrough
	case DW_OP_const4u:
		fallthrough
	case DW_OP_const8u:
		t = append(t, unsignedInt)
	case DW_OP_const1s:
		fallthrough
	case DW_OP_const2s:
		fallthrough
	case DW_OP_const4s:
		fallthrough
	case DW_OP_const8s:
		t = append(t, signedInt)
	case DW_OP_bregx:
		t = append(t, unsignedLEBI, signedLEBI)
	}
	return
}

func (op Opcode) NoOperand() (no uint) {
	switch {
	case DW_OP_const1u <= op && op <= DW_OP_const8s:
		fallthrough
	case op >= DW_OP_breg0 && op <= DW_OP_breg31:
		fallthrough
	case op == DW_OP_constu:
		fallthrough
	case op == DW_OP_consts:
		fallthrough
	case op == DW_OP_addr:
		no = 1
	case op == DW_OP_bregx:
		no = 2
	}

	return
}

func (op Opcode) operation(s *stack, operands []Value, regs Registers) {
	var res Value
	if op >= DW_OP_lit0 && op <= DW_OP_lit31 {
		value := uint64(op - DW_OP_lit0)
		res.Uvalue = value
	} else if op >= DW_OP_breg0 && op <= DW_OP_breg31 {
		reg := uint64(op - DW_OP_breg0)
		regVal := int64(regs.GetRegister(reg))
		offset := operands[0]
		fmt.Printf("Offset %+v\n", offset)
		address := uint64(regVal + offset.Int64())
		fmt.Println("address", address)
		res.Uvalue = address
	}

	switch op {
	case DW_OP_addr:
		fallthrough
	case DW_OP_const1u, DW_OP_const2u, DW_OP_const4u, DW_OP_const8u:
		fallthrough
	case DW_OP_const1s, DW_OP_const2s, DW_OP_const4s, DW_OP_const8s:
		fallthrough
	case DW_OP_constu, DW_OP_consts:
		res = operands[0]
	case DW_OP_fbreg:
		val := operands[0]
		base := int64(regs.GetFrameBase())
		fmt.Println("offset in here", val.Int64())
		fmt.Println("Base", base)
		res.Uvalue = uint64(val.Int64() + base)
	case DW_OP_bregx:
		reg := operands[0]
		offset := operands[1]
		regVal := int64(regs.GetRegister(reg.Uvalue))
		res.Uvalue = uint64(regVal + offset.Int64())
	}
	s.push(res)
}
