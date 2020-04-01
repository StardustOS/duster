package main

import (
	"debug/dwarf"
	"debug/elf"
	"fmt"
)

func main() {
	file, err := elf.Open("./file/a.out")
	if err != nil {
		fmt.Println(err)
	}
	d, err := file.DWARF()
	if err != nil {
		fmt.Println(err)
	}

	reader := d.Reader()
	entry, _ := reader.SeekPC(0x66f)
	lineReader, _ := d.LineReader(entry)
	var e dwarf.LineEntry
	lineReader.SeekPC(0x66f, &e)
	fmt.Printf("%+v\n", e)
}
