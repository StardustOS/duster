package xen

import (
	"reflect"
	"testing"
	"unsafe"
)

const (
	startingAddress = 5034
)

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
	if start != startingAddress {
		t.Errorf("Error: starting address should be %d not %d", startingAddress, start)
	}

	if end != startingAddress+4096 {
		t.Errorf("Error: ending address should be %d not %d", startingAddress+4096, end)
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
