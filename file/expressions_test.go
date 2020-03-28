package file

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"testing"

	"github.com/AtomicMalloc/go-leb128"
)

const (
	pc           = 0x91a
	reg dummyReg = 40
)

type dummyReg uint

func (r dummyReg) GetRegister(v uint64) uint64 {
	if v == uint64(reg) {
		return 0x08
	}
	return 0x102004
}

func (r dummyReg) GetFrameBase() uint64 {
	return pc
}

type exprTests struct {
	Input *bytes.Reader
	Res   Value
}

var tests = []exprTests{
	exprTests{
		Input: bytes.NewReader([]byte{byte(DW_OP_lit0)}),
		Res:   Value{Uvalue: uint64(0)},
	},
	exprTests{
		Input: bytes.NewReader([]byte{byte(DW_OP_addr), 0x10, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0}),
		Res:   Value{Uvalue: binary.BigEndian.Uint64([]byte{0x10, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0})},
	},
	exprTests{
		Input: bytes.NewReader([]byte{byte(DW_OP_const1u), 0x24}),
		Res:   Value{Uvalue: 0x24},
	},
	exprTests{
		Input: bytes.NewReader([]byte{byte(DW_OP_const2u), 0x36, 0x42}),
		Res:   Value{Uvalue: 13890},
	},
	exprTests{
		Input: bytes.NewReader([]byte{byte(DW_OP_const4u), 0x46, 0x65, 0x78, 0x87}),
		Res:   Value{Uvalue: 1181055111},
	},
	exprTests{
		Input: bytes.NewReader([]byte{byte(DW_OP_const8u), 0x20, 0x0, 0x7a, 0x0, 0x0, 0x0, 0x0, 0x0}),
		Res:   Value{Uvalue: 2305977149632282624},
	},
	exprTests{
		Input: bytes.NewReader([]byte{byte(DW_OP_const1s), 255}),
		Res:   Value{Svalue: -1, Signed: true},
	},
	exprTests{
		Input: bytes.NewReader([]byte{byte(DW_OP_const2s), 255, 255}),
		Res:   Value{Svalue: -1, Signed: true},
	},
	exprTests{
		Input: bytes.NewReader([]byte{byte(DW_OP_const4s), 255, 255, 255, 255}),
		Res:   Value{Svalue: -1, Signed: true},
	},
	exprTests{
		Input: bytes.NewReader([]byte{byte(DW_OP_const8s), 255, 255, 255, 255, 255, 255, 255, 255}),
		Res:   Value{Svalue: -1, Signed: true},
	},
	exprTests{
		Input: bytes.NewReader(append([]byte{byte(DW_OP_constu)}, leb128.FromUInt64(2^40)...)),
		Res:   Value{Uvalue: 2 ^ 40},
	},
	exprTests{
		Input: bytes.NewReader(append([]byte{byte(DW_OP_fbreg)}, leb128.FromUInt64(2^20)...)),
		Res:   Value{Uvalue: 2 ^ 20 + pc},
	},
}

func init() {
	for i := 1; i < 32; i++ {
		input := bytes.NewReader([]byte{byte(DW_OP_lit0 + Opcode(i))})
		output := Value{Uvalue: uint64(i)}
		e := exprTests{input, output}
		tests = append(tests, e)
	}

	for i := 0; i < 32; i++ {
		inputSlice := []byte{byte(DW_OP_breg0 + Opcode(i))}
		inputSlice = append(inputSlice, leb128.FromUInt64(2^20+2)...)

		input := bytes.NewReader(inputSlice)
		output := Value{Uvalue: uint64(2 ^ 20 + 2 + 0x102004)}
		e := exprTests{input, output}
		tests = append(tests, e)
	}

	inputSlice := []byte{byte(DW_OP_bregx)}
	fmt.Println("HERE", uint64(reg))
	inputSlice = append(inputSlice, leb128.FromUInt64(uint64(reg))...)
	inputSlice = append(inputSlice, leb128.FromUInt64(128)...)
	input := bytes.NewReader(inputSlice)
	output := Value{Uvalue: uint64(0x08 + 128)}
	e := exprTests{input, output}
	tests = append(tests, e)
}

func TestParse(t *testing.T) {
	for _, test := range tests {
		p := Parser{Input: test.Input, StackPointer: pc, Regs: dummyReg(0)}
		err := p.Parse()
		if err != nil {
			t.Error(err)
		}
		res, err := p.Result()
		if err != nil {
			t.Error(err)
		}
		if res != test.Res {
			t.Errorf("Expected %+v but got %+v", test.Res, res)
		}
	}
}
