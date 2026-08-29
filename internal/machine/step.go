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

	if _, ok := binaryOperations[c.IR.Opcode]; ok {
		if err := c.execBinaryOp(c.IR.Opcode); err != nil {
			return err
		}
	} else if _, ok := unaryOperations[c.IR.Opcode]; ok {
		if err := c.execUnaryOp(c.IR.Opcode); err != nil {
			return err
		}
	} else {

		switch c.IR.Opcode {
		case isa.HALT:
			c.halted = true
		case isa.NOP:

		case isa.JMP:
			c.PC = uint32(c.IR.Operand)

		case isa.LD_IMM:
			c.ACC = int32(c.IR.Operand)
		case isa.LD_ADDR:
			word, err := c.readWord(uint32(c.IR.Operand))
			if err != nil {
				return err
			}
			c.ACC = int32(word)
		case isa.ST_ADDR:
			err := c.writeWord(uint32(c.IR.Operand), uint32(c.ACC))
			if err != nil {
				return err
			}
		case isa.LD_IND:
			ptr, err := c.readWord(uint32(c.IR.Operand))
			if err != nil {
				return err
			}

			word, err = c.readWord(ptr)
			if err != nil {
				return err
			}

			c.ACC = int32(word)
		case isa.ST_IND:
			ptr, err := c.readWord(uint32(c.IR.Opcode))
			if err != nil {
				return err
			}

			err = c.writeWord(ptr, uint32(c.ACC))
			if err != nil {
				return err
			}

		default:
			return fmt.Errorf("unknown opcode: %v", c.IR.Opcode)
		}
	}
	c.tickCounter += c.instructionTick(c.IR.Opcode)
	return nil
}
