package isa

import (
	"encoding/binary"
	"os"
	"fmt"
)

func ReadProgram(filename string) ([]uint32, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	if len(data)%4 != 0 {
		return nil, fmt.Errorf(
			"data size must be a multiple of 4 bytes: %d received",
			len(data),
		)
	}
	
	return fromByteToUint32(data), nil
}

func WriteProgram(filename string, data []uint32) error {
	bytes := fromUint32ToByte(data)
	err := os.WriteFile(filename, bytes, 0644)	
	if err != nil {
		return err
	}
	return nil
}

func EncodeProgram(instrs []Instruction) []uint32 {
	program := []uint32{}
	for _, instr := range instrs {
		program = append(program, EncodeInstruction(instr))
	}
	return program
}

func DecodeProgram(program []uint32) []Instruction {
	instrs := []Instruction{}
	for _, instr := range program {
		instrs = append(instrs, DecodeInstruction(instr))
	}
	return instrs
}


func EncodeInstruction(instr Instruction) uint32 {
	return uint32(instr.Opcode & 0xFF)<<24 | uint32(instr.Operand & 0xFFFFFF)
}

func DecodeInstruction(instr uint32) Instruction {
	return Instruction{
		Opcode:  Opcode(instr >> 24),
		Operand: Operand(instr & 0xFFFFFF),
	}
}

func fromUint32ToByte(data []uint32) []byte {
	result := make([]byte, len(data)*4)

	for i, word := range data {
		binary.BigEndian.PutUint32(result[i*4:(i+1)*4], word)
	}

	return result
}

func fromByteToUint32(data []byte) []uint32 {
	result := make([]uint32, len(data)/4)

	for i:= range result {
		result[i] = binary.BigEndian.Uint32(data[i*4 : (i+1)*4])
	}

	return result
}