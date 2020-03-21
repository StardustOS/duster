package file

import (
	"debug/dwarf"
)

type SymbolError int

const (
	InvalidDWARF SymbolError = 0
)

func (err SymbolError) Error() string {
	switch err {
	case InvalidDWARF:
		return "Error: not a valid dwarf file"
	}
	return ""
}

type Symbol struct {
	pc    uint64
	Data  *dwarf.Data
	entry *dwarf.Entry
}

func (sym *Symbol) Update(pc uint64) error {
	if sym.entry == nil {
		reader := sym.Data.Reader()
		_, err := reader.SeekPC(pc)
		if err != nil {
			return err
		}
		for entry, err := reader.Next(); entry != nil; entry, err = reader.Next() {
			if err != nil {
				return err
			}
			if entry.Tag == dwarf.TagSubprogram || entry.Tag == dwarf.TagLexDwarfBlock {
				lowPC := entry.AttrField(dwarf.AttrLowpc)
				highPC := entry.AttrField(dwarf.AttrHighpc)
				if lowPC != nil && highPC != nil {
					lowPC, ok := lowPC.Val.(uint64)

					if !ok {
						return InvalidDWARF
					}
					highPC := uint64(highPC.Val.(int64))
					if !ok {
						return InvalidDWARF
					}
					highPC += uint64(lowPC)
					if lowPC <= pc && highPC >= pc {
						sym.entry = entry
						sym.pc = pc
						return nil
					}
				}
			}
		}
	}
	return nil
}

func (sym *Symbol) Symbols() ([]string, error) {
	var symbols []string
	reader := sym.Data.Reader()
	reader.Seek(sym.entry.Offset)
	reader.Next()
	for child, err := reader.Next(); child != nil && child.Tag == dwarf.TagVariable; child, err = reader.Next() {
		if err != nil {
			return nil, err
		}
		field := child.AttrField(dwarf.AttrName)
		if field != nil {
			name, ok := field.Val.(string)
			if !ok {
				return nil, InvalidDWARF
			}
			symbols = append(symbols, name)
		}
	}
	return symbols, nil
}
