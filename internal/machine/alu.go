package machine

import (
	"cpu-simulator/internal/isa"
	"fmt"
)

type BinaryOperation func(left, right int32) (int32, error)

var binaryOperations = map[isa.Opcode]BinaryOperation{
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

type UnaryOperation func(value int32) (int32, error)

var unaryOperations = map[isa.Opcode]UnaryOperation{
	isa.INC: func(value int32) (int32, error) {
		return value + 1, nil
	},
	isa.DEC: func(value int32) (int32, error) {
		return value - 1, nil
	},
	isa.NOT: func(value int32) (int32, error) {
		return ^value, nil
	},
}

func (c *CPU) updateZN(value int32) {
	c.setFlag(isa.Z, value == 0)
	c.setFlag(isa.N, value < 0)
}

func (c *CPU) execBinaryOp(opcode isa.Opcode) error {
	operation, ok := binaryOperations[opcode]
	if !ok {
		return fmt.Errorf("unknown alu opcode: %v", opcode)
	}

	right, err := c.readOperand()
	if err != nil {
		return err
	}

	result, err := operation(c.ACC, right)
	if err != nil {
		return err
	}

	c.ACC = result
	c.updateZN(result)
	return nil
}

func (c *CPU) execUnaryOp(opcode isa.Opcode) error {
	operation, ok := unaryOperations[opcode]
	if !ok {
		return fmt.Errorf("unknown unary alu opcode: %v", opcode)
	}

	result, err := operation(c.ACC)
	if err != nil {
		return err
	}

	c.ACC = result
	c.updateZN(result)
	return nil
}

func (c *CPU) execCmp(value int32) {
	result := c.ACC - value
	c.updateZN(result)
}
