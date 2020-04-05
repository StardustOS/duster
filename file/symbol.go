package file

import (
	"debug/dwarf"
	"errors"
)

type SymbolError int

const (
	InvalidDWARF   SymbolError = 0
	SymbolNotFound SymbolError = 1
	NoLoctionFound SymbolError = 2
	NoName SymbolError = 3
)

func (err SymbolError) Error() string {
	switch err {
	case InvalidDWARF:
		return "Error: not a valid dwarf file"
	case SymbolNotFound:
		return "Error: symbol not found"
	case NoLoctionFound:
		return "Error: not found location"
	}
	return ""
}

type SymbolTable struct {
	parent   *SymbolTable
	symbols  map[string]*Variable
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
			return c
		}
	}
	return sym
}

func (sym *SymbolTable) Parent() *SymbolTable {
	return sym.parent
}

func (sym *SymbolTable) Lookup(symbolName string) (*Variable, error) {
	if variable, ok := sym.symbols[symbolName]; ok {
		return variable, nil
	} else if sym.parent != nil {
		return sym.parent.Lookup(symbolName)
	}
	return nil, SymbolNotFound
}

func (sym *SymbolTable) AddVariable(variable *Variable) {
	if sym.symbols == nil {
		sym.symbols = make(map[string]*Variable)
	}
	sym.symbols[variable.Name] = variable
}

func (sym *SymbolTable) AddParent(parent *SymbolTable) {
	sym.parent = parent
}

func parseVariable(entry *dwarf.Entry, manager *TypeManager) (*Variable, error) {
	field := entry.AttrField(dwarf.AttrName)
	variable := new(Variable)
	if field == nil {
		return nil, NoName
	}
	name := field.Val.(string)
	variable.Name = name
	field = entry.AttrField(dwarf.AttrType)
	if field == nil {
		return nil, errors.New("Error: could not find type")
	}
	offset := field.Val.(dwarf.Offset)
	variable.typeVar = manager.getType(offset)

	field = entry.AttrField(dwarf.AttrLocation)
	if field == nil {
		return nil, NoLoctionFound
	}
	variable.location = field.Val.([]byte)
	return variable, nil
}

type SymbolManager struct {
	typemanager  *TypeManager
	currentTable *SymbolTable
	rootTable    *SymbolTable
}

func (manager *SymbolManager) ParseDwarfEntry(entry *dwarf.Entry) error {
	if manager.rootTable == nil {
		manager.rootTable = new(SymbolTable)
		manager.currentTable = manager.rootTable
	}

	switch entry.Tag {
	case dwarf.TagVariable:
		variable, err := parseVariable(entry, manager.typemanager)
		if err != nil {
			if err == NoLoctionFound || err == NoName {
				return nil
			}
			return err
		}
		manager.currentTable.AddVariable(variable)
	case dwarf.TagSubprogram:
		lowPC, highPC, err := parsePC(entry)
		if err != nil {
			return err
		}
		parent := manager.rootTable
		newTable := &SymbolTable{LowerPC: lowPC, UpperPC: highPC}
		parent.AddChild(newTable)
		newTable.AddParent(parent)
		manager.currentTable = newTable

	case dwarf.TagLexDwarfBlock:
		lowPC, highPC, err := parsePC(entry)
		if err != nil {
			return err
		}

		newTable := &SymbolTable{LowerPC: lowPC, UpperPC: highPC}
		if manager.currentTable.PCInStack(lowPC) {
			manager.currentTable.AddChild(newTable)
			newTable.AddParent(manager.currentTable)
		} else {
			parent := manager.currentTable.Parent()
			parent.AddChild(newTable)
			newTable.AddParent(parent)
		}
		manager.currentTable = newTable
	}
	return nil
}

func (manager *SymbolManager) GetSymbol(pc uint64, name string) (*Variable, error) {
	table := manager.rootTable.GetNextTable(pc)
	variable, err := table.Lookup(name)
	return variable, err
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
