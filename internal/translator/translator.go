package translator

import (
	"cpu-simulator/internal/isa"
	"errors"
	"fmt"
)

type Program struct {
	Code             []isa.Instruction
	Data             map[isa.Operand]int32
	InterruptVectors map[uint8]isa.Operand
}

func (p Program) CodeWords() ([]uint32, error) {
	words := isa.EncodeProgram(p.Code)

	for vector, addr := range p.InterruptVectors {
		if vector >= isa.InterruptVectors {
			return nil, fmt.Errorf("interrupt vector %d is outside vector table", vector)
		}

		index := int(isa.InterruptVectorAddress(vector) / isa.Operand(isa.WordSize))
		if index >= len(words) {
			return nil, fmt.Errorf("interrupt vector %d is outside code image", vector)
		}
		words[index] = uint32(addr)
	}

	return words, nil
}

func (p Program) MemoryImage() ([]uint32, error) {
	maxLen := len(p.Code)

	for addr := range p.Data {
		if addr%isa.Operand(isa.WordSize) != 0 {
			return nil, fmt.Errorf("unaligned data address 0x%06X", addr)
		}
		if uint32(addr) >= isa.MemSize {
			return nil, fmt.Errorf("data address 0x%06X is outside memory", addr)
		}

		index := int(addr / isa.Operand(isa.WordSize))
		if index < len(p.Code) {
			return nil, fmt.Errorf("data address 0x%06X overlaps code", addr)
		}
		if index+1 > maxLen {
			maxLen = index + 1
		}
	}

	for i := range p.Code {
		addr := uint32(i) * isa.WordSize
		if addr >= isa.MemSize {
			return nil, fmt.Errorf("code address 0x%06X is outside memory", addr)
		}
	}

	code, err := p.CodeWords()
	if err != nil {
		return nil, err
	}

	image := make([]uint32, maxLen)
	for i, word := range code {
		image[i] = word
	}
	for addr, value := range p.Data {
		image[int(addr/isa.Operand(isa.WordSize))] = uint32(value)
	}

	return image, nil
}

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
	printRuntime := printRuntime{
		valueAddr:          isa.Operand(isa.MemSize - 4*isa.WordSize),
		divisorAddr:        isa.Operand(isa.MemSize - 5*isa.WordSize),
		digitAddr:          isa.Operand(isa.MemSize - 6*isa.WordSize),
		startedAddr:        isa.Operand(isa.MemSize - 7*isa.WordSize),
		initialDivisorAddr: isa.Operand(isa.MemSize - 8*isa.WordSize),
		tenAddr:            isa.Operand(isa.MemSize - 9*isa.WordSize),
		asciiZeroAddr:      isa.Operand(isa.MemSize - 10*isa.WordSize),
	}

	data := map[isa.Operand]int32{
		mainScratch.tmpAddr:             0,
		interruptScratch.tmpAddr:        0,
		printRuntime.valueAddr:          0,
		printRuntime.divisorAddr:        0,
		printRuntime.digitAddr:          0,
		printRuntime.startedAddr:        0,
		printRuntime.initialDivisorAddr: 1000000000,
		printRuntime.tenAddr:            10,
		printRuntime.asciiZeroAddr:      48,
		zeroAddr:                        0,
	}
	translator := &Translator{
		code:             code,
		mainScratch:      mainScratch,
		interruptScratch: interruptScratch,
		printRuntime:     printRuntime,
		nextDataAddr:     printRuntime.asciiZeroAddr,
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

func (t *Translator) validateSymbolName(name string) error {
	if _, ok := builtins[name]; ok {
		return errors.New("name conflicts with builtin word")
	}

	if isSpecialWord(name) {
		return errors.New("name conflicts with special word")
	}

	if _, ok := t.lookupVariable(name); ok {
		return fmt.Errorf("name %s already defined as variable", name)
	}

	if _, ok := t.lookupWord(name); ok {
		return fmt.Errorf("name %s already defined as word", name)
	}

	return nil
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

		err := sortToken(t, translator)
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

func sortToken(tkn Token, trlr *Translator) error {
	switch tkn.Kind {
	case TokenNumber:
		trlr.emitPushImm(tkn.Value)
		return nil
	case TokenString:
		addr, err := trlr.defineString(tkn.Text)
		if err != nil {
			return err
		}
		trlr.emitPushImm(int32(addr))
		return nil
	case TokenWord:
		if builtin, ok := builtins[tkn.Text]; ok {
			return builtin(trlr)
		}

		if addr, ok := trlr.lookupVariable(tkn.Text); ok {
			trlr.emitPushImm(int32(addr))
			return nil
		}

		if addr, ok := trlr.lookupWord(tkn.Text); ok {
			trlr.emitCall(addr)
			return nil
		}

		if trlr.insideWord {
			trlr.emitUnresolvedCall(tkn.Text)
			return nil
		}

		return fmt.Errorf("unknown word: %s", tkn.Text)
	default:
		return errors.New("unknown token type")
	}
}
