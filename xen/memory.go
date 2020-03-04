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
void write_memory(void* map, void* buffer, int offset, int length) {
	memcpy((map + offset), buffer, length);
}
*/
import "C"

import (
	"fmt"
	"unsafe"
)

//PageError type representing error codes
type PageError int

const (
	OffsetTooLarge          PageError = 0
	NotEnoughBytes          PageError = 1
	SizeTooLarge            PageError = 2
	MismatchingNoBytesWrite PageError = 3
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
//size the amont of data to be read
func (page *Page) Read(offset, size uint16) ([]byte, error) {
	err := validateOffsetAndSize(offset, size, page.pageSize)
	if err != nil {
		return nil, err
	}

	buffer := C.GoBytes(page.memory, C.int(page.pageSize))
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
	page.start = address
	page.end = address + uint64(page.pageSize)
	page.memory = memory
	return page
}

type Memory struct {
	key      *C.xenforeignmemory_handle
	ctrl     *Xenctrl
	memories map[uint64]unsafe.Pointer
}

func (mem *Memory) Init(ctrl *Xenctrl) error {
	mem.key = C.xenforeignmemory_open(nil, 0)
	mem.ctrl = ctrl
	return nil
}

func (mem *Memory) Close() error {
	C.xenforeignmemory_close(mem.key)
	return nil
}

func (mem *Memory) Map(address uint64, domid uint32, size uint32, vcpu uint32) error {
	if mem.memories == nil {
		mem.memories = make(map[uint64]unsafe.Pointer)
	}
	memory := C.map_memory(mem.key, mem.ctrl.Key(), C.uint32_t(domid), C.int(vcpu), C.ulong(address), C.ulong(size), C.PROT_READ)
	if memory == nil {
		fmt.Println("SOMETHING WENT WRONG")
	}
	mem.memories[address] = memory
	return nil
}

func (mem *Memory) Read(address uint64, size uint32) []byte {
	if mem.memories == nil {
		return nil
	}
	if val, ok := mem.memories[address]; ok {
		buffer := C.GoBytes(val, C.int(size))
		return buffer
	}
	return nil
}

func (mem *Memory) Write(address uint64, bytes []byte) error {
	if v, ok := mem.memories[address]; ok {
		fmt.Println(v)
		//C.write_memory(v, bytes, C.int(len(bytes)))
	}
	return nil
}

func (mem *Memory) UnMap(address uint64, pages uint32) error {
	err := C.xenforeignmemory_unmap(mem.key, mem.memories[address], 1)
	fmt.Println("Error in unmap", err)
	return nil
}
