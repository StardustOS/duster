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
	printf("DOMAIN ID: %d\n", domid);
	printf("VCPU: %d\n", vcpu);
	printf("vaddr: %lu\n", vaddr);
	printf("num_bytes: %lu\n", num_bytes);
	printf("perm: %d\n", perm);
	const size_t num_pages = (num_bytes + XC_PAGE_SIZE - 1) >> XC_PAGE_SHIFT;
	printf("%lu\n", num_pages);
	printf("%lu\n", num_bytes);
	xen_pfn_t*pages= (xen_pfn_t*) malloc(num_pages* sizeof(xen_pfn_t));
	int* errors = (int*) malloc(num_pages* sizeof(int));
	printf("Addrewss before: %lu\n", vaddr);
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
	printf("%lu\n", XC_PAGE_SIZE);
	unsigned long offset = vaddr % XC_PAGE_SIZE;
	printf("%lu\n", offset);
	long* data = (long*)mem;
	data =(data + offset);
	printf("YOU %lu\n", *data);
	return mem;
}
void write_memory(void* map, char* buffer, int length) {
	memcpy(map, (void*)buffer, length);
}
*/
import "C"

import (
	"fmt"
	"unsafe"
)

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
		C.write_memory(v, C.CString(string(bytes)), C.int(len(bytes)))
	}
	return nil
}

func (mem *Memory) UnMap(address uint64, pages uint32) error {
	err := C.xenforeignmemory_unmap(mem.key, mem.memories[address], 1)
	fmt.Println("Error in unmap", err)
	return nil
}
