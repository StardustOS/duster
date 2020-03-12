package xen

/*
#cgo CFLAGS: -g -Wall
#cgo LDFLAGS: -lxenforeignmemory -lxenctrl
#include <xenctrl.h>
#include <string.h>
#include <sys/mman.h>
#include <fcntl.h>
#include <stdlib.h>
#include <xenforeignmemory.h>
void* map_memory(xenforeignmemory_handle*fmem, xc_interface* xch, uint32_t domid, int vcpu, uint64_t vaddr, size_t num_bytes, int perm) {
	const size_t num_pages = (num_bytes + XC_PAGE_SIZE - 1) >> XC_PAGE_SHIFT;
	xen_pfn_t*pages= (xen_pfn_t*) malloc(num_pages* sizeof(xen_pfn_t));
	int* errors = (int*) malloc(num_pages* sizeof(int));
	const xen_pfn_t base_gfn = xc_translate_foreign_address(xch, domid, vcpu, vaddr);
	printf("Addresses after: %lu\n", base_gfn);
	for (size_t i = 0;i < num_pages; ++i)
		pages[i] = base_gfn + i;
	void* mem = xenforeignmemory_map(fmem, domid, perm, num_pages, pages, errors);
	if (!mem) {
		puts("NO MEMORY MAPPED!!!");
	}
	for (size_t i = 0;i < num_pages; ++i) {
		if(errors[i]) {
			puts("SOMETHING WENT WRONG DURING THE MAPPING!");
			return NULL;
		}
	}
	return mem;
}

int unmap(xenforeignmemory_handle* fmem, void* address, unsigned long pages) {
	return xenforeignmemory_unmap(fmem, address, pages);
}
void write_memory(void* map, void* buffer, int offset, int length) {
	//char* b = (char*)map;
	memcpy((map + offset), buffer, length);
}

*/
import "C"

import (
	"errors"
	"fmt"
	"unsafe"
)

//PageError type representing error codes
type PageError int
type MemoryError int

const (
	OffsetTooLarge          PageError   = 0
	NotEnoughBytes          PageError   = 1
	SizeTooLarge            PageError   = 2
	MismatchingNoBytesWrite PageError   = 3
	CouldNotGetMemoryHandle MemoryError = 0
)

func (err PageError) Error() string {
	switch err {
	case OffsetTooLarge:
		return "The offset is larger than the page size would allow for"
	case NotEnoughBytes:
		return "There is not enough bytes in the page left to read."
	case SizeTooLarge:
		return "Trying to read more bytes than available in a page"
	case MismatchingNoBytesWrite:
		return "Mismatching number of bytes in Write (i.e. the byte slice has more or less than what the size var say)"
	}
	return "Unknown error"
}

func (err MemoryError) Error() string {
	switch err {
	case CouldNotGetMemoryHandle:
		return "Error: could not get memory handle. Please try running as sudo and make sure the domain is correct"
	}
	return ""
}

func validateOffsetAndSize(offset, size, pageSize uint16) error {
	if offset > pageSize {
		return OffsetTooLarge
	} else if size > pageSize {
		return SizeTooLarge
	} else if pageSize-offset < size {
		return NotEnoughBytes
	}
	return nil
}

//Page represents a single page in memory for
//VM
type Page struct {
	memory   unsafe.Pointer
	start    uint64
	end      uint64
	pageSize uint16
}

//Range returns the lower and upperbound for addresses in the page
func (page *Page) Range() (uint64, uint64) {
	return page.start, page.end
}

//Read returns a byte array containing the data stored at the page
//offset is the offset we start reading at (i.e. set to zero if we want to read from the beginning)
//size the amount of data to be read
func (page *Page) Read(offset, size uint16) ([]byte, error) {
	err := validateOffsetAndSize(offset, size, page.pageSize)
	if err != nil {
		return nil, err
	}
	fmt.Println("offset", offset)
	buffer := C.GoBytes(page.memory, C.int(page.pageSize))
	fmt.Println("Length of offset", offset+size)
	return buffer[offset:(offset + size)], nil
}

//Write writes the buffer passed at the speficied offset
func (page *Page) Write(offset, size uint16, bytes []byte) error {
	err := validateOffsetAndSize(offset, size, page.pageSize)
	if err != nil {
		return err
	}
	if len(bytes) != int(size) {
		return MismatchingNoBytesWrite
	}
	fmt.Println("Offset before write", offset)
	C.write_memory(page.memory, C.CBytes(bytes), C.int(offset), C.int(size))
	return nil
}

//CalculateOffset works out where in the page the value is stored
func (page *Page) CalculateOffset(address uint64) uint16 {
	offset := address % uint64(page.pageSize)
	return uint16(offset)
}

//CreatePage constructor for the Page struct
func CreatePage(address uint64, memory unsafe.Pointer) *Page {
	page := new(Page)
	page.pageSize = uint16(C.XC_PAGE_SIZE)
	page.start = address - uint64(page.CalculateOffset(address))
	page.end = page.start + uint64(page.pageSize)
	page.memory = memory
	return page
}

type Map struct {
	memories []*Page
	pointer  unsafe.Pointer
}

