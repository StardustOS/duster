package file

import (
	"encoding/binary"
	"math/big"

	"github.com/filecoin-project/go-leb128"
)

type Opcode byte

const (
	// DW_OP_litx just means push the literal value x on to the stack
	// (e.g. DW_OP_lit0 just means push zero onto the stack)
	DW_OP_lit0    Opcode = 0x30
	DW_OP_lit1    Opcode = 0x31
	DW_OP_lit2    Opcode = 0x32
	DW_OP_lit3    Opcode = 0x33
	DW_OP_lit4    Opcode = 0x34
	DW_OP_lit5    Opcode = 0x35
	DW_OP_lit6    Opcode = 0x36
	DW_OP_lit7    Opcode = 0x37
	DW_OP_lit8    Opcode = 0x38
	DW_OP_lit9    Opcode = 0x39
	DW_OP_lit10   Opcode = 0x3a
	DW_OP_lit11   Opcode = 0x3b
	DW_OP_lit12   Opcode = 0x3c
	DW_OP_lit13   Opcode = 0x3d
	DW_OP_lit14   Opcode = 0x3e
	DW_OP_lit15   Opcode = 0x3f
	DW_OP_lit16   Opcode = 0x40
	DW_OP_lit17   Opcode = 0x41
	DW_OP_lit18   Opcode = 0x42
	DW_OP_lit19   Opcode = 0x43
	DW_OP_lit20   Opcode = 0x44
	DW_OP_lit21   Opcode = 0x45
	DW_OP_lit22   Opcode = 0x46
	DW_OP_lit23   Opcode = 0x47
	DW_OP_lit24   Opcode = 0x48
	DW_OP_lit25   Opcode = 0x49
	DW_OP_lit26   Opcode = 0x4a
	DW_OP_lit27   Opcode = 0x4b
	DW_OP_lit28   Opcode = 0x4c
	DW_OP_lit29   Opcode = 0x4d
	DW_OP_lit30   Opcode = 0x4e
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
)

var operands = map[Opcode]int{
	DW_OP_addr:    8,
	DW_OP_const1u: 1,
	DW_OP_const2u: 2,
	DW_OP_const4u: 4,
	DW_OP_const8u: 8,
	DW_OP_const1s: 1,
	DW_OP_const2s: 2,
	DW_OP_const4s: 4,
	DW_OP_const8s: 8,
}

func (op Opcode) operation(s *stack, operand []byte, regs Registers) {
	if op >= DW_OP_lit0 && op <= DW_OP_lit31 {
		value := uint64(op - DW_OP_lit0)
		element := item{uVal: value, signed: false}
		s.push(element)
	} else if op >= DW_OP_breg0 && op <= DW_OP_breg31 {
		regVal := regs.GetRegister(uint(op - DW_OP_breg0))
		bigInt := leb128.ToBigInt(operand)
		if bigInt.IsUint64() {
			val := bigInt.Uint64()
			element := item{uVal: uint64(int64(val) + int64(regVal))}
			s.push(element)
		}
	}

	switch op {
	case DW_OP_addr:
		address := binary.BigEndian.Uint64(operand)
		element := item{uVal: address, signed: false}
		s.push(element)
	case DW_OP_const1u:
		unsignedConst := uint64(uint8(operand[0]))
		element := item{uVal: unsignedConst, signed: false}
		s.push(element)
	case DW_OP_const2u:
		unsignedConst := uint64(binary.BigEndian.Uint16(operand))
		element := item{uVal: unsignedConst, signed: false}
		s.push(element)
	case DW_OP_const4u:
		unsignedConst := uint64(binary.BigEndian.Uint32(operand))
		element := item{uVal: unsignedConst, signed: false}
		s.push(element)
	case DW_OP_const8u:
		unsignedConst := binary.BigEndian.Uint64(operand)
		element := item{uVal: unsignedConst, signed: false}
		s.push(element)
	case DW_OP_const1s:
		signedConst := int64(int8(operand[0]))
		element := item{sVal: signedConst, signed: true}
		s.push(element)
	case DW_OP_const2s:
		signedConst := int64(int16(binary.BigEndian.Uint16(operand)))
		element := item{sVal: signedConst, signed: true}
		s.push(element)
	case DW_OP_const4s:
		signedConst := int64(int32(binary.BigEndian.Uint32(operand)))
		element := item{sVal: signedConst, signed: true}
		s.push(element)
	case DW_OP_const8s:
		signedConst := int64(binary.BigEndian.Uint64(operand))
		element := item{sVal: signedConst, signed: true}
		s.push(element)
	case DW_OP_constu:
		bigInt := leb128.ToBigInt(operand)
		var element item
		if bigInt.IsUint64() {
			element = item{uVal: bigInt.Uint64()}
		} else {
			mask := big.NewInt(2 ^ 64)
			val := big.NewInt(0)
			val.And(bigInt, mask)
			element = item{uVal: val.Uint64()}
		}
		s.push(element)
	case DW_OP_consts:
		bigInt := leb128.ToBigInt(operand)
		var element item
		if bigInt.IsInt64() {
			element = item{sVal: bigInt.Int64(), signed: true}
		} else {
			mask := big.NewInt(2 ^ 64)
			val := big.NewInt(0)
			val.And(bigInt, mask)
			element = item{sVal: val.Int64(), signed: true}
		}
		s.push(element)
	case DW_OP_fbreg:
		bigInt := leb128.ToBigInt(operand)
		if bigInt.IsUint64() {
			val := bigInt.Uint64()
			element := item{uVal: val}
			s.push(element)
		}
		//TODO: error stuff
	}
}
