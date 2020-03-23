package file

import (
	"debug/elf"
	"fmt"
	"testing"
)

type pcPos struct {
	PC        uint64
	Variables []Variable
}

type test struct {
	Filename  string
	Positions []pcPos
}

var positions = []test{
	test{
		Filename: "testfiles/simple",
		Positions: []pcPos{
			pcPos{
				PC: 0x63a,
				Variables: []Variable{
					Variable{Name: "x"},
					Variable{Name: "y"},
				},
			},
		},
	},
	test{
		Filename: "testfiles/globalvars",
		Positions: []pcPos{
			pcPos{
				PC: 0x63a,
				Variables: []Variable{
					Variable{Name: "x"},
					Variable{Name: "y"},
					Variable{Name: "k"},
					Variable{Name: "hello"},
				},
			},
		},
	},
}

func TestUpdateSymbolSimple(t *testing.T) {
	for _, test := range positions {
		file, err := elf.Open(test.Filename)
		if err != nil {
			t.Error(err)
		}
		d, err := file.DWARF()
		if err != nil {
			t.Error(err)
		}
		for _, pos := range test.Positions {
			name := fmt.Sprintf("%s:%d", test.Filename, pos.PC)
			t.Run(name, func(t *testing.T) {
				sym := Symbol{Data: d}
				err = sym.Update(pos.PC)
				if err != nil {
					t.Error(err)
				}
				for _, expected := range pos.Variables {
					variable, err := sym.GetSymbol(expected.Name)
					if err != nil {
						t.Error(err)
					}
					if variable != expected {
						t.Errorf("Expected %+v but got %+v", expected, variable)
					}
				}
			})
		}
	}

}

// func TestUpdateSymbols(t *testing.T) {
// 	file, err := elf.Open("testfiles/variable_data")
// 	if err != nil {
// 		t.Error(err)
// 	}
// 	d, err := file.DWARF()
// 	if err != nil {
// 		t.Error(err)
// 	}
// 	sym := Symbol{Data: d}
// 	err = sym.Update(0x710)
// 	if err != nil {
// 		t.Error(err)
// 	}
// 	variable, err := sym.GetSymbol("a")
// 	if err != nil {
// 		t.Error(err)
// 	}
// 	expected := Variable{Name: "a"}
// 	if variable != expected {
// 		t.Errorf("Expected %+v but got %+v", expected, variable)
// 	}
// }
