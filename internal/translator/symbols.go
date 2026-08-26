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

	addr := t.nextDataAddr
	t.variables[name] = addr
	t.data[addr] = 0
	t.nextDataAddr += isa.Operand(isa.WordSize)
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
