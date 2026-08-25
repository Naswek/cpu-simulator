package isa

import (
	"fmt"
	"os"
	"strings"
)

func DisassembleInstruction(instr Instruction) string {
	info, ok := OpcodeTable[instr.Opcode]
	if !ok {
		return fmt.Sprintf(".word opcode=0x%06X operand=0x%06X", instr.Opcode, instr.Operand)
	}

	kind := info.OperandKind
	switch kind {
		case none:
			return fmt.Sprintf("%v", info.Mnemonic)
		case addr:
			return fmt.Sprintf("%v 0x%06X", info.Mnemonic, instr.Operand)
		case imm:
			return fmt.Sprintf("%v #%v", info.Mnemonic, instr.Operand)
		case stack_offset:
			return fmt.Sprintf("%v %v", info.Mnemonic, instr.Operand)
		case port:
			deviceInfo, ok := PortsTable[uint8(instr.Operand)]
			if !ok {
				return fmt.Sprintf("%v port %v", info.Mnemonic, instr.Operand)
			}
			return fmt.Sprintf("%v %s ; %v", info.Mnemonic, instr.Operand, deviceInfo.Device)
		default:
			return fmt.Sprintf("Unknown opcode: %s, operand: %v", instr.Opcode, instr.Operand)
	}	
}

func DisassembleWord(word uint32) string {
	instr := DecodeInstruction(word)
	return DisassembleInstruction(instr)
}

func DisassembleProgram(program []uint32) []string {
	lines := []string{}
	for i, word := range program {
		address := i * int(WordSize)
		lines = append(
			lines, 
			fmt.Sprintf("%08X - %08X - %s", address, word, DisassembleWord(word)),)
	}
	return lines
}

func WriteDisasm(filename string, program []uint32) error {
	lines := DisassembleProgram(program)
	data := []byte(strings.Join(lines, "\n")) 
	data = append(data, '\n')
	return os.WriteFile(filename, data, 0644)
}
