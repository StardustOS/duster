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

type Attribute struct {
	FieldName string
	Offset    int
	base      Type
}

type Struct struct {
	Name       string
	attributes []*Attribute
	needType   map[dwarf.Offset]*Attribute
}

func (s *Struct) AddAtribute(attr *Attribute) {
	s.attributes = append(s.attributes, attr)
}

func (s *Struct) AddNeedType(attr *Attribute, offset dwarf.Offset) {
	if s.needType == nil {
		s.needType = make(map[dwarf.Offset]*Attribute)
	}
	s.needType[offset] = attr
}

func (s *Struct) AddType(offset dwarf.Offset, t Type) {
	if attr, ok := s.needType[offset]; ok {
		attr.base = t
	}
}

func (s *Struct) Size() int {
	size := 0
	for _, attr := range s.attributes {
		size += attr.Offset
	}
	return size
}

func (s *Struct) Parse(bytes []byte, endianess binary.ByteOrder) (string, error) {
	fmt.Println(len(bytes))
	str := "{"
	for _, val := range s.attributes {
		start := val.Offset
		end := val.Offset + val.base.Size()
		attributeData := bytes[start:end]
		out, err := val.base.Parse(attributeData, endianess)
		if err != nil {
			return "", err
		}
		str = fmt.Sprintf("%s %s: %s", str, val.FieldName, out)
	}
	str = fmt.Sprintf("%s }", str)
	return str, nil
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

func parseStruct(entry *dwarf.Entry) (*Struct, error) {
	newStruct := new(Struct)
	field := entry.AttrField(dwarf.AttrName)
	if field == nil {
		return nil, errors.New("no name for struct")
	}
	newStruct.Name = field.Val.(string)
	return newStruct, nil
}

func parseMember(entry *dwarf.Entry, manager *TypeManager) (*Attribute, error) {
	newAttribute := new(Attribute)
	field := entry.AttrField(dwarf.AttrName)
	if field == nil {
		return nil, errors.New("No name for attribute")
	}
	name := field.Val.(string)
	newAttribute.FieldName = name
	field = entry.AttrField(dwarf.AttrType)
	if field == nil {
		return nil, errors.New("No type for attribute")
	}
	offset := field.Val.(dwarf.Offset)
	t := manager.getType(offset)
	if t == nil {
		manager.currentStruct.AddNeedType(newAttribute, offset)
		manager.addWaiting(offset, manager.currentStruct)
	} else {
		newAttribute.base = t
	}

	field = entry.AttrField(dwarf.AttrDataMemberLoc)
	if field == nil {
		return nil, errors.New("No memeber location")
	}
	newAttribute.Offset = int(field.Val.(int64))

	return newAttribute, nil
}

type TypeManager struct {
	Endianess     binary.ByteOrder
	types         map[dwarf.Offset]Type
	waitingDef    map[dwarf.Offset][]Type
	currentStruct *Struct
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
			default:
				element = typeToAdd
			case *TypeDef:
				t := element.(*TypeDef)
				t.Base = typeToAdd
			case *Struct:
				t := element.(*Struct)
				t.AddType(offset, typeToAdd)
			}
		}
		delete(manager.waitingDef, offset)
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
	fmt.Println(t)
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
	case dwarf.TagStructType:
		newStruct, err := parseStruct(entry)
		if err != nil {
			return err
		}
		manager.currentStruct = newStruct
		manager.addType(entry.Offset, newStruct)
		added = true
	case dwarf.TagMember:
		memeber, err := parseMember(entry, manager)
		if err != nil {
			return err
		}
		manager.currentStruct.AddAtribute(memeber)
	}

	if added {
		manager.removeWaitingList(entry.Offset)
	}
	return nil
}
