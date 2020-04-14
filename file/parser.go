package file

import (
	"debug/dwarf"
	"debug/elf"
	"encoding/binary"
	"errors"

	"github.com/AtomicMalloc/debugger/debugger"
)

type Parser struct {
	data      *dwarf.Data
	types     *TypeManager
	symbols   *SymbolManager
	cu        *dwarf.Entry
	endianess binary.ByteOrder
}

func (p *Parser) parseTypes(reader *dwarf.Reader) error {
	p.types = new(TypeManager)
	p.types.Endianess = p.endianess
	for entry, err := reader.Next(); entry != nil; entry, err = reader.Next() {
		if err != nil {
			return err
		}
		err = p.types.ParseDwarfEntry(entry)
		if err != nil {
			return err
		}
	}
	return nil
}

func (p *Parser) parseVariables(reader *dwarf.Reader) error {
	p.symbols = new(SymbolManager)
	p.symbols.typemanager = p.types
	for entry, err := reader.Next(); entry != nil; entry, err = reader.Next() {
		if err != nil {
			return err
		}
		err = p.symbols.ParseDwarfEntry(entry)
		if err != nil {
			return err
		}
	}
	return nil
}

func (p *Parser) Parse(pc uint64) error {
	reader := p.data.Reader()
	entry, err := reader.SeekPC(pc)
	if err != nil {
		return err
	}

	if p.cu == entry {
		return nil
	}
	p.cu = entry
	err = p.parseTypes(reader)
	if err != nil {
		return err
	}
	reader.SeekPC(pc)
	err = p.parseVariables(reader)
	return err
}

func (p *Parser) GetSymbol(name string, rip uint64) (debugger.Variable, error) {
	err := p.Parse(rip)
	if err != nil {
		return nil, err
	}
	return p.symbols.GetSymbol(rip, name)
}

func (p *Parser) IsPointer(variable debugger.Variable) bool {
	v := variable.(*Variable)
	switch v.typeVar.(type) {
	default:
		return false
	case *Pointer:
		return true
	}
}

func (p *Parser) ParsePointer(variable debugger.Variable, bytes []byte, endianess binary.ByteOrder) (string, error) {
	v := variable.(*Variable)
	switch v.typeVar.(type) {
	default:
		return "", errors.New("Not a pointer")
	case *Pointer:
		p := v.typeVar.(*Pointer)
		return p.typeOfPointer.Parse(bytes, endianess)
	}
}

func (p *Parser) GetPointContentSize(variable debugger.Variable) int {
	v := variable.(*Variable)
	k := v.typeVar.(*Pointer)
	return k.typeOfPointer.Size()
}

func (p *Parser) SymbolManager() *SymbolManager {
	return p.symbols
}

func NewParser(filename string, endianess binary.ByteOrder) (*Parser, error) {
	file, err := elf.Open(filename)
	if err != nil {
		return nil, err
	}
	d, err := file.DWARF()
	if err != nil {
		return nil, err
	}
	newParse := new(Parser)
	newParse.data = d
	newParse.endianess = endianess
	return newParse, nil
}
