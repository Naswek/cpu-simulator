package tests

import (
	"cpu-simulator/internal/machine"
	"cpu-simulator/internal/translator"
	"os"
	"path/filepath"
	"testing"
)

const (
	inputStartTick uint64 = 1000
	inputTickStep  uint64 = 100
)

type goldenCase struct {
	name     string
	maxTicks int
	wantErr  string
}

func TestGoldenPrograms(t *testing.T) {
	cases := []goldenCase{
		{name: "hello", maxTicks: 10000},
		{name: "cat", maxTicks: 1500, wantErr: "tick limit exceeded"},
		{name: "hello_user_name", maxTicks: 20000},
		{name: "prob2", maxTicks: 50000},
		{name: "sort", maxTicks: 20000},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			source := readFile(t, path(tc.name, "source.forth"))
			input := readOptionalFile(t, path(tc.name, "input.txt"))

			program, err := translator.Translate(string(source))
			if err != nil {
				t.Fatalf("translate: %v", err)
			}

			image, err := program.MemoryImage()
			if err != nil {
				t.Fatalf("memory image: %v", err)
			}

			cpu := machine.NewCPUWithInputEvents(image, buildInputEvents(input))
			err = cpu.Run(tc.maxTicks)

			gotErr := ""
			if err != nil {
				gotErr = err.Error()
			}
			if gotErr != tc.wantErr {
				t.Fatalf("error = %q, want %q", gotErr, tc.wantErr)
			}

			compareFile(t, path(tc.name, "expected.output"), cpu.Output())
		})
	}
}

func readFile(t *testing.T, filePath string) []byte {
	t.Helper()

	data, err := os.ReadFile(filePath)
	if err != nil {
		t.Fatalf("read %s: %v", filePath, err)
	}
	return data
}

func readOptionalFile(t *testing.T, filePath string) []byte {
	t.Helper()

	data, err := os.ReadFile(filePath)
	if os.IsNotExist(err) {
		return nil
	}
	if err != nil {
		t.Fatalf("read %s: %v", filePath, err)
	}
	return data
}

func compareFile(t *testing.T, filePath string, actual string) {
	t.Helper()

	expected := readFile(t, filePath)
	if string(expected) != actual {
		t.Fatalf("%s mismatch\ngot:\n%q\nwant:\n%q", filePath, actual, string(expected))
	}
}

func buildInputEvents(input []byte) []machine.InputEvent {
	events := make([]machine.InputEvent, 0, len(input))
	for i, b := range input {
		events = append(events, machine.InputEvent{
			Tick:  inputStartTick + uint64(i)*inputTickStep,
			Value: int32(b),
		})
	}
	return events
}

func path(caseName string, fileName string) string {
	return filepath.Join("golden", caseName, fileName)
}
