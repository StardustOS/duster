package file

import (
	"bytes"
	"encoding/binary"
	"testing"

	"github.com/filecoin-project/go-leb128"
)

const (
	pc = 0x91a
)

type exprTests struct {
	Input *bytes.Buffer
	Res   Result
}

var tests = []exprTests{
	exprTests{
		Input: bytes.NewBuffer([]byte{byte(DW_OP_lit0)}),
		Res:   Result{Uvalue: uint64(0)},
	},
	exprTests{
		Input: bytes.NewBuffer([]byte{byte(DW_OP_addr), 0x10, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0}),
		Res:   Result{Uvalue: binary.BigEndian.Uint64([]byte{0x10, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0})},
	},
	exprTests{
		Input: bytes.NewBuffer([]byte{byte(DW_OP_const1u), 0x24}),
		Res:   Result{Uvalue: 0x24},
	},
	exprTests{
		Input: bytes.NewBuffer([]byte{byte(DW_OP_const2u), 0x36, 0x42}),
		Res:   Result{Uvalue: 13890},
	},
	exprTests{
		Input: bytes.NewBuffer([]byte{byte(DW_OP_const4u), 0x46, 0x65, 0x78, 0x87}),
		Res:   Result{Uvalue: 1181055111},
	},
	exprTests{
		Input: bytes.NewBuffer([]byte{byte(DW_OP_const8u), 0x20, 0x0, 0x7a, 0x0, 0x0, 0x0, 0x0, 0x0}),
		Res:   Result{Uvalue: 2305977149632282624},
	},
	exprTests{
		Input: bytes.NewBuffer([]byte{byte(DW_OP_const1s), 85}),
		Res:   Result{Svalue: 85, Signed: true},
	},
	exprTests{
		Input: bytes.NewBuffer([]byte{byte(DW_OP_const2s), 255, 255}),
		Res:   Result{Svalue: -1, Signed: true},
	},
	exprTests{
		Input: bytes.NewBuffer([]byte{byte(DW_OP_const4s), 255, 255, 255, 255}),
		Res:   Result{Svalue: -1, Signed: true},
	},
	exprTests{
		Input: bytes.NewBuffer([]byte{byte(DW_OP_const8s), 255, 255, 255, 255, 255, 255, 255, 255}),
		Res:   Result{Svalue: -1, Signed: true},
	},
	exprTests{
		Input: bytes.NewBuffer(append([]byte{byte(DW_OP_constu)}, leb128.FromUInt64(2^40)...)),
		Res:   Result{Uvalue: 2 ^ 40},
	},
	exprTests{
		Input: bytes.NewBuffer(append([]byte{byte(DW_OP_fbreg)}, leb128.FromUInt64(2^20)...)),
		Res:   Result{Uvalue: 2 ^ 20 + pc},
	},
}

func init() {
	for i := 1; i < 32; i++ {
		input := bytes.NewBuffer([]byte{byte(DW_OP_lit0 + Opcode(i))})
		output := Result{Uvalue: uint64(i)}
		e := exprTests{input, output}
		tests = append(tests, e)
	}
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
