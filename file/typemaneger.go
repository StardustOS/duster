package file

import (
	"debug/dwarf"
	"encoding/binary"
	"errors"
	"fmt"
	"math"
	"strings"
)
//Defines the Dwarf base type 
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
	AnnoymousStruct  ErrorD = 2
	NoBoundary       ErrorD = 3
	NeedParseLoction ErrorD = 4
)

func (e ErrorD) Error() string {
	switch e {
	case NoAssociatedType:
		return "No associated Type"
	case AnnoymousStruct:
		return "Anonymous struct"
	}
	return ""
}

//Type outlines the methods for interacting the 
//with types
type Type interface {
	//Size returns the number bytes of required to 
	//represent this type
	Size() int

	//Parse returns a human readable string of that 
	//type (i.e. take the raw bytes passed and convert 
	// them into something readable) 
	Parse([]byte, binary.ByteOrder) (string, error)
}

//Array represents the array type (i.e. char[] or int[])
//NOTE: this not represent array of the type char* or int*
type Array struct {
	typeArray Type
	noElement int
	Location  []byte
}

//SetSize is used to set the number of elements
//(only used when a non-constant value is used in the array initilisation)
func (arr *Array) SetSize(size int) {
	arr.noElement = size
}

//Size returns the total number of bytes used to represent the array
//(i.e. number of elements multipled by the size of the type)
func (arr *Array) Size() int {
	return arr.noElement * (arr.typeArray.Size() + 1)
}

//Parse returns a human readable string of the array
func (arr *Array) Parse(bytes []byte, endianess binary.ByteOrder) (string, error) {
	if arr.Location != nil {
		return "", NeedParseLoction
	}
	str := ""
	for start := 0; start < len(bytes); start += arr.typeArray.Size() {
		end := start + arr.typeArray.Size()
		element := bytes[start:end]
		strElement, err := arr.typeArray.Parse(element, endianess)
		if err != nil {
			return "", err
		}
		str = fmt.Sprintf("%s%s ", str, strElement)
	}
	str = strings.TrimSpace(str)
	return str, nil
}

//parseArrayEntry - parses the array dwarf entry and returns an Array
func parseArrayEntry(entry *dwarf.Entry, manager *TypeManager) (*Array, error) {
	arr := new(Array)
	field := entry.AttrField(dwarf.AttrType)
	if field == nil {
		return nil, errors.New("No type")
	}
	offset := field.Val.(dwarf.Offset)
	t := manager.getType(offset)
	if t == nil {
		manager.addWaiting(offset, arr)
	} else {
		arr.typeArray = t
	}
	return arr, nil
}

//parseArrayRange parses the array range from the dwarf
//and set it in the struct
func parseArrayRange(entry *dwarf.Entry, arr *Array) error {
	field := entry.AttrField(dwarf.AttrUpperBound)
	if field == nil {
		return NoBoundary
	}
	upperBound, ok := field.Val.(int64)
	if ok {
		arr.noElement = int(upperBound)
	} else {
		location := field.Val.([]byte)
		arr.Location = location
	}
	return nil
}

//Pointer represents the pointer in a C program
type Pointer struct {
	size          int
	typeOfPointer Type
}

//Size returns the size of the Pointer 
//(not the size of the content)
func (p *Pointer) Size() int {
	return p.size
}

//Parse returns a human string for the pointer 
func (p *Pointer) Parse(bytes []byte, endianess binary.ByteOrder) (string, error) {
	address := parseUinteger(bytes, endianess)
	var name string
	switch p.typeOfPointer.(type) {
	default:
		name = "void"
	case *BaseType:
		val := p.typeOfPointer.(*BaseType)
		name = val.Name
	case *Struct:
		val := p.typeOfPointer.(*Struct)
		name = val.Name
	case *TypeDef:
		val := p.typeOfPointer.(*TypeDef)
		name = val.Name
	case *VolatileType:
		name = "volatile"
	}
	return fmt.Sprintf("(%s*) 0x%x", name, address), nil
}

//Type returns the type of a pointer
func (p *Pointer) Type() Type {
	return p.typeOfPointer
}

//TypeDef represents the typedef type in C
type TypeDef struct {
	//Name represents the name of the typedef (e.g. for size_t)
	Name string
	//Base represents the type used to define the typedef (i.e. typedef uint_t unsigned int)
	Base Type
}

//Size returns the size of the type
func (t *TypeDef) Size() int {
	return t.Base.Size()
}

//Parse returns a human readable of the data in bytes
func (t *TypeDef) Parse(bytes []byte, endianess binary.ByteOrder) (string, error) {
	return t.Base.Parse(bytes, endianess)
}

