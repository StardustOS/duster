package file

import "encoding/binary"

type Variable struct {
	Name     string
	typeVar  Type
	location []byte
}

func (variable *Variable) Parse(bytes []byte, endianess binary.ByteOrder) (string, error) {
	return variable.typeVar.Parse(bytes, endianess)
}

func (variable *Variable) Size() int {
	return variable.typeVar.Size()
}

func (variable *Variable) Location() []byte {
	return variable.location
}
