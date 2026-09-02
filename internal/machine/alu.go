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

func (c *CPU) updateALUFlags(opcode isa.Opcode, left, right, result int32) {
	c.updateZN(result)

	switch opcode {
	case isa.ADD, isa.INC:
		c.setFlag(isa.C, addCarry(left, right))
		c.setFlag(isa.V, addOverflow(left, right, result))
	case isa.SUB, isa.CMP, isa.DEC:
		c.setFlag(isa.C, subBorrow(left, right))
		c.setFlag(isa.V, subOverflow(left, right, result))
	default:
		c.setFlag(isa.C, false)
		c.setFlag(isa.V, false)
	}
}

func addCarry(left, right int32) bool {
	return uint32(left)+uint32(right) < uint32(left)
}

func addOverflow(left, right, result int32) bool {
	return (left >= 0 && right >= 0 && result < 0) || (left < 0 && right < 0 && result >= 0)
}

func subBorrow(left, right int32) bool {
	return uint32(left) < uint32(right)
}

func subOverflow(left, right, result int32) bool {
	return (left < 0 && right >= 0 && result >= 0) || (left >= 0 && right < 0 && result < 0)
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

	left := c.ACC
	result, err := operation(left, right)
	if err != nil {
		return err
	}

	c.ACC = result
	c.updateALUFlags(opcode, left, right, result)
	return nil
}

func (c *CPU) execUnaryOp(opcode isa.Opcode) error {
	operation, ok := unaryOperations[opcode]
	if !ok {
		return fmt.Errorf("unknown unary alu opcode: %v", opcode)
	}

	left := c.ACC
	result, err := operation(left)
	if err != nil {
		return err
	}

	c.ACC = result
	c.updateALUFlags(opcode, left, unaryRightOperand(opcode), result)
	return nil
}

func unaryRightOperand(opcode isa.Opcode) int32 {
	switch opcode {
	case isa.INC, isa.DEC:
		return 1
	default:
		return 0
	}
}

func (c *CPU) execCmp(value int32) {
	left := c.ACC
	result := left - value
	c.updateALUFlags(isa.CMP, left, value, result)
}
