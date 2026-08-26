package translator

import (
	"cpu-simulator/internal/isa"
	"fmt"
)

func (t *Translator) defineVariable(name string) error {
	if _, ok := t.lookupVariable(name); ok {
		return fmt.Errorf("variable %s already exists", name)
	}

	if _, ok := t.lookupWord(name); ok {
		return fmt.Errorf("variable %s conflicts with word", name)
	}

	addr, err := t.allocateDataCell(0)
	if err != nil {
		return err
	}

	t.variables[name] = addr
	return nil
}

func (t *Translator) lookupVariable(name string) (isa.Operand, bool) {
	addr, ok := t.variables[name]
	if !ok {
		return 0, false
	}

	return addr, true
}

func (t *Translator) defineWord(name string) error {
	if _, ok := t.lookupWord(name); ok {
		return fmt.Errorf("word %s already exists", name)
	}

	if _, ok := t.lookupVariable(name); ok {
		return fmt.Errorf("word %s conflicts with variable", name)
	}

	t.words[name] = isa.Operand(t.currentAddr())
	return nil
}

func (t *Translator) lookupWord(name string) (isa.Operand, bool) {
	addr, ok := t.words[name]
	if !ok {
		return 0, false
	}

	return addr, true
}

func (t *Translator) resolveCalls() error {
	for _, patch := range t.unresolvedCalls {
		addr, ok := t.lookupWord(patch.name)
		if !ok {
			return fmt.Errorf("unknown word: %s", patch.name)
		}
		t.code[patch.index].Operand = addr
	}

	return nil
}

func (t *Translator) defineString(value string) (isa.Operand, error) {
	startAddr := t.nextDataAddr

	for _, ch := range []byte(value) {
		if _, err := t.allocateDataCell(int32(ch)); err != nil {
			return 0, err
		}
	}

	if _, err := t.allocateDataCell(0); err != nil {
		return 0, err
	}

	return startAddr, nil
}

func (t *Translator) allocateDataCell(value int32) (isa.Operand, error) {
	addr := t.nextDataAddr
	if uint32(addr) >= isa.MemSize {
		return 0, fmt.Errorf("data memory overflow at address 0x%06X", addr)
	}

	if _, ok := t.data[addr]; ok {
		return 0, fmt.Errorf("data address 0x%06X is already used", addr)
	}

	t.data[addr] = value
	t.nextDataAddr += isa.Operand(isa.WordSize)
	return addr, nil
}
