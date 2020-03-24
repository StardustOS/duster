package file

import (
	"debug/dwarf"
	"fmt"
)

type SymbolError int

const (
	InvalidDWARF   SymbolError = 0
	SymbolNotFound SymbolError = 1
)

func (err SymbolError) Error() string {
	switch err {
	case InvalidDWARF:
		return "Error: not a valid dwarf file"
	case SymbolNotFound:
		return "Error: symbol not found"
	}
	return ""
}

type Variable struct {
	Name string
}

type SymbolTable struct {
	parent   *SymbolTable
	symbols  map[string]Variable
	children []*SymbolTable
	LowerPC  uint64
	UpperPC  uint64
}

func (sym *SymbolTable) PCInStack(pc uint64) bool {
	return pc >= sym.LowerPC && pc < sym.UpperPC
}

func (sym *SymbolTable) AddChild(table *SymbolTable) {
	sym.children = append(sym.children, table)
}

func (sym *SymbolTable) GetNextTable(pc uint64) *SymbolTable {
	for _, child := range sym.children {
		if child.PCInStack(pc) {
			c := child.GetNextTable(pc)
			fmt.Printf("%+v\n", c)
			return c
		}
	}
	return sym
}

func (sym *SymbolTable) Parent() *SymbolTable {
	return sym.parent
}

func (sym *SymbolTable) Lookup(symbolName string) (Variable, error) {
	if variable, ok := sym.symbols[symbolName]; ok {
		return variable, nil
	} else if sym.parent != nil {
		return sym.parent.Lookup(symbolName)
	}
	return Variable{}, SymbolNotFound
}

func (sym *SymbolTable) AddVariable(name string) {
	if sym.symbols == nil {
		sym.symbols = make(map[string]Variable)
	}
	sym.symbols[name] = Variable{Name: name}
}

func (sym *SymbolTable) AddParent(parent *SymbolTable) {
	sym.parent = parent
}

type Symbol struct {
	pc    uint64
	Data  *dwarf.Data
	entry *dwarf.Entry
	root  *SymbolTable
}

func parsePC(entry *dwarf.Entry) (lower, upper uint64, err error) {
	lowPC := entry.AttrField(dwarf.AttrLowpc)
	highPC := entry.AttrField(dwarf.AttrHighpc)
	if lowPC != nil && highPC != nil {
		var ok bool
		lowPC, ok := lowPC.Val.(uint64)
		lower = lowPC
		if !ok {
			err = InvalidDWARF
			return
		}
		highPC, ok := highPC.Val.(int64)
		if !ok {
			err = InvalidDWARF
			return
		}
		upper = uint64(highPC)
		upper += uint64(lowPC)
	}
	return
}

func (sym *Symbol) parse(cu *dwarf.Entry) error {
	if sym.entry == cu {
		return nil
	}
	reader := sym.Data.Reader()
	sym.entry = cu
	sym.root = new(SymbolTable)
	current := sym.root
	sym.root.parent = nil
	for entry, err := reader.Next(); entry != nil; entry, err = reader.Next() {
		if err != nil {
			return err
		}
		switch entry.Tag {
		case dwarf.TagVariable:
			field := entry.AttrField(dwarf.AttrName)
			if field != nil {
				name, ok := field.Val.(string)
				if !ok {
					return InvalidDWARF
				}
				current.AddVariable(name)
			}
		case dwarf.TagSubprogram:
			lowPC, highPC, err := parsePC(entry)
			if err != nil {
				return err
			}
			parent := sym.root
			newTable := &SymbolTable{LowerPC: lowPC, UpperPC: highPC}
			parent.AddChild(newTable)
			newTable.AddParent(parent)
			current = newTable

		case dwarf.TagLexDwarfBlock:
			lowPC, highPC, err := parsePC(entry)
			if err != nil {
				return err
			}

			newTable := &SymbolTable{LowerPC: lowPC, UpperPC: highPC}
			if current.PCInStack(lowPC) {
				current.AddChild(newTable)
				newTable.AddParent(current)
			} else {
				parent := current.Parent()
				parent.AddChild(newTable)
				newTable.AddParent(parent)
			}
			current = newTable
		}
	}
	return nil
}

func (sym *Symbol) GetSymbol(pc uint64, name string) (Variable, error) {
	reader := sym.Data.Reader()
	compileUnit, err := reader.SeekPC(pc)
	if err != nil {
		return Variable{}, err
	}
	err = sym.parse(compileUnit)
	if err != nil {
		return Variable{}, err
	}
	current := sym.root.GetNextTable(pc)
	return current.Lookup(name)
}
