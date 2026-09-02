package translator

import (
	"cpu-simulator/internal/isa"
	"errors"
)

func newTranslator() *Translator {
	code := make([]isa.Instruction, 1+isa.InterruptVectors)
	code[0] = isa.Instruction{Opcode: isa.JMP, Operand: 0}

	zeroAddr := isa.Operand(isa.MemSize - isa.WordSize)
	mainScratch := scratch{
		tmpAddr: isa.Operand(isa.MemSize - 2*isa.WordSize),
	}
	interruptScratch := scratch{
		tmpAddr: isa.Operand(isa.MemSize - 3*isa.WordSize),
	}
	printRuntime := newPrintRuntimeBelow(interruptScratch.tmpAddr)

	data := map[isa.Operand]int32{
		mainScratch.tmpAddr:      0,
		interruptScratch.tmpAddr: 0,
		zeroAddr:                 0,
	}
	for addr, value := range printRuntime.initialData() {
		data[addr] = value
	}
	translator := &Translator{
		code:             code,
		mainScratch:      mainScratch,
		interruptScratch: interruptScratch,
		printRuntime:     printRuntime,
		nextDataAddr:     printRuntime.lowestAddr(),
		zeroAddr:         zeroAddr,
		variables:        make(map[string]isa.Operand),
		data:             data,
		words:            make(map[string]isa.Operand),
		interruptVectors: make(map[uint8]isa.Operand),
		entryJumpIndex:   0,
	}
	translator.emitPrintNumberRoutine()
	return translator
}

func (t *Translator) startMain() {
	if t.mainStarted {
		return
	}

	t.patchOperand(t.entryJumpIndex, t.currentAddr())
	t.mainStarted = true
}

func Translate(source string) (Program, error) {
	translator := newTranslator()

	tokens, err := tokenize(source)
	if err != nil {
		return Program{}, err
	}

	for i := 0; i < len(tokens); i++ {
		t := tokens[i]
		if t.Kind == TokenWord {
			if handler, ok := specialHandlers[t.Text]; ok {
				i, err = handler(translator, tokens, i)
				if err != nil {
					return Program{}, err
				}
				continue
			}
		}

		if !translator.insideWord {
			translator.startMain()
		}

		err := translator.compileToken(t)
		if err != nil {
			return Program{}, err
		}
	}

	if translator.insideWord {
		return Program{}, errors.New("unclosed word definition")
	}

	if len(translator.controlStack) != 0 {
		return Program{}, errors.New("unclosed control frame")
	}

	if err = translator.resolveCalls(); err != nil {
		return Program{}, err
	}

	translator.resolveInterruptVectors()
	translator.startMain()
	translator.emitNoArg(isa.HALT)

	return Program{
		Code:             translator.code,
		Data:             translator.data,
		InterruptVectors: translator.interruptVectors,
	}, nil
}
