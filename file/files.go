package file

import (
	"debug/dwarf"
	"debug/elf"
	"fmt"
	"strings"
)

type info struct {
	lineToAddressInfo map[int][]uint64
}

func (i *info) addLineInformation(entry *dwarf.LineEntry) {
	if i.lineToAddressInfo == nil {
		i.lineToAddressInfo = make(map[int][]uint64)
	}
	if strings.Contains(entry.File.Name, "startup.c") {
		// fmt.Println("Line no", entry.Line)
		// fmt.Printf("%+v\n", entry)
	}
	list, _ := i.lineToAddressInfo[entry.Line]
	i.lineToAddressInfo[entry.Line] = append(list, entry.Address)
}

func (i *info) getAddress(line int) uint64 {
	//fmt.Println(i.lineToAddressInfo[line])
	if val, ok := i.lineToAddressInfo[line]; ok {
		return val[0]
	}
	return 0
}

type File struct {
	Name            string
	data            *dwarf.Data
	information     map[uint64]dwarf.LineEntry
	informationFile map[string]info
	currentLine     int
	currentFile     string
}

func (f *File) Init() error {
	f.information = make(map[uint64]dwarf.LineEntry)
	f.informationFile = make(map[string]info)

	file, err := elf.Open(f.Name)
	defer file.Close()

	if err != nil {
		return err
	}
	d, err := file.DWARF()
	if err != nil {
		return err
	}
	f.data = d
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
			if _, ok := f.information[lineEntry.Address]; !ok {
				f.information[lineEntry.Address] = *lineEntry
			}
			filename := lineEntry.File.Name
			path := strings.Split(filename, "/")
			name := path[len(path)-1]
			if strings.Compare(name, "startup.c") == 0 {
				fmt.Printf("%+v\n", lineEntry)
			}
			i, _ := f.informationFile[name]
			i.addLineInformation(lineEntry)
			f.informationFile[name] = i
		}
	}

	return nil
}

func (f *File) GetAddress(filename string, line int) uint64 {
	if info, ok := f.informationFile[filename]; ok {
		address := info.getAddress(line)
		return address
	}
	return 0
}

func (f *File) CurrentLine() (string, int) {
	return f.currentFile, f.currentLine
}

func (f *File) UpdateLine(rip uint64) (changed bool) {
	reader := f.data.Reader()
	entry, err := reader.SeekPC(rip)
	// fmt.Println(err)
	// fmt.Println(entry)
	if entry == nil {
		fmt.Println("NIL!")
	}
	lineReader, _ := f.data.LineReader(entry)
	// fmt.Println(lineReader)
	var lineEntry dwarf.LineEntry
	err = lineReader.SeekPC(rip, &lineEntry)
	if err == nil {
		if lineEntry.Line != f.currentLine || strings.Compare(lineEntry.File.Name, f.currentFile) != 0 {
			changed = true
		}
		f.currentLine = lineEntry.Line
		f.currentFile = lineEntry.File.Name
	}
	return
}
