package machine

import (
	"cpu-simulator/internal/isa"
	"fmt"
)

func (c *CPU) pushData(value int32) error {
	if len(c.dataStack) >= int(isa.DataStackSize) {
		return fmt.Errorf("data stack overflow")
	}

	c.dataStack = append(c.dataStack, value)
	c.DSP++
	return nil
}

func (c *CPU) popData() (int32, error) {
	if len(c.dataStack) == 0 {
		return 0, fmt.Errorf("data stack underflow")
	}

	data := c.dataStack[len(c.dataStack)-1]
	c.dataStack = c.dataStack[:len(c.dataStack)-1]
	c.DSP--
	return data, nil
}

func (c *CPU) pushReturn(addr uint32) error {
	if len(c.returnStack) >= int(isa.ReturnStackSize) {
		return fmt.Errorf("return stack overflow")
	}

	c.returnStack = append(c.returnStack, addr)
	c.RSP++
	return nil
}

func (c *CPU) popReturn() (uint32, error) {
	if len(c.returnStack) == 0 {
		return 0, fmt.Errorf("return stack underflow")
	}

	addr := c.returnStack[len(c.returnStack)-1]
	c.returnStack = c.returnStack[:len(c.returnStack)-1]
	c.RSP--
	return addr, nil
}

func (c *CPU) peekData(offset uint32) (int32, error) {
	if offset >= uint32(len(c.dataStack)) {
		return 0, fmt.Errorf("data stack offset out of range: %d", offset)
	}

	index := len(c.dataStack) - 1 - int(offset)
	return c.dataStack[index], nil
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
