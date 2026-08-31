package machine

import (
	"cpu-simulator/internal/isa"
	"fmt"
	"strings"
)

type LogEntry struct {
	Tick        int
	PC          uint32
	ACC         int32
	SR          uint8
	IR          isa.Instruction
	InInterrupt bool
	Output      string
}

func (c *CPU) appendLog() {
	c.trace = append(c.trace, LogEntry{
		Tick:        c.tickCounter,
		PC:          c.PC,
		ACC:         c.ACC,
		SR:          c.SR,
		IR:          c.IR,
		InInterrupt: c.flag(isa.IM),
		Output:      c.Output(),
	})
}

func (c *CPU) Log() string {
	var b strings.Builder

	for _, entry := range c.trace {
		fmt.Fprintf(
			&b,
			"tick=%d pc=0x%08X acc=%d sr=0x%02X ir=0x%08X interrupt=%t output=%q\n",
			entry.Tick,
			entry.PC,
			entry.ACC,
			entry.SR,
			entry.IR,
			entry.InInterrupt,
			entry.Output,
		)
	}

	return b.String()
}