//ComplexType represents the methods used to interact with complex types such
//as structs or unions
type ComplexType interface {
	//Must include all the methods of primitive type
	Type

	//AddAtribute method for adding attribute to the type.
	//It should be noted that how the attribute is interpreted depends
	//on the implementation and is not something the interface prescribes. 
	AddAtribute(*Attribute)

	//AddNeedType method for saying that attrubute is waiting for a type 
	//definition
	AddNeedType(*Attribute, dwarf.Offset)

	//AddType method for adding type definition to attributes missing said
	//type
	AddType(dwarf.Offset, Type)
}

//Attribute represents an attribute in a struct
type Attribute struct {
	//FieldName the name of the attribute
	FieldName string
	//Offset in the buffer where this attribute start
	Offset    int
	//The type of the attribute
	base      Type
}

//Struct represents a struct time in C
type Struct struct {
	//Name of the struct
	Name       string
	attributes []*Attribute
	needType   map[dwarf.Offset][]*Attribute
}

//AddAtribute adds an attribute to the struct
func (s *Struct) AddAtribute(attr *Attribute) {
	s.attributes = append(s.attributes, attr)
}

//AddNeedType - in DWARF types can be referred to before they are
//actually defined. Hence this function is used to indicate that an 
//attribute is waiting for type definition to be defined
func (s *Struct) AddNeedType(attr *Attribute, offset dwarf.Offset) {
	if s.needType == nil {
		s.needType = make(map[dwarf.Offset][]*Attribute)
	}
	s.needType[offset] = append(s.needType[offset], attr)
}

//AddType add the missing type defintion to an attribute
func (s *Struct) AddType(offset dwarf.Offset, t Type) {
	if attrs, ok := s.needType[offset]; ok {
		for _, attr := range attrs {
			attr.base = t
		}
	}
}

//Size calculates the size in bytes. Note the way DWARF
//encodes this information makes sure the calculation accounts
//for the offset due to C's alignment rules
func (s *Struct) Size() int {
	size := 0
	for _, attr := range s.attributes {
		size += attr.Offset
	}
	return size
}

