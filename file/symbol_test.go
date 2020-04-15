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
				PC: 0x649,
				Variables: []Variable{
					Variable{name: "x"},
					Variable{name: "y"},
				},
			},
		},
	},
	test{
		Filename: "testfiles/globalvars",
		Positions: []pcPos{
			pcPos{
				PC: 0x656,
				Variables: []Variable{
					Variable{name: "x"},
					Variable{name: "y"},
					Variable{name: "k"},
					Variable{name: "hello"},
				},
			},
		},
	},
	test{
		Filename: "testfiles/different-scopes",
		Positions: []pcPos{
			pcPos{
				PC: 0x5fe,
				Variables: []Variable{
					Variable{name: "j"},
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
				PC: 0x626,
				Variables: []Variable{
					Variable{name: "k"},
					Variable{name: "i"},
					Variable{name: "factor"},
				},
			},
			pcPos{
				PC: 0x628,
				Variables: []Variable{
					Variable{name: "k"},
					Variable{name: "i"},
					Variable{name: "factor"},
					Variable{name: "j"},
				},
			},
			pcPos{
				PC: 0x635,
				Variables: []Variable{
					Variable{name: "k"},
					Variable{name: "i"},
					Variable{name: "factor"},
				},
			},
		},
	},
	test{
		Filename: "testfiles/variable_data",
		Positions: []pcPos{
			pcPos{
				PC: 0x71b,
				Variables: []Variable{
					Variable{name: "clean_a"},
					Variable{name: "z"},
					Variable{name: "x"},
				},
			},
			pcPos{
				PC: 0x6aa,
				Variables: []Variable{
					Variable{name: "i"},
					Variable{name: "a"},
				},
			},
			pcPos{
				PC: 0x701,
				Variables: []Variable{
					Variable{name: "clean_a"},
					Variable{name: "z"},
					Variable{name: "a"},
					Variable{name: "meh"},
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
				PC: 0x6b5,
				Variables: []Variable{
					Variable{name: "clean_a"},
					Variable{name: "val1"},
					Variable{name: "val2"},
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
			name := fmt.Sprintf("%s:0x%x", test.Filename, pos.PC)
			t.Run(name, func(t *testing.T) {
				err := parser.Parse(pos.PC)
				if err != nil {
					t.Error(err)
				}
				sym := parser.SymbolManager()
				for _, expected := range pos.Variables {
					variable, err := sym.GetSymbol(pos.PC, expected.Name())
					if err != nil && !test.ExpectedError {
						t.Error(err)
					} else if test.ExpectedError {
						if err != test.Err {
							t.Errorf("Error: expected to get %s not %s", test.Err, err)
						}
					} else if variable.Name() != expected.Name() && !test.ExpectedError {
						t.Errorf("Expected %+v but got %+v", expected, variable)
					}
				}
			})
		}
	}

}
