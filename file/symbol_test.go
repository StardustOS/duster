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
	Filename      string
	Positions     []pcPos
	ExpectedError bool
	Err           error
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
	test{
		Filename: "testfiles/different-scopes",
		Positions: []pcPos{
			pcPos{
				PC: 0x635,
				Variables: []Variable{
					Variable{Name: "j"},
				},
			},
		},
		ExpectedError: true,
		Err:           SymbolNotFound,
	},
	test{
		Filename: "testfiles/different-scopes",
		Positions: []pcPos{
			pcPos{
				PC: 0x60e,
				Variables: []Variable{
					Variable{Name: "k"},
					Variable{Name: "i"},
					Variable{Name: "factor"},
				},
			},
			pcPos{
				PC: 0x628,
				Variables: []Variable{
					Variable{Name: "k"},
					Variable{Name: "i"},
					Variable{Name: "factor"},
					Variable{Name: "j"},
				},
			},
			pcPos{
				PC: 0x635,
				Variables: []Variable{
					Variable{Name: "k"},
					Variable{Name: "i"},
					Variable{Name: "factor"},
				},
			},
		},
	},
}

func TestGetSymbol(t *testing.T) {
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
				for _, expected := range pos.Variables {
					variable, err := sym.GetSymbol(pos.PC, expected.Name)

					if err != nil && !test.ExpectedError {
						t.Error(err)
					} else if test.ExpectedError {
						if err != test.Err {
							t.Errorf("Error: expected to get %s not %s", test.Err, err)
						}
					}
					if variable != expected && !test.ExpectedError {
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
