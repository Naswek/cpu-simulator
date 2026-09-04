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
	c.SR = statusWithZN(c.SR, value)
}

func (c *CPU) updateALUFlags(opcode isa.Opcode, left, right, result int32) {
	c.SR = statusWithALUFlags(c.SR, opcode, left, right, result)
}

func statusWithALUFlags(status uint8, opcode isa.Opcode, left, right, result int32) uint8 {
	status = statusWithZN(status, result)
	switch opcode {
	case isa.ADD, isa.INC:
		status = statusWithFlag(status, isa.C, addCarry(left, right))
		return statusWithFlag(status, isa.V, addOverflow(left, right, result))
	case isa.SUB, isa.CMP, isa.DEC:
		status = statusWithFlag(status, isa.C, subBorrow(left, right))
		return statusWithFlag(status, isa.V, subOverflow(left, right, result))
	default:
		status = statusWithFlag(status, isa.C, false)
		return statusWithFlag(status, isa.V, false)
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
	if _, ok := binaryOperations[opcode]; !ok {
		return fmt.Errorf("unknown alu opcode: %v", opcode)
	}

	c.appendStages(c.directReadStages()...)
	c.appendStages(c.stageSetACCFromBinaryOp(opcode))
	return nil
}

func (c *CPU) execUnaryOp(opcode isa.Opcode) error {
	if _, ok := unaryOperations[opcode]; !ok {
		return fmt.Errorf("unknown unary alu opcode: %v", opcode)
	}

	c.appendStages(
		c.stageSetACCFromUnaryOp(opcode),
	)
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

func (c *CPU) execCmpInstruction() error {
	c.appendStages(c.directReadStages()...)
	c.appendStages(c.stageSetSR(func(c *CPU) uint8 {
		left := c.ACC
		right := int32(c.DR)
		return statusWithALUFlags(c.SR, isa.CMP, left, right, left-right)
	}))
	return nil
}

func (c *CPU) stageSetACCFromBinaryOp(opcode isa.Opcode) stageOperation {
	return func(c *CPU) error {
		operation := binaryOperations[opcode]
		left := c.ACC
		right := int32(c.DR)
		result, err := operation(left, right)
		if err != nil {
			return err
		}

		c.ACC = result
		c.SR = statusWithALUFlags(c.SR, opcode, left, right, result)
		c.tick()
		return nil
	}
}

func (c *CPU) stageSetACCFromUnaryOp(opcode isa.Opcode) stageOperation {
	return func(c *CPU) error {
		operation := unaryOperations[opcode]
		left := c.ACC
		result, err := operation(left)
		if err != nil {
			return err
		}

		c.ACC = result
		c.SR = statusWithALUFlags(c.SR, opcode, left, unaryRightOperand(opcode), result)
		c.tick()
		return nil
	}
}
