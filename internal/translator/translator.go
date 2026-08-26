package translator

import (
	"cpu-simulator/internal/isa"
	"errors"
	"fmt"
)

const (
	startTmpAddr      = 0x400
	startTmpValueAddr = 0x404
	zeroAddr          = 0x408
	startNextDataAddr = 0x40C
)

var specialWords = map[string]struct{}{
	"variable": {},
	":":        {},
	";":        {},
	"if":       {},
	"then":     {},
	"else":     {},
	"begin":    {},
	"until":    {},
	"again":    {},
	"while":    {},
	"repeat":   {},
}

type Program struct {
	Code []isa.Instruction
	Data map[isa.Operand]int32
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

	image := make([]uint32, maxLen)
	for i, instr := range p.Code {
		image[i] = isa.EncodeInstruction(instr)
	}
	for addr, value := range p.Data {
		image[int(addr/isa.Operand(isa.WordSize))] = uint32(value)
	}

	return image, nil
}

func newTranslator() *Translator {
	code := []isa.Instruction{
		{Opcode: isa.JMP, Operand: 0},
	}
	data := map[isa.Operand]int32{
		startTmpAddr:      0,
		startTmpValueAddr: 0,
		zeroAddr:          0,
	}
	return &Translator{
		code:           code,
		tmpAddr:        startTmpAddr,
		tmpValueAddr:   startTmpValueAddr,
		nextDataAddr:   startNextDataAddr,
		zeroAddr:       zeroAddr,
		variables:      make(map[string]isa.Operand),
		data:           data,
		words:          make(map[string]isa.Operand),
		entryJumpIndex: 0,
	}
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

	if _, ok := specialWords[name]; ok {
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
			switch t.Text {
			case "variable":
				if translator.insideWord {
					return Program{}, errors.New("variable definition inside word is not supported")
				}
				if i+1 >= len(tokens) || tokens[i+1].Kind != TokenWord {
					return Program{}, errors.New("syntax error: variable definitions cannot be without variable name, variable cannot be number")
				}
				name := tokens[i+1]
				if err = translator.validateSymbolName(name.Text); err != nil {
					return Program{}, fmt.Errorf("cannot define variable %q: %w", name.Text, err)
				}
				if err = translator.defineVariable(name.Text); err != nil {
					return Program{}, fmt.Errorf("cannot define variable %q: %w", name.Text, err)
				}
				i++
				continue
			case ":":
				if translator.insideWord {
					return Program{}, errors.New("nested word definitions are not supported")
				}
				if translator.mainStarted {
					return Program{}, errors.New("word definitions must appear before executable code")
				}
				if len(translator.controlStack) != 0 {
					return Program{}, errors.New("word definition inside control frame is not supported")
				}
				if i+1 >= len(tokens) || tokens[i+1].Kind != TokenWord {
					return Program{}, errors.New("syntax error: word definition cannot be without word name")
				}
				name := tokens[i+1]
				if err = translator.validateSymbolName(name.Text); err != nil {
					return Program{}, fmt.Errorf("cannot define word %q: %w", name.Text, err)
				}
				if err = translator.defineWord(name.Text); err != nil {
					return Program{}, fmt.Errorf("cannot define word %q: %w", name.Text, err)
				}
				translator.insideWord = true
				i++
				continue
			case ";":
				if !translator.insideWord {
					return Program{}, errors.New("semicolon without matching word definition")
				}
				if len(translator.controlStack) != 0 {
					return Program{}, errors.New("unclosed control frame inside word definition")
				}
				translator.emitReturn()
				translator.insideWord = false
				continue
			case "if":
				if !translator.insideWord {
					translator.startMain()
				}
				translator.emitIf()
				continue
			case "then":
				if !translator.insideWord {
					translator.startMain()
				}
				err = translator.emitThen()
				if err != nil {
					return Program{}, err
				}
				continue
			case "else":
				if !translator.insideWord {
					translator.startMain()
				}
				err = translator.emitElse()
				if err != nil {
					return Program{}, err
				}
				continue
			case "begin":
				if !translator.insideWord {
					translator.startMain()
				}
				translator.emitBegin()
				continue
			case "until":
				if !translator.insideWord {
					translator.startMain()
				}
				err = translator.emitUntil()
				if err != nil {
					return Program{}, err
				}
				continue
			case "again":
				if !translator.insideWord {
					translator.startMain()
				}
				err = translator.emitAgain()
				if err != nil {
					return Program{}, err
				}
				continue
			case "while":
				if !translator.insideWord {
					translator.startMain()
				}
				err = translator.emitWhile()
				if err != nil {
					return Program{}, err
				}
				continue
			case "repeat":
				if !translator.insideWord {
					translator.startMain()
				}
				err = translator.emitRepeat()
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

	translator.startMain()
	translator.emitNoArg(isa.HALT)

	return Program{
		Code: translator.code,
		Data: translator.data,
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
