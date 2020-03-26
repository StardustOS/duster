package file

import "bytes"

type Error int

const (
	Empty Error = 1
)

func (e Error) Error() string {
	switch e {
	case Empty:
		return "Stack is empty"
	}
	return ""
}

type stack struct {
	stack []uint64
}

func (s *stack) push(element uint64) {
	s.stack = append([]uint64{element}, s.stack...)
}

func (s *stack) pop() (uint64, error) {
	if s.stack == nil {
		return 0, Empty
	}
	element := s.stack[0]
	s.stack = s.stack[1:]
	return element, nil
}

type Result struct {
	Value uint64
}

type Parser struct {
	Input *bytes.Buffer
	stack stack
}

func (p *Parser) Parse() error {
	for p.Input.Len() > 0 {
		rawInput := p.Input.Next(1)
		op := Opcode(rawInput[0])
		if size, ok := operands[op]; ok {
			operand := p.Input.Next(size)
			op.operation(&p.stack, operand)
		} else {
			op.operation(&p.stack, nil)
		}
	}
	return nil
}

func (p *Parser) Result() (Result, error) {
	result, err := p.stack.pop()
	if err != nil {
		return Result{}, err
	}
	return Result{result}, nil
}
