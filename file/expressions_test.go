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
