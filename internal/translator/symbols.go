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
	bytes := []byte(value)
	startAddr, err := t.allocateDataBlock(len(bytes) + 1)
	if err != nil {
		return 0, err
	}

	for i, ch := range bytes {
		addr := startAddr + isa.Operand(i)*isa.Operand(isa.WordSize)
		t.data[addr] = int32(ch)
	}

	terminatorAddr := startAddr + isa.Operand(len(bytes))*isa.Operand(isa.WordSize)
	t.data[terminatorAddr] = 0

	return startAddr, nil
}

func (t *Translator) allocateDataCell(value int32) (isa.Operand, error) {
	addr, err := t.allocateDataBlock(1)
	if err != nil {
		return 0, err
	}

	t.data[addr] = value
	return addr, nil
}

func (t *Translator) allocateDataBlock(cells int) (isa.Operand, error) {
	if cells <= 0 {
		return t.nextDataAddr, nil
	}

	size := isa.Operand(cells) * isa.Operand(isa.WordSize)
	if t.nextDataAddr < size {
		return 0, fmt.Errorf("data memory overflow: cannot allocate %d cells", cells)
	}

	startAddr := t.nextDataAddr - size
	for i := 0; i < cells; i++ {
		addr := startAddr + isa.Operand(i)*isa.Operand(isa.WordSize)
		if uint32(addr) >= isa.MemSize {
			return 0, fmt.Errorf("data address 0x%06X is outside memory", addr)
		}
		if _, ok := t.data[addr]; ok {
			return 0, fmt.Errorf("data address 0x%06X is already used", addr)
		}
	}

	t.nextDataAddr = startAddr
	return startAddr, nil
}
