package machine

import (
	"cpu-simulator/internal/isa"
)

type CPU struct {
	ACC         int32
	PC          uint32
	IR          isa.Instruction
	SR          uint8
	DSP         int32
	RSP         int32
	memoryImage []uint32
	tickCounter int
	halted      bool
}

func NewCPU(program []uint32) *CPU {
	mem := make([]uint32, isa.MemSize/isa.WordSize)
	copy(mem, program)

	return &CPU{
		memoryImage: mem,
		PC:          0,
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

func (c *CPU) instructionTick(opcode isa.Opcode) int {
	info, ok := isa.OpcodeTable[opcode]
	if !ok {
		return 0
	}
	return info.Ticks
}
