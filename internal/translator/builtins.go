package translator

import (
	"cpu-simulator/internal/isa"
)

type BuiltinEmitter func(t *Translator) error

var builtins = map[string]BuiltinEmitter{
	"+": func(t *Translator) error {
		t.emitBinaryOp(isa.ADD)
		return nil
	},
	"-": func(t *Translator) error {
		t.emitBinaryOp(isa.SUB)
		return nil
	},
	"*": func(t *Translator) error {
		t.emitBinaryOp(isa.MUL)
		return nil
	},
	"/": func(t *Translator) error {
		t.emitBinaryOp(isa.DIV)
		return nil
	},
	"mod": func(t *Translator) error {
		t.emitBinaryOp(isa.MOD)
		return nil
	},

	"dup": func(t *Translator) error {
		t.emitDup()
		return nil
	},
	"drop": func(t *Translator) error {
		t.emitDrop()
		return nil
	},

	"emit": func(t *Translator) error {
		t.emitEmit()
		return nil
	},
	"key": func(t *Translator) error {
		t.emitKey()
		return nil
	},

	"ei": func(t *Translator) error {
		t.emitNoArg(isa.EI)
		return nil
	},
	"di": func(t *Translator) error {
		t.emitNoArg(isa.DI)
		return nil
	},
	"iret": func(t *Translator) error {
		t.emitNoArg(isa.IRET)
		return nil
	},
	"@": func(t *Translator) error {
		t.emitFetch()
		return nil
	},
	"!": func(t *Translator) error {
		t.emitStore()
		return nil
	},
}
