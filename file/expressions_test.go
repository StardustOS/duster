package file

import (
	"bytes"
	"encoding/binary"
	"testing"
)

type exprTests struct {
	Input *bytes.Buffer
	Res   Result
}

var tests = []exprTests{
	exprTests{
		Input: bytes.NewBuffer([]byte{byte(DW_OP_lit0)}),
		Res:   Result{Value: uint64(0)},
	},
	exprTests{
		Input: bytes.NewBuffer([]byte{byte(DW_OP_addr), 0x10, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0}),
		Res:   Result{Value: binary.BigEndian.Uint64([]byte{0x10, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0})},
	},
	exprTests{
		Input: bytes.NewBuffer([]byte{byte(DW_OP_const1u), 0x24}),
		Res:   Result{Value: 0x24},
	},
	exprTests{
		Input: bytes.NewBuffer([]byte{byte(DW_OP_const2u), 0x36, 0x42}),
		Res:   Result{Value: 13890},
	},
	exprTests{
		Input: bytes.NewBuffer([]byte{byte(DW_OP_const4u), 0x46, 0x65, 0x78, 0x87}),
		Res:   Result{Value: 1181055111},
	},
	exprTests{
		Input: bytes.NewBuffer([]byte{byte(DW_OP_const8u), 0x20, 0x0, 0x7a, 0x0, 0x0, 0x0, 0x0, 0x0}),
		Res:   Result{Value: 2305977149632282624},
	},
}

func init() {
	for i := 1; i < 32; i++ {
		input := bytes.NewBuffer([]byte{byte(DW_OP_lit0 + Opcode(i))})
		output := Result{Value: uint64(i)}
		e := exprTests{input, output}
		tests = append(tests, e)
	}
}

func TestParse(t *testing.T) {
	for _, test := range tests {
		p := Parser{Input: test.Input}
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
