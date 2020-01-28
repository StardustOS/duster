package xen

type WordSize uint

const (
	SixtyFourBit WordSize = 8
	ThirtyTwoBit WordSize = 4
)

type Register struct {
	Type WordSize
	registers map[string][]byte
}

func (register *Register) AddRegister(name string, content []byte) error {
	if register.registers == nil {
		register.registers = make(map[string][]byte)
	}
	register.registers[name] = content
	return nil
}

func (register *Register) GetRegister(name string) ([]byte, error) {
	if content, ok := register.registers[name]; ok {
		return content, nil 
	}
	return nil, nil
}