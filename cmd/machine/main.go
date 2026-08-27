package main

import (
	"cpu-simulator/internal/isa"
	"cpu-simulator/internal/machine"
	"log"
	"os"
	"strconv"
)

func main() {
	if len(os.Args) != 3 {
		log.Fatalf("not enough args\nusage: machine <program.bin> <max_ticks>")
	}

	programFile := os.Args[1]
	limit, err := strconv.Atoi(os.Args[2])
	if err != nil {
		log.Fatalf("second arg must be number")
	}

	program, err := isa.ReadProgram(programFile)
	if err != nil {
		log.Fatalf("unexpected error: %v", err)
	}

	cpu := machine.NewCPU(program)
	err = cpu.Run(limit)

	if err != nil {
		log.Fatalf("something went wrong: %v", err)
	}
}
