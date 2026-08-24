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
	t.code = append(t.code, isa.Instruction {
		Opcode: opcode,
		Operand: 0,
	})
}

func (t *Translator) emitPushImm(value int) {
	t.code = append(t.code, isa.Instruction{
		Opcode: isa.LD_IMM,
		Operand: isa.Operand(value),
	})
	t.code = append(t.code, isa.Instruction{
		Opcode: isa.PUSH,
		Operand: 0,
	})
}

func (t *Translator) emitBinaryOp(opcode isa.Opcode) {
	t.emitDrop()
	t.code = append(t.code, isa.Instruction{
		Opcode: isa.ST_ADDR,
		Operand: t.tmpAddr,
	})
	t.emitDrop()
	t.code = append(t.code, isa.Instruction{
		Opcode: opcode,
		Operand: t.tmpAddr,
	})
	t.code = append(t.code, isa.Instruction{
		Opcode: isa.PUSH,
		Operand: 0,
	})
}

func (t *Translator) emitDup() {
	t.code = append(t.code, isa.Instruction{
		Opcode: isa.LD_SP_N,
		Operand: 0,
	})
	t.code = append(t.code, isa.Instruction{
		Opcode: isa.PUSH,
		Operand: 0,
	})
}

func (t *Translator) emitDrop() {
	t.code = append(t.code, isa.Instruction{
		Opcode: isa.POP,
		Operand: 0,
	})
}

func (t *Translator) emitKey() {
	t.code = append(t.code, isa.Instruction{
		Opcode: isa.IN,
		Operand: 0,
	})
	t.code = append(t.code, isa.Instruction{
		Opcode: isa.PUSH,
		Operand: 0,
	})
}

func (t *Translator) emitEmit() {
	
}