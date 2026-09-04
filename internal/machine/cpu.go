package machine

import (
	"cpu-simulator/internal/isa"
)

type stageOperation func(*CPU) error

type CPU struct {
	ACC            int32
	PC             uint32
	IR             isa.Instruction
	SR             uint8
	AR             uint32
	DR             uint32
	DSP            uint32
	dataStack      []int32
	returnStack    []uint32
	RSP            uint32
	memoryImage    []uint32
	tickCounter    int
	step           int
	stages         []stageOperation
	pendingAddr    uint32
	pendingFrame   InterruptFrame
	halted         bool
	io             *IOController
	trace          []LogEntry
	interruptStack []InterruptFrame
}

func NewCPU(program []uint32) *CPU {
	return NewCPUWithInputEvents(program, nil)
}

func NewCPUWithInputEvents(program []uint32, inputEvents []InputEvent) *CPU {
	mem := make([]uint32, isa.MemSize/isa.WordSize)
	copy(mem, program)

	return &CPU{
		memoryImage:    mem,
		PC:             0,
		dataStack:      make([]int32, 0),
		returnStack:    make([]uint32, 0),
		DSP:            0,
		RSP:            0,
		io:             newIOControllerWithEvents(inputEvents),
		interruptStack: make([]InterruptFrame, 0),
	}
}

func (c *CPU) setFlag(flag isa.Status, value bool) {
	c.SR = statusWithFlag(c.SR, flag, value)
}

func (c *CPU) flag(flag isa.Status) bool {
	return c.SR&(1<<flag) != 0
}

func statusWithFlag(status uint8, flag isa.Status, value bool) uint8 {
	mask := uint8(1 << flag)
	if value {
		return status | mask
	}
	return status &^ mask
}

func statusWithZN(status uint8, value int32) uint8 {
	status = statusWithFlag(status, isa.Z, value == 0)
	return statusWithFlag(status, isa.N, value < 0)
}
