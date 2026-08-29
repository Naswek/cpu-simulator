package machine

import (
	"cpu-simulator/internal/isa"
	"fmt"
)

var loadOpcodes = map[isa.Opcode]struct{}{
	isa.LD_IMM:  {},
	isa.LD_ADDR: {},
	isa.LD_IND:  {},
	isa.LD_SP_N: {},
}

var storeOpcodes = map[isa.Opcode]struct{}{
	isa.ST_ADDR: {},
	isa.ST_IND:  {},
}

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

	if _, ok := loadOpcodes[c.IR.Opcode]; ok {
		err = c.execLoad()
	} else if _, ok := storeOpcodes[c.IR.Opcode]; ok {
		err = c.writeOperand(c.ACC)
	} else if _, ok := binaryOperations[c.IR.Opcode]; ok {
		err = c.execBinaryOp(c.IR.Opcode)
	} else if _, ok := unaryOperations[c.IR.Opcode]; ok {
		err = c.execUnaryOp(c.IR.Opcode)
	} else if _, ok := jumpOpcodes[c.IR.Opcode]; ok {
		err = c.execJump(c.IR.Opcode)
	} else if operation, ok := controlOpcodes[c.IR.Opcode]; ok {
		err = operation(c)
	} else {

		switch c.IR.Opcode {
		case isa.HALT:
			c.halted = true
		case isa.NOP:

		case isa.POP:
			var value int32
			value, err = c.popData()
			if err == nil {
				c.ACC = value
				c.updateZN(c.ACC)
			}
		case isa.PUSH:
			err = c.pushData(c.ACC)
		case isa.CMP:
			var value int32
			value, err = c.readOperand()
			if err == nil {
				c.execCmp(value)
			}
		default:
			return fmt.Errorf("unknown opcode: %v", c.IR.Opcode)
		}
	}
	if err != nil {
		return err
	}
	c.tickCounter += c.instructionTick(c.IR.Opcode)
	return nil
}
