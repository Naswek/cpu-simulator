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

type InterruptFrame struct {
	PC  uint32
	ACC int32
	SR  uint8
}

func (c *CPU) execEI() error {
	c.appendStages(c.stageSetSR(func(c *CPU) uint8 {
		return statusWithFlag(c.SR, isa.IE, true)
	}))
	return nil
}

func (c *CPU) execDI() error {
	c.appendStages(c.stageSetSR(func(c *CPU) uint8 {
		return statusWithFlag(c.SR, isa.IE, false)
	}))
	return nil
}

func (c *CPU) execIRET() error {
	c.appendStages(
		(*CPU).stageReadInterruptFrame,
		c.stageSetACC(func(c *CPU) int32 {
			return c.pendingFrame.ACC
		}),
		c.stageSetPC(func(c *CPU) uint32 {
			return c.pendingFrame.PC
		}),
		c.stageSetSR(func(c *CPU) uint8 {
			return statusWithFlag(c.pendingFrame.SR, isa.IM, false)
		}),
	)
	return nil
}

func (c *CPU) hasPendingInterrupt() bool {
	return c.flag(isa.IE) && !c.flag(isa.IM) && c.io.IRQPending()
}

func (c *CPU) startInterrupt() {
	vector := c.io.IRQVector()
	c.stages = []stageOperation{
		c.stageSetAR(func(c *CPU) uint32 {
			return uint32(isa.InterruptVectorAddress(vector))
		}),
		(*CPU).stageReadMemory,
		(*CPU).stagePushInterruptFrame,
		c.stageSetSR(func(c *CPU) uint8 {
			return statusWithFlag(c.SR, isa.IM, true)
		}),
		c.stageSetPC(func(c *CPU) uint32 {
			return c.DR
		}),
		(*CPU).stageClearIRQ,
	}
	c.step = 0
}

func (c *CPU) stagePushInterruptFrame() error {
	if c.DR == 0 {
		return fmt.Errorf("interrupt vector %d is not initialized", c.io.IRQVector())
	}

	frame := InterruptFrame{
		PC:  c.PC,
		ACC: c.ACC,
		SR:  c.SR,
	}
	if err := c.pushFrame(frame); err != nil {
		return err
	}

	c.tick()
	return nil
}

func (c *CPU) stageReadInterruptFrame() error {
	frame, err := c.popFrame()
	if err != nil {
		return err
	}

	c.pendingFrame = frame
	c.tick()
	return nil
}

func (c *CPU) stageClearIRQ() error {
	c.io.ClearIRQ()
	c.tick()
	return nil
}
