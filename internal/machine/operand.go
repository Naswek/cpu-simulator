package machine

import (
	"cpu-simulator/internal/isa"
	"fmt"
)

func (c *CPU) execLoad() error {
	switch c.IR.Opcode {
	case isa.LD_IMM:
		c.appendStages(
			c.stageSetACCAndZN(func(c *CPU) int32 {
				return int32(c.IR.Operand)
			}),
		)
	case isa.LD_ADDR:
		c.appendStages(c.directReadStages()...)
		c.appendStages(
			c.stageSetACCAndZN(func(c *CPU) int32 {
				return int32(c.DR)
			}),
		)
	case isa.LD_IND:
		c.appendStages(c.indirectReadStages()...)
		c.appendStages(
			c.stageSetACCAndZN(func(c *CPU) int32 {
				return int32(c.DR)
			}),
		)
	case isa.LD_SP_N:
		c.appendStages(
			c.stageSetAR(func(c *CPU) uint32 {
				offset := uint32(c.IR.Operand)
				if offset >= c.DSP {
					return 0
				}
				return (c.DSP - 1 - offset) * isa.WordSize
			}),
			(*CPU).stageReadDataStack,
			c.stageSetACCAndZN(func(c *CPU) int32 {
				return int32(c.DR)
			}),
		)
	default:
		return fmt.Errorf("unknown load opcode: %v", c.IR.Opcode)
	}
	return nil
}

func (c *CPU) directReadStages() []stageOperation {
	return []stageOperation{
		c.stageSetAR(func(c *CPU) uint32 {
			return uint32(c.IR.Operand)
		}),
		(*CPU).stageReadMemory,
	}
}

func (c *CPU) indirectReadStages() []stageOperation {
	return []stageOperation{
		c.stageSetAR(func(c *CPU) uint32 {
			return uint32(c.IR.Operand)
		}),
		(*CPU).stageReadMemory,
		c.stageSetAR(func(c *CPU) uint32 {
			return c.DR
		}),
		(*CPU).stageReadMemory,
	}
}

func (c *CPU) execStore() error {
	switch c.IR.Opcode {
	case isa.ST_ADDR:
		c.appendStages(
			c.stageSetAR(func(c *CPU) uint32 {
				return uint32(c.IR.Operand)
			}),
			(*CPU).stageSetDRFromACC,
			(*CPU).stageWriteMemory,
		)
	case isa.ST_IND:
		c.appendStages(
			c.stageSetAR(func(c *CPU) uint32 {
				return uint32(c.IR.Operand)
			}),
			(*CPU).stageReadMemory,
			c.stageSetAR(func(c *CPU) uint32 {
				return c.DR
			}),
			(*CPU).stageSetDRFromACC,
			(*CPU).stageWriteMemory,
		)
	default:
		return fmt.Errorf("unknown store opcode: %v", c.IR.Opcode)
	}
	return nil
}
