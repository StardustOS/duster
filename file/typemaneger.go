package file

import (
	"debug/dwarf"
	"encoding/binary"
	"errors"
	"fmt"
	"math"
)

type DType uint
type ErrorD int

const (
	_ DType = iota
	Address
	Boolean
	ComplexFloat
	Float
	Sinteger
	Schar
	Uinteger
	Uchar
	Ifloat
	Pdecimal
	NumericString
	EditedString
	SfixedPointInteger
	UfixedPointInteger
	WrongSize ErrorD = 0
)

type parse func()

type Type interface {
	Size() int
	Parse([]byte, binary.ByteOrder) (string, error)
}

type BaseType struct {
	size     int
	Encoding DType
	Name     string
}

func parseInteger(bytes []byte, endianess binary.ByteOrder) int64 {
	var integer int64
	length := len(bytes)
	switch length {
	case 1:
		val := int8(bytes[0])
		integer = int64(val)
		fmt.Println("INTEGER", integer)
	case 2:
		val := int16(endianess.Uint16(bytes))
		integer = int64(val)
	case 4:
		val := int32(endianess.Uint32(bytes))
		integer = int64(val)
	case 8:
		integer = int64(endianess.Uint64(bytes))
	}
	return integer
}

func parseUinteger(bytes []byte, endianess binary.ByteOrder) uint64 {
	var val uint64
	length := len(bytes)
	switch length {
	case 1:
		val = uint64(uint8(bytes[0]))
	case 2:
		val = uint64(endianess.Uint16(bytes))
	case 4:
		val = uint64(endianess.Uint32(bytes))
	case 8:
		val = uint64(endianess.Uint64(bytes))
	}
	return val
}

func (t *BaseType) Parse(bytes []byte, endianess binary.ByteOrder) (string, error) {
	if len(bytes) != t.size {
		return "", errors.New("MEH")
	}
	var output string
	switch t.Encoding {
	case Address:
		integer := uint64(parseInteger(bytes, endianess))
		output = fmt.Sprintf("%x", integer)
	case Boolean:
		boolean := parseInteger(bytes, endianess) == 1
		output = fmt.Sprintf("%t", boolean)
	case Float:
		integer := parseInteger(bytes, endianess)
		if len(bytes) == 4 {
			float := math.Float32frombits(uint32(integer))
			output = fmt.Sprintf("%f", float)
		} else {
			float := math.Float64frombits(uint64(integer))
			output = fmt.Sprintf("%f", float)
		}
	case Sinteger, Schar:
		integer := parseInteger(bytes, endianess)
		output = fmt.Sprintf("%d", integer)
	case Uinteger, Uchar:
		integer := uint64(parseUinteger(bytes, endianess))
		output = fmt.Sprintf("%d", integer)
	}
	return output, nil
}

func (t *BaseType) Size() int {
	return t.size
}

func parseBaseEntry(entry *dwarf.Entry) (*BaseType, error) {
	base := new(BaseType)
	field := entry.AttrField(dwarf.AttrByteSize)
	if field == nil {
		return nil, errors.New("Oops")
	}
	base.size = int(field.Val.(int64))
	field = entry.AttrField(dwarf.AttrEncoding)
	if field == nil {
		return nil, errors.New("Oops")
	}
	base.Encoding = DType(field.Val.(int64))
	field = entry.AttrField(dwarf.AttrName)
	if field == nil {
		return nil, errors.New("Oops")
	}
	base.Name = field.Val.(string)
	return base, nil
}

type TypeManager struct {
	Endianess binary.ByteOrder
	types     map[dwarf.Offset]Type
}

func (manager *TypeManager) addType(offset dwarf.Offset, t Type) {
	if manager.types == nil {
		manager.types = make(map[dwarf.Offset]Type)
	}
	manager.types[offset] = t
}

func (manager *TypeManager) getType(offset dwarf.Offset) Type {
	return manager.types[offset]
}

func (manager *TypeManager) Size(offset dwarf.Offset) int {
	t := manager.types[offset]
	if t != nil {
		return t.Size()
	}
	return 0
}

func (manager *TypeManager) ParseBytes(offset dwarf.Offset, bytes []byte) (string, error) {
	t := manager.getType(offset)
	fmt.Printf("%d\n", offset)
	fmt.Println(t)
	str, err := t.Parse(bytes, manager.Endianess)
	return str, err
}

func (manager *TypeManager) ParseDwarfEntry(entry *dwarf.Entry) error {
	switch entry.Tag {
	case dwarf.TagBaseType:
		fmt.Printf("%+v\n", entry)
		base, err := parseBaseEntry(entry)
		if err != nil {
			return err
		}
		manager.addType(entry.Offset, base)
	}
	return nil
}