//Parse returns a human readable string of struct
func (s *Struct) Parse(bytes []byte, endianess binary.ByteOrder) (string, error) {
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

//Union represents unions in C
type Union struct {
	Name       string
	attributes []*Attribute
	needType   map[dwarf.Offset][]*Attribute
}

//AddAtribute add an attribute (i.e. a potential of intrepreting the data)
func (s *Union) AddAtribute(attr *Attribute) {
	s.attributes = append(s.attributes, attr)
}

//AddNeedType works as stated in the interface comment
func (s *Union) AddNeedType(attr *Attribute, offset dwarf.Offset) {
	if s.needType == nil {
		s.needType = make(map[dwarf.Offset][]*Attribute)
	}
	s.needType[offset] = append(s.needType[offset], attr)
}

//AddNeedType works as stated in the interface comment
func (s *Union) AddType(offset dwarf.Offset, t Type) {
	if attrs, ok := s.needType[offset]; ok {
		for _, attr := range attrs {
			attr.base = t
		}
	}
}

//Size calculates the size of the union (i.e. returns the largest
//size of the potential types it can hold)
func (union *Union) Size() int {
	var largest int
	for _, attr := range union.attributes {
		if largest < attr.base.Size() {
			largest = attr.base.Size()
		}
	}
	return largest
}

//Parse returns a human readable string in the format {a1: v1, a2: v2} where 
//ax represents a potential way of interpreting the data
func (union *Union) Parse(bytes []byte, endianess binary.ByteOrder) (string, error) {
	str := "{"
	for _, attr := range union.attributes {
		val, err := attr.base.Parse(bytes[:attr.base.Size()], endianess)
		if err != nil {
			return "", err
		}
		str = fmt.Sprintf("%s %s : %s", str, attr.FieldName, val)
	}
	return fmt.Sprintf("%s }", str), nil
}

//parseInteger is helper function for parsing signed integers
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

//parseInteger is helper function for parsing unsigned integers
//note we cannot just use the above functionality and cast to uint64
//as the conversion doesn't correctly overflow when required.
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


//BaseType represent the basic types defined by the DWARF file format. These 
//are used to implement every other type in this file
type BaseType struct {
	size     int
	Encoding DType
	Name     string
}

//Parse returns a human readable string of the bytes interpreted as that type
func (t *BaseType) Parse(bytes []byte, endianess binary.ByteOrder) (string, error) {
	if len(bytes) != t.size {
		return "", fmt.Errorf("Error: type %s expects %d bytes but got %d", t.Name, t.Size(), len(bytes))
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

//Size returns the size of the type
func (t *BaseType) Size() int {
	return t.size
}

//parseBaseEntry return a parse BaseType from the dwarf entry
//Note: the fields specifying the encoding, size and name are required
func parseBaseEntry(entry *dwarf.Entry) (*BaseType, error) {
	base := new(BaseType)
	field := entry.AttrField(dwarf.AttrByteSize)
	if field == nil {
		return nil, errors.New("Error: no bytes size for the base types")
	}
	base.size = int(field.Val.(int64))
	field = entry.AttrField(dwarf.AttrEncoding)
	if field == nil {
		return nil, errors.New("Error: no encoding field")
	}
	base.Encoding = DType(field.Val.(int64))
	field = entry.AttrField(dwarf.AttrName)
	if field == nil {
		return nil, errors.New("Error: no type name")
	}
	base.Name = field.Val.(string)
	return base, nil
}

//parseTypeDef parses the typedef from the DWARF
func parseTypeDef(entry *dwarf.Entry, manager *TypeManager) (*TypeDef, error) {
	typedef := new(TypeDef)
	field := entry.AttrField(dwarf.AttrName)
	if field == nil {
		return nil, errors.New("Error: no attribute name in for the typedef")
	}
	typedef.Name = field.Val.(string)

	field = entry.AttrField(dwarf.AttrType)
	if field == nil {
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

//parses the union from the dwarf
func parseUnion(entry *dwarf.Entry) (*Union, error) {
	newUnion := new(Union)
	field := entry.AttrField(dwarf.AttrName)
	if field == nil {
		return nil, errors.New("Error: Union has no name")
	}
	newUnion.Name = field.Val.(string)
	return newUnion, nil
}

//parses the struct from the DWARF
func parseStruct(entry *dwarf.Entry) (*Struct, error) {
	newStruct := new(Struct)
	field := entry.AttrField(dwarf.AttrName)
	if field == nil {
		return nil, AnnoymousStruct
	}
	newStruct.Name = field.Val.(string)
	return newStruct, nil
}

//parses the member or attribute of a struct
func parseMember(entry *dwarf.Entry, manager *TypeManager) (*Attribute, error) {
	newAttribute := new(Attribute)
	field := entry.AttrField(dwarf.AttrName)
	if field == nil {
		return nil, nil //errors.New("No name for attribute")
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
		manager.current.AddNeedType(newAttribute, offset)
		manager.addWaiting(offset, manager.current)
	} else {
		newAttribute.base = t
	}

	field = entry.AttrField(dwarf.AttrDataMemberLoc)
	if field == nil {
		return newAttribute, nil //errors.New("No memeber location")
	}
	newAttribute.Offset = int(field.Val.(int64))

	return newAttribute, nil
}

//Parses a pointer 
func parsePointer(entry *dwarf.Entry, manager *TypeManager) (*Pointer, error) {
	pointer := new(Pointer)
	field := entry.AttrField(dwarf.AttrByteSize)
	if field == nil {
		return nil, errors.New("No byte size attribute")
	}
	pointer.size = int(field.Val.(int64))
	field = entry.AttrField(dwarf.AttrType)
	if field == nil {
		return pointer, nil
	}
	offset := field.Val.(dwarf.Offset)
	t := manager.getType(offset)
	if t == nil {
		manager.addWaiting(offset, pointer)
	} else {
		pointer.typeOfPointer = t
	}
	return pointer, nil
}

//TypeManager handles the parsing and interpretation of the 
//types
type TypeManager struct {
	Endianess    binary.ByteOrder
	types        map[dwarf.Offset]Type
	waitingDef   map[dwarf.Offset][]Type
	current      ComplexType
	currentArray *Array
}

//addWaiting - waiting list for any type that needs another type to be
//defined before it is fully parsed (e.g. typedef not having its base type
//defined yet)
func (manager *TypeManager) addWaiting(offset dwarf.Offset, t Type) {
	if manager.waitingDef == nil {
		manager.waitingDef = make(map[dwarf.Offset][]Type)
	}
	manager.waitingDef[offset] = append(manager.waitingDef[offset], t)
}

//removeWaitingList - remove types from the waiting list. Please note
//for this function to work the newly parsed type must have been added 
//via the addType method for this to work.
func (manager *TypeManager) removeWaitingList(offset dwarf.Offset) {
	typeToAdd := manager.getType(offset)
	if list, ok := manager.waitingDef[offset]; ok {
		for _, element := range list {
			switch element.(type) {
			default:
				return
			case *TypeDef:
				t := element.(*TypeDef)
				t.Base = typeToAdd
			case *Struct:
				t := element.(*Struct)
				t.AddType(offset, typeToAdd)
			case *Pointer:
				t := element.(*Pointer)
				t.typeOfPointer = typeToAdd
			case *Array:
				t := element.(*Array)
				t.typeArray = typeToAdd
			case *VolatileType:
				t := element.(*VolatileType)
				t.t = typeToAdd
			}
		}
		delete(manager.waitingDef, offset)
	}
}

//addType - adds a newly parsed type to the struct. Please note
//you can add a new type and place it into the waiting list
func (manager *TypeManager) addType(offset dwarf.Offset, t Type) {
	if manager.types == nil {
		manager.types = make(map[dwarf.Offset]Type)
	}
	manager.types[offset] = t
}

//Helper function for abstracting over the simple map used 
func (manager *TypeManager) getType(offset dwarf.Offset) Type {
	return manager.types[offset]
}

//Size returns the size based off the dwarf.Offset 
//(which is unique to each type)
func (manager *TypeManager) Size(offset dwarf.Offset) int {
	t := manager.types[offset]
	if t != nil {
		return t.Size()
	}
	return 0
}

//ParseBytes takes raw memory bytes and return a human readable string of that type
func (manager *TypeManager) ParseBytes(offset dwarf.Offset, bytes []byte) (string, error) {
	t := manager.getType(offset)
	str, err := t.Parse(bytes, manager.Endianess)
	return str, err
}

//VolatileType represents any type that has been 
//defined has volatile int c;
type VolatileType struct {
	t Type
}

//Size the size of the volatile 
func (v *VolatileType) Size() int {
	return v.t.Size()
}

//Parse parses the volatile type
func (v *VolatileType) Parse(bytes []byte, endianess binary.ByteOrder) (string, error) {
	return v.t.Parse(bytes, endianess)
}

//parses the dwarf entry for volatile type
func parseVolatile(entry *dwarf.Entry, manager *TypeManager) (*VolatileType, error) {
	volatile := new(VolatileType)
	field := entry.AttrField(dwarf.AttrType)
	if field == nil {
		return nil, nil
	}
	offset := field.Val.(dwarf.Offset)
	t := manager.getType(offset)
	if t == nil {
		manager.addWaiting(offset, volatile)
	} else {
		volatile.t = t
	}
	return volatile, nil
}

//ConstType represents the type defined as 
//const int k;
type ConstType struct {
	t Type
}

//Size return the number of bytes to represent the type
func (c *ConstType) Size() int {
	return c.t.Size()
}

//Parse returns a pretty string of the value
func (c *ConstType) Parse(bytes []byte, endianess binary.ByteOrder) (string, error) {
	return c.t.Parse(bytes, endianess)
}

//Parse constant type
func parseConst(entry *dwarf.Entry, manager *TypeManager) (*ConstType, error) {
	constant := new(ConstType)
	field := entry.AttrField(dwarf.AttrType)
	if field == nil {
		return nil, nil
	}
	offset := field.Val.(dwarf.Offset)
	t := manager.getType(offset)
	if t == nil {
		manager.addWaiting(offset, constant)
	} else {
		constant.t = t
	}
	return constant, nil
}

//ParseDwarfEntry parses a dwarf entry and adds it the typemanager struct
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
			if err == AnnoymousStruct {
				return nil
			}
			return err
		}
		manager.current = newStruct
		manager.addType(entry.Offset, newStruct)
		added = true
	case dwarf.TagMember:
		memeber, err := parseMember(entry, manager)
		if err != nil {
			return err
		}
		manager.current.AddAtribute(memeber)
	case dwarf.TagPointerType:
		pointer, err := parsePointer(entry, manager)
		if err != nil {
			return err
		}
		manager.addType(entry.Offset, pointer)
		added = true
	case dwarf.TagArrayType:
		arr, err := parseArrayEntry(entry, manager)
		if err != nil {
			return err
		}
		manager.addType(entry.Offset, arr)
		added = true
		manager.currentArray = arr
	case dwarf.TagSubrangeType:
		if manager.currentArray == nil {
			return nil
		}
		err := parseArrayRange(entry, manager.currentArray)
		if err != nil && err != NoBoundary {
			return err
		}
	case dwarf.TagUnionType:
		union, err := parseUnion(entry)
		if err != nil {
			return nil
		}
		manager.current = union
		manager.addType(entry.Offset, union)
		added = true
	case dwarf.TagVolatileType:
		volatile, err := parseVolatile(entry, manager)
		if err != nil {
			return err
		}
		manager.addType(entry.Offset, volatile)
		added = true
	case dwarf.TagConstType:
		constant, err := parseConst(entry, manager)
		if err != nil {
			return err
		}
		manager.addType(entry.Offset, constant)
		added = true
	}

	if added {
		manager.removeWaitingList(entry.Offset)
	}
	return nil
}
