package machine

import (
	"cpu-simulator/internal/isa"
	"fmt"
)

type CPU struct {
	ACC int32
	PC uint32
	IR isa.Instruction
	SR uint8
	DSP int32
	RSP int32
	memoryImage []uint32
	tickCounter int
	halted bool
	
}

func NewCPU(program []uint32) *CPU {
	mem := make([]uint32, isa.MemSize/isa.WordSize)
	copy(mem, program)

	return &CPU{
		memoryImage: mem,
		PC: 0,
	}
}

func (c *CPU) setFlag(flag isa.Status, value bool) {
	mask := uint8(1 << flag)
	if value {
		c.SR |= mask
	} else {
		c.SR &^= mask 
	} 
}

func (c *CPU) flag(flag isa.Status) bool {
	return c.SR&(1<<flag) != 0
}

func (c *CPU) readWord(addr uint32) (uint32, error) {
	if addr%isa.WordSize != 0 {
		return 0, fmt.Errorf("unlighned memory address: 0x%X", addr)
	}

	index := addr / isa.WordSize
	if index >= uint32(len(c.memoryImage)) {
		return 0, fmt.Errorf("memory addres out of reange: 0x%X", addr)
	}
	return c.memoryImage[index], nil
}

func (c *CPU) instructionTick(opcode isa.Opcode) int {
	info, ok := isa.OpcodeTable[isa.Opcode(opcode)]
	if !ok {
		return 0
	}
	return info.Ticks
}
