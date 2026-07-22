package isa

type Opcode uint8
type Operand uint32
type Register string

const (
	WordSize uint32 = 4
	MemSize  uint32 = 2048
	DataStackSize uint32 = MemSize
	ReturnStackSize uint32 = 1500
	IOPortCount uint32 = 4 
)



const (
	HALT = 0x01
	NOP  = 0x02
	EI	 = 0x03
	DI	 = 0x04
	IRET = 0x05

	JMP = 0x10
	JZ	= 0x11
	JN  = 0x12
	JNZ = 0x13
	CALL = 0x14
	CALL_IND = 0x15
	RET = 0x16
	PUSH = 0x17
	POP = 0x18

	LD_IMM = 0x20
	LD_ADDR = 0x21
	LD_SP = 0x22
	LD_IND = 0x23
	LD_SP_N = 0x24
	ST_ADDR = 0x25
	ST_SP = 0x26
	ST_IND = 0x27
	ST_SP_N = 0x28

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
	HasOperant 	bool
	Ticks		int
}

type Operand struct {
	Value int
}

var OpcodeTable = map[Opcode]OpcodeInfo {
	HALT : {"halt", false, },
	NOP	 : {"nop", false, },
	EI	 : {"ei", false, },
	DI   : {"di", false, },
	IRET : {"iret", false, },

	JMP  : {"jmp", true, },
	JZ   : {"jz", true, },
	JN 	 : {"jn", true, },
	JNZ	 : {"jnz", true, },
	CALL : {"call", true, },
	CALL_IND : {"call_ind", false, },
	RET : {"ret", false, },

	PUSH : {"push", false, },
	POP  : {"pop", false, },

	LD_IMM  : {"ld_imm", true, },
	LD_ADDR : {"ld_addr", true, },
	LD_SP   : {"ld_sp", true, },
	LD_IND  : {"ld_ind", true, },
	LD_SP_N : {"ld_sp_n", true, },
	ST_ADDR : {"st_addr", true, },
	ST_SP   : {"st_sp", true, },
	ST_IND  : {"st_ind", true, },
	ST_SP_N : {"st_sp_n", true, },

	ADD : {"add", true, },
	SUB : {"sub", true, },
	MUL : {"mul", true, },
	MOD : {"mod", true, },
	DIV : {"div", true, },
	INC : {"inc", true, },
	DEC : {"dec", true, },
	CMP : {"cmp", true, },
	AND : {"and", true, },
	OR  : {"or", true, },
	NOT : {"not", true, },

	IN  : {"in", true, },
	OUT : {"out", true, },
	
}

var RegisterTable = map[Register]uint32 {
	
}


type Instruction struct {
	Opcode
	Operand
}

type IntVector struct {
	
}

