package translator

import (
	"cpu-simulator/internal/isa"
	"errors"
)

type scratch struct {
	tmpAddr isa.Operand
}

type printRuntime struct {
	entryAddr          isa.Operand
	valueAddr          isa.Operand
	divisorAddr        isa.Operand
	digitAddr          isa.Operand
	startedAddr        isa.Operand
	initialDivisorAddr isa.Operand
	tenAddr            isa.Operand
	asciiZeroAddr      isa.Operand
}

type Translator struct {
	code                   []isa.Instruction
	mainScratch            scratch
	interruptScratch       scratch
	printRuntime           printRuntime
	zeroAddr               isa.Operand
	variables              map[string]isa.Operand
	nextDataAddr           isa.Operand
	controlStack           []ControlFrame
	data                   map[isa.Operand]int32
	words                  map[string]isa.Operand
	interruptVectors       map[uint8]isa.Operand
	unresolvedCalls        []CallPatch
	insideWord             bool
	insideInterruptHandler bool
	mainStarted            bool
	entryJumpIndex         int
}

type CallPatch struct {
	name  string
	index int
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

func (t *Translator) emitPrintNumber() {
	t.emitCall(t.printRuntime.entryAddr)
}

func (t *Translator) emitPrintNumberRoutine() {
	rt := &t.printRuntime
	rt.entryAddr = isa.Operand(t.currentAddr())

	t.emitDrop()
	t.emit(isa.ST_ADDR, rt.valueAddr)
	t.emit(isa.CMP, t.zeroAddr)
	zeroJumpIndex := len(t.code)
	t.emit(isa.JZ, 0)

	t.emit(isa.LD_ADDR, rt.initialDivisorAddr)
	t.emit(isa.ST_ADDR, rt.divisorAddr)
	t.emit(isa.LD_IMM, 0)
	t.emit(isa.ST_ADDR, rt.startedAddr)

	loopAddr := t.currentAddr()
	t.emit(isa.LD_ADDR, rt.divisorAddr)
	t.emit(isa.CMP, t.zeroAddr)
	doneJumpIndex := len(t.code)
	t.emit(isa.JZ, 0)

	t.emit(isa.LD_ADDR, rt.valueAddr)
	t.emit(isa.DIV, rt.divisorAddr)
	t.emit(isa.ST_ADDR, rt.digitAddr)

	t.emit(isa.LD_ADDR, rt.digitAddr)
	t.emit(isa.CMP, t.zeroAddr)
	printJumpIndex := len(t.code)
	t.emit(isa.JNZ, 0)

	t.emit(isa.LD_ADDR, rt.startedAddr)
	t.emit(isa.CMP, t.zeroAddr)
	skipPrintJumpIndex := len(t.code)
	t.emit(isa.JZ, 0)

	printAddr := t.currentAddr()
	t.patchOperand(printJumpIndex, printAddr)
	t.emit(isa.LD_ADDR, rt.digitAddr)
	t.emit(isa.ADD, rt.asciiZeroAddr)
	t.emit(isa.OUT, isa.Operand(isa.PortOutput))
	t.emit(isa.LD_IMM, 1)
	t.emit(isa.ST_ADDR, rt.startedAddr)

	skipPrintAddr := t.currentAddr()
	t.patchOperand(skipPrintJumpIndex, skipPrintAddr)
	t.emit(isa.LD_ADDR, rt.valueAddr)
	t.emit(isa.MOD, rt.divisorAddr)
	t.emit(isa.ST_ADDR, rt.valueAddr)
	t.emit(isa.LD_ADDR, rt.divisorAddr)
	t.emit(isa.DIV, rt.tenAddr)
	t.emit(isa.ST_ADDR, rt.divisorAddr)
	t.emit(isa.JMP, isa.Operand(loopAddr))

	doneAddr := t.currentAddr()
	t.patchOperand(doneJumpIndex, doneAddr)
	t.emitReturn()

	zeroAddr := t.currentAddr()
	t.patchOperand(zeroJumpIndex, zeroAddr)
	t.emit(isa.LD_ADDR, rt.asciiZeroAddr)
	t.emit(isa.OUT, isa.Operand(isa.PortOutput))
	t.emitReturn()
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
