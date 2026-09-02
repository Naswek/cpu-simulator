package translator

import (
	"cpu-simulator/internal/isa"
	"fmt"
)

type Program struct {
	Code             []isa.Instruction
	Data             map[isa.Operand]int32
	InterruptVectors map[uint8]isa.Operand
}

func (p Program) CodeWords() ([]uint32, error) {
	words := isa.EncodeProgram(p.Code)

	for vector, addr := range p.InterruptVectors {
		if vector >= isa.InterruptVectors {
			return nil, fmt.Errorf("interrupt vector %d is outside vector table", vector)
		}

		index := int(isa.InterruptVectorAddress(vector) / isa.Operand(isa.WordSize))
		if index >= len(words) {
			return nil, fmt.Errorf("interrupt vector %d is outside code image", vector)
		}
		words[index] = uint32(addr)
	}

	return words, nil
}

func (p Program) MemoryImage() ([]uint32, error) {
	maxLen := len(p.Code)

	for addr := range p.Data {
		if addr%isa.Operand(isa.WordSize) != 0 {
			return nil, fmt.Errorf("unaligned data address 0x%06X", addr)
		}
		if uint32(addr) >= isa.MemSize {
			return nil, fmt.Errorf("data address 0x%06X is outside memory", addr)
		}

		index := int(addr / isa.Operand(isa.WordSize))
		if index < len(p.Code) {
			return nil, fmt.Errorf("data address 0x%06X overlaps code", addr)
		}
		if index+1 > maxLen {
			maxLen = index + 1
		}
	}

	for i := range p.Code {
		addr := uint32(i) * isa.WordSize
		if addr >= isa.MemSize {
			return nil, fmt.Errorf("code address 0x%06X is outside memory", addr)
		}
	}

	code, err := p.CodeWords()
	if err != nil {
		return nil, err
	}

	image := make([]uint32, maxLen)
	for i, word := range code {
		image[i] = word
	}
	for addr, value := range p.Data {
		image[int(addr/isa.Operand(isa.WordSize))] = uint32(value)
	}

	return image, nil
}
