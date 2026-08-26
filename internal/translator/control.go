package translator

import (
	"cpu-simulator/internal/isa"
)

type ControlFrame struct {
	kind      ControlKind
	jumpIndex int
	startAddr uint32
}

type ControlKind uint8

const (
	If ControlKind = iota
	Else
	Begin
	While
)

func (t *Translator) currentAddr() uint32 {
	return uint32(len(t.code)) * isa.WordSize
}

func (t *Translator) patchOperand(index int, addr uint32) {
	t.code[index].Operand = isa.Operand(addr)
}
