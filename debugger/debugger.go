package debugger

import "C"

import (
	"bufio"
	"encoding/binary"
	"fmt"
	"os"

	"github.com/go-delve/delve/pkg/dwarf/op"
)

type DebuggingError int

func (err DebuggingError) Error() string {
	switch err {
	case NotPaused:
		return "Error: Domain is not paused"
	case NotPointer:
		return "Error: Not pointer type"
	}
	return ""
}

const (
	singleStep uint64         = 0x00000100
	//NotPaused is returned when a call is made which 
	//modifies the VM states but it is still running
	NotPaused  DebuggingError = iota

	//NotPointer is returned when dereference is run 
	//non pointer type.
	NotPointer
)

//Registers is an interface that defines how the debugger
//will interact with register at a given point in time.
type Registers interface {
	//GetRegister given a string, such as "rip", return the value stored
	//in that register
	GetRegister(string) (uint64, error)

	//SetRegister given a string, such as "rip", set that register to the new
	//value. However, do not write that value back to the virtual machine.
	//That is not the responsability of this interface.
	SetRegister(string, uint64) error

	//DwarfRegisters return a representation of the registers in the 
	//format that the Vendor has specified for the current hardware.
	DwarfRegisters() *op.DwarfRegisters
}

//RegisterHandler provides functionality for retrieving and setting
//the register from the VM
type RegisterHandler interface {
	//GetRegister should get the register of a certain VCPU
	GetRegisters(uint32) (Registers, error)

	//SetRegisters should write the register back to a certain VCPU. 
	//You may assume that the struct returned from GetRegister is the 
	//same as the one passed here (i.e. the debugger will not use internal
	//struct instead)
	SetRegisters(uint32, Registers) error
}

//Control define interface for pausing and unpausing the VM
//along with checking it is paused (please note an implementation
//where the functionality of IsPaused is implemented as a consequence
//of the other two will NOT work. As the the IsPaused function is used
//to check whether a breakpoint has been hit).
type Control interface {
	
	//IsPaused checks whether the VM is paused please read above note.
	IsPaused() bool
	
	//Pauses the VM
	Pause() error
	
	//Unpause the VM
	Unpause() error
}

//LineInformation defines interface for getting information about the 
//source information. This interface does not handle symbols.
type LineInformation interface {
	
	//CurrentLine, gets the current line. Please note
	//the return types are assumed to be the filename and 
	//the line number. DO NOT RETURN THE LINE ITSELF, the caller
	//will handle getting this. 
	CurrentLine() (string, int)
	
	//IsNewLine this method will be passed the current program counter (PC).
	//Return true when we start a new source line, otherwise false.
	IsNewLine(uint64) bool
	
	//Address this function takes a filename and a line number then returns an
	//address. Please note the filename will be in form file.c not a complete path
	//it is your responsability to handle issues because of this.
	Address(string, int) uint64


	//AddressToLine takes an address and converts it into line information.
	//Please note this operation cannot affect the operation of the other 
	//methods.
	AddressToLine(address uint64) (string, int, error)
}

//MemoryAccess defines API that will be used to for reading and writing to memory
type MemoryAccess interface {
	//Read reads the memory at the passed address
	Read(address uint64, size uint) ([]byte, error)
	//Write write a slice of bytes to an address
	Write(address uint64, bytes []byte, size uint) error
}

//Variable interface defines how the debugger will interact 
//variables.
type Variable interface {
	//Name returns the name of the variable
	Name() string

	//Parse should return a HUMAN readable string of
	//the content of that Variable. The bytes passed 
	//are the bytes stored in that variable's memory location.
	//It is your responsability to validate them (i.e. right size and so-on).  
	Parse([]byte, binary.ByteOrder) (string, error)

	//Location should return the DWARF expression associated with a variable.
	Location() []byte

	//Size should return the number of bytes in memory used to represent the type.
	//In the case of pointers return the size of the point NOT the size of the type
	//we're pointing to
	Size() int
}

