package translator

import (
	"cpu-simulator/internal/isa"
)

type Translator struct {
	code []isa.Instruction
	tmpAddr isa.Operand
}

func (t *Translator) emit(opcode isa.Opcode, operand isa.Operand) {
	t.code = append(t.code, isa.Instruction{
		Opcode: opcode,
		Operand: operand,
	})
}
func (t *Translator) emitNoArg(opcode isa.Opcode) {
	t.emit(opcode, 0)
}

func (t *Translator) emitPushImm(value int) {
	t.emit(isa.LD_IMM, isa.Operand(value))
	t.emitNoArg(isa.PUSH)
}

func (t *Translator) emitBinaryOp(opcode isa.Opcode) {
	t.emitDrop()
	t.emit(isa.ST_ADDR, t.tmpAddr)
	t.emitDrop()
	t.emit(opcode, t.tmpAddr)
	t.emitNoArg(isa.PUSH)
}

func (t *Translator) emitDup() {
	t.emit(isa.LD_SP_N, 0)
	t.emitNoArg(isa.PUSH)
}

func (t *Translator) emitDrop() {
	t.emitNoArg(isa.POP)
}

func (t *Translator) emitKey(port uint8) {
	t.emit(isa.IN, isa.Operand(port))
	t.emitNoArg(isa.PUSH)
}

func (t *Translator) emitEmit(port uint8) {
	t.emitDrop()
	t.emit(isa.OUT, isa.Operand(port))
}