package main

import (
	"debug/elf"
	"fmt"
)

func main() {
	//file, _ := elf.Open("../stardust-experimental/build/mini-os")
	file, _ := elf.Open("./file/testfiles/globalvars")

	for _, section := range file.Sections {
		fmt.Println(section.Name)
	}

}
