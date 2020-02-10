package main

import (
	"debug/dwarf"
	"debug/elf"
	"fmt"
)

func main() {
	file, err := elf.Open("../../stardust-experimental/build/mini-os")
	pc := uint64(5996)
	variable := "i"
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
	reader.SeekPC(pc)
	for entry, _ := reader.Next(); entry != nil; entry, _ = reader.Next() {
		// fmt.Printf("%+v\n", entry)
		// fmt.Println(entry.Tag)
		if entry.Tag == dwarf.TagCompileUnit {
			break
		} else if entry.Tag == dwarf.TagVariable {
			filed := entry.AttrField(dwarf.AttrName)
			if filed != nil {
				name, ok := filed.Val.(string)
				if ok {
					if name == variable {
						fmt.Println("FOUND ENTRY")
						fmt.Printf("%+v\n", entry)
						filed := entry.AttrField(dwarf.AttrType)
						typeOffset, ok := filed.Val.(dwarf.Offset)
						if ok {
							currentOffset := entry.Offset
							reader.Seek(typeOffset)
							entry, _ := reader.Next()
							fmt.Println(entry)
							reader.Seek(currentOffset)
						}
						filed = entry.AttrField(dwarf.AttrLocation)
						fmt.Println(filed)
						location, ok := filed.Val.(int64)
						ok = true

						if ok {
							location := dwarf.Offset(location)
							currentOffset := entry.Offset
							reader.Seek(location)
							entry, _ := reader.Next()
							fmt.Printf("%+v\n", entry)
							reader.Seek(currentOffset)
						}
					}
					reader.Next()
				}
			}
		}
	}
	// reader.Seek(874276)
	// entry, _ := reader.Next()
	// fmt.Println(entry)
	// entry, _ = reader.Next()
	// fmt.Println(entry)
	// for entry, _ := reader.Next(); entry != nil; entry, _ = reader.Next() {
	// 	fmt.Printf("%+v\n", entry)
	// 	fmt.Println(entry.Tag)
	// 	if entry.Tag == dwarf.TagCompileUnit {
	// 		break
	// 	}
	// }
}
