package file

import (
	"debug/dwarf"
	"debug/elf"
	"strings"
)

type info struct {
	lineToAddressInfo map[int][]uint64
}

func (i *info) addLineInformation(entry *dwarf.LineEntry) {
	if !entry.IsStmt {
		return 
	}
	if i.lineToAddressInfo == nil {
		i.lineToAddressInfo = make(map[int][]uint64)
	} 
	
	list, _ := i.lineToAddressInfo[entry.Line]
	i.lineToAddressInfo[entry.Line] = append(list, entry.Address)
}

func (i *info) getAddress(line int) uint64 {
	if val, ok := i.lineToAddressInfo[line]; ok {
		return val[0]
	}
	return 0
}

//LineInformation - handles the translation between source line information
//and addresses in the executable
type LineInformation struct {
	// Name of the image (i.e. my-os and so-)
	Name            string
	data            *dwarf.Data
	pcToLineEntry     map[uint64]dwarf.LineEntry
	filenameToAddress map[string]info
	currentLine     int
	currentFile     string
}

//Init - sets up the File struct. This must be run before any
//of the other methods are run 
func (lineInfo *LineInformation) Init() error {
	lineInfo.pcToLineEntry = make(map[uint64]dwarf.LineEntry)
	lineInfo.filenameToAddress = make(map[string]info)

	file, err := elf.Open(lineInfo.Name)
	defer file.Close()
	if err != nil {
		return err
	}

	d, err := file.DWARF()
	if err != nil {
		return err
	}
	lineInfo.data = d
	compilationReader := d.Reader()

	for entry, err := compilationReader.Next(); entry != nil && err == nil; entry, err = compilationReader.Next() {
		lineReader, err := d.LineReader(entry)
		if err != nil {
			return err
		}
		//Compilation units may contain nothing so we just skip it
		if lineReader == nil {
			continue
		}

		lineEntry := new(dwarf.LineEntry)
		for err := lineReader.Next(lineEntry); err == nil; err = lineReader.Next(lineEntry) {
			if _, ok := lineInfo.pcToLineEntry[lineEntry.Address]; !ok {
				lineInfo.pcToLineEntry[lineEntry.Address] = *lineEntry
			}
			filename := lineEntry.File.Name
			path := strings.Split(filename, "/")
			name := path[len(path)-1]
			i, _ := lineInfo.filenameToAddress[name]
			i.addLineInformation(lineEntry)
			lineInfo.filenameToAddress[name] = i
		}
	}

	return nil
}

//Address - gets the address of a place in a file and line
func (lineInfo *LineInformation) Address(filename string, line int) uint64 {
	if info, ok := lineInfo.filenameToAddress[filename]; ok {
		address := info.getAddress(line)
		return address
	}
	return 0
}

//AddressToLine - translates an address to the filename and line number
func (lineInfo *LineInformation) AddressToLine(address uint64) (string, int, error) {
	reader := lineInfo.data.Reader()
	entry, err := reader.SeekPC(address)
	if err != nil {
		return "", 0, err
	}

	lineReader, _ := lineInfo.data.LineReader(entry)
	var lineEntry dwarf.LineEntry
	err = lineReader.SeekPC(address, &lineEntry)
	if err != nil {
		return "", 0, err
	}
	return lineEntry.File.Name, lineEntry.Line, nil 
}

func (lineInfo *LineInformation) CurrentLine() (string, int) {
	return lineInfo.currentFile, lineInfo.currentLine
}

//IsNewLine - takes a PC counter and returns whether the new program
//counter is on a new source line or not
func (lineInfo *LineInformation) IsNewLine(rip uint64) (changed bool) {
	reader := lineInfo.data.Reader()
	entry, err := reader.SeekPC(rip)
	if entry == nil {
		return false
	}
	lineReader, _ := lineInfo.data.LineReader(entry)
	var lineEntry dwarf.LineEntry
	err = lineReader.SeekPC(rip, &lineEntry)
	if err == nil {
		if lineEntry.Line != lineInfo.currentLine || strings.Compare(lineEntry.File.Name, lineInfo.currentFile) != 0 {
			changed = true
		}
		lineInfo.currentLine = lineEntry.Line
		lineInfo.currentFile = lineEntry.File.Name
	} 
	return
}
