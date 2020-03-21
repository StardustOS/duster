package file

import (
	"fmt"
	"path/filepath"
	"strings"
	"testing"
)

func TestGetAddress(t *testing.T) {
	file := File{Name: "testfiles/test"}
	err := file.Init()
	if err != nil {
		t.Error(err)
	}

	address := file.GetAddress("test.c", 6)
	if address != 0x654 {
		t.Errorf("Error: expecting 0x64d but got %d\n", address)
	}
}

func TestGetLine(t *testing.T) {
	file := File{Name: "testfiles/test"}
	err := file.Init()
	if err != nil {
		t.Error(err)
	}

	changed := file.UpdateLine(0x660)
	if !changed {
		t.Error("Error should indicated changed")
	}

	dir, err := filepath.Abs(filepath.Dir("testfiles/test.c"))
	if err != nil {
		t.Error(err)
	}
	expected := fmt.Sprintf("%s/test.c", dir)
	str, line := file.CurrentLine()

	if strings.Compare(str, expected) != 0 {
		t.Errorf("Expected file to be %s not %s", expected, str)
	}

	if line != 7 {
		t.Errorf("Address occurs in line 7 not %d", line)
	}

}
