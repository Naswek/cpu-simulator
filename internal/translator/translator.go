package translator

import (
	"cpu-simulator/internal/isa"
	"errors"
	"fmt"
)

func newTranslator() *Translator {
	return &Translator{
		code: make([]isa.Instruction, 0),
		tmpAddr: 0x400,
	}
}

func Translate(source string) ([]isa.Instruction, error) {
	translator := newTranslator()

	tokens, err := tokenize(source)
	if err != nil {
		return nil, err
	}

	for _, t := range tokens {
		err := sortToken(t, translator)
		if err != nil {
			return nil, err
		}
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
				return fmt.Errorf("unknown word: %s", tkn.Text)
			}
		default:
			return errors.New("unknown token type")
	}
}