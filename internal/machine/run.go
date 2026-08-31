package machine

import (
	"fmt"
	"os"
)

func (c *CPU) Run(maxTicks int, logFile string) error {
	if logFile != "" {
		if err := os.WriteFile(logFile, []byte(c.Log()), 0644); err != nil {
			return err
		}

	}

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