func (mappedMemory *Map) AddPointer(pointer unsafe.Pointer, size uint32, address uint64) {
	mappedMemory.pointer = pointer
	page := CreatePage(address, pointer)
	mappedMemory.memories = append(mappedMemory.memories, page)
}

func (mappedMemory *Map) NoPages() uint64 {
	return uint64(len(mappedMemory.memories))
}

func (mappedMemory *Map) Destroy() unsafe.Pointer {
	mappedMemory.memories = nil
	pointer := mappedMemory.pointer
	return pointer
}

func (mappedMemory *Map) findPage(address uint64) (*Page, error) {
	for _, page := range mappedMemory.memories {
		lower, upper := page.Range()
		if lower <= address && upper >= address {
			return page, nil
		}
	}
	return nil, errors.New("NOT FOUND")
}

func (mappedMemory *Map) Read(address uint64, size uint16) (bytes []byte, err error) {
	page, err := mappedMemory.findPage(address)
	if err != nil {
		return nil, err
	}
	bytes, err = page.Read(page.CalculateOffset(address), size)
	return
}

func (mappedMemory *Map) Write(address uint64, bytes []byte, size uint16) error {
	page, err := mappedMemory.findPage(address)
	if err != nil {
		return err
	}
	err = page.Write(page.CalculateOffset(address), size, bytes)
	return err
}

func (mappedMemory *Map) StartingAddress() (address uint64) {
	if mappedMemory.memories != nil {
		address, _ = mappedMemory.memories[0].Range()
	}
	return
}

func (mappedMemory *Map) EndingAddress() (address uint64) {
	if mappedMemory.memories != nil {
		_, address = mappedMemory.memories[len(mappedMemory.memories)-1].Range()
	}
	return
}

func (mappedMemory *Map) IsInMap(address uint64) bool {
	startAddress := mappedMemory.StartingAddress()
	endAddress := mappedMemory.EndingAddress()
	fmt.Printf("Start address: %d and end address: %d (address: %d)\n", startAddress, endAddress, address)
	return startAddress <= address && endAddress >= address

}

//Memory is used for interacting with the memory
//of the virtual machine
type Memory struct {
	key  *C.xenforeignmemory_handle
	ctrl *Xenctrl
	maps []Map
}

//Init must be called first gets handles for accessing memories
func (mem *Memory) Init(ctrl *Xenctrl) error {
	mem.key = C.xenforeignmemory_open(nil, 0)
	if mem.key == nil {
		return CouldNotGetMemoryHandle
	}
	mem.ctrl = ctrl
	return nil
}

//Close called when operation are done. None of the
//methods, except Init, will work
func (mem *Memory) Close() error {
	C.xenforeignmemory_close(mem.key)
	return nil
}

//Map - maps an area of memory that allows use to read from it
func (mem *Memory) Map(address uint64, domid uint32, size uint32, vcpu uint32) error {
	memory := C.map_memory(mem.key, mem.ctrl.Key(), C.uint32_t(domid), C.int(vcpu), C.ulong(address), C.ulong(size), C.PROT_READ|C.PROT_WRITE)
	if memory == nil {
		fmt.Println("SOMETHING WENT WRONG")
	}
	newMap := Map{}
	newMap.AddPointer(memory, size, address)

	startNewMap := newMap.StartingAddress()
	for index, currMap := range mem.maps {
		start := currMap.StartingAddress()
		if startNewMap < start {
			mem.maps = append(mem.maps[:index], append([]Map{newMap}, mem.maps[index:]...)...)
			return nil
		}
	}
	mem.maps = append(mem.maps, newMap)
	fmt.Println(mem.maps)
	return nil
}

func (mem *Memory) getMap(address uint64) (Map, int) {
	for index, current := range mem.maps {
		if current.IsInMap(address) {
			return current, index
		}
	}
	return Map{}, -1
}

func (mem *Memory) Read(address uint64, size uint32) ([]byte, error) {
	mapToRead, index := mem.getMap(address)
	if index == -1 {
		return nil, nil
	}
	fmt.Println("size", size)
	bytes, err := mapToRead.Read(address, uint16(size))
	if err != nil {
		return nil, err
	}
	return bytes, nil
}

func (mem *Memory) Write(address uint64, size uint32, bytes []byte) error {
	mapToWrite, index := mem.getMap(address)
	if index == -1 {
		return errors.New("Not found")
	}
	//FIX THIS
	err := mapToWrite.Write(address, bytes, uint16(size))
	return err
}

//UnMap - &safe_place_to_write cleans up the memory once it's has been finished being used
func (mem *Memory) UnMap(address uint64) error {
	mapToRemove, index := mem.getMap(address)
	if index == -1 {
		return errors.New("Not found")
	}
	mem.maps = append(mem.maps[:index], mem.maps[index+1:]...)

	noPages := C.ulong(mapToRemove.NoPages())
	err := C.unmap(mem.key, mapToRemove.Destroy(), noPages) // C.size_t(mapToRemove.NoPages()))

	if err != 0 {
		fmt.Println(err)
		return errors.New("Error occured when unmapping")
	}

	return nil
}
