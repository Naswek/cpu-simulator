package translator

import (
	"errors"
	"fmt"
)

type SpecialHandler func(t *Translator, tokens []Token, i int) (int, error)

var specialHandlers = map[string]SpecialHandler{
	"variable": handleVariable,
	":":        handleWordStart,
	";":        handleWordEnd,
	"if":       handleIf,
	"then":     handleThen,
	"else":     handleElse,
	"begin":    handleBegin,
	"until":    handleUntil,
	"again":    handleAgain,
	"while":    handleWhile,
	"repeat":   handleRepeat,
}

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

func isSpecialWord(name string) bool {
	_, ok := specialWords[name]
	return ok
}

func handleVariable(t *Translator, tokens []Token, i int) (int, error) {
	if t.insideWord {
		return i, errors.New("variable definition inside word is not supported")
	}

	name, err := nextWord(tokens, i, "syntax error: variable definitions cannot be without variable name, variable cannot be number")
	if err != nil {
		return i, err
	}

	if err := t.validateSymbolName(name.Text); err != nil {
		return i, fmt.Errorf("cannot define variable %q: %w", name.Text, err)
	}
	if err := t.defineVariable(name.Text); err != nil {
		return i, fmt.Errorf("cannot define variable %q: %w", name.Text, err)
	}

	return i + 1, nil
}

func handleWordStart(t *Translator, tokens []Token, i int) (int, error) {
	if t.insideWord {
		return i, errors.New("nested word definitions are not supported")
	}
	if t.mainStarted {
		return i, errors.New("word definitions must appear before executable code")
	}
	if len(t.controlStack) != 0 {
		return i, errors.New("word definition inside control frame is not supported")
	}

	name, err := nextWord(tokens, i, "syntax error: word definition cannot be without word name")
	if err != nil {
		return i, err
	}

	if err := t.validateSymbolName(name.Text); err != nil {
		return i, fmt.Errorf("cannot define word %q: %w", name.Text, err)
	}
	if err := t.defineWord(name.Text); err != nil {
		return i, fmt.Errorf("cannot define word %q: %w", name.Text, err)
	}

	t.insideWord = true
	t.insideInterruptHandler = name.Text == "handle_input"

	return i + 1, nil
}

func handleWordEnd(t *Translator, _ []Token, i int) (int, error) {
	if !t.insideWord {
		return i, errors.New("semicolon without matching word definition")
	}
	if len(t.controlStack) != 0 {
		return i, errors.New("unclosed control frame inside word definition")
	}

	t.emitReturn()
	t.insideWord = false
	t.insideInterruptHandler = false

	return i, nil
}

func handleIf(t *Translator, _ []Token, i int) (int, error) {
	t.startMainIfNeeded()
	t.emitIf()
	return i, nil
}

func handleThen(t *Translator, _ []Token, i int) (int, error) {
	t.startMainIfNeeded()
	return i, t.emitThen()
}

func handleElse(t *Translator, _ []Token, i int) (int, error) {
	t.startMainIfNeeded()
	return i, t.emitElse()
}

func handleBegin(t *Translator, _ []Token, i int) (int, error) {
	t.startMainIfNeeded()
	t.emitBegin()
	return i, nil
}

func handleUntil(t *Translator, _ []Token, i int) (int, error) {
	t.startMainIfNeeded()
	return i, t.emitUntil()
}

func handleAgain(t *Translator, _ []Token, i int) (int, error) {
	t.startMainIfNeeded()
	return i, t.emitAgain()
}

func handleWhile(t *Translator, _ []Token, i int) (int, error) {
	t.startMainIfNeeded()
	return i, t.emitWhile()
}

func handleRepeat(t *Translator, _ []Token, i int) (int, error) {
	t.startMainIfNeeded()
	return i, t.emitRepeat()
}

func nextWord(tokens []Token, i int, message string) (Token, error) {
	if i+1 >= len(tokens) || tokens[i+1].Kind != TokenWord {
		return Token{}, errors.New(message)
	}

	return tokens[i+1], nil
}

func (t *Translator) startMainIfNeeded() {
	if !t.insideWord {
		t.startMain()
	}
}
