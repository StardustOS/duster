package file

import "bytes"

type Error int

const (
	Empty Error = 1
)

type Registers interface {
	GetRegister(uint) uint64
}

func (e Error) Error() string {
	switch e {
	case Empty:
		return "Stack is empty"
	}
	return ""
}

type item struct {
	signed bool
	uVal   uint64
	sVal   int64
}

type stack struct {
	stack []item
}

func (s *stack) push(element item) {
	s.stack = append([]item{element}, s.stack...)
}

func (s *stack) pop() (item, error) {
	if s.stack == nil {
		return item{}, Empty
	}
	element := s.stack[0]
	s.stack = s.stack[1:]
	return element, nil
}

type Result struct {
	Signed bool
	Uvalue uint64
	Svalue int64
}

type Parser struct {
	Input        *bytes.Buffer
	StackPointer uint64
	stack        stack
	Regs         Registers
}

func readLEBI128Integer(input *bytes.Buffer) []byte {
	found := false
	var bytes []byte
	for !found {
		b := input.Next(1)
		currByte := b[0]
		bytes = append(bytes, currByte)
		found = currByte&0x80 == 0
	}
	return bytes
}

func (p *Parser) Parse() error {
	for p.Input.Len() > 0 {
		rawInput := p.Input.Next(1)
		op := Opcode(rawInput[0])
		if size, ok := operands[op]; ok {
			operand := p.Input.Next(size)
			op.operation(&p.stack, operand, p.Regs)
		} else {
			if op == DW_OP_constu || op == DW_OP_consts || op == DW_OP_fbreg || (op >= DW_OP_breg0 && op <= DW_OP_breg31) || op == DW_OP_bregx {
				operand := readLEBI128Integer(p.Input)
				if op == DW_OP_bregx {
					operand = append(operand, readLEBI128Integer(p.Input)...)
				}
				op.operation(&p.stack, operand, p.Regs)
				if op == DW_OP_fbreg {
					element, _ := p.stack.pop()
					element.uVal = element.uVal + p.StackPointer
					p.stack.push(element)
				}
			} else {
				op.operation(&p.stack, nil, p.Regs)
			}
		}
	}
	return nil
}

func (p *Parser) Result() (Result, error) {
	result, err := p.stack.pop()
	if err != nil {
		return Result{}, err
	}
	return Result{result.signed, result.uVal, result.sVal}, nil
}
