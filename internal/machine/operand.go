package machine

import (
	"cpu-simulator/internal/isa"
	"fmt"
)

type OperandReader func(c *CPU) (int32, error)
type OperandWriter func(c *CPU, value int32) error

var operandReaders = map[isa.Opcode]OperandReader{
	isa.LD_IMM:  readImmediate,
	isa.LD_ADDR: readDirect,
	isa.LD_IND:  readIndirect,
	isa.LD_SP_N: readStackOffset,

	isa.ADD: readDirect,
	isa.SUB: readDirect,
	isa.MUL: readDirect,
	isa.DIV: readDirect,
	isa.MOD: readDirect,
	isa.CMP: readDirect,
	isa.AND: readDirect,
	isa.OR:  readDirect,
}

func readImmediate(c *CPU) (int32, error) {
	return int32(c.IR.Operand), nil
}

func readDirect(c *CPU) (int32, error) {
	word, err := c.readWord(uint32(c.IR.Operand))
	if err != nil {
		return 0, err
	}

	return int32(word), nil
}

func readIndirect(c *CPU) (int32, error) {
	ptr, err := c.readWord(uint32(c.IR.Operand))
	if err != nil {
		return 0, err
	}

	word, err := c.readWord(ptr)
	if err != nil {
		return 0, err
	}
	return int32(word), nil
}

func readStackOffset(c *CPU) (int32, error) {
	return c.peekData(uint32(c.IR.Operand))
}

var operandWriters = map[isa.Opcode]OperandWriter{
	isa.ST_ADDR: writeDirect,
	isa.ST_IND:  writeIndirect,
}

func writeDirect(c *CPU, value int32) error {
	return c.writeWord(uint32(c.IR.Operand), uint32(value))
}

func writeIndirect(c *CPU, value int32) error {
	ptr, err := c.readWord(uint32(c.IR.Operand))
	if err != nil {
		return err
	}

	return c.writeWord(ptr, uint32(value))
}

func (c *CPU) readOperand() (int32, error) {
	reader, ok := operandReaders[c.IR.Opcode]
	if !ok {
		return 0, fmt.Errorf("no operand reader for opcode: %v", c.IR.Opcode)
	}
	return reader(c)
}

func (c *CPU) writeOperand(value int32) error {
	writer, ok := operandWriters[c.IR.Opcode]
	if !ok {
		return fmt.Errorf("no operand writer for opcode: %v", c.IR.Opcode)
	}

	return writer(c, value)
}

func (c *CPU) execLoad() error {
	value, err := c.readOperand()
	if err != nil {
		return err
	}

	c.ACC = value
	c.updateZN(c.ACC)
	return nil
}
