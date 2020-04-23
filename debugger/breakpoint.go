package debugger

import (
	"fmt"
)

type errorType int

const (
	breakInt byte      = 0xCC
	NotFound errorType = iota
	AlreadyBreakpointSet
)

type BreakPointError struct {
	Address uint64
	errType errorType
}

func (e BreakPointError) Error() string {
	switch e.errType {
	case NotFound:
		return fmt.Sprintf("Error: no breakpoint at %d", e.Address)
	case AlreadyBreakpointSet:
		return fmt.Sprintf("Error: breakpoint already at %d", e.Address)
		
	}
	return ""
}

//Breakpoints manages the setting and remove of breakpoints
type Breakpoints struct {
	breakpoints        map[uint64]byte
	mem                MemoryAccess
	restoreBreakpoints []uint64
}

//NewBreakpointManager create a new strcut for handling the creation, and removal
//of breakpoints
func NewBreakpointManager(mem MemoryAccess) *Breakpoints {
	bp := new(Breakpoints)
	bp.breakpoints = make(map[uint64]byte)
	bp.mem = mem
	return bp
}

//Add - writes a the break instruction to specified address
//returns error if there is already breakpoint at that location or something goes
//wrong while writing to memory
func (point *Breakpoints) Add(address uint64) error {
	if _, ok := point.breakpoints[address]; ok {
		return BreakPointError{Address: address, errType: AlreadyBreakpointSet}
	}
	bytes, err := point.mem.Read(address, 1)
	if err != nil {
		return err
	}
	point.breakpoints[address] = bytes[0]
	err = point.mem.Write(address, []byte{breakInt}, 1)
	return err
}

//Remove - deletes a breakpoint and puts the memory back in original state
func (point *Breakpoints) Remove(address uint64) error {
	if origByte, ok := point.breakpoints[address]; ok {
		err := point.mem.Write(address, []byte{origByte}, 1)
		if err != nil {
			return err
		}
		delete(point.breakpoints, address)
		return nil
	}
	return BreakPointError{Address: address, errType: NotFound}
}

//Addresses - returns a list of addresses
func (point *Breakpoints) Addresses() []uint64 {
	var list []uint64
	for address, _ := range point.breakpoints {
		list = append(list, address)
	}
	return list
}

//RestoreInstruction - writes back the original byte that was overwritten 
//with the breakpoint
func (point *Breakpoints) RestoreInstruction(address uint64) error {
	if point.AddressIsBreakpoint(address) {
		err := point.Remove(address)
		if err != nil {
			return err
		}
		point.restoreBreakpoints = append(point.restoreBreakpoints, address)
		return nil
	}
	return BreakPointError{address, NotFound}
}

//RestoreBreakpoint - puts breakpoints back in their place (only the ones
//remove by RestoreInstruction, NOT Remove)
func (point *Breakpoints) RestoreBreakpoint() error {
	for _, addressRestoreInstruction := range point.restoreBreakpoints {
		err := point.Add(addressRestoreInstruction)
		if err != nil {
			return err
		}
	} 
	point.restoreBreakpoints = nil
	return nil
}

//AddressIsBreakpoint - returns the whether the address has a breakpoint
//at it
func (point *Breakpoints) AddressIsBreakpoint(address uint64) bool {
	_, ok := point.breakpoints[address]
	return ok
}
