package machine

import (
	"fmt"
)

func (c *CPU) Run(maxTicks int) error {
	for !c.halted {
		if c.tickCounter >= maxTicks {
			return fmt.Errorf("tick limit exceeded")
		}

		if err := c.Step(); err != nil {
			return err
		}
	}
	return nil
}
