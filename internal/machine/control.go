package machine

import (
	"cpu-simulator/internal/isa"
	"fmt"
)

type JumpOperation func(c *CPU, addr isa.Operand) error

var jumpOpcodes = map[isa.Opcode]JumpOperation{
	isa.JMP: func(c *CPU, addr isa.Operand) error {
		return c.jump(addr)
	},
	isa.JZ: func(c *CPU, addr isa.Operand) error {
		return c.jumpIf(c.flag(isa.Z), addr)
	},
	isa.JNZ: func(c *CPU, addr isa.Operand) error {
		return c.jumpIf(!c.flag(isa.Z), addr)
	},
	isa.JN: func(c *CPU, addr isa.Operand) error {
		return c.jumpIf(c.flag(isa.N), addr)
	},
	isa.JP: func(c *CPU, addr isa.Operand) error {
		return c.jumpIf(!c.flag(isa.N) && !c.flag(isa.Z), addr)
	},
}

func (c *CPU) jump(addr isa.Operand) error {
	if uint32(addr)%isa.WordSize != 0 {
		return fmt.Errorf("unaligned jump address: 0x%X", addr)
	}
	if uint32(addr) >= isa.MemSize {
		return fmt.Errorf("jump address out of range: 0x%X", addr)
	}

	c.PC = uint32(addr)
	return nil
}

func (c *CPU) jumpIf(condition bool, addr isa.Operand) error {
	if !condition {
		return nil
	}

	return c.jump(addr)
}

func (c *CPU) execJump(opcode isa.Opcode) error {
	jumpOperation, ok := jumpOpcodes[opcode]
	if !ok {
		return fmt.Errorf("unknown jump opcode: %v", opcode)
	}

	return jumpOperation(c, c.IR.Operand)
}
