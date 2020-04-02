package file

import (
	"debug/dwarf"
	"encoding/binary"
	"fmt"
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

func (e ErrorD) Error() string {
	switch e {
	case WrongSize:
		return "Mismatching size"
	}
	return ""
}

type Type struct {
	Size     int64
	Encoding DType
	Name     string
}

func (t *Type) Entry(entry *dwarf.Entry) error {
	if entry.Tag == dwarf.TagBaseType {

		field := entry.AttrField(dwarf.AttrByteSize)
		if field != nil {
			t.Size = field.Val.(int64)
		}
		field = entry.AttrField(dwarf.AttrEncoding)
		if field != nil {
			t.Encoding = DType(field.Val.(int64))
		}
		field = entry.AttrField(dwarf.AttrName)
		if field != nil {
			t.Name = field.Val.(string)
		}
	}
	return nil
}

func (t *Type) Parse(bytes []byte) (string, error) {
	fmt.Printf("TYPE: %+v\n", t)
	fmt.Println("BYTES: ", bytes)
	fmt.Println(Boolean)
	switch t.Encoding {
	case Boolean:
		if int64(len(bytes)) != t.Size {
			return "", WrongSize
		}
		i := int(bytes[0])
		if i == 1 {
			return "true", nil
		} else {
			return "false", nil
		}
	case Sinteger:
		if int64(len(bytes)) != t.Size {
			return "", WrongSize
		}
		integer := int64(binary.BigEndian.Uint64(bytes))
		return fmt.Sprintf("%d", integer), nil
	case Uinteger:
		integer := binary.LittleEndian.Uint64(bytes)
		return fmt.Sprintf("%d", integer), nil

	}
	return "", nil
}
