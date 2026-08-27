package main

import (
	"cpu-simulator/internal/isa"
	"cpu-simulator/internal/translator"
	"log"
	"os"
	"path/filepath"
	"strings"
)

func main() {

	if len(os.Args) != 3 {
		log.Fatal("not enough args\nusage: translator <input.forth> <output.bin>")
	}

	inputFile := os.Args[1]
	outputFile := os.Args[2]

	bytes, err := os.ReadFile(inputFile)
	if err != nil {
		log.Fatal(err)
	}

	source := string(bytes)

	program, err := translator.Translate(source)
	if err != nil {
		log.Fatal(err)
	}

	image, err := program.MemoryImage()
	if err != nil {
		log.Fatal(err)
	}

	err = isa.WriteProgram(outputFile, image)
	if err != nil {
		log.Fatal(err)
	}

	code, err := program.CodeWords()
	if err != nil {
		log.Fatal(err)
	}

	err = isa.WriteDisasm(makeDisasm(outputFile), code)
	if err != nil {
		log.Fatalf("something went wrong with disasm log: %v", err)
	}
}

func makeDisasm(outputFile string) string {
	ext := filepath.Ext(outputFile)
	return strings.TrimSuffix(outputFile, ext) + ".disasm"
}
