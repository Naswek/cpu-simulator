package machine

import (
	"cpu-simulator/internal/isa"
	"fmt"
)

var interruptOpcodes = map[isa.Opcode]ControlOperation{
	isa.EI:   (*CPU).execEI,
	isa.DI:   (*CPU).execDI,
	isa.IRET: (*CPU).execIRET,
}

func (c *CPU) execEI() error {
	c.setFlag(isa.IE, true)
	return nil
}

func (c *CPU) execDI() error {
	c.setFlag(isa.IE, false)
	return nil
}

func (c *CPU) execIRET() error {
	addr, err := c.popReturn()
	if err != nil {
		return err
	}
	c.setFlag(isa.IM, false)
	c.PC = addr
	return nil
}

func (c *CPU) interruptHandler() error {
	if c.flag(isa.IE) && !c.flag(isa.IM) && c.io.IRQPending() {
		addr, err := c.readWord(uint32(isa.InterruptVectorAddress(c.io.IRQVector())))
		if err != nil {
			return err
		}
		if addr == 0 {
			return fmt.Errorf("interrupt vector %d is not initialized", c.io.IRQVector())
		}
		if err := c.pushReturn(c.PC); err != nil {
			return err
		}
		c.setFlag(isa.IM, true)
		c.io.ClearIRQ()
		return c.jump(isa.Operand(addr))
	}
	return nil
}
