package translator

import (
	"errors"
	"fmt"
)

func (t *Translator) compileToken(tkn Token) error {
	switch tkn.Kind {
	case TokenNumber:
		t.emitPushImm(tkn.Value)
		return nil
	case TokenString:
		addr, err := t.defineString(tkn.Text)
		if err != nil {
			return err
		}
		t.emitPushImm(int32(addr))
		return nil
	case TokenWord:
		return t.compileWord(tkn.Text)
	default:
		return errors.New("unknown token type")
	}
}

func (t *Translator) compileWord(word string) error {
	if builtin, ok := builtins[word]; ok {
		return builtin(t)
	}

	if addr, ok := t.lookupVariable(word); ok {
		t.emitPushImm(int32(addr))
		return nil
	}

	if addr, ok := t.lookupWord(word); ok {
		t.emitCall(addr)
		return nil
	}

	if t.insideWord {
		t.emitUnresolvedCall(word)
		return nil
	}

	return fmt.Errorf("unknown word: %s", word)
}
