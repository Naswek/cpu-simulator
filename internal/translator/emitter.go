package translator

import (
	"cpu-simulator/internal/isa"
	"errors"
)

type Translator struct {
	code         []isa.Instruction
	tmpAddr      isa.Operand
	tmpValueAddr isa.Operand
	zeroAddr     isa.Operand
	variables    map[string]isa.Operand
	nextDataAddr isa.Operand
	controlStack []ControlFrame
}

func (t *Translator) emit(opcode isa.Opcode, operand isa.Operand) {
	t.code = append(t.code, isa.Instruction{
		Opcode:  opcode,
		Operand: operand,
	})
}
func (t *Translator) emitNoArg(opcode isa.Opcode) {
	t.emit(opcode, 0)
}

func (t *Translator) emitPushImm(value int32) {
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

func (t *Translator) emitKey() {
	t.emit(isa.IN, isa.Operand(isa.PortInput))
	t.emitNoArg(isa.PUSH)
}

func (t *Translator) emitEmit() {
	t.emitDrop()
	t.emit(isa.OUT, isa.Operand(isa.PortOutput))
}

func (t *Translator) emitFetch() {
	t.emitDrop()
	t.emit(isa.ST_ADDR, t.tmpAddr)
	t.emit(isa.LD_IND, t.tmpAddr)
	t.emitNoArg(isa.PUSH)
}

func (t *Translator) emitStore() {
	t.emitDrop()
	t.emit(isa.ST_ADDR, t.tmpAddr)
	t.emitDrop()
	t.emit(isa.ST_IND, t.tmpAddr)
}

func (t *Translator) emitIf() {
	t.emitDrop()
	t.emit(isa.CMP, t.zeroAddr)

	jumpIndex := len(t.code)
	t.emit(isa.JZ, 0)

	t.controlStack = append(t.controlStack, ControlFrame{
		kind:      If,
		jumpIndex: jumpIndex,
	})
}

func (t *Translator) emitThen() error {
	if len(t.controlStack) == 0 {
		return errors.New("then without matching if")
	}

	frame := t.controlStack[len(t.controlStack)-1]
	if frame.kind != If && frame.kind != Else {
		return errors.New("then closes wrong frame")
	}

	t.controlStack = t.controlStack[:len(t.controlStack)-1]
	addr := t.currentAddr()
	t.patchOperand(frame.jumpIndex, addr)
	return nil
}

func (t *Translator) emitElse() error {
	if len(t.controlStack) == 0 {
		return errors.New("else without matching if")
	}

	frame := t.controlStack[len(t.controlStack)-1]
	if frame.kind != If {
		return errors.New("else closes wrong frame")
	}

	oldJumpIndex := frame.jumpIndex
	newJumpIndex := len(t.code)
	t.emit(isa.JMP, 0)
	addr := t.currentAddr()
	t.patchOperand(oldJumpIndex, addr)
	t.controlStack[len(t.controlStack)-1] = ControlFrame{
		kind:      Else,
		jumpIndex: newJumpIndex,
	}
	return nil
}

func (t *Translator) emitBegin() {
	startAddr := t.currentAddr()
	t.controlStack = append(t.controlStack, ControlFrame{
		kind:      Begin,
		startAddr: startAddr,
	})
}

func (t *Translator) emitUntil() error {
	if len(t.controlStack) == 0 {
		return errors.New("until without matching begin")
	}

	frame := t.controlStack[len(t.controlStack)-1]
	if frame.kind != Begin {
		return errors.New("until closes wrong frame")
	}

	t.controlStack = t.controlStack[:len(t.controlStack)-1]

	t.emitDrop()
	t.emit(isa.CMP, t.zeroAddr)
	t.emit(isa.JZ, isa.Operand(frame.startAddr))
	return nil
}