//Symbol interface defines how the debugger will interact with symbolic 
//information. This may come from DWARF or some other format.
//Please note the underlying type of Variable returned from GetSymbol may be 
//assumed to be the same type that is passed to IsPointer or GetPointContentSize.
type Symbol interface {
	//GetSymbol give a variable name and the current PC return a variable.
	//Note not finding the variable is considered an error and must return an error
	//DO NOT JUST RETURN NIL
	GetSymbol(string, uint64) (Variable, error)

	//IsPointer return true if the variable is a pointer otherwise false
	IsPointer(Variable) bool

	//GetPointContentSize return the size of the type that is pointed by a pointer
	GetPointContentSize(Variable) int

	//ParsePointer pretty print the content memory pointed by a pointer.
	//Note the variable passed is the pointer
	ParsePointer(Variable, []byte, binary.ByteOrder) (string, error)
}

//Debugger struct carries out the debugging
type Debugger struct {
	endianess         binary.ByteOrder
	domainid          uint32
	breakpointManager *Breakpoints
	registers         RegisterHandler
	controller        Control
	memory            MemoryAccess
	lineInfo          LineInformation
	symbols           Symbol
}

//NewDebugger - constructor the debugger struct
func NewDebugger(memory MemoryAccess, controller Control, lineInfo LineInformation, registers RegisterHandler, symbols Symbol) *Debugger {
	debugger := new(Debugger)
	debugger.breakpointManager = NewBreakpointManager(memory)
	debugger.controller = controller
	debugger.lineInfo = lineInfo
	debugger.registers = registers
	debugger.symbols = symbols
	debugger.memory = memory
	debugger.endianess = binary.LittleEndian
	return debugger
}

func (debugger *Debugger) singleStep(vcpu uint32, start bool) error {
	registers, err := debugger.registers.GetRegisters(vcpu)
	if err != nil {
		return err
	}

	rflags, err := registers.GetRegister("rflags")
	if err != nil {
		return err
	}
	if start {
		rflags |= singleStep
	} else {
		rflags &= ^singleStep
	}

	err = registers.SetRegister("rflags", rflags)
	if err != nil {
		return err
	}

	err = debugger.registers.SetRegisters(vcpu, registers)
	if err != nil {
		return err
	}
	return nil
}

//GetLineInformation returns the filename, line number and line content
//for the current point in the execution
func (debugger *Debugger) GetLineInformation() string {
	filename, lineNo := debugger.lineInfo.CurrentLine()
	file, _ := os.Open(filename)
	reader := bufio.NewReader(file)
	var line string
	//Go to the current line (the last iteration will yield the line)
	for i := 0; i < lineNo; i += 1 {
		line, _ = reader.ReadString('\n')
	}
	return fmt.Sprintf("%s:%d - %s", filename, lineNo, line)
}

//Step - moves the program to the next source line
//Note only works when the process has been paused 
//and will put into single step mode. 
func (debugger *Debugger) Step(vcpu uint32) error {
	if !debugger.controller.IsPaused() {
		return NotPaused
	}

	err := debugger.singleStep(vcpu, true)
	if err != nil {
		return err
	}

	isNewline := false

	for !isNewline {
		//Busy wait until the process is no longer running
		for !debugger.controller.IsPaused() {
		}

		registers, err := debugger.registers.GetRegisters(vcpu)
		if err != nil {
			return err
		}

		rip, err := registers.GetRegister("rip")
		if err != nil {
			return err
		}


		isNewline = debugger.lineInfo.IsNewLine(rip)

		//Check whether a breakpoint exists here (if we've went through one
		//we need to fix the instruction that we've broken at before moving on).
		if debugger.breakpointManager.AddressIsBreakpoint(rip) {

			//We need to restore the instruction to its original state before 
			//the breakpoint was inserted because the instruction we over wrote 
			//may modify some state of the VM.
			err = debugger.breakpointManager.RestoreInstruction(rip)
			if err != nil {
				return err
			}

			//We want to rollback to run the instruction that was over written with
			//the breakpoint 
			rip -= 1
			registers.SetRegister("rip", rip)
			debugger.registers.SetRegisters(vcpu, registers)
		}

		err = debugger.breakpointManager.RestoreBreakpoint()
		if err != nil {
			return err
		}

		err = debugger.controller.Unpause()
		if err != nil {
			return err
		}
	}
	return nil
}

