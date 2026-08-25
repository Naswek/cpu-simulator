package isa

const (
	WordSize        uint32 = 4
	MemSize         uint32 = 2048
	DataStackSize   uint32 = MemSize
	ReturnStackSize uint32 = 1500
	PortInput       uint8  = 0
	PortOutput      uint8  = 1
)

type IODevice string

const (
	DeviceInput  IODevice = "input"
	DeviceOutput IODevice = "output"
)

type IODeviceInfo struct {
	Device          IODevice
	InterruptVector uint8
}

var PortsTable = map[uint8]IODeviceInfo{
	0: {Device: DeviceInput, InterruptVector: 0},
	1: {Device: DeviceOutput, InterruptVector: 1},
}

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
	EI   = 0x03
	DI   = 0x04
	IRET = 0x05

	JMP      = 0x10
	JZ       = 0x11
	JC       = 0x12
	JNC      = 0x13
	JP       = 0x14
	JN       = 0x15
	JNZ      = 0x16
	CALL     = 0x17
	CALL_IND = 0x18
	RET      = 0x19
	PUSH     = 0x1A
	POP      = 0x1B

	LD_IMM  = 0x20
	LD_ADDR = 0x21
	LD_IND  = 0x22
	LD_SP_N = 0x23
	ST_ADDR = 0x24
	ST_IND  = 0x25
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
	OR  = 0x39
	NOT = 0x3A

	IN  = 0x40
	OUT = 0x41
)

type OpcodeInfo struct {
	Mnemonic    string
	OperandKind OpcodeKind
	Ticks       int
}

var OpcodeTable = map[Opcode]OpcodeInfo{
	HALT: {"halt", none, 1},
	NOP:  {"nop", none, 2},
	EI:   {"ei", none, 1},
	DI:   {"di", none, 1},
	IRET: {"iret", none, 1},

	JMP:      {"jmp", addr, 1},
	JZ:       {"jz", addr, 1},
	JN:       {"jn", addr, 1},
	JNZ:      {"jnz", addr, 1},
	JNC:      {"jnc", addr, 1},
	CALL:     {"call", addr, 1},
	JP:       {"jp", addr, 1},
	JC:       {"jc", addr, 1},
	CALL_IND: {"call_ind", none, 1},
	RET:      {"ret", none, 1},

	PUSH: {"push", none, 1},
	POP:  {"pop", none, 1},

	LD_IMM:  {"ld_imm", imm, 1},
	LD_ADDR: {"ld_addr", addr, 1},
	LD_IND:  {"ld_ind", addr, 1},
	LD_SP_N: {"ld_sp_n", stack_offset, 1},
	ST_ADDR: {"st_addr", addr, 1},
	ST_IND:  {"st_ind", addr, 1},
	ST_SP_N: {"st_sp_n", stack_offset, 1},

	ADD: {"add", addr, 1},
	SUB: {"sub", addr, 1},
	MUL: {"mul", addr, 1},
	MOD: {"mod", addr, 1},
	DIV: {"div", addr, 1},
	INC: {"inc", none, 1},
	DEC: {"dec", none, 1},
	CMP: {"cmp", addr, 1},
	AND: {"and", addr, 1},
	OR:  {"or", addr, 1},
	NOT: {"not", none, 1},

	IN:  {"in", port, 1},
	OUT: {"out", port, 1},
}

type Instruction struct {
	Opcode
	Operand
}
