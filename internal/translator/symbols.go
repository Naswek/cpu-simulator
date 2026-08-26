package translator

import (
	"cpu-simulator/internal/isa"
	"fmt"
)

func (t *Translator) defineVariable(name string) error {
	if _, ok := t.lookupVariable(name); !ok {
		t.variables[name] = t.nextDataAddr
		addr := t.nextDataAddr
		t.variables[name] = addr
		t.data[addr] = 0
		t.nextDataAddr += isa.Operand(isa.WordSize)
		return nil
	} else {
		return fmt.Errorf("variable %s already exists", name)
	}
}

func (t *Translator) lookupVariable(name string) (isa.Operand, bool) {
	addr, ok := t.variables[name]
	if !ok {
		return 0, false
	}

	return addr, true
}
