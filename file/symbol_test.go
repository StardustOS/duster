package file

import (
	"encoding/binary"
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
				PC: 0x401126,
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
				PC: 0x401106,
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
				PC: 0x40111a,
				Variables: []Variable{
					Variable{Name: "k"},
					Variable{Name: "i"},
					Variable{Name: "factor"},
				},
			},
			pcPos{
				PC: 0x401134,
				Variables: []Variable{
					Variable{Name: "k"},
					Variable{Name: "i"},
					Variable{Name: "factor"},
					Variable{Name: "j"},
				},
			},
			pcPos{
				PC: 0x401141,
				Variables: []Variable{
					Variable{Name: "k"},
					Variable{Name: "i"},
					Variable{Name: "factor"},
				},
			},
		},
	},
	test{
		Filename: "testfiles/variable_data",
		Positions: []pcPos{
			pcPos{
				PC: 0x401194,
				Variables: []Variable{
					Variable{Name: "clean_a"},
					Variable{Name: "z"},
					Variable{Name: "x"},
				},
			},
			pcPos{
				PC: 0x401136,
				Variables: []Variable{
					Variable{Name: "i"},
					Variable{Name: "a"},
				},
			},
			pcPos{
				PC: 0x401189,
				Variables: []Variable{
					Variable{Name: "clean_a"},
					Variable{Name: "z"},
					Variable{Name: "a"},
					Variable{Name: "meh"},
				},
			},
		},
		ExpectedError: true,
		Err:           SymbolNotFound,
	},
	test{
		Filename: "testfiles/variable_data",
		Positions: []pcPos{
			pcPos{
				PC: 0x401136,
				Variables: []Variable{
					Variable{Name: "clean_a"},
					Variable{Name: "val1"},
					Variable{Name: "val2"},
				},
			},
		},
	},
}

func TestGetSymbol(t *testing.T) {
	for _, test := range positions {
		parser, err := NewParser(test.Filename, binary.LittleEndian)
		if err != nil {
			t.Error(err)
		}
		for _, pos := range test.Positions {
			name := fmt.Sprintf("%s:%d", test.Filename, pos.PC)
			t.Run(name, func(t *testing.T) {
				err := parser.Parse(pos.PC)
				if err != nil {
					t.Error(err)
				}
				sym := parser.SymbolManager()
				for _, expected := range pos.Variables {
					variable, err := sym.GetSymbol(pos.PC, expected.Name)
					fmt.Println(variable)
					fmt.Println(err)
					if err != nil && !test.ExpectedError {
						t.Error(err)
					} else if test.ExpectedError {
						if err != test.Err {
							t.Errorf("Error: expected to get %s not %s", test.Err, err)
						}
					} else if variable.Name != expected.Name && !test.ExpectedError {
						t.Errorf("Expected %+v but got %+v", expected, variable)
					}
				}
			})
		}
	}

}
