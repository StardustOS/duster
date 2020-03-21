package file

import (
	"debug/elf"
	"reflect"
	"testing"
)

func TestUpdateSymbols(t *testing.T) {
	file, err := elf.Open("testfiles/variable_data")
	if err != nil {
		t.Error(err)
	}
	d, err := file.DWARF()
	if err != nil {
		t.Error(err)
	}
	sym := Symbol{Data: d}
	err = sym.Update(0x63a)
	if err != nil {
		t.Error(err)
	}
	symbols, err := sym.Symbols()
	if err != nil {
		t.Error(err)
	}
	expected := []string{"a", "hello_world", "meh"}
	if !reflect.DeepEqual(symbols, expected) {
		t.Errorf("Expected %v but got %v", expected, symbols)
	}
}
