package file

import (
	"debug/dwarf"
	"debug/elf"
	"encoding/binary"
	"fmt"

	"github.com/StardustOS/debugger/debugger"
)

//SymbolicInformation represents the type and variables 
//in the program. Implements the Symbol interface defined in the
//debugger package
type SymbolicInformation struct {
	data      *dwarf.Data
	types     *TypeManager
	symbols   *SymbolManager
	cu        *dwarf.Entry
	endianess binary.ByteOrder
}

//parseType - parses the type information the DWARF file
func (symbolicInfo *SymbolicInformation) parseTypes(reader *dwarf.Reader) error {
	symbolicInfo.types = new(TypeManager)
	symbolicInfo.types.Endianess = symbolicInfo.endianess
	for entry, err := reader.Next(); entry != nil; entry, err = reader.Next() {
		if err != nil {
			return err
		}
		err = symbolicInfo.types.ParseDwarfEntry(entry)
		if err != nil {
			return err
		}
	}
	return nil
}

//parseVariables - parses the variable information in the dwarf
func (symbolicInfo *SymbolicInformation) parseVariables(reader *dwarf.Reader) error {
	symbolicInfo.symbols = new(SymbolManager)
	symbolicInfo.symbols.typemanager = symbolicInfo.types
	for entry, err := reader.Next(); entry != nil; entry, err = reader.Next() {
		if err != nil {
			return err
		}
		err = symbolicInfo.symbols.ParseDwarfEntry(entry)
		if err != nil {
			return err
		}
	}
	return nil
}

//Parse - parses the compile unit that contains the information (both types and variables) that 
//will be required for the current program counter
func (symbolicInfo *SymbolicInformation) Parse(pc uint64) error {
	reader := symbolicInfo.data.Reader()
	entry, err := reader.SeekPC(pc)
	if err != nil {
		return err
	}

	//Checks whether we've already parsed this particularly Compile
	//unit and then skips if we have
	if symbolicInfo.cu == entry {
		return nil
	}

	symbolicInfo.cu = entry
	err = symbolicInfo.parseTypes(reader)
	if err != nil {
		return err
	}
	reader.SeekPC(pc)
	err = symbolicInfo.parseVariables(reader)
	return err
}

//GetSymbol - takes a variable name and the current program counter and returns
//the variable
func (symbolicInfo *SymbolicInformation) GetSymbol(name string, rip uint64) (debugger.Variable, error) {
	err := symbolicInfo.Parse(rip)
	if err != nil {
		return nil, err
	}
	return symbolicInfo.symbols.GetSymbol(rip, name)
}

//IsPointer - takes a variable and returns whether it is a pointer  
func (symbolicInfo *SymbolicInformation) IsPointer(variable debugger.Variable) bool {
	v := variable.(*Variable)
	switch v.typeVar.(type) {
	default:
		return false
	case *Pointer:
		return true
	}
}

//ParsePointer - takes a variable (which a pointer type) and parses the bytes of that type
func (symbolicInfo *SymbolicInformation) ParsePointer(variable debugger.Variable, bytes []byte, endianess binary.ByteOrder) (string, error) {
	v := variable.(*Variable)
	switch v.typeVar.(type) {
	default:
		name := v.Name()
		return "", fmt.Errorf("Error: %s is not a pointer", name)
	case *Pointer:
		pointer := v.typeVar.(*Pointer)
		return pointer.typeOfPointer.Parse(bytes, endianess)
	}
}

//GetPointContentSize - takes a pointer variable and returns the size of that type
//(i.e. if we have int* var; then it would return 4 (assuming an int is 4 bytes))
func (symbolicInfo *SymbolicInformation) GetPointContentSize(variable debugger.Variable) int {
	variableOrignalStruct := variable.(*Variable)
	pointer := variableOrignalStruct.typeVar.(*Pointer)
	return pointer.typeOfPointer.Size()
}

//SymbolManager - returns the symbol manager
func (symbolicInfo *SymbolicInformation) SymbolManager() *SymbolManager {
	return symbolicInfo.symbols
}

//NewSymbolicInformation - this is the constructor for SymbolicInformation struct
func NewSymbolicInformation(filename string, endianess binary.ByteOrder) (*SymbolicInformation, error) {
	file, err := elf.Open(filename)
	if err != nil {
		return nil, err
	}
	dwarfData, err := file.DWARF()
	if err != nil {
		return nil, err
	}
	symbolicInfo := new(SymbolicInformation)
	symbolicInfo.data = dwarfData
	symbolicInfo.endianess = endianess
	return symbolicInfo, nil
}
