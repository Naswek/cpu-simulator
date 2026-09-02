package translator

import "cpu-simulator/internal/isa"

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

func newPrintRuntimeBelow(addr isa.Operand) printRuntime {
	cellBelow := func(offset int) isa.Operand {
		return addr - isa.Operand(offset)*isa.Operand(isa.WordSize)
	}

	return printRuntime{
		valueAddr:          cellBelow(1),
		divisorAddr:        cellBelow(2),
		digitAddr:          cellBelow(3),
		startedAddr:        cellBelow(4),
		initialDivisorAddr: cellBelow(5),
		tenAddr:            cellBelow(6),
		asciiZeroAddr:      cellBelow(7),
	}
}

func (rt printRuntime) initialData() map[isa.Operand]int32 {
	return map[isa.Operand]int32{
		rt.valueAddr:          0,
		rt.divisorAddr:        0,
		rt.digitAddr:          0,
		rt.startedAddr:        0,
		rt.initialDivisorAddr: 1000000000,
		rt.tenAddr:            10,
		rt.asciiZeroAddr:      48,
	}
}

func (rt printRuntime) lowestAddr() isa.Operand {
	return rt.asciiZeroAddr
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
