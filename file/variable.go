package file

import (
	"encoding/binary"
	"fmt"
)

type Variable struct {
	name     string
	typeVar  Type
	location []byte
}

func (variable *Variable) Parse(bytes []byte, endianess binary.ByteOrder) (string, error) {
	return variable.typeVar.Parse(bytes, endianess)
}

func (variable *Variable) Size() int {
	fmt.Println(variable.typeVar)
	return variable.typeVar.Size()
}

func (variable *Variable) Location() []byte {
	return variable.location
}

func (variable *Variable) Type() Type {
	return variable.typeVar
}

func (variable *Variable) Name() string {
	return variable.name
}
