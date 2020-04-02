package file

import (
	"debug/dwarf"
	"debug/elf"
	"encoding/binary"
	"fmt"
	"strings"

	"github.com/go-delve/delve/pkg/dwarf/op"
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
	Name       string
	Location   []byte
	frame      uint64
	typeVar    Type
	typeOffset dwarf.Offset
	data       *dwarf.Data
}

func (variable *Variable) ParseVal(bytes []byte) (string, error) {
	reader := variable.data.Reader()
	reader.Seek(variable.typeOffset)
	entry, _ := reader.Next()
	if entry.Tag == dwarf.TagTypedef {
		field := entry.AttrField(dwarf.AttrType)
		offset := field.Val.(dwarf.Offset)
		reader := variable.data.Reader()
		reader.Seek(offset)
		entry, _ = reader.Next()
	}
	var newBase Type
	newBase.Entry(entry)
	fmt.Printf("%+v\n", entry)
	return newBase.Parse(bytes)
}

func (variable *Variable) Size() int64 {
	reader := variable.data.Reader()
	reader.Seek(variable.typeOffset)
	entry, _ := reader.Next()
	if entry.Tag == dwarf.TagTypedef {
		field := entry.AttrField(dwarf.AttrType)
		offset := field.Val.(dwarf.Offset)
		reader := variable.data.Reader()
		reader.Seek(offset)
		entry, _ = reader.Next()
	}
	fmt.Printf("%+v\n", entry)
	var newBase Type
	newBase.Entry(entry)
	return newBase.Size

}

func (variable *Variable) Address(regs op.DwarfRegisters) int64 {
	//	regs.FrameBase = int64(variable.frame)

	address, m, b := op.ExecuteStackProgram(regs, variable.Location)
	fmt.Println("M ", m)
	fmt.Println("B", b)
	return address
}

type SymbolTable struct {
	parent   *SymbolTable
	symbols  map[string]Variable
	children []*SymbolTable
	LowerPC  uint64
	UpperPC  uint64
	Regs     Registers
	types    map[dwarf.Offset]Type
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

func (sym *SymbolTable) Lookup(symbolName string) (Variable, error) {
	if variable, ok := sym.symbols[symbolName]; ok {
		return variable, nil
	} else if sym.parent != nil {
		return sym.parent.Lookup(symbolName)
	}
	return Variable{}, SymbolNotFound
}

func (sym *SymbolTable) AddVariable(name string, location []byte, upper uint64, offset dwarf.Offset, data *dwarf.Data) {
	if sym.symbols == nil {
		sym.symbols = make(map[string]Variable)
	}
	sym.symbols[name] = Variable{Name: name, Location: location, frame: upper, typeOffset: offset, data: data}
}

func (sym *SymbolTable) AddType(offset dwarf.Offset, varType Type) {
	if sym.types == nil {
		sym.types = make(map[dwarf.Offset]Type)
	}
	sym.types[offset] = varType
}

func (sym *SymbolTable) GetType(offset dwarf.Offset) (Type, error) {
	if val, ok := sym.types[offset]; ok {
		return val, nil
	} else if sym.parent != nil {
		return sym.parent.GetType(offset)
	}
	return Type{}, SymbolNotFound
}

func (sym *SymbolTable) AddParent(parent *SymbolTable) {
	sym.parent = parent
}

type Symbol struct {
	pc    uint64
	RSP   uint64
	Data  *dwarf.Data
	entry *dwarf.Entry
	root  *SymbolTable
	types map[dwarf.Offset]Type
}

func (sym *Symbol) Init(filename string) error {
	file, err := elf.Open(filename)
	if err != nil {
		return err
	}
	sym.Data, err = file.DWARF()
	return err
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

func parseFramebase(entry *dwarf.Entry) (v uint64) {
	field := entry.AttrField(dwarf.AttrFrameBase)

	if field != nil {
		framebase := field.Val.([]uint8)
		bytes := []byte(framebase)
		length := len(bytes)
		switch length {
		case 1:
			val := int8(bytes[0])
			v = uint64(val)
		case 2:
			val := int16(binary.BigEndian.Uint16(bytes))
			v = uint64(val)
		case 4:
			val := int32(binary.BigEndian.Uint32(bytes))
			v = uint64(val)
		case 8:
			val := int64(binary.BigEndian.Uint64(bytes))
			v = uint64(val)
		}
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
				if strings.Compare(name, "val") == 0 {
					fmt.Println(entry)
				}
				if !ok {
					return InvalidDWARF
				}
				//fmt.Printf("%+v\n", entry)

				field = entry.AttrField(dwarf.AttrLocation)
				if field != nil {
					bytes := field.Val.([]byte)
					field = entry.AttrField(dwarf.AttrType)
					offset := field.Val.(dwarf.Offset)
					current.AddVariable(name, bytes, current.UpperPC, offset, sym.Data)
				}
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
		case dwarf.TagBaseType:
			var newBaseType Type
			newBaseType.Entry(entry)
			current.AddType(entry.Offset, newBaseType)
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
	variable, err := current.Lookup(name)

	if err != nil {
		return Variable{}, err
	}
	// fmt.Println(variable)
	// p := Parser{Input: bytes.NewReader(variable.Location), StackPointer: sym.RSP}
	// err = p.Parse()
	// if err != nil {
	// 	return Variable{}, err
	// }
	// val, err := p.Result()
	// if err != nil {
	// 	return Variable{}, err
	// }
	// variable.Address = val.Uvalue
	return variable, nil
}
