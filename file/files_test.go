package file

import (
	"fmt"
	"testing"
)

func TestGetAddress(t *testing.T) {
	file := File{Name: "./a.out"}
	err := file.Init()
	if err != nil {
		t.Error(err)
	}

	address := file.GetAddress("test.c", 6)
	if address != 0x659 {
		t.Errorf("Error: expecting 0x64d but got %d\n", address)
	}
}

func TestGetLine(t *testing.T) {
	file := File{Name: "./a.out"}
	err := file.Init()
	if err != nil {
		t.Error(err)
	}

	changed := file.UpdateLine(0x659)
	if !changed {
		t.Error("Error should indicated changed")
	}
	str, line := file.CurrentLine()
	fmt.Printf("File: %s\nLine No: %d\n", str, line)
}