//Helper function for reading the contents of variables from memory
//Note we need the registers in the DWARF format, because we'll need to 
//evaluate a DWARF expression
func (debugger *Debugger) readMemory(variable Variable, regs *op.DwarfRegisters) ([]byte, error) {
	address, piece, err := op.ExecuteStackProgram(*regs, variable.Location())
	if err != nil {
		return nil, err
	}
	if piece == nil {
		size := uint(variable.Size())
		bytes, err := debugger.memory.Read(uint64(address), size)
		return bytes, err
	}
	return nil, nil
}

//GetVariable returns a pretty printed string represent the content 
//of that variable
func (debugger *Debugger) GetVariable(name string) (string, error) {
	if !debugger.controller.IsPaused() {
		return "", NotPaused
	}

	registers, err := debugger.registers.GetRegisters(0)
	if err != nil {
		return "", err
	}

	rip, err := registers.GetRegister("rip")
	if err != nil {
		return "", err
	}

	variable, err := debugger.symbols.GetSymbol(name, rip)
	if err != nil {
		return "", err
	}
	dregs := registers.DwarfRegisters()
	bytes, err := debugger.readMemory(variable, dregs)

	if err != nil {
		return "", err
	}

	val, err := variable.Parse(bytes, debugger.endianess)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%s = %s", name, val), nil
}


//Dereference returns the content of a point in pretty printed string
func (debugger *Debugger) Dereference(vcpu uint32, name string) (string, error) {
	if !debugger.controller.IsPaused() {
		return "", NotPaused
	}

	registers, err := debugger.registers.GetRegisters(0)
	if err != nil {
		return "", err
	}

	rip, err := registers.GetRegister("rip")
	if err != nil {
		return "", err
	}

	variable, err := debugger.symbols.GetSymbol(name, rip)
	if err != nil {
		return "", err
	}

	if !debugger.symbols.IsPointer(variable) {
		return "", NotPointer
	}

	dregs := registers.DwarfRegisters()
	bytes, err := debugger.readMemory(variable, dregs)
	if err != nil {
		return "", err
	}

	//Since we know the variable is a pointer, we know the content 
	//of its memory is an address 
	address := debugger.endianess.Uint64(bytes)
	size := debugger.symbols.GetPointContentSize(variable)

	bytes, err = debugger.memory.Read(address, uint(size))
	if err != nil {
		return "", err
	}

	val, err := debugger.symbols.ParsePointer(variable, bytes, debugger.endianess)
	if err != nil {
		return "", err 
	}
	return fmt.Sprintf("*%s = %s", name, val), nil 
}

//Continues to the next breakpoint or until the VM terminates
func (debugger *Debugger) Continue(vcpu uint32) error {
	if !debugger.controller.IsPaused() {
		return NotPaused
	}

	//Makes sure we don't break on each new instruction
	debugger.singleStep(vcpu, false)

	err := debugger.controller.Unpause()
	if err != nil {
		return err
	}

	//Busy wait until we hit the next breakpoint
	for !debugger.controller.IsPaused() {
	}
	return nil
}

//SetBreakpoint sets a breakpoint at specific point in the source 
func (debugger *Debugger) SetBreakpoint(filename string, line int, vcpu uint32) error {
	if !debugger.controller.IsPaused() {
		return NotPaused
	}

	address := debugger.lineInfo.Address(filename, line)

	//If the address is zero we're either trying to set a breakpoint 
	//on an empty line or during the preamble. 
	if address == 0 {
		return fmt.Errorf("Error: could not set breakpoint @ %s:%d (most likely empty line or comment)", filename, line)
	}

	err := debugger.breakpointManager.Add(address)
	return err
}

//RemoveBreakpoints removes a breakpoint from the VM
func (debugger *Debugger) RemoveBreakpoint(filename string, line int, vcpu uint32) error {
	if !debugger.controller.IsPaused() {
		return NotPaused
	}
	address := debugger.lineInfo.Address(filename, line)
	err := debugger.breakpointManager.Remove(address)
	return err
}

//ListBreakpoints - returns a formatted list of the breakpoints that have been set
func (debugger *Debugger) ListBreakpoints() string {
	addresses := debugger.breakpointManager.Addresses()
	var formattedList string
	for _, address := range addresses {
		filename, line, _ := debugger.lineInfo.AddressToLine(address)
		formattedList = fmt.Sprintf("%s0x%x (%s:%d)\n", formattedList, address, filename, line)
	}
	if len(formattedList) == 0 {
		formattedList = "No breakpoints have been set!"
	}
	return formattedList
}