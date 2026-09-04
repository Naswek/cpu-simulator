package machine

import (
	"cpu-simulator/internal/isa"
	"fmt"
)

func (c *CPU) execPush() error {
	c.appendStages(
		(*CPU).stageSetDRFromACC,
		(*CPU).stageWriteDataStack,
		(*CPU).stageIncrementDSP,
	)
	return nil
}

func (c *CPU) execPop() error {
	c.appendStages(
		c.stageSetAR(func(c *CPU) uint32 {
			if c.DSP == 0 {
				return 0
			}
			return (c.DSP - 1) * isa.WordSize
		}),
		(*CPU).stageReadDataStack,
		c.stageSetACCAndZN(func(c *CPU) int32 {
			return int32(c.DR)
		}),
		(*CPU).stagePopDataStack,
	)
	return nil
}

func (c *CPU) peekData(offset uint32) (int32, error) {
	if offset >= uint32(len(c.dataStack)) {
		return 0, fmt.Errorf("data stack offset out of range: %d", offset)
	}

	index := len(c.dataStack) - 1 - int(offset)
	return c.dataStack[index], nil
}

func (c *CPU) stageWriteDataStack() error {
	if len(c.dataStack) >= int(isa.DataStackSize) {
		return fmt.Errorf("data stack overflow")
	}

	c.dataStack = append(c.dataStack, int32(c.DR))
	c.tick()
	return nil
}

func (c *CPU) stageReadDataStack() error {
	value, err := c.peekData(uint32(c.IR.Operand))
	if err != nil {
		return err
	}

	c.DR = uint32(value)
	c.tick()
	return nil
}

func (c *CPU) stagePopDataStack() error {
	if len(c.dataStack) == 0 {
		return fmt.Errorf("data stack underflow")
	}

	c.dataStack = c.dataStack[:len(c.dataStack)-1]
	c.DSP--
	c.tick()
	return nil
}

func (c *CPU) stageIncrementDSP() error {
	c.DSP++
	c.tick()
	return nil
}

func (c *CPU) stageWriteReturnStack() error {
	if len(c.returnStack) >= int(isa.ReturnStackSize) {
		return fmt.Errorf("return stack overflow")
	}

	c.returnStack = append(c.returnStack, c.PC)
	c.tick()
	return nil
}

func (c *CPU) stageReadReturnStack() error {
	if len(c.returnStack) == 0 {
		return fmt.Errorf("return stack underflow")
	}

	c.pendingAddr = c.returnStack[len(c.returnStack)-1]
	c.tick()
	return nil
}

func (c *CPU) stageSetDRFromPendingAddr() error {
	c.DR = c.pendingAddr
	c.tick()
	return nil
}

func (c *CPU) stageIncrementRSP() error {
	c.RSP++
	c.tick()
	return nil
}

func (c *CPU) stagePopReturnStack() error {
	if len(c.returnStack) == 0 {
		return fmt.Errorf("return stack underflow")
	}

	c.returnStack = c.returnStack[:len(c.returnStack)-1]
	c.RSP--
	c.tick()
	return nil
}

func (c *CPU) pushFrame(frame InterruptFrame) error {
	if len(c.interruptStack) >= int(isa.InterruptStackSize) {
		return fmt.Errorf("interrupt stack overflow")
	}

	c.interruptStack = append(c.interruptStack, frame)
	return nil
}

func (c *CPU) popFrame() (InterruptFrame, error) {
	if len(c.interruptStack) == 0 {
		return InterruptFrame{}, fmt.Errorf("interrupt stack underflow")
	}

	frame := c.interruptStack[len(c.interruptStack)-1]
	c.interruptStack = c.interruptStack[:len(c.interruptStack)-1]
	return frame, nil
}
