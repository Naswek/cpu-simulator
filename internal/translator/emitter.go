package translator

import (
	"cpu-simulator/internal/isa"
	"errors"
)

func (t *Translator) emit(opcode isa.Opcode, operand isa.Operand) {
	t.code = append(t.code, isa.Instruction{
		Opcode:  opcode,
		Operand: operand,
	})
}

func (t *Translator) emitNoArg(opcode isa.Opcode) {
	t.emit(opcode, 0)
}

func (t *Translator) emitCall(addr isa.Operand) {
	t.emit(isa.CALL, addr)
}

func (t *Translator) emitUnresolvedCall(name string) {
	index := len(t.code)
	t.emitCall(0)
	t.unresolvedCalls = append(t.unresolvedCalls, CallPatch{
		name:  name,
		index: index,
	})
}

func (t *Translator) emitReturn() {
	t.emitNoArg(isa.RET)
}

func (t *Translator) emitPushImm(value int32) {
	t.emit(isa.LD_IMM, isa.Operand(value))
	t.emitNoArg(isa.PUSH)
}

func (t *Translator) currentTmpAddr() isa.Operand {
	if t.insideInterruptHandler {
		return t.interruptScratch.tmpAddr
	}
	return t.mainScratch.tmpAddr
}

func (t *Translator) emitBinaryOp(opcode isa.Opcode) {
	tmpAddr := t.currentTmpAddr()
	t.emitDrop()
	t.emit(isa.ST_ADDR, tmpAddr)
	t.emitDrop()
	t.emit(opcode, tmpAddr)
	t.emitNoArg(isa.PUSH)
}

func (t *Translator) emitCompareOp(jumpIfTrue isa.Opcode) {
	tmpAddr := t.currentTmpAddr()
	t.emitDrop()
	t.emit(isa.ST_ADDR, tmpAddr)
	t.emitDrop()
	t.emit(isa.CMP, tmpAddr)
	t.emitBoolFromJump(jumpIfTrue)
}

func (t *Translator) emitZeroCompare(jumpIfTrue isa.Opcode) {
	t.emitDrop()
	t.emit(isa.CMP, t.zeroAddr)
	t.emitBoolFromJump(jumpIfTrue)
}

func (t *Translator) emitBoolFromJump(jumpIfTrue isa.Opcode) {
	trueJumpIndex := len(t.code)
	t.emit(jumpIfTrue, 0)

	t.emitPushImm(0)
	endJumpIndex := len(t.code)
	t.emit(isa.JMP, 0)

	trueAddr := t.currentAddr()
	t.patchOperand(trueJumpIndex, trueAddr)
	t.emitPushImm(1)

	endAddr := t.currentAddr()
	t.patchOperand(endJumpIndex, endAddr)
}

func (t *Translator) emitInc() {
	t.emitDrop()
	t.emitNoArg(isa.INC)
	t.emitNoArg(isa.PUSH)
}

func (t *Translator) emitCellInc() {
	t.emitDrop()
	for i := uint32(0); i < isa.WordSize; i++ {
		t.emitNoArg(isa.INC)
	}
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
	tmpAddr := t.currentTmpAddr()
	t.emitDrop()
	t.emit(isa.ST_ADDR, tmpAddr)
	t.emit(isa.LD_IND, tmpAddr)
	t.emitNoArg(isa.PUSH)
}

func (t *Translator) emitStore() {
	tmpAddr := t.currentTmpAddr()
	t.emitDrop()
	t.emit(isa.ST_ADDR, tmpAddr)
	t.emitDrop()
	t.emit(isa.ST_IND, tmpAddr)
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

func (t *Translator) emitAgain() error {
	if len(t.controlStack) == 0 {
		return errors.New("again without matching begin")
	}

	frame := t.controlStack[len(t.controlStack)-1]
	if frame.kind != Begin {
		return errors.New("again closes wrong frame")
	}

	t.controlStack = t.controlStack[:len(t.controlStack)-1]
	t.emit(isa.JMP, isa.Operand(frame.startAddr))
	return nil
}

func (t *Translator) emitWhile() error {
	if len(t.controlStack) == 0 {
		return errors.New("while without matching begin")
	}

	frame := t.controlStack[len(t.controlStack)-1]
	if frame.kind != Begin {
		return errors.New("while closes wrong frame")
	}

	t.emitDrop()
	t.emit(isa.CMP, t.zeroAddr)

	jumpIndex := len(t.code)
	t.emit(isa.JZ, 0)

	t.controlStack[len(t.controlStack)-1] = ControlFrame{
		kind:      While,
		jumpIndex: jumpIndex,
		startAddr: frame.startAddr,
	}
	return nil
}

func (t *Translator) emitRepeat() error {
	if len(t.controlStack) == 0 {
		return errors.New("repeat without matching begin")
	}

	frame := t.controlStack[len(t.controlStack)-1]
	if frame.kind != While {
		return errors.New("repeat closes wrong frame")
	}

	t.controlStack = t.controlStack[:len(t.controlStack)-1]
	t.emit(isa.JMP, isa.Operand(frame.startAddr))
	addr := t.currentAddr()
	t.patchOperand(frame.jumpIndex, addr)
	return nil
}
