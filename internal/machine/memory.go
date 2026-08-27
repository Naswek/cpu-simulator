package machine

import (
	"cpu-simulator/internal/isa"
	"fmt"
)

func (c *CPU) readWord(addr uint32) (uint32, error) {
	if addr%isa.WordSize != 0 {
		return 0, fmt.Errorf("unaligned memory address: 0x%X", addr)
	}

	index := addr / isa.WordSize
	if index >= uint32(len(c.memoryImage)) {
		return 0, fmt.Errorf("memory address out of range: 0x%X", addr)
	}
	return c.memoryImage[index], nil
}

func (c *CPU) writeWord(addr uint32, value uint32) error {
	if addr%isa.WordSize != 0 {
		return fmt.Errorf("unaligned memory address: 0x%X", addr)
	}

	index := addr / isa.WordSize
	if index >= uint32(len(c.memoryImage)) {
		return fmt.Errorf("memory address out of range: 0x%X", addr)
	}
	c.memoryImage[index] = value
	return nil
}
