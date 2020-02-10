package main

import (
	"debug/dwarf"
	"debug/elf"
	"fmt"
)

func main() {
	file, err := elf.Open("../../stardust-experimental/build/mini-os")
	if err != nil {
		fmt.Println(err)
		fmt.Println("Something went wrong")
		return
	}
	d, err := file.DWARF()
	if err != nil {
		fmt.Println("Something went wrong getting the dwarf")
		return
	}

	reader := d.Reader()

	for entry, err := reader.Next(); entry != nil && err == nil; entry, err = reader.Next() {
		fmt.Println(entry.Tag.GoString())
		lineReader, err := d.LineReader(entry)
		if err != nil {
			fmt.Println(err)
		}
		if lineReader == nil {
			continue
		}

		lineEntry := new(dwarf.LineEntry)
		for err := lineReader.Next(lineEntry); err == nil; err = lineReader.Next(lineEntry) {
			//fmt.Println(lineEntry)
			fmt.Printf("The name of the file is: %s and the line number is %d. Address: %d\n", lineEntry.File.Name, lineEntry.Line, lineEntry.Address)
		}
		// entry, err = reader.Next()
		// if err != nil {

		// 	fmt.Println(err)
		// }
	}

}
