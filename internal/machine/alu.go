package machine

import (
	"cpu-simulator/internal/isa"
	"fmt"
)

type Operation func(left, right int32) (int32, error)

var operations = map[isa.Opcode]Operation{
	isa.ADD: func(left, right int32) (int32, error) {
		return left + right, nil
	},
	isa.SUB: func(left, right int32) (int32, error) {
		return left - right, nil
	},
	isa.MUL: func(left, right int32) (int32, error) {
		return left * right, nil
	},
	isa.DIV: func(left, right int32) (int32, error) {
		if right == 0 {
			return 0, fmt.Errorf("division by zero")
		}
		return left / right, nil
	},
	isa.MOD: func(left, right int32) (int32, error) {
		if right == 0 {
			return 0, fmt.Errorf("division by zero")
		}
		return left % right, nil
	},
	isa.AND: func(left, right int32) (int32, error) {
		return left & right, nil
	},
	isa.OR: func(left, right int32) (int32, error) {
		return left | right, nil
	},
}

func (c *CPU) updateZN(value int32) {
	c.setFlag(isa.Z, value == 0)
	c.setFlag(isa.N, value < 0)
}

func (c *CPU) readOperand() (int32, error) {
	word, err := c.readWord(uint32(c.IR.Operand))
	if err != nil {
		return 0, err
	}
	return int32(word), nil
}

func (c *CPU) execOp(opcode isa.Opcode) error {
	operation, ok := operations[opcode]
	if !ok {
		return fmt.Errorf("unknown opcode: %v", opcode)
	}

	right, err := c.readOperand()
	if err != nil {
		return err
	}

	result, err := operation(c.ACC, right)
	if err != nil {
		return nil
	}

	c.ACC = result
	c.updateZN(result)
	return nil
}

func (c *CPU) increment() {
	c.ACC++
	c.updateZN(c.ACC)
}

func (c *CPU) decrement() {
	c.ACC--
	c.updateZN(c.ACC)
}

func (c *CPU) execCmp(value int32) {
	result := c.ACC - value
	c.updateZN(result)
}

func (c *CPU) execNot() {
	result := ^c.ACC
	c.updateZN(result)
	c.ACC = result
}
