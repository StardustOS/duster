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
	offset   dwarf.Offset
	data     []byte
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

	var basicTypes = []val{
		val{offset: 0x00000038, data: integerBytes, expected: fmt.Sprintf("%d", uint64(val64))},
		val{offset: 0x0000004d, data: fourBytes, expected: fmt.Sprintf("%d", uint32(val32))},
		val{offset: 0x00000062, data: neg, expected: "-1"},
		val{offset: 0x0000008e, data: []byte{byte(255)}, expected: "-1"},
		val{offset: 0x00000371, data: []byte{byte(1)}, expected: "true"},
		val{offset: 0x00000371, data: []byte{byte(0)}, expected: "false"},
		val{offset: 0x0000036a, data: floatBytes, expected: "1.300000"},
		val{offset: 0x00000378, data: float64Bytes, expected: "123.121000"},
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

func TestTypedef(t *testing.T) {
	integerBytes := make([]byte, 8)
	val64 := math.Pow(2, 64)
	binary.LittleEndian.PutUint64(integerBytes, uint64(val64))

	bytes := make([]byte, 4)
	val32 := uint32(math.Pow(2, 32)) - 1
	binary.LittleEndian.PutUint32(bytes, val32)
	var typedefed = []val{
		val{offset: 0x0000002d, data: integerBytes, expected: fmt.Sprintf("%d", uint64(val64))}, //Case where only one link
		val{offset: 0x00000353, data: bytes, expected: fmt.Sprintf("%d", val32)},                //Case where multiple links back to original type
	}

	reader := setup("./testfiles/typedef", t)
	var manager TypeManager
	manager.Endianess = binary.LittleEndian
	for entry, _ := reader.Next(); entry != nil; entry, _ = reader.Next() {
		err := manager.ParseDwarfEntry(entry)
		if err != nil {
			t.Fatalf(err.Error())
		}
	}

	for _, v := range typedefed {
		str, err := manager.ParseBytes(v.offset, v.data)
		if err != nil {
			t.Fatalf(err.Error())
		}
		if strings.Compare(str, v.expected) != 0 {
			t.Errorf("Expected %s but got %s", v.expected, str)
		}
	}
}

func TestStruct(t *testing.T) {
	intBytes := make([]byte, 4)
	val32 := uint32(math.Pow(2, 32)) - 1
	binary.LittleEndian.PutUint32(intBytes, val32)
	charBytes := make([]byte, 4)
	charBytes[0] = byte(255)
	floatBytes := make([]byte, 4)
	i := math.Float32bits(1.3)
	binary.LittleEndian.PutUint32(floatBytes, i)
	var bytes []byte
	bytes = append(bytes, intBytes...)
	bytes = append(bytes, charBytes...)
	bytes = append(bytes, floatBytes...)

	intBytes64 := make([]byte, 8)
	val64 := uint64(math.Pow(2, 64)) - 1
	binary.LittleEndian.PutUint64(intBytes64, val64)
	charBytes2 := []byte{byte(255)}

	var bytes2 []byte

	bytes2 = append(bytes2, intBytes64...)
	bytes2 = append(bytes2, bytes...)
	bytes2 = append(bytes2, charBytes2...)
	var tests = []val{
		val{offset: 0x000002ff, data: bytes, expected: "{ v: -1 c: -1 f: 1.300000 }"},
		val{offset: 0x0000032f, data: bytes2, expected: "{ m: 9223372036854775807 meh: { v: -1 c: -1 f: 1.300000 } b: -1 }"},
	}

	reader := setup("./testfiles/structs", t)
	var manager TypeManager
	manager.Endianess = binary.LittleEndian
	for entry, _ := reader.Next(); entry != nil; entry, _ = reader.Next() {
		err := manager.ParseDwarfEntry(entry)
		if err != nil {
			t.Fatalf(err.Error())
		}
	}
	for _, v := range tests {
		str, err := manager.ParseBytes(v.offset, v.data)
		if err != nil {
			t.Fatalf(err.Error())
		}
		if strings.Compare(str, v.expected) != 0 {
			t.Errorf("Expected %s but got %s", v.expected, str)
		}
	}
}

func TestPointer(t *testing.T) {
	pointer := make([]byte, 8)
	binary.LittleEndian.PutUint32(pointer, 0x10294)
	var tests = []val{
		val{offset: 0x0000035c, data: pointer, expected: "(int*) 0x10294"},
		val{offset: 0x00000362, data: pointer, expected: "(k*) 0x10294"},
	}

	reader := setup("./testfiles/pointer", t)
	var manager TypeManager
	manager.Endianess = binary.LittleEndian
	for entry, _ := reader.Next(); entry != nil; entry, _ = reader.Next() {
		err := manager.ParseDwarfEntry(entry)
		if err != nil {
			t.Fatalf(err.Error())
		}
	}
	for _, v := range tests {
		str, err := manager.ParseBytes(v.offset, v.data)
		if err != nil {
			t.Fatalf(err.Error())
		}
		if strings.Compare(str, v.expected) != 0 {
			t.Errorf("Expected %s but got %s", v.expected, str)
		}
	}

}


func TestArray(t *testing.T) {
	expected := "Hola el mundo"
	array := []byte(expected)
	e := fmt.Sprintf("%v", array)
	expected = e[1 : len(e)-1]
	var tests = []val{
		val{offset: 0x0000034d, data: array, expected: expected},
	}

	reader := setup("./testfiles/arrays", t)
	var manager TypeManager
	manager.Endianess = binary.LittleEndian
	for entry, _ := reader.Next(); entry != nil; entry, _ = reader.Next() {
		err := manager.ParseDwarfEntry(entry)
		if err != nil {
			t.Fatalf(err.Error())
		}
	}
	for _, v := range tests {
		str, err := manager.ParseBytes(v.offset, v.data)
		if err != nil {
			t.Fatalf(err.Error())
		}
		if strings.Compare(str, v.expected) != 0 {
			t.Errorf("Expected %s but got %s", v.expected, str)
		}
	}

}

// func TestDynamicArray(t *testing.T) {
// 	var tests = []val{
// 		val{offset: 0x00000341},
// 	}

// 	reader := setup("./testfiles/dynamic-array", t)
// 	var manager TypeManager
// 	manager.Endianess = binary.LittleEndian
// 	for entry, _ := reader.Next(); entry != nil; entry, _ = reader.Next() {
// 		err := manager.ParseDwarfEntry(entry)
// 		if err != nil {
// 			t.Fatalf(err.Error())
// 		}
// 	}
// 	for _, v := range tests {
// 		_, err := manager.ParseBytes(v.offset, v.data)
// 		if err != NeedParseLoction {
// 			t.Fatalf("Need to returns error (we don't know the size of the data)")
// 		}
// 	}

// }

func TestRecusriveStruct(t *testing.T) {
	p1 := make([]byte, 8)
	p2 := make([]byte, 8)

	binary.LittleEndian.PutUint64(p1, 0x10294)
	binary.LittleEndian.PutUint64(p2, 0x102945)
	data := append(p1, p2...)
	var tests = []val{
		val{offset: 0x000002ff, data: data, expected: "{ next: (list_head*) 0x10294 prev: (list_head*) 0x102945 }"},
	}
	reader := setup("./testfiles/recursive_struct", t)
	var manager TypeManager
	manager.Endianess = binary.LittleEndian
	for entry, _ := reader.Next(); entry != nil; entry, _ = reader.Next() {
		err := manager.ParseDwarfEntry(entry)
		if err != nil {
			t.Fatalf(err.Error())
		}
	}
	for _, v := range tests {
		str, _ := manager.ParseBytes(v.offset, v.data)
		if strings.Compare(str, v.expected) != 0 {
			t.Errorf("Expected %s but got %s", str, v.expected)
		}
	}

}

func TestUnion(t *testing.T) {
	data := make([]byte, 8)

	binary.LittleEndian.PutUint64(data, 254)

	var tests = []val{
		val{offset: 0x000002ff, data: data, expected: "{ hello : 254 c : -2 }"},
	}

	reader := setup("./testfiles/unions", t)
	var manager TypeManager
	manager.Endianess = binary.LittleEndian
	for entry, _ := reader.Next(); entry != nil; entry, _ = reader.Next() {
		err := manager.ParseDwarfEntry(entry)
		if err != nil {
			t.Fatalf(err.Error())
		}
	}
	for _, v := range tests {
		str, _ := manager.ParseBytes(v.offset, v.data)
		if strings.Compare(str, v.expected) != 0 {
			t.Errorf("Expected %s but got %s", str, v.expected)
		}
	}

}

func TestVolatile(t *testing.T) {
	data := make([]byte, 4)
	binary.LittleEndian.PutUint32(data, 100)
	var tests = []val{
		val{offset: 0x00000069, data: data, expected: "100"},
	}
	reader := setup("./testfiles/volatile", t)
	var manager TypeManager
	manager.Endianess = binary.LittleEndian
	for entry, _ := reader.Next(); entry != nil; entry, _ = reader.Next() {
		err := manager.ParseDwarfEntry(entry)
		if err != nil {
			t.Fatalf(err.Error())
		}
	}
	for _, v := range tests {
		str, _ := manager.ParseBytes(v.offset, v.data)
		if strings.Compare(str, v.expected) != 0 {
			t.Errorf("Expected %s but got %s", str, v.expected)
		}
	}

}

func TestConst(t *testing.T) {
	data := make([]byte, 4)
	binary.LittleEndian.PutUint32(data, 100)
	var tests = []val{
		val{offset: 0x00000069, data: data, expected: "100"},
	}
	reader := setup("./testfiles/constant", t)
	var manager TypeManager
	manager.Endianess = binary.LittleEndian
	for entry, _ := reader.Next(); entry != nil; entry, _ = reader.Next() {
		err := manager.ParseDwarfEntry(entry)
		if err != nil {
			t.Fatalf(err.Error())
		}
	}
	for _, v := range tests {
		str, _ := manager.ParseBytes(v.offset, v.data)
		if strings.Compare(str, v.expected) != 0 {
			t.Errorf("Expected %s but got %s", str, v.expected)
		}
	}

}


func TestStatic(t *testing.T) {
	binary.LittleEndian.PutUint32(data, 100)
	var tests = []val{
		val{offset: 0x00000069, data: data, expected: "100"},
	}

}