package translator

import "cpu-simulator/internal/isa"

type scratch struct {
	tmpAddr isa.Operand
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
