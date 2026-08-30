package machine

import (
	"cpu-simulator/internal/isa"
)

type InputEvent struct {
	Tick  uint64
	Value int32
}

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
	c.setFlag(isa.IM, false)
	addr, err := c.popReturn()
	if err != nil {
		return err
	}
	c.PC = addr
	return nil
}

func (c *CPU) interruptHandler() error {
	if c.flag(isa.IE) && !c.flag(isa.IM) && c.io.IRQPending() {
		addr, err := c.readWord(uint32(isa.InterruptVectorAddress(c.io.IRQVector())))
		if err != nil {
			return err
		}
		if err := c.pushReturn(c.PC); err != nil {
			return err
		}
		c.setFlag(isa.IM, true)
		c.io.ClearIRQ()
		c.PC = addr
	}
	return nil
}
