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

func newTranslator() *Translator {
	return &Translator{
		code:         make([]isa.Instruction, 0),
		tmpAddr:      startTmpAddr,
		tmpValueAddr: startTmpValueAddr,
		nextDataAddr: startNextDataAddr,
		zeroAddr:     zeroAddr,
		variables:    make(map[string]isa.Operand),
	}
}

func Translate(source string) ([]isa.Instruction, error) {
	translator := newTranslator()

	tokens, err := tokenize(source)
	if err != nil {
		return nil, err
	}

	for i := 0; i < len(tokens); i++ {
		t := tokens[i]
		isWord := t.Kind == TokenWord
		if isWord && t.Text == "variable" {
			if i+1 >= len(tokens) || tokens[i+1].Kind != TokenWord {
				return nil, errors.New("syntax error: variable definitions cannot be without variable name, variable cannot be number")
			}
			name := tokens[i+1]
			if _, builtinsHas := builtins[name.Text]; builtinsHas || name.Text == "variable" {
				return nil, errors.New("variable cannot have name like function words")
			}
			err = translator.defineVariable(name.Text)
			if err != nil {
				return nil, fmt.Errorf("cannot define variable %q: %w", name.Text, err)
			}
			i++
			continue
		}

		if isWord && t.Text == "if" {
			translator.emitIf()
			continue
		}

		if isWord && t.Text == "then" {
			err = translator.emitThen()
			if err != nil {
				return nil, err
			}
			continue
		}

		if isWord && t.Text == "else" {
			err = translator.emitElse()
			if err != nil {
				return nil, err
			}
			continue
		}

		err := sortToken(t, translator)
		if err != nil {
			return nil, err
		}
	}

	if len(translator.controlStack) != 0 {
		return nil, errors.New("unclosed control frame")
	}

	return translator.code, nil
}

func sortToken(tkn Token, trlr *Translator) error {
	switch tkn.Kind {
	case TokenNumber:
		trlr.emitPushImm(tkn.Value)
		return nil
	case TokenString:
		return errors.New("strings are not implemented")
	case TokenWord:
		if builtin, ok := builtins[tkn.Text]; ok {
			return builtin(trlr)
		} else {
			if addr, ok := trlr.lookupVariable(tkn.Text); ok {
				trlr.emitPushImm(int32(addr))
				return nil
			}
		}
		return fmt.Errorf("unknown word: %s", tkn.Text)
	default:
		return errors.New("unknown token type")
	}
}
