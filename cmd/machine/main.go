package main

import (
	"cpu-simulator/internal/isa"
	"cpu-simulator/internal/machine"
	"fmt"
	"log"
	"os"
	"strconv"
)

const inputTickStep uint64 = 10

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
	events := buildInputEvents(program, inputTickStep)

	cpu := machine.NewCPUWithInputEvents(program, events)
	err = cpu.Run(limit)

	if err != nil {
		log.Fatalf("something went wrong: %v", err)
	}

	fmt.Print(cpu.Output())
}

func buildInputEvents(data []uint32, step uint64) []machine.InputEvent {
	events := make([]machine.InputEvent, 0, len(data))

	for i, b := range data {
		events = append(events, machine.InputEvent{
			Tick: uint64(i) * step,
			Value: int32(b),
		})
	}
	return events
}