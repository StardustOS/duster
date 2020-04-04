package file 

import (
	"debug/dwarf"
	"debug/elf"
	"encoding/binary"
	"fmt"
	"math"
	"strings"
	"testing"
)

func setup(filename string, t *testing.T) *dwarf.Reader {
	file, err := elf.Open(filename)
	if err != nil {
		t.Fatalf(err.Error())
	}
	d, err := file.DWARF()
	if err != nil {
		t.Fatalf(err.Error())
	}
	reader := d.Reader()
	return reader
}

type val struct {
	offset dwarf.Offset
	data []byte 
	expected string
}
func TestBasicType(t *testing.T) {
	integerBytes := make([]byte, 8)
	val64 := math.Pow(2, 64)
	val32 := math.Pow(2, 32)
	binary.LittleEndian.PutUint64(integerBytes, uint64(val64))
	fourBytes := make([]byte, 4)
	binary.LittleEndian.PutUint32(integerBytes, uint32(val32))
	neg := make([]byte, 4)
	binary.LittleEndian.PutUint32(neg, uint32(val32-1))
	chars := make([]byte, 2)
	binary.LittleEndian.PutUint16(chars, 254)
	floatBytes := make([]byte, 4)
	i := math.Float32bits(1.3)
	binary.LittleEndian.PutUint32(floatBytes, i)
	k := math.Float64bits(123.121)
	float64Bytes := make([]byte, 8)
	binary.LittleEndian.PutUint64(float64Bytes, k)

	var basicTypes = []val {
		val{offset: 0x00000039, data: integerBytes, expected:fmt.Sprintf("%d", uint64(val64))},
		val{offset: 0x00000040, data: fourBytes, expected: fmt.Sprintf("%d", uint32(val32))},
		val{offset: 0x00000065, data: neg, expected: "-1"},
		val{offset: 0x00000049, data:[]byte{byte(255)}, expected: "255"},
		val{offset: 0x00000091, data: []byte{byte(255)}, expected: "-1"},
		val{offset: 0x0000035b, data: []byte{byte(1)}, expected: "true"},
		val{offset: 0x0000035b, data: []byte{byte(0)}, expected: "false"},
		val{offset: 0x00000354, data: floatBytes, expected: "1.300000"},
		val{offset: 0x00000362, data: float64Bytes, expected: "123.121000"},
	}
	reader := setup("./testfiles/basicType", t)
	var manager TypeManager
	manager.Endianess = binary.LittleEndian
	for entry, _ := reader.Next(); entry != nil; entry, _ = reader.Next() {
		err := manager.ParseDwarfEntry(entry)
		if err != nil {
			t.Fatalf(err.Error())
		}
	}

	for _, v := range basicTypes {
		str, err := manager.ParseBytes(v.offset, v.data)
		if err != nil {
			t.Fatalf(err.Error())
		}
		if strings.Compare(str, v.expected) != 0 {
			t.Errorf("Expected %s but got %s", v.expected, str)
		}
	}
}