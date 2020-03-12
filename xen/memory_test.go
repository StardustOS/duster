package xen

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"os/exec"
	"reflect"
	"strconv"
	"strings"
	"testing"
	"unsafe"
)

const (
	startingAddress = 5034
)

func testSetup() int {
	cmd := exec.Command("./start.sh")
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Run()
	output := out.String()
	rows := strings.Split(output, "\n")
	for _, row := range rows {
		if strings.Contains(row, "stardust") {
			values := strings.Fields(row)
			id := values[1]
			num, _ := strconv.Atoi(id)
			return num
		}
	}
	return 0
}

func testTeardown() {
	cmd := exec.Command("xl", "destroy", "stardust")
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Run()
}

func CreateDummyBuffer() unsafe.Pointer {
	buffer := make([]byte, 4096)
	return unsafe.Pointer(&buffer[0])
}

func TestCreatePage(t *testing.T) {
	dummyBuffer := CreateDummyBuffer()
	page := CreatePage(startingAddress, dummyBuffer)
	if page == nil {
		t.Error("Error: page was nil when should be non-nil value")
		return
	}
	start, end := page.Range()
	expectedStart := uint64(startingAddress - (startingAddress % 4096))
	expectedEnd := uint64(expectedStart + 4096)
	fmt.Printf("Expected %d\n", expectedEnd)
	if start != expectedStart {
		t.Errorf("Error: starting address should be %d not %d", expectedStart, start)
	}

	if end != expectedEnd {
		t.Errorf("Error: ending address should be %d not %d", expectedEnd, end)
	}
}

func TestCalculateOffset(t *testing.T) {
	dummyBuffer := CreateDummyBuffer()
	page := CreatePage(startingAddress, dummyBuffer)
	offset := page.CalculateOffset(884816)
	if offset != 80 {
		t.Errorf("Error: offset should equal %d not %d", 80, offset)
	}
}

func TestReadWrite(t *testing.T) {
	toWrite := []byte{1, 2, 3, 4}
	dummyBuffer := CreateDummyBuffer()
	page := CreatePage(startingAddress, dummyBuffer)
	addressWrite := uint64(startingAddress + 80)
	err := page.Write(page.CalculateOffset(addressWrite), 4, toWrite)
	if err != nil {
		t.Error(err)
	}
	bytes, err := page.Read(page.CalculateOffset(addressWrite), 4)
	if err != nil {
		t.Error(err)
	}

	if !reflect.DeepEqual(bytes, toWrite) {
		t.Errorf("Error: %+v should match %+v", bytes, toWrite)
	}
}

type Args struct {
	Offset   uint16
	Size     uint16
	Expected error
}

var values = []Args{
	Args{Offset: 6100, Size: 30, Expected: OffsetTooLarge},
	Args{Offset: 4000, Size: 101, Expected: NotEnoughBytes},
	Args{Offset: 0, Size: 4097, Expected: SizeTooLarge},
}

func TestReadWriteError(t *testing.T) {
	for _, val := range values {
		t.Run("Testing", func(t *testing.T) {
			dummyBuffer := CreateDummyBuffer()
			page := CreatePage(startingAddress, dummyBuffer)
			bytes, err := page.Read(val.Offset, val.Size)
			if bytes != nil {
				t.Error("Error: bytes should be nil")
			}
			if err != val.Expected {
				t.Errorf("Error: error produced should be %s not %s", val.Expected, err)
			}
			page = CreatePage(startingAddress, dummyBuffer)
			err = page.Write(val.Offset, val.Size, nil)

			if err != val.Expected {
				t.Errorf("Error: error produced should be %s not %s", val.Expected, err)
			}

		})
	}
}

func TestInitWorks(t *testing.T) {
	testSetup()
	memory := Memory{}
	err := memory.Init(nil)
	if err != nil {
		t.Error(err)
	}
	testTeardown()
}

func TestMapMemory(t *testing.T) {
	id := uint32(testSetup())
	memory := &Memory{}
	domain := &Xenctrl{}
	err := memory.Init(domain)
	if err != nil {
		t.Error(err)
	}
	err = domain.Init()
	if err != nil {
		t.Error(err)
	}
	err = memory.Map(1900632, id, 8, 0)
	if err != nil {
		t.Error("Error with mapping", err)
	}
	err = memory.UnMap(1900632)
	if err != nil {
		t.Errorf("Error with unmapping %s", err)
	}
	testTeardown()
}

func TestReadMemory(t *testing.T) {
	id := uint32(testSetup())
	memory := &Memory{}
	domain := &Xenctrl{}

	err := memory.Init(domain)
	if err != nil {
		t.Error(err)
	}

	err = domain.Init()
	if err != nil {
		t.Error(err)
	}

	err = memory.Map(884816, id, 8, 0)
	if err != nil {
		t.Error("Error with mapping", err)
	}
	bytes, err := memory.Read(884816, 8)
	if err != nil {
		t.Error(err)
		return
	}
	fmt.Println(bytes)
	data := binary.LittleEndian.Uint64(bytes)
	fmt.Println(data)
	if data != 131072 {
		t.Error("Read the wrong value from memory")
	}

	err = memory.UnMap(884816)
	if err != nil {
		t.Errorf("Error with unmapping %s", err)
	}
	testTeardown()
}

func TestReadWriteMemory(t *testing.T) {
	id := uint32(testSetup())
	memory := &Memory{}
	domain := &Xenctrl{}

	err := memory.Init(domain)
	if err != nil {
		t.Error(err)
	}

	err = domain.Init()
	if err != nil {
		t.Error(err)
	}

	err = memory.Map(884848, id, 8, 0)
	if err != nil {
		t.Error("Error with mapping", err)
	}
	bs := make([]byte, 8)
	binary.LittleEndian.PutUint64(bs, 884848)
	err = memory.Write(884848, 8, bs)
	if err != nil {
		t.Error(err)
		return
	}

	bytes, err := memory.Read(884848, 8)
	if err != nil {
		t.Error(err)
		return
	}
	data := binary.LittleEndian.Uint64(bytes)

	fmt.Println(data)
	if data != 884848 {
		t.Error("Read the wrong value from memory")
	}

	err = memory.UnMap(884848)
	if err != nil {
		t.Errorf("Error with unmapping %s", err)
	}
	testTeardown()
}
