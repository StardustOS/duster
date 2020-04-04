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
	WrongSize        ErrorD = 0
	NoAssociatedType ErrorD = 1
)

func (e ErrorD) Error() string {
	switch e {
	case NoAssociatedType:
		return "No associated TYpe"
	}
	return ""
}

type Type interface {
	Size() int
	Parse([]byte, binary.ByteOrder) (string, error)
}

type TypeDef struct {
	Name string
	Base Type
}

func (t *TypeDef) Size() int {
	return t.Base.Size()
}

func (t *TypeDef) Parse(bytes []byte, endianess binary.ByteOrder) (string, error) {
	return t.Base.Parse(bytes, endianess)
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

func parseTypeDef(entry *dwarf.Entry, manager *TypeManager) (*TypeDef, error) {
	typedef := new(TypeDef)
	field := entry.AttrField(dwarf.AttrName)
	if field == nil {
		return nil, errors.New("No attrName")
	}
	typedef.Name = field.Val.(string)

	field = entry.AttrField(dwarf.AttrType)
	if field == nil {
		fmt.Printf("%+v\n", entry)
		return nil, NoAssociatedType
	}
	offset := field.Val.(dwarf.Offset)
	t := manager.getType(offset)
	if t == nil {
		manager.addWaiting(offset, typedef)
	} else {
		typedef.Base = t
	}
	return typedef, nil
}

type TypeManager struct {
	Endianess  binary.ByteOrder
	types      map[dwarf.Offset]Type
	waitingDef map[dwarf.Offset][]Type
}

func (manager *TypeManager) addWaiting(offset dwarf.Offset, t Type) {
	if manager.waitingDef == nil {
		manager.waitingDef = make(map[dwarf.Offset][]Type)
	}
	manager.waitingDef[offset] = append(manager.waitingDef[offset], t)
}

func (manager *TypeManager) removeWaitingList(offset dwarf.Offset) {
	typeToAdd := manager.getType(offset)
	if list, ok := manager.waitingDef[offset]; ok {
		for _, element := range list {
			switch element.(type) {
			case *TypeDef:
				t := element.(*TypeDef)
				t.Base = typeToAdd
			}
		}
	}
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
	str, err := t.Parse(bytes, manager.Endianess)
	return str, err
}

func (manager *TypeManager) ParseDwarfEntry(entry *dwarf.Entry) error {
	var added bool
	switch entry.Tag {
	case dwarf.TagBaseType:
		base, err := parseBaseEntry(entry)
		if err != nil {
			return err
		}
		manager.addType(entry.Offset, base)
		added = true
	case dwarf.TagTypedef:
		typedef, err := parseTypeDef(entry, manager)
		if err != nil {
			if err == NoAssociatedType {
				return nil
			}
			return err
		}
		manager.addType(entry.Offset, typedef)
		added = true
	}

	if added {
		manager.removeWaitingList(entry.Offset)
	}
	return nil
}
