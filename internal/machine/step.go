package machine

import (
	"cpu-simulator/internal/isa"
	"fmt"
)

func (c *CPU) Step() error {
	if c.halted {
		return nil
	}

	word, err := c.readWord(c.PC)
	if err != nil {
		return err
	}

	c.IR = isa.DecodeInstruction(word)
	c.PC += isa.WordSize

	switch c.IR.Opcode{
		case isa.HALT:
			c.halted = true
		case isa.NOP:
			
		case isa.JMP:
			c.PC = uint32(c.IR.Operand)
		default:
		return fmt.Errorf("unknown opcode: %v", c.IR.Opcode)
	}

	c.tickCounter += c.instructionTick(c.IR.Opcode)
	return nil
}

