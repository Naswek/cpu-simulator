package isa

const (
	WordSize uint32 = 4
	MemSize  uint32 = 2048
	DataStackSize uint32 = MemSize
	ReturnStackSize uint32 = 1500
	IOPortCount uint32 = 4 
	PORT_IN = 0
	PORT_OUT = 1
	InterruptVectorInput = 0
	InterruptVectorOutput = 1
)


type Register uint32
const (
	ACC Register = iota
	PC
	IR
	DSP
	RS
	SR
	AR
	DR
)



type Status uint8
const (
	Z Status = iota
	N
	C
	V
	IE
	IM
)

type Opcode uint8
type OpcodeKind uint8
const (
	none OpcodeKind = iota
	imm
	addr
	stack_offset
	port
)

type Operand uint32
const (
	HALT = 0x01
	NOP  = 0x02
	EI	 = 0x03
	DI	 = 0x04
	IRET = 0x05

	JMP = 0x10
	JZ	= 0x11
	JС  = 0x12
	JNС = 0x13
	JP = 0x14
	JN = 0x15
	JNZ = 0x16
	CALL = 0x17
	CALL_IND = 0x18
	RET = 0x19
	PUSH = 0x1A
	POP = 0x1B

	LD_IMM = 0x20
	LD_ADDR = 0x21
	LD_IND = 0x22
	LD_SP_N = 0x23
	ST_ADDR = 0x24
	ST_IND = 0x25
	ST_SP_N = 0x26

	ADD = 0x30
	SUB = 0x31
	MUL = 0x32
	MOD = 0x33
	DIV = 0x34
	INC = 0x35
	DEC = 0x36
	CMP = 0x37
	AND = 0x38
	OR = 0x39
	NOT = 0x3A

	IN = 0x40
	OUT = 0x41
)

type OpcodeInfo struct {
	Mnemonic 	string
	OperandKind OpcodeKind
	Ticks		int
}


var OpcodeTable = map[Opcode]OpcodeInfo {
	HALT : {"halt", none, },
	NOP	 : {"nop", none, },
	EI	 : {"ei", none, },
	DI   : {"di", none, },
	IRET : {"iret", none, },

	JMP  : {"jmp", addr, },
	JZ   : {"jz", addr, },
	JN 	 : {"jn", addr, },
	JNZ	 : {"jnz", addr, },
	CALL : {"call", addr, },
	JP	 : {"jp", addr, },
	JC	 : {"jc", addr, },
	CALL_IND : {"call_ind", none, },
	RET : {"ret", none, },

	PUSH : {"push", none, },
	POP  : {"pop", none, },

	LD_IMM  : {"ld_imm", imm, },
	LD_ADDR : {"ld_addr", addr, },
	LD_IND  : {"ld_ind", addr, },
	LD_SP_N : {"ld_sp_n", stack_offset, },
	ST_ADDR : {"st_addr", addr, },
	ST_IND  : {"st_ind", addr, },
	ST_SP_N : {"st_sp_n",imm, },

	ADD : {"add", addr, },
	SUB : {"sub", addr, },
	MUL : {"mul", addr, },
	MOD : {"mod", addr, },
	DIV : {"div", addr, },
	INC : {"inc", none, },
	DEC : {"dec", none, },
	CMP : {"cmp", addr, },
	AND : {"and", addr, },
	OR  : {"or", addr, },
	NOT : {"not", none, },

	IN  : {"in", port, },
	OUT : {"out", port, },
	
}

type Instruction struct {
	Opcode
	Operand
}



