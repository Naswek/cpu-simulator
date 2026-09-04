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

	if err := c.io.Advance(c.tickCounter); err != nil {
		return err
	}

	if len(c.stages) == 0 {
		if c.hasPendingInterrupt() {
			c.startInterrupt()
		} else {
			c.startInstruction()
		}
	}

	if c.step < 0 || c.step >= len(c.stages) {
		return fmt.Errorf("unknown cpu step: %d", c.step)
	}

	operation := c.stages[c.step]
	if err := operation(c); err != nil {
		return err
	}

	c.step++
	if c.step >= len(c.stages) {
		c.step = 0
		c.stages = nil
		c.clearPending()
	}
	return nil
}

func (c *CPU) startInstruction() {
	c.stages = []stageOperation{
		(*CPU).fetchInstruction,
	}
	c.step = 0
}

func (c *CPU) fetchInstruction() error {
	word, err := c.readWord(c.PC)
	if err != nil {
		return err
	}

	c.IR = isa.DecodeInstruction(word)
	c.PC += isa.WordSize
	c.tick()
	return c.appendExecuteStages()
}

func (c *CPU) appendExecuteStages() error {
	var err error
	if _, ok := loadOpcodes[c.IR.Opcode]; ok {
		err = c.execLoad()
	} else if _, ok := storeOpcodes[c.IR.Opcode]; ok {
		err = c.execStore()
	} else if _, ok := binaryOperations[c.IR.Opcode]; ok {
		err = c.execBinaryOp(c.IR.Opcode)
	} else if _, ok := unaryOperations[c.IR.Opcode]; ok {
		err = c.execUnaryOp(c.IR.Opcode)
	} else if _, ok := jumpOpcodes[c.IR.Opcode]; ok {
		err = c.execJump(c.IR.Opcode)
	} else if operation, ok := controlOpcodes[c.IR.Opcode]; ok {
		err = operation(c)
	} else if operation, ok := ioOpcodes[c.IR.Opcode]; ok {
		err = operation(c)
	} else if operation, ok := interruptOpcodes[c.IR.Opcode]; ok {
		err = operation(c)
	} else {

		switch c.IR.Opcode {
		case isa.HALT:
			err = c.execHalt()
		case isa.NOP:
			err = c.execNop()
		case isa.POP:
			err = c.execPop()
		case isa.PUSH:
			err = c.execPush()
		case isa.CMP:
			err = c.execCmpInstruction()
		default:
			return fmt.Errorf("unknown opcode: %v", c.IR.Opcode)
		}
	}
	return err
}

func (c *CPU) tick() {
	c.tickCounter++
	c.appendLog()
}

func (c *CPU) appendStages(stages ...stageOperation) {
	c.stages = append(c.stages, stages...)
}

func (c *CPU) clearPending() {
	c.pendingAddr = 0
	c.pendingFrame = InterruptFrame{}
}

func (c *CPU) execHalt() error {
	c.appendStages((*CPU).stageSetHalted)
	return nil
}

func (c *CPU) execNop() error {
	c.appendStages((*CPU).stageNoop)
	return nil
}

func (c *CPU) stageSetAR(addr func(*CPU) uint32) stageOperation {
	return func(c *CPU) error {
		c.AR = addr(c)
		c.tick()
		return nil
	}
}

func (c *CPU) stageReadMemory() error {
	word, err := c.readWord(c.AR)
	if err != nil {
		return err
	}

	c.DR = word
	c.tick()
	return nil
}

func (c *CPU) stageWriteMemory() error {
	if err := c.writeWord(c.AR, c.DR); err != nil {
		return err
	}

	c.tick()
	return nil
}

func (c *CPU) stageSetDRFromACC() error {
	c.DR = uint32(c.ACC)
	c.tick()
	return nil
}

func (c *CPU) stageSetACC(value func(*CPU) int32) stageOperation {
	return func(c *CPU) error {
		c.ACC = value(c)
		c.tick()
		return nil
	}
}

func (c *CPU) stageSetACCAndZN(value func(*CPU) int32) stageOperation {
	return func(c *CPU) error {
		c.ACC = value(c)
		c.SR = statusWithZN(c.SR, c.ACC)
		c.tick()
		return nil
	}
}

func (c *CPU) stageSetSR(status func(*CPU) uint8) stageOperation {
	return func(c *CPU) error {
		c.SR = status(c)
		c.tick()
		return nil
	}
}

func (c *CPU) stageSetPC(addr func(*CPU) uint32) stageOperation {
	return func(c *CPU) error {
		if err := c.jump(isa.Operand(addr(c))); err != nil {
			return err
		}

		c.tick()
		return nil
	}
}

func (c *CPU) stageSetHalted() error {
	c.halted = true
	c.tick()
	return nil
}

func (c *CPU) stageNoop() error {
	c.tick()
	return nil
}
