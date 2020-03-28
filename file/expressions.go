package file

import (
	"bytes"
	"encoding/binary"
	"math/big"

	"github.com/AtomicMalloc/go-leb128"
)

type Error int

const (
	Empty Error = 1
)

type Registers interface {
	GetRegister(uint64) uint64
	GetFrameBase() uint64
}

func (e Error) Error() string {
	switch e {
	case Empty:
		return "Stack is empty"
	}
	return ""
}

type Value struct {
	Signed bool
	Uvalue uint64
	Svalue int64
}

func (val Value) Int64() int64 {
	if val.Signed {
		return val.Svalue
	} else {
		return int64(val.Uvalue)
	}
}

type stack struct {
	stack []Value
}

func (s *stack) push(element Value) {
	s.stack = append([]Value{element}, s.stack...)
}

func (s *stack) pop() (Value, error) {
	if s.stack == nil {
		return Value{}, Empty
	}
	element := s.stack[0]
	s.stack = s.stack[1:]
	return element, nil
}

type Parser struct {
	Input        *bytes.Reader
	StackPointer uint64
	stack        stack
	Regs         Registers
}

func getLast64Bits(no *big.Int) *big.Int {
	mask := big.NewInt(2 ^ 64)
	val := big.NewInt(0)
	return val.And(mask, no)
}

func extendSlice(bytes []byte, bigEndian bool) []byte {
	length := len(bytes)
	if length < 8 {
		extend := make([]byte, 8-length)
		if bigEndian {
			bytes = append(extend, bytes...)
		} else {
			bytes = append(bytes, extend...)
		}
	}
	return bytes
}

func parseUnsignedInt(input *bytes.Reader, opcode Opcode) uint64 {
	length := opcode.OperandSize()
	bytes := make([]byte, length)
	input.Read(bytes)
	bytes = extendSlice(bytes, true)

	integer := binary.BigEndian.Uint64(bytes)
	return integer
}

func (p *Parser) parseSignedInt(opcode Opcode) (v int64) {
	length := opcode.OperandSize()
	bytes := make([]byte, length)
	p.Input.Read(bytes)
	switch length {
	case 1:
		val := int8(bytes[0])
		v = int64(val)
	case 2:
		val := int16(binary.BigEndian.Uint16(bytes))
		v = int64(val)
	case 4:
		val := int32(binary.BigEndian.Uint32(bytes))
		v = int64(val)
	case 8:
		val := int64(binary.BigEndian.Uint64(bytes))
		v = int64(val)
	}
	return
}

func (p *Parser) parseLEBI(signed bool) (val Value, err error) {
	no, err := leb128.ToBigInt(p.Input)
	if err != nil {
		return Value{}, err
	}
	if signed {
		if !no.IsInt64() {
			no = getLast64Bits(no)
		}
		val.Svalue = no.Int64()
		val.Signed = true
	} else {
		if !no.IsUint64() {
			no = getLast64Bits(no)
		}
		val.Uvalue = no.Uint64()
	}
	return
}

func (p *Parser) Parse() error {
	for p.Input.Len() > 0 {
		rawOpcode, err := p.Input.ReadByte()
		if err != nil {
			return err
		}
		opcode := Opcode(rawOpcode)
		noOperands := opcode.NoOperand()
		var operands []Value

		operandTypes := opcode.OperandType()
		for _, operandType := range operandTypes {
			var val Value
			switch operandType {
			case signedLEBI:
				no, err := leb128.ToBigInt(p.Input)
				if err != nil {
					return err
				}
				if !no.IsInt64() {
					no = getLast64Bits(no)
				}
				val.Svalue = no.Int64()
				val.Signed = true
			case unsignedLEBI:
				no, err := leb128.ToBigInt(p.Input)
				if err != nil {
					return err
				}
				if !no.IsUint64() {
					no = getLast64Bits(no)
				}
				val.Uvalue = no.Uint64()
			case unsignedInt:
				integer := uint64(p.parseSignedInt(opcode))
				val.Uvalue = integer
			case signedInt:
				integer := p.parseSignedInt(opcode)
				val.Svalue = integer
				val.Signed = true
			}
			operands = append(operands, val)
		}

		opcode.operation(&p.stack, operands, p.Regs)

	}
	return nil
}

func (p *Parser) Result() (Value, error) {
	result, err := p.stack.pop()
	return result, err
}
