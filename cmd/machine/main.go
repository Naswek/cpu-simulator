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
	if len(os.Args) != 5 {
		log.Fatalf("not enough args\nusage: machine <program.bin> <input.txt> <max_ticks> <log.txt>")
	}

	programFile := os.Args[1]
	inputFile := os.Args[2]
	limit, err := strconv.Atoi(os.Args[3])
	logFile := os.Args[4]
	if err != nil {
		log.Fatalf("max_ticks must be number")
	}

	program, err := isa.ReadProgram(programFile)
	if err != nil {
		log.Fatalf("unexpected error with reading .bin file.: %v", err)
	}

	data, err := os.ReadFile(inputFile)
	if err != nil {
		log.Fatalf("unexpected error with reading .txt file: %v", err)
	}
	events := buildInputEvents(data, inputTickStep)

	cpu := machine.NewCPUWithInputEvents(program, events)
	err = cpu.Run(limit)

	if logFile != "" {
		if errWrite := os.WriteFile(logFile, []byte(cpu.Log()), 0644); errWrite != nil {
			fmt.Printf("failed to save log file: %v", errWrite)
		}

	}

	fmt.Print(cpu.Output())

	if err != nil {
		log.Fatalf("something went wrong: %v", err)
	}
}

func buildInputEvents(data []byte, step uint64) []machine.InputEvent {
	events := make([]machine.InputEvent, 0, len(data))

	for i, b := range data {
		events = append(events, machine.InputEvent{
			Tick:  uint64(i) * step,
			Value: int32(b),
		})
	}
	return events
}
